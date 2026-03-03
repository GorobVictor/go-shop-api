# shop-api

Commands:

1. update swagger
```bash
swag init -g cmd/main.go -o internal/docs
```

## Debug in Docker (VS Code)

1. Open this repository in VS Code.
2. Start debug config `Go: Attach to Docker (auto up/down)` from Run and Debug.
3. Set breakpoints in files under `src/`.

What is used:
- `.vscode/tasks.json` runs `docker compose -f docker-compose.debug.yml up --build -d` and `down`.
- `.vscode/launch.json` attaches to Delve on port `2345`.
- `docker-compose.debug.yml` starts the debug container with `src` mounted to `/app`.

Manual mode:
```bash
docker compose -f docker-compose.debug.yml up --build -d
```
Then run debug config `Go: Attach to Docker (manual)`.
