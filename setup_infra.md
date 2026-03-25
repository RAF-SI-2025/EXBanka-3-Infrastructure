# EXBanka 3 Infrastructure Notes

This repository is no longer the active application runtime.

## Current role

`EXBanka-3-Infrastructure` now acts as a support repository for:

- historical proto-generation references
- seeding/helpers that informed the current microservice split
- legacy notes from the earlier monolithic phase

## Important clarification

The active application is started from the project root with:

```powershell
cd C:\Dev\Projects\SI
docker compose up -d --build
```

That root stack launches:

- frontend nginx gateway
- shared backend docs/health container
- auth, employee, client, account, transfer, payment, exchange, and loan services
- PostgreSQL
- Mailhog

## Do not use the old monolith instructions as runtime guidance

Older notes in this repository may refer to:

- `cmd/server/main.go` as the main application server
- a monolithic HTTP + gRPC runtime on `:8080` / `:9090`
- setup flows that seed and run everything from this repo alone

Those instructions are outdated for the current codebase shape.

## Where to look instead

- backend runtime/docs:
  [EXBanka-3-Backend/README.md](/C:/Dev/Projects/SI/EXBanka-3-Backend/README.md)
- frontend/gateway runtime:
  [EXBanka-3-Frontend/README.md](/C:/Dev/Projects/SI/EXBanka-3-Frontend/README.md)
- root stack wiring:
  [docker-compose.yml](/C:/Dev/Projects/SI/docker-compose.yml)
- gateway routing:
  [EXBanka-3-Frontend/nginx.conf](/C:/Dev/Projects/SI/EXBanka-3-Frontend/nginx.conf)
