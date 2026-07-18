# Inspírate Inventory

Sistema web para administrar el inventario y las ventas de una perfumería familiar en Bogotá, Colombia. Dos módulos principales: **gestión de inventario** y **registro de ventas**, usado por un administrador (dueña del negocio) y vendedoras.

> Este repositorio está en el **módulo 0**: setup inicial del monorepo. Todavía no hay lógica de negocio, tablas de base de datos ni autenticación implementadas.

## Stack

- **Backend**: Go 1.22+, chi, sqlc, goose, PostgreSQL 15+, `slog`.
- **Frontend**: Next.js 16 (App Router) + React 19 + TypeScript, Tailwind CSS v4, TanStack Query, Zustand, React Hook Form + Zod.
- **CI**: GitHub Actions.

Ver [`CLAUDE.md`](./CLAUDE.md) para las convenciones y reglas completas del proyecto.

## Prerrequisitos

- Go 1.22+
- Node 20+
- Docker (Postgres local)
- `make`

## Herramientas requeridas

| Herramienta | Uso | Instalación |
|---|---|---|
| Go 1.22+ | Compilar y correr el backend | https://go.dev/dl/ |
| Node 20+ | Correr el frontend | https://nodejs.org/ |
| Docker + Docker Compose | Postgres local | https://docs.docker.com/get-docker/ |
| GNU Make | Targets de `backend/Makefile` | gestor de paquetes del sistema (`apt install make`, etc.) |
| [goose](https://github.com/pressly/goose) | Migraciones SQL | `go install github.com/pressly/goose/v3/cmd/goose@latest` |
| [golangci-lint](https://golangci-lint.run/) | Lint del backend | `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest` |
| [sqlc](https://sqlc.dev/) | Generación de código desde SQL | `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest` |

`cd backend && make install-tools` instala goose, golangci-lint y sqlc en un solo paso.

## Arrancar el ambiente local

1. **Backend**

   ```bash
   cd backend
   cp .env.example .env
   make dev        # db-up + levanta la API en http://localhost:8080
   ```

   Verificar: `curl http://localhost:8080/api/v1/health`

2. **Frontend** (en otra terminal)

   ```bash
   cd frontend
   cp .env.example .env.local
   npm install
   npm run dev     # http://localhost:3000
   ```

## Estructura del monorepo

```
inspirate-inventory/
├── backend/    API en Go (chi + sqlc + goose)
├── frontend/   App Next.js (App Router + TypeScript)
├── docs/       ADRs y diagramas
└── .github/    Workflows de CI
```

## Comandos útiles

| Dónde | Comando | Qué hace |
|---|---|---|
| `backend/` | `make run` | Levanta la API |
| `backend/` | `make dev` | `db-up` + levanta la API en un solo comando |
| `backend/` | `make test` | Corre los tests |
| `backend/` | `make lint` | Corre `golangci-lint` |
| `backend/` | `make db-up` / `make db-down` | Postgres local |
| `backend/` | `make install-tools` | Instala goose, golangci-lint y sqlc |
| `frontend/` | `npm run dev` | Servidor de desarrollo |
| `frontend/` | `npm run build` | Build de producción |
| `frontend/` | `npm run lint` | ESLint |

Cada carpeta (`backend/`, `frontend/`) tiene su propio README con más detalle.
