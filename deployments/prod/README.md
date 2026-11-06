# LabFlow Go - Production Deployment (labflow) - No reverse proxy

> Deploy dir on VPS is `~/deployment/labflow/` - **single directory** contains all deployment files: `docker-compose.yml`, `.env`, `backups/`. Stack is **app only** (postgres, redis, s3 (SeaweedFS), backend, frontend). No nginx/certbot - you handle TLS/reverse proxy yourself.

## Exposed Ports

- **Backend**: `127.0.0.1:8080` → `backend:8080` (`/api`, `/auth`, `/health`)
- **Frontend**: `127.0.0.1:3000` → `frontend:3000` (SSR app)

Bind is `127.0.0.1` for security. If your reverse proxy is on same host (recommended), it can proxy to `127.0.0.1:3000` and `127.0.0.1:8080`. If you want public direct access, change to `0.0.0.0:8080:8080` etc in `docker-compose.yml`.

**You need to configure your own reverse proxy** (example below) to map:
- `/api/` and `/auth/` → `http://127.0.0.1:8080`
- `/` → `http://127.0.0.1:3000`

---

## Quick Start (First time on VPS)

No git clone needed for deploys - CI builds images, deploy workflow SCPs compose file and updates VPS `.env`.

### Minimal setup (no git)

```bash
mkdir -p ~/deployment/labflow
cd ~/deployment/labflow
curl -fsSL https://raw.githubusercontent.com/<owner>/labflow/main/deployments/prod/docker-compose.yml -o docker-compose.yml
curl -fsSL https://raw.githubusercontent.com/<owner>/labflow/main/deployments/prod/.env.example -o .env.example
cp .env.example .env
nano .env
```

### Clone (first time convenience)

```bash
cd ~
git clone https://github.com/<owner>/labflow.git labflow
cd labflow
cp deployments/prod/docker-compose.yml ~/deployment/labflow/
cp deployments/prod/.env.example ~/deployment/labflow/
cd ~/deployment/labflow
cp .env.example .env
nano .env
```

**Must set:**
- `POSTGRES_PASSWORD`, `MINIO_ROOT_PASSWORD`
- `JWT_SECRET`, `COOKIE_SECRET` (`openssl rand -base64 32`)
- `MINIO_ROOT_USER` (object storage access key; use letters and digits only for it and `MINIO_ROOT_PASSWORD`)
- `BACKEND_IMAGE`, `FRONTEND_IMAGE` - set defaults `ghcr.io/<owner>/labflow-backend:latest`, deploy workflow will **update these values in VPS .env** to chosen tag
- `SMTP_*` optional

All required vars use `${VAR:?required}` - missing = fail fast.

### Start

```bash
docker compose up -d
docker compose ps                   # wait backend healthy
```

### Seed (first time)

```bash
docker compose exec backend ./main seed
```

### Access

- Direct: `http://VPS_IP:3000` (frontend), `http://VPS_IP:8080/health` (backend) if you bound to 0.0.0.0
- Via your reverse proxy: `https://your-domain.com`

---

## CI/CD Overview

- **CI** (`ci.yml`): Runs on PRs and non-main pushes - vet, build, test, typecheck, hadolint, compose validate (no image push)
- **Build Images** (`build.yml`): Runs on push to `main`, tags `v*`, `backend-v*`, `frontend-v*`, and manual dispatch. Builds backend/frontend with buildx cache and pushes to GHCR with tags:
  - `latest` (main only)
  - `${{ github.sha }}` (raw SHA, main only)
  - `sha-<short>`, `sha-<long>`
  - `ref` branch name (e.g., `feat-xyz`)
  - `ref` tag (e.g., `v1.0.0`, `backend-v1.2.0`) + semver
- **Deploy** (`deploy.yml`): **Manual only** via GitHub UI → Actions → Deploy → Run workflow
  - Inputs: `backend_tag`, `frontend_tag` (e.g., `latest`, `<sha>`, `v1.0.0`, branch name), `branch` for compose file source
  - SCPs `docker-compose.yml` to `~/deployment/labflow/` (single dir, all files inside)
  - **Updates `BACKEND_IMAGE` and `FRONTEND_IMAGE` in VPS `.env`** via `sed`
  - Pulls exact images, backs up DB to `~/deployment/labflow/backups/`, up -d, health wait, smoke tests, prune 72h

After first `.env`, no git needed on VPS. All deployment files live in `~/deployment/labflow/`.

## Deployment - Manual with tag selection

### Via GitHub UI (recommended)

1. Go to Actions → `Deploy` → Run workflow
2. Fill:
   - `backend_tag`: e.g., `latest`, or `abc123def...` (SHA), or `v1.0.0`, or `main`, or `feat-branch`
   - `frontend_tag`: same
   - `branch`: which branch to take `docker-compose.yml` from (default `main`)
3. Run → workflow will:
   - Set `BACKEND_IMAGE=ghcr.io/<owner>/labflow-backend:<backend_tag>`
   - Set `FRONTEND_IMAGE=ghcr.io/<owner>/labflow-frontend:<frontend_tag>`
   - SSH to VPS dir `~/deployment/labflow/`, **update `.env`** with those values, pull and redeploy

