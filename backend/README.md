# Inspírate Inventory — Backend

API en Go para el sistema de inventario y ventas de Inspírate Inventory.

## Stack

- Go 1.22+, router [chi](https://github.com/go-chi/chi)
- [sqlc](https://sqlc.dev/) para queries tipadas (sin ORM)
- [goose](https://github.com/pressly/goose) para migraciones SQL
- PostgreSQL 15+ (Docker local / Supabase en producción)
- `log/slog` para logging estructurado

## Requisitos

- Go 1.22+
- Docker (para Postgres local)

## Desarrollo local

```bash
cp .env.example .env
make db-up     # levanta Postgres local
make run       # levanta la API en :8080
```

Verificar:

```bash
curl http://localhost:8080/api/v1/health
```

## Comandos

| Comando | Descripción |
|---|---|
| `make run` | Corre la API con `go run` |
| `make build` | Compila el binario a `bin/api` |
| `make test` | Corre los tests |
| `make test-coverage` | Corre los tests con reporte de cobertura |
| `make lint` | Corre `golangci-lint` |
| `make sqlc-generate` | Genera código a partir de `db/queries` y `db/migrations` |
| `make migrate-up` | Aplica migraciones pendientes |
| `make migrate-down` | Revierte la última migración |
| `make migrate-status` | Muestra el estado de las migraciones |
| `make migrate-create name=nombre` | Crea una nueva migración SQL |
| `make db-up` / `make db-down` | Levanta / detiene Postgres local |
| `make db-reset` | Reinicia Postgres local (borra el volumen) |

## Estructura

```
cmd/api/          entrypoint del servidor HTTP
internal/domain/       entidades y reglas de negocio puras
internal/usecase/      casos de uso, orquestación
internal/repository/   interfaces + implementaciones sqlc
internal/http/         handlers, middleware, router, respuestas JSON
internal/platform/     config, db, mailer, jwt, logger
db/migrations/         migraciones SQL (goose)
db/queries/            queries SQL (sqlc)
```
