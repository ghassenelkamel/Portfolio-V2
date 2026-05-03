# Portfolio V2

Portfolio web application with a Go backend and a SvelteKit frontend, designed with Docker Compose.

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

### API service

- `ADDR`: API listen address (example: `:8080`)
- `GITHUB_USER`: GitHub username used by API integrations
- `GITHUB_TOKEN`: Optional GitHub token (helps with rate limits)
- `CONTENT_REPO_OWNER`: Owner of the portfolio content repository
- `CONTENT_REPO_NAME`: Repository name for structured site content
- `CONTENT_REPO_REF`: Branch/tag/commit to read (example: `main`)
- `CONTENT_REPO_DIR`: Folder in content repo containing JSON files (example: `content`)
- `CONTENT_CACHE_TTL`: API cache TTL in seconds for content fetches
- `CONTACT_FORWARD_URL`: Optional endpoint to forward contact submissions

### Web service

- `VITE_API_PROXY_TARGET`: API target used by the frontend proxy (compose default usually points to `http://api:8080`)

## Running Without Docker (Optional)

### Backend

```bash
cd backend
go run ./cmd/portfolio
```

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

## License

This project is licensed under the MIT License.