Example: deploy specific commit that CI built:
- Build on main produced tags `latest` + `<sha>` (e.g., `a1b2c3d...`)
- In Deploy workflow, set `backend_tag: a1b2c3d...`, `frontend_tag: a1b2c3d...` → VPS `.env` at `~/deployment/labflow/.env` updated to those SHA images

Example: rollback to previous version:
- Set `backend_tag: v1.0.0`, `frontend_tag: v1.0.0` → deploys tagged release to `~/deployment/labflow/`

### Manual without GitHub (direct SSH)

```bash
export VPS_USER=youruser VPS_HOST=YOUR_VPS_IP
export BACKEND_IMAGE=ghcr.io/<owner>/labflow-backend:<tag>
export FRONTEND_IMAGE=ghcr.io/<owner>/labflow-frontend:<tag>
scp deployments/prod/docker-compose.yml $VPS_USER@$VPS_HOST:~/deployment/labflow/
ssh $VPS_USER@$VPS_HOST "
  cd ~/deployment/labflow &&
  sed -i \"s|^BACKEND_IMAGE=.*|BACKEND_IMAGE=$BACKEND_IMAGE|\" .env &&
  sed -i \"s|^FRONTEND_IMAGE=.*|FRONTEND_IMAGE=$FRONTEND_IMAGE|\" .env &&
  docker compose pull && docker compose up -d
"
```

### Tagging Backend and Frontend Differently

You can build and deploy backend and frontend with **different versions**. `build.yml` now splits into `build-backend` and `build-frontend` jobs:

- Push `backend-v*` → builds **only backend**, skips frontend
- Push `frontend-v*` → builds **only frontend**, skips backend
- Push `v*` → builds **both** with same version
- Push to `main` → builds both with `latest` + SHA
- Manual dispatch allows different custom tags per image

**Git commands to tag separately:**

```bash
# Backend only release 1.2.0
git tag backend-v1.2.0
git push origin backend-v1.2.0
# Result: ghcr.io/<owner>/labflow-backend:backend-v1.2.0 + 1.2.0 + 1.2

# Frontend only release 2.0.0
git tag frontend-v2.0.0
git push origin frontend-v2.0.0
# Result: ghcr.io/<owner>/labflow-frontend:frontend-v2.0.0 + 2.0.0 + 2.0

# Both together same version
git tag v1.0.0
git push origin v1.0.0
# Result: both images get v1.0.0, 1.0

# Delete a tag if you made a mistake
git tag -d backend-v1.2.0
git push origin :refs/tags/backend-v1.2.0
```

**Manual build with different tags (recommended for independent versions):**

1. GitHub → Actions → **Build Images** → Run workflow
2. Inputs:
   - `backend_tag: v1.2.0` (or `backend-v1.2.0`)
   - `frontend_tag: v3.0.0` (or `frontend-v2.0.0`)
3. Run → pushes:
   - `ghcr.io/<owner>/labflow-backend:v1.2.0` + `latest` + SHA
   - `ghcr.io/<owner>/labflow-frontend:v3.0.0` + `latest` + SHA

**Deploy different tags:**

Actions → **Deploy** → Run workflow
- `backend_tag: backend-v1.2.0` or `1.2.0` or SHA `abc123`
- `frontend_tag: frontend-v2.0.0` or `2.0.0` or different SHA `def456`
- `branch: main` (compose file source)

This updates VPS `.env` to:
```env
BACKEND_IMAGE=ghcr.io/<owner>/labflow-backend:backend-v1.2.0
FRONTEND_IMAGE=ghcr.io/<owner>/labflow-frontend:frontend-v2.0.0
```

Check deployed images:
```bash
cat ~/deployment/labflow/.env | grep IMAGE
docker compose images
```

---

## Reverse Proxy Example (bring your own)

```nginx
upstream backend { server 127.0.0.1:8080; }
upstream frontend { server 127.0.0.1:3000; }

server {
  listen 80;
  server_name your-domain.com;

  location /api/ {
    proxy_pass http://backend;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
  }
  location /auth/ {
    proxy_pass http://backend;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
  }
  location / {
    proxy_pass http://frontend;
    proxy_set_header Host $host;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
  }
}
```

For TLS, `sudo certbot --nginx -d your-domain.com`.

---

## Useful Commands

```bash
docker compose logs -f backend
docker compose logs -f s3
docker compose ps
docker compose exec postgres psql -U labflow -d labflow
docker compose restart backend
docker compose down
docker compose down -v  # DANGER: deletes data
ls ~/deployment/labflow/backups/
docker compose exec backend ls -la /app/sql/migrations
curl -f http://127.0.0.1:8080/health
curl -f http://127.0.0.1:3000/
cat .env | grep IMAGE  # check current deployed images
```

## Improvements implemented

- Dockerfiles: cache mounts, ldflags, wget healthcheck, non-root frontend, dumb-init
- .dockerignore: fixed *.sql bug
- Compose: env_file, required vars, logging rotation, resource limits, SeaweedFS object storage (buckets created by the backend), healthchecks, ports 8080/3000 on 127.0.0.1, no nginx/certbot
- CI: ci.yml PR checks, build.yml auto builds on main/tags, deploy.yml manual with tag selection and VPS .env update, gha cache, SHA tags, SCP only compose
- dependabot weekly
