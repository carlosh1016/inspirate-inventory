# Inspírate Inventory

Sistema web para administrar el inventario y las ventas de una perfumería familiar en Bogotá, Colombia. Dos módulos principales: **gestión de inventario** y **registro de ventas**, usado por un administrador (dueña del negocio) y vendedoras.

> Este repositorio está en el **módulo 0**: setup inicial del monorepo. Todavía no hay lógica de negocio, tablas de base de datos ni autenticación implementadas.

## Stack

- **Backend**: Go 1.22+, chi, sqlc, goose, PostgreSQL 15+, `slog`.
- **Frontend**: Next.js (App Router) + TypeScript, Tailwind CSS, TanStack Query, Zustand, React Hook Form + Zod.
- **CI**: GitHub Actions.

Ver [`CLAUDE.md`](./CLAUDE.md) para las convenciones y reglas completas del proyecto.

## Prerrequisitos

- Go 1.22+
- Node 20+
- Docker (Postgres local)
- `make`

## Arrancar el ambiente local

1. **Backend**

   ```bash
   cd backend
   cp .env.example .env
   make db-up      # levanta Postgres local (Docker)
   make run        # levanta la API en http://localhost:8080
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
| `backend/` | `make test` | Corre los tests |
| `backend/` | `make lint` | Corre `golangci-lint` |
| `backend/` | `make db-up` / `make db-down` | Postgres local |
| `frontend/` | `npm run dev` | Servidor de desarrollo |
| `frontend/` | `npm run build` | Build de producción |
| `frontend/` | `npm run lint` | ESLint |

Cada carpeta (`backend/`, `frontend/`) tiene su propio README con más detalle.
