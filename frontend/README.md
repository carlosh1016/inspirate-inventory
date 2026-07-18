# Inspírate Inventory — Frontend

Aplicación web (Next.js App Router + TypeScript) para el sistema de inventario y ventas de Inspírate Inventory.

## Stack

- Next.js (App Router) + TypeScript
- Tailwind CSS
- TanStack Query (estado del servidor)
- Zustand (estado global)
- React Hook Form + Zod (formularios)
- axios (cliente HTTP)

## Requisitos

- Node 20+

## Desarrollo local

```bash
cp .env.example .env.local
npm install
npm run dev
```

Abrir [http://localhost:3000](http://localhost:3000).

## Comandos

| Comando | Descripción |
|---|---|
| `npm run dev` | Levanta el servidor de desarrollo |
| `npm run build` | Compila la aplicación para producción |
| `npm run start` | Sirve el build de producción |
| `npm run lint` | Corre ESLint |

## Estructura

```
src/app/           App Router: layouts, páginas, estilos globales
src/components/    componentes reutilizables
src/features/      módulos por dominio
src/lib/           utilidades (cliente API, formatters)
src/hooks/         custom hooks reutilizables
src/stores/        stores de Zustand
src/types/         tipos TypeScript compartidos
```
