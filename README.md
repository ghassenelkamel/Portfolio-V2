# Portfolio V2

Portfolio web application with a Go backend and a SvelteKit frontend, designed with Docker Compose.

## Demo

Short walkthrough video of the portfolio experience.

[![Portfolio Demo](docs/media/portfolio-terminal.gif)](https://github.com/user-attachments/assets/0d171656-f655-4275-8379-22b1d3b5b65a)

[Download demo file](https://github.com/ghassenelkamel/Portfolio-V2/raw/refs/heads/main/docs/media/portfolio-demo.mp4)

## Stack

- Backend: Go
- Frontend: SvelteKit + Vite + TypeScript
- Styling: Tailwind CSS
- Containers: Docker + Docker Compose
- Package managers: Go modules (backend), Bun (frontend)

## Project Structure

```text
.
├── backend/                # Go API service
├── frontend/               # SvelteKit app
├── docker-compose.dev.yml  # Development compose setup
├── docker-compose.yml      # Compose setup for non-dev runs
└── README.md
```

## Prerequisites

- Docker Engine
- Docker Compose (v2)
- Node.js 24+ (for local frontend tooling outside Docker)
- Bun 1.3+ (recommended for frontend commands)

## Run Locally (Development)

1. Clone the repository.
2. Enter the project directory:

```bash
cd <repo-directory>
```

3. Start development services:

```bash
docker compose -f docker-compose.dev.yml up
```

4. Open:

- Frontend: http://localhost:5173
- Backend API: http://localhost:8080

## Stop Local Environment

If running in the foreground, press `Ctrl+C`.

Or stop and remove containers from another terminal:

```bash
docker compose -f docker-compose.dev.yml down
```

## Useful Docker Commands

Start in background:

```bash
docker compose -f docker-compose.dev.yml up -d
```

View logs:

```bash
docker compose -f docker-compose.dev.yml logs -f
```

Rebuild after Dockerfile/dependency changes:

```bash
docker compose -f docker-compose.dev.yml up --build
```

Reset containers + anonymous volumes:

```bash
docker compose -f docker-compose.dev.yml down -v
```

## Environment Variables

Configure values in `docker-compose.dev.yml` (or via your environment):

```bash
cp .env.example .env
```

Then set your real SMTP credentials in `.env`.

### API service

- `ADDR`: API listen address (example: `:8080`)
- `LOG_LEVEL`: Backend log level (`debug`, `info`, `warn`, `error`; default `info`)
- `GITHUB_USER`: GitHub username used by API integrations
- `GITHUB_TOKEN`: Optional GitHub token (helps with rate limits)
- `CONTENT_REPO_OWNER`: Owner of the portfolio content repository
- `CONTENT_REPO_NAME`: Repository name for structured site content
- `CONTENT_REPO_REF`: Branch/tag/commit to read (example: `main`)
- `CONTENT_REPO_DIR`: Folder in content repo containing JSON files (example: `content`)
- `CONTENT_CACHE_TTL`: API cache TTL in seconds for content fetches
- `SMTP_HOST`: SMTP server hostname for direct email delivery
- `SMTP_PORT`: SMTP server port (default: `587`)
- `SMTP_USER`: SMTP username/login
- `SMTP_PASS`: SMTP password/app password
- `SMTP_FROM`: Sender address used by backend (example: `no-reply@yourdomain.com`)
- `SMTP_TO`: Recipient address where contact form messages are delivered
- `SMTP_REQUIRE_TLS`: Require STARTTLS before auth/send (`true` by default)
- `CONTACT_FORWARD_URL`: Optional fallback endpoint (used only when SMTP is not configured)

### Web service

- `VITE_API_PROXY_TARGET`: API target used by the frontend proxy (compose default usually points to `http://api:8080`)

## Log Retention

- Docker Compose services use rotating container logs (`max-size` + `max-file`) so old logs are automatically cleaned.
- Follow logs live with `docker compose -f docker-compose.dev.yml logs -f`.

## Running Without Docker (Optional)

### Backend

```bash
cd backend
go run ./cmd/portfolio
```

The backend auto-loads `.env` when present (from `backend/.env` or repo-root `../.env`) unless variables are already set in the process environment.

### Frontend

```bash
cd frontend
bun install
bun dev
```

## Troubleshooting

- Port conflict: ensure `5173` and `8080` are free.
- Build/cache issues: run `docker compose -f docker-compose.dev.yml down -v` then start again.
- If one service fails, inspect logs with `docker compose -f docker-compose.dev.yml logs -f <service-name>`.

## Dynamic Content Source

The About, Skills, Experience, and Contact pages are loaded from:

- `/api/content/about`
- `/api/content/skills`
- `/api/content/experience`
- `/api/content/contact`

Back-end content endpoints pull JSON from your public content repository (`CONTENT_REPO_*` variables) and keep a last-known-good cache fallback.

## Contact Form Delivery

`POST /api/contact` supports two delivery modes:

1. SMTP (preferred): set `SMTP_HOST`, `SMTP_FROM`, and `SMTP_TO` (plus auth vars if required)
2. Forwarding fallback: set `CONTACT_FORWARD_URL` if SMTP is not configured

If neither SMTP nor forwarding is configured, backend accepts the submission but does not send it (dev-friendly no-op).

## License

This project is licensed under the MIT License.
