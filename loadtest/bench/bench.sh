#!/usr/bin/env bash
# Compare backend builds under the same k6 load on the same machine.
#
#   loadtest/bench/bench.sh [-r rounds] [-w warmup-rounds] [-R rate] [-d duration] label=git-ref ...
#   loadtest/bench/bench.sh -w 1 -r 4 gin-gorm=main nethttp=HEAD
#
# Each variant is built from kernel/ at its git ref. Then, for every round and
# every variant (interleaved, so machine drift affects all variants equally),
# after any warm-up rounds, which run the same way but are not summarized:
#   1. reset Postgres and Redis, then migrate and seed
#   2. start the backend with fixed CPU and memory limits
#   3. warm up, then run loadtest/k6/lifecycle.js
#   4. sample `docker stats` for the backend container during the run
# Results go to loadtest/results/bench/<run-id>/. Summarize them with
#   python3 loadtest/bench/summarize.py loadtest/results/bench/<run-id>
set -euo pipefail

ROUNDS=3
WARMUP=0
RATE=100
DURATION=2m
while getopts "r:w:R:d:" opt; do
  case $opt in
    r) ROUNDS=$OPTARG ;;
    w) WARMUP=$OPTARG ;;
    R) RATE=$OPTARG ;;
    d) DURATION=$OPTARG ;;
    *) exit 2 ;;
  esac
