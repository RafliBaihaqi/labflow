# LabFlow Go - Dev Infrastructure (shared VPS)

Runs Postgres, Redis, and S3-compatible object storage (SeaweedFS) in Docker for
development. These can run on the **same VPS as the prod stack** without
conflict: dev uses `labflow-dev-*` container names, `dev_*` volumes and its own
network. Ports bind to `127.0.0.1` on the VPS.

## Before you start: stop the old dev containers

If you previously ran dev services that reuse prod's names/ports, stop them
first so they don't clash with the prod stack:

```bash
docker ps
# stop/remove any old dev containers using prod names (labflow-postgres, etc.)
```

## Start

```bash
cd deployments/dev
docker compose up -d
docker compose ps                    # all should be healthy
```

The backend creates the buckets `certificates`, `reservations` and `users` on
startup. Storage credentials come from `MINIO_ROOT_USER` / `MINIO_ROOT_PASSWORD`
(default `minioadmin` / `minioadmin`).

## Connect from your laptop (SSH tunnel)

Because the ports bind to `127.0.0.1` on the VPS, tunnel them over SSH:

```bash
ssh -L 5432:127.0.0.1:5432 \
    -L 6379:127.0.0.1:6379 \
    -L 9000:127.0.0.1:9000 \
    user@YOUR_VPS_IP
```

Then point `kernel/.env` at localhost (traffic flows through the tunnel):

```dotenv
DATABASE_URL="postgres://postgres:database@localhost:5432/labflow"
MIGRATIONS_URL="postgres://postgres:database@localhost:5432/labflow?sslmode=disable"
REDIS_URL="redis://localhost:6379/0"
MINIO_ENDPOINT="localhost:9000"
MINIO_ACCESS="minioadmin"
MINIO_SECRET="minioadmin"
MINIO_BUCKET_CERTIFICATE=certificates
MINIO_BUCKET_RESERVATION=reservations
MINIO_BUCKET_USER=users
JWT_SECRET=dev-jwt-secret-min-32-chars-long
COOKIE_SECRET=dev-cookie-secret-min-32-chars-long

# Turnstile - dev bypass stays true (prod must NOT bypass)
ENV=development
BYPASS_TURNSTILE=true
# For dev you can use Cloudflare test keys that always pass:
TURNSTILE_SECRET=1x0000000000000000000000000000000AA
```

And for frontend dev, create `app/.env`:

```dotenv
VITE_TURNSTILE_SITE_KEY=1x00000000000000000000AA
```

Test keys from Cloudflare docs: sitekey `1x00000000000000000000AA` always passes, secret `1x0000000000000000000000000000000AA`.

> If you cannot use a tunnel, change the dev `ports` from `127.0.0.1:` to
> `0.0.0.0:` and restrict access with a VPS firewall to your IP only. Do not
> leave these ports open to the internet.

For frontend dev, set `BACKEND_URL=http://localhost:8080` or run backend locally.

## Useful commands

```bash
docker compose ps              # status + health
docker compose logs -f         # logs all
docker compose logs -f s3
docker compose down            # stop (keeps data)
docker compose down -v         # stop + delete all dev data
docker compose exec postgres psql -U postgres -d labflow
```

## Isolation from prod

| Resource        | Dev                  | Prod (`deployments/prod`) |
|-----------------|----------------------|---------------------------|
| Container names | `labflow-dev-*`      | `labflow-*`               |
| Volumes         | `dev_*`              | `labflow_*_data`          |
| DB/cache ports  | 5432 / 6379 / 9000   | not published             |
| Bucket creation | backend on startup   | backend on startup        |
| nginx           | not used             | 80/443                    |

Both stacks can run on the same VPS at the same time.

## Improvements applied

- Replaced MinIO with SeaweedFS (`chrislusf/seaweedfs`), because the `minio/minio`
  and `minio/mc` images are no longer published. The backend creates buckets itself.
- Vite proxy now uses `process.env.BACKEND_URL || localhost:8080` instead of hardcoded WSL IP.
