// LabFlow API load test.
//
// Three scenarios run in parallel, modelled on how the lab is actually used:
//   browse - customers and staff reading dashboards (most traffic)
//   submit - customers creating new test requests
//   login  - sign-ins (bcrypt-bound, so measured on its own)
//
//   k6 run loadtest/k6/lifecycle.js
//   BASE_URL=http://host:8080 RATE=200 k6 run loadtest/k6/lifecycle.js
//
// Needs a seeded database (`main seed`) and DISABLE_RATE_LIMIT=true on the backend.
import http from 'k6/http';
import { check, fail } from 'k6';
import { Counter } from 'k6/metrics';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const PASSWORD = __ENV.SEED_PASSWORD || 'password123';
const RATE = Number(__ENV.RATE || 150); // browse iterations per second at peak
const DURATION = __ENV.DURATION || '2m';

const created = new Counter('reservations_created');
const JSON_HEADERS = { 'Content-Type': 'application/json' };

export const options = {
  scenarios: {
    browse: {
      executor: 'ramping-arrival-rate',
      exec: 'browse',
      startRate: 10,
      timeUnit: '1s',
      preAllocatedVUs: 50,
      maxVUs: 300,
      stages: [
        { target: RATE, duration: '30s' },
        { target: RATE, duration: DURATION },
        { target: 0, duration: '10s' },
      ],
    },
    submit: {
      executor: 'constant-arrival-rate',
      exec: 'submit',
      rate: Math.max(1, Math.round(RATE / 15)),
      timeUnit: '1s',
      duration: DURATION,
      startTime: '30s',
      preAllocatedVUs: 20,
      maxVUs: 100,
    },
    login: {
      executor: 'constant-arrival-rate',
      exec: 'login',
      rate: 5,
      timeUnit: '1s',
      duration: DURATION,
      startTime: '30s',
      preAllocatedVUs: 10,
      maxVUs: 50,
    },
  },
  thresholds: {
    'http_req_failed': ['rate<0.01'],
    'http_req_duration{scenario:browse}': ['p(95)<100', 'p(99)<250'],
    'http_req_duration{scenario:submit}': ['p(95)<250'],
    'http_req_duration{scenario:login}': ['p(95)<800'],
    'checks': ['rate>0.99'],
  },
  summaryTrendStats: ['avg', 'med', 'p(90)', 'p(95)', 'p(99)', 'max'],
};

function signIn(identifier) {
  const res = http.post(
    `${BASE_URL}/auth/login`,
    JSON.stringify({ identifier, password: PASSWORD, turnstileToken: 'loadtest' }),
    { headers: JSON_HEADERS, tags: { name: 'POST /auth/login' } },
  );
  const ok = check(res, { 'login accepted': (r) => r.status === 202 || r.status === 200 });
  const cookie = res.cookies.Authorization && res.cookies.Authorization[0];
  return ok && cookie ? cookie.value : null;
}

// Log in once per role up front; scenarios reuse the session cookies.
export function setup() {
  const sessions = {};
  for (const user of ['customer', 'admin', 'qa']) {
    const token = signIn(user);
    if (!token) fail(`could not sign in seeded user "${user}" - is the database seeded?`);
    sessions[user] = token;
  }
  return sessions;
}

function authed(token, name) {
  return { headers: { ...JSON_HEADERS, Cookie: `Authorization=${token}` }, tags: { name } };
}

export function browse(sessions) {
  const role = ['customer', 'customer', 'admin', 'qa'][__ITER % 4];
  const token = sessions[role];
  const category = 1 + (__ITER % 3);
  const responses = http.batch([
    ['GET', `${BASE_URL}/api/auth/me`, null, authed(token, 'GET /api/auth/me')],
    ['GET', `${BASE_URL}/api/reservations?page=1&limit=10`, null, authed(token, 'GET /api/reservations')],
    ['GET', `${BASE_URL}/api/reservations/stats`, null, authed(token, 'GET /api/reservations/stats')],
    ['GET', `${BASE_URL}/api/tests/catalog/${category}`, null, authed(token, 'GET /api/tests/catalog/:id')],
  ]);
  for (const r of responses) {
    check(r, { 'browse 200': (res) => res.status === 200 });
  }
}

export function submit(sessions) {
  const body = JSON.stringify({
    sample_name: `Load sample ${__VU}-${__ITER}`,
    sample_amount: 3,
    special_instructions: 'k6 load test',
    sample_type: 'Liquid',
    sample_taken_date: new Date().toISOString(),
    sample_taken_by: 'k6',
    tests: [
      { test_list_id: 1, method: 'USP <61>', criteria: '< 100 CFU/g', test_unit: 'CFU/g' },
      { test_list_id: 11, method: 'USP <791>', criteria: '6.5 - 7.5', test_unit: 'pH' },
    ],
  });
  const res = http.post(`${BASE_URL}/api/reservations`, body, authed(sessions.customer, 'POST /api/reservations'));
  if (check(res, { 'reservation created': (r) => r.status === 201 || r.status === 200 })) {
    created.add(1);
  } else {
    console.warn(`submit failed: ${res.status} ${String(res.body).slice(0, 200)}`);
  }
}

// The login scenario signs in as accounts the browse/submit scenarios don't use,
// so its sessions never interfere with theirs.
const LOGIN_USERS = ['micro_analyst', 'chem_analyst', 'micro_qc', 'chem_qc', 'lab_head', 'mr'];

export function login() {
  signIn(LOGIN_USERS[__ITER % LOGIN_USERS.length]);
}

export function handleSummary(data) {
  const out = __ENV.SUMMARY_PATH || 'loadtest/results/summary.json';
  return {
    [out]: JSON.stringify(data, null, 2),
    stdout: textSummary(data),
  };
}

function textSummary(data) {
  const m = data.metrics;
  const fmt = (v) => (v === undefined ? '-' : `${v.toFixed(1)}ms`);
  const lines = ['', 'scenario  p50       p95       p99'];
  for (const s of ['browse', 'submit', 'login']) {
    const t = m[`http_req_duration{scenario:${s}}`];
    if (t) lines.push(`${s.padEnd(9)} ${fmt(t.values.med).padEnd(9)} ${fmt(t.values['p(95)']).padEnd(9)} ${fmt(t.values['p(99)'])}`);
  }
  lines.push(`requests: ${m.http_reqs.values.count} (${m.http_reqs.values.rate.toFixed(0)}/s)  failed: ${(m.http_req_failed.values.rate * 100).toFixed(2)}%`);
  lines.push(`reservations created: ${m.reservations_created ? m.reservations_created.values.count : 0}`);
  const failed = Object.entries(m).filter(([, v]) => v.thresholds && Object.values(v.thresholds).some((t) => !t.ok)).map(([k]) => k);
  lines.push(failed.length ? `THRESHOLDS FAILED: ${failed.join(', ')}` : 'all thresholds passed');
  return lines.join('\n') + '\n';
}
