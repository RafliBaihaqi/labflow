"""Compare two running builds of the API response by response.

    python3 loadtest/bench/parity.py http://127.0.0.1:18080 http://127.0.0.1:18081

Both builds must share one seeded database and Redis (so a session cookie from
one is valid on the other). Reads must return identical JSON. Writes are sent
to each build separately and compared after masking generated values
(ids, xids, timestamps, request numbers).
"""

import json
import re
import sys
import urllib.error
import urllib.request

OLD, NEW = sys.argv[1], sys.argv[2]
PASSWORD = "password123"
MASK = re.compile(r"(^|_)(id|xid|created_at|updated_at|reservation_date|test_request_number|approved_date)$")


def call(base, method, path, token=None, body=None):
    data = json.dumps(body).encode() if body is not None else None
    req = urllib.request.Request(base + path, data=data, method=method)
    req.add_header("Content-Type", "application/json")
    if token:
        req.add_header("Cookie", f"Authorization={token}")
    try:
        with urllib.request.urlopen(req) as resp:
            return resp.status, resp.headers, json.loads(resp.read() or b"null")
    except urllib.error.HTTPError as err:
        return err.code, err.headers, json.loads(err.read() or b"null")


def login(base, user):
    status, headers, _ = call(base, "POST", "/auth/login",
                              body={"identifier": user, "password": PASSWORD, "turnstileToken": "x"})
    assert status in (200, 202), f"login {user} on {base}: {status}"
    cookie = headers.get("Set-Cookie", "")
    return re.search(r"Authorization=([^;]+)", cookie).group(1)


def normalized(value, key=""):
    # GetMe builds permissions from a Go map, so their order is random in every build.
    if isinstance(value, dict):
        return {k: normalized(v, k) for k, v in value.items()}
    if isinstance(value, list):
        items = [normalized(v, key) for v in value]
        return sorted(items) if key == "permissions" else items
    return value


def masked(value, key=""):
    if isinstance(value, dict):
        return {k: masked(v, k) for k, v in value.items()}
    if isinstance(value, list):
        return [masked(v, key) for v in value]
    if MASK.search(key) and value not in (None, 0, ""):
        return "<generated>"
    return value


failures = []


def compare(label, a, b, mask=False):
    a, b = (a[0], normalized(a[1])), (b[0], normalized(b[1]))
    if mask:
        a, b = (a[0], masked(a[1])), (b[0], masked(b[1]))
    if a != b:
        failures.append(label)
        print(f"DIFF {label}\n  old: {json.dumps(a)[:600]}\n  new: {json.dumps(b)[:600]}")
    else:
        print(f"same {label}")


tokens = {user: login(OLD, user) for user in ("customer", "admin", "qa")}

reservation = {
    "sample_name": "Parity sample", "sample_amount": 3, "special_instructions": "parity",
    "sample_type": "Liquid", "sample_taken_date": "2026-11-09T01:00:00Z", "sample_taken_by": "parity",
    "tests": [
        {"test_list_id": 1, "method": "USP <61>", "criteria": "< 100 CFU/g", "test_unit": "CFU/g"},
        {"test_list_id": 11, "method": "USP <791>", "criteria": "6.5 - 7.5", "test_unit": "pH"},
    ],
}
# Distinct names so that sorting by sample_name has no ties.
created = {}
for name, base in (("old", OLD), ("new", NEW)):
    status, _, body = call(base, "POST", "/api/reservations", tokens["customer"],
                           {**reservation, "sample_name": f"Parity sample {name}"})
    created[name] = (status, {**body, "sample_name": "<per build>"})
compare("POST /api/reservations", created["old"], created["new"], mask=True)
xid = created["old"][1]["xid"]

transitions = {}
for name, base in (("old", OLD), ("new", NEW)):
    target = created[name][1]["xid"]
    transitions[name] = call(base, "PUT", f"/api/reservations/{target}/status", tokens["qa"],
                             {"status": "rejected"})[::2]
compare("PUT /api/reservations/{xid}/status", transitions["old"], transitions["new"])

reads = [
    ("customer", "/api/auth/me"), ("admin", "/api/auth/me"), ("qa", "/api/auth/me"),
    ("customer", "/api/reservations?page=1&limit=10"),
    ("admin", "/api/reservations?page=1&limit=5&sort_by=sample_name&sort_order=asc"),
    ("admin", "/api/reservations?page=2&limit=1"),
    ("admin", "/api/reservations?search=Parity&status=rejected"),
    ("admin", "/api/reservations?status=pending&status=rejected"),
    ("admin", "/api/reservations/stats"),
    ("admin", f"/api/reservations/{xid}"),
    ("admin", "/api/reservations/00000000-0000-0000-0000-000000000000"),
    ("admin", "/api/tests/catalog/1"), ("admin", "/api/tests/catalog/2"), ("admin", "/api/tests/catalog/3"),
    ("admin", "/api/test-lists?page=1&limit=20"), ("admin", "/api/test-lists?search=pH&category_id=2"),
    ("admin", "/api/quotations"), ("customer", "/api/quotations"),
    ("admin", "/api/lifecycle/transitions"),
]
for user, path in reads:
    compare(f"GET {path} as {user}", call(OLD, "GET", path, tokens[user])[::2],
            call(NEW, "GET", path, tokens[user])[::2])

print(f"\n{len(failures)} difference(s)")
sys.exit(1 if failures else 0)