done
shift $((OPTIND - 1))
if [ $# -lt 1 ]; then
  echo "usage: $0 [-r rounds] [-w warmup-rounds] [-R rate] [-d duration] label=git-ref ..." >&2
  exit 2
fi

ROOT=$(git -c safe.directory="*" rev-parse --show-toplevel)
PROJECT=lfbench
COMPOSE=(docker compose -p "$PROJECT" -f "$ROOT/deployments/loadtest/docker-compose.yml")
NET=${PROJECT}_default
API=${PROJECT}-api
K6_IMAGE=grafana/k6:2.3.0
API_CPUS=2
API_MEMORY=512m
RUN_ID=$(date +%Y%m%d-%H%M%S)
OUT=$ROOT/loadtest/results/bench/$RUN_ID
mkdir -p "$OUT"

# Same settings as the backend service in deployments/loadtest/docker-compose.yml.
API_ENV=(
  -e ENV=loadtest
  -e "DATABASE_URL=postgres://labflow:labflow@postgres:5432/labflow?sslmode=disable"
  -e "MIGRATIONS_URL=postgres://labflow:labflow@postgres:5432/labflow?sslmode=disable"
  -e REDIS_URL=redis://redis:6379/0
  -e MINIO_ENDPOINT=s3:9000
  -e MINIO_ACCESS=minioadmin
  -e MINIO_SECRET=minioadmin
  -e JWT_SECRET=loadtest-jwt-secret-not-for-production-use
  -e COOKIE_SECRET=loadtest-cookie-secret-not-for-production
  -e BYPASS_TURNSTILE=true
  -e DISABLE_RATE_LIMIT=true
  -e SMTP_HOST=mailpit
  -e SMTP_PORT=1025
)

log() { printf '[%s] %s\n' "$(date +%H:%M:%S)" "$*"; }

cleanup() { docker rm -f "$API" >/dev/null 2>&1 || true; }
trap cleanup EXIT

log "starting infrastructure ($PROJECT)"
"${COMPOSE[@]}" up -d --wait postgres redis s3 mailpit >/dev/null

LABELS=()
{
  echo "{"
  echo "  \"run_id\": \"$RUN_ID\", \"rounds\": $ROUNDS, \"warmup_rounds\": $WARMUP, \"rate\": $RATE, \"duration\": \"$DURATION\","
  echo "  \"api_cpus\": $API_CPUS, \"api_memory\": \"$API_MEMORY\", \"k6_image\": \"$K6_IMAGE\","
  echo "  \"host\": \"$(nproc) threads, $(free -m | awk '/Mem:/ {print $2}') MiB RAM, $(uname -r)\","
  echo "  \"variants\": {"
} >"$OUT/meta.json"
first=1
for spec in "$@"; do
  label=${spec%%=*}
  ref=${spec#*=}
  sha=$(git -c safe.directory="*" -C "$ROOT" rev-parse --short "$ref")
  image=labflow-bench:$label
  log "building $label from $ref ($sha)"
  git -c safe.directory="*" -C "$ROOT" archive --format=tar "$sha:kernel" | docker build --network host -q -t "$image" - >/dev/null
  image_bytes=$(docker image inspect --format '{{.Size}}' "$image")
  binary_bytes=$(docker run --rm --entrypoint stat "$image" -c %s /app/main)
  [ $first -eq 1 ] || echo "," >>"$OUT/meta.json"
  first=0
  printf '    "%s": {"ref": "%s", "sha": "%s", "image_bytes": %s, "binary_bytes": %s}' \
    "$label" "$ref" "$sha" "$image_bytes" "$binary_bytes" >>"$OUT/meta.json"
  LABELS+=("$label")
done
printf '\n  }\n}\n' >>"$OUT/meta.json"

run_k6() {
  local rate=$1 duration=$2 summary=$3
  docker run --rm --network host -v "$ROOT/loadtest:/loadtest" \
    -e BASE_URL=http://127.0.0.1:8080 -e RATE="$rate" -e DURATION="$duration" \
    -e SUMMARY_PATH="$summary" "$K6_IMAGE" run --quiet /loadtest/k6/lifecycle.js
}

# Warm-up rounds are named w<n> so summarize.py ignores them.
ROUND_IDS=()
for i in $(seq 1 "$WARMUP"); do ROUND_IDS+=("w$i"); done
for i in $(seq 1 "$ROUNDS"); do ROUND_IDS+=("r$i"); done

for round in "${ROUND_IDS[@]}"; do
  for label in "${LABELS[@]}"; do
    image=labflow-bench:$label
    prefix=$OUT/$label-$round
    log "round $round: $label"

    docker exec "${PROJECT}-postgres-1" psql -q -U labflow -d labflow \
      -c 'SET client_min_messages TO warning; DROP SCHEMA public CASCADE; CREATE SCHEMA public;' >/dev/null
    docker exec "${PROJECT}-redis-1" redis-cli FLUSHALL >/dev/null
    docker run --rm --network "$NET" "${API_ENV[@]}" "$image" seed >"$prefix-seed.log" 2>&1

    docker run -d --name "$API" --network "$NET" -p 127.0.0.1:8080:8080 \
      --cpus "$API_CPUS" --memory "$API_MEMORY" "${API_ENV[@]}" "$image" >/dev/null
    for _ in $(seq 1 60); do
      curl -sf http://127.0.0.1:8080/health >/dev/null && break
      sleep 1
    done
    curl -sf http://127.0.0.1:8080/health >/dev/null

    run_k6 20 10s /tmp/warmup.json >/dev/null 2>&1 || true

    (
      while true; do
        docker stats --no-stream --format '{{.CPUPerc}};{{.MemUsage}}' "$API" 2>/dev/null || true
        sleep 1
      done
    ) >"$prefix-stats.csv" &
    sampler=$!
    run_k6 "$RATE" "$DURATION" "/loadtest/results/bench/$RUN_ID/$label-$round.json" \
      >"$prefix-k6.txt" 2>&1 || true
    kill "$sampler" 2>/dev/null || true
    wait "$sampler" 2>/dev/null || true

    docker logs --tail 50 "$API" >"$prefix-api.log" 2>&1
    docker rm -f "$API" >/dev/null
    tail -n 6 "$prefix-k6.txt"
  done
done

log "done: $OUT"
python3 "$ROOT/loadtest/bench/summarize.py" "$OUT"
