# CLAUDE.md — Inspírate Inventory

## Descripción del proyecto

Inspírate Inventory es un sistema web para administrar una perfumería familiar en Bogotá, Colombia. Tiene dos módulos principales: gestión de inventario y registro de ventas. Los usuarios son un administrador (dueña del negocio) y vendedoras (personal no técnico, sin conocimientos de sistemas). La interfaz y los mensajes de error deben ser claros y simples para usuarios no técnicos.

## Stack técnico

- **Backend**: Go 1.22+, router [chi](https://github.com/go-chi/chi).
- **Query builder**: [sqlc](https://sqlc.dev/) (NO ORM).
- **Migraciones**: [goose](https://github.com/pressly/goose), formato SQL puro. Se usa como binario CLI (`~/go/bin/goose` vía `go install`), no como dependencia del módulo — no se importa desde código Go.
- **Base de datos**: PostgreSQL 15+ (Supabase en producción, Docker local en desarrollo).
- **Frontend**: Next.js 16 (App Router) + React 19 + TypeScript. Config de Next en `next.config.ts` (no `.mjs`).
- **Estilos**: Tailwind CSS v4 + shadcn/ui. Tailwind v4 se configura por CSS (`@theme` en `src/app/globals.css`), no con `tailwind.config.ts`. Al instalar shadcn/ui, usar la versión compatible con Tailwind v4.
- **Linter frontend**: ESLint con flat config (`eslint.config.mjs`), no `.eslintrc.json`.
- **Estado del servidor**: TanStack Query.
- **Estado global**: Zustand.
- **Formularios**: React Hook Form + Zod.
- **Testing backend**: paquete `testing` estándar de Go.
- **Linters**: golangci-lint (backend), ESLint + Prettier (frontend).
- **CI**: GitHub Actions (`.github/workflows/backend.yml`, `.github/workflows/frontend.yml`).

## Convenciones de código

- **PostgreSQL y JSON**: `snake_case` (columnas, claves JSON).
- **TypeScript**: `camelCase` para variables/funciones, `PascalCase` para tipos y componentes React.
- **Go**: `PascalCase` para tipos y funciones exportadas, `camelCase` para no exportadas.

## Estructura de carpetas

```
backend/
  cmd/api/            entrypoint del servidor HTTP
  internal/domain/    entidades y reglas de negocio puras, sin dependencias externas
  internal/usecase/   casos de uso, orquestación
  internal/repository/ interfaces + implementaciones sqlc
  internal/http/       handlers, middleware, router, respuestas JSON estandarizadas
  internal/platform/   config, db (pgxpool), mailer (Resend), jwt, logger (slog)
  db/migrations/       migraciones SQL (goose)
  db/queries/          queries SQL (sqlc)
  pkg/                 utilidades reutilizables

frontend/
  src/app/         App Router de Next.js
  src/components/  componentes reutilizables
  src/features/    módulos por dominio
  src/lib/         utilidades (cliente API, formatters)
  src/hooks/       custom hooks reutilizables
  src/stores/      stores de Zustand
  src/types/       tipos TypeScript compartidos

docs/
  adr/         Architecture Decision Records
  diagrams/    diagramas ERD y similares
```

## Reglas obligatorias

- **NO usar GORM ni ningún ORM.** Solo sqlc + SQL puro.
- **NO usar Prisma** en el frontend/backend.
- **NO usar librerías de UI** diferentes a shadcn/ui.
- **NO agregar dependencias sin justificar** — explicar por qué se necesita antes de instalarla.
- **NO commitear secretos.** Solo `.env.example` se versiona; `.env`/`.env.local` están en `.gitignore`.
- **Cada endpoint debe validar autorización explícitamente** dentro del handler o caso de uso — no confiar solo en middleware genérico.
- **Cada movimiento de inventario debe ejecutarse dentro de una transacción** junto con la actualización de `stock_actual`, para evitar inconsistencias.

## Herramientas externas requeridas

Además de Go y Node, el backend depende de binarios CLI que no son dependencias del módulo (no aparecen en `go.mod`):

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

`make install-tools` (en `backend/`) corre los tres comandos anteriores.

## Comandos comunes

### Backend (`cd backend`)

| Comando | Descripción |
|---|---|
| `make run` | Levanta la API con `go run` |
| `make dev` | `db-up` + levanta la API en un solo comando |
| `make build` | Compila el binario |
| `make test` | Corre los tests |
| `make test-coverage` | Tests con cobertura |
| `make lint` | Corre `golangci-lint` |
| `make sqlc-generate` | Genera código desde `db/queries` y `db/migrations` |
| `make migrate-up` / `make migrate-down` / `make migrate-status` | Migraciones (goose) |
| `make migrate-create name=nombre` | Crea una migración nueva |
| `make db-up` / `make db-down` / `make db-reset` | Postgres local (Docker) |
| `make install-tools` | Instala goose, golangci-lint y sqlc como binarios CLI |

### Frontend (`cd frontend`)

| Comando | Descripción |
|---|---|
| `npm run dev` | Servidor de desarrollo (puerto 3000) |
| `npm run build` | Build de producción |
| `npm run lint` | ESLint |

## Antes de dar por terminada una tarea

1. Backend: correr `make lint` y `make test` en `backend/` — ambos deben pasar sin errores.
2. Frontend: correr `npm run lint` y `npm run build` en `frontend/` — ambos deben pasar sin errores.
3. Si el cambio afecta un flujo visible en la UI, levantar `npm run dev` y probarlo en el navegador antes de reportar como completo.
