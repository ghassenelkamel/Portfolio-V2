# ---------- Frontend build ----------
FROM oven/bun:1.3.6 AS web
WORKDIR /app
COPY frontend/package.json frontend/bun.lock ./frontend/
COPY frontend/svelte.config.js frontend/vite.config.ts frontend/tsconfig.json ./frontend/
COPY frontend/src ./frontend/src
COPY frontend/static ./frontend/static
RUN cd frontend && bun install --frozen-lockfile && bun run build

# ---------- Backend build ----------
FROM golang:1.25.5-bookworm AS api
WORKDIR /app
COPY backend/go.mod ./backend/
COPY backend/cmd ./backend/cmd
COPY backend/internal ./backend/internal
RUN cd backend && go build -trimpath -ldflags="-s -w" -o /out/portfolio ./cmd/portfolio

# ---------- Runtime ----------
FROM debian:bookworm-slim
WORKDIR /app
RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates \
  && rm -rf /var/lib/apt/lists/* \
  && useradd -r -u 10001 appuser
COPY --from=api /out/portfolio /app/portfolio
COPY --from=web /app/frontend/build /app/web
ENV ADDR=:8080
ENV WEB_DIR=/app/web
USER appuser
EXPOSE 8080
CMD ["/app/portfolio"]
