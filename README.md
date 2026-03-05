# shop-api

Commands:

1. update swagger
```bash
swag init -g cmd/main.go -o internal/docs
```

## Release (Docker)

Build and run release stack (app + postgres):
```bash
docker compose up --build -d
```

Stop release stack:
```bash
docker compose down
```

Release files:
- `src/Dockerfile` - production image (multi-stage, distroless)
- `docker-compose.yml` - release compose (without debug)

## Debug in Docker (VS Code)

1. Open this repository in VS Code.
2. Start debug config `Go: Attach to Docker (auto up/down)` from Run and Debug.
3. Set breakpoints in files under `src/`.

What is used:
- `.vscode/tasks.json` runs `docker compose -f docker-compose.debug.yml up --build -d` and `down`.
- `.vscode/launch.json` attaches to Delve on port `2345` using `debugAdapter: legacy` (compatible with `dlv debug --headless`).
- `docker-compose.debug.yml` starts the debug container with `src` mounted to `/app` and also starts `postgres`.

Manual mode:
```bash
docker compose -f docker-compose.debug.yml up --build -d
```
Then run debug config `Go: Attach to Docker (manual)`.

Postgres in debug compose:
- Host: `localhost`
- Port: `5432`
- DB: `shop_db`
- User: `shop_user`
- Password: `shop_password`
- Compose will create the `postgres` container and `postgres_data` volume automatically if they do not exist.
