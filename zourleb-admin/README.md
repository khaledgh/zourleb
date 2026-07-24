# Zourleb Admin

React + TypeScript dashboard for the Zourleb platform — a single SPA serving two
role-based experiences: **Super Admin** (platform management) and **Agency**
(portal). Built with Vite, React Router, React Query, Tailwind v3,
react-hook-form, and react-i18next.

## Stack

- **Vite + React 18 + TypeScript** (strict)
- **React Router v6** with auth + role guards
- **React Query v5** for all server state
- **Tailwind v3** with logical properties (`start`/`end`) for RTL
- **react-hook-form** for forms, **zustand** for auth/toast stores
- **axios** instance with bearer injection, `Accept-Language`, and a
  single-flight **refresh-on-401** interceptor
- **react-i18next** hydrated from the backend `GET /i18n/:locale` bundle, with
  built-in ar/en fallback and automatic `dir=rtl`

## Running

The dev server proxies `/api` and `/uploads` to the Go backend on `:8080`, so
run the API first (`cd ../zourleb-api && make run`), then:

```bash
npm install
npm run dev          # http://localhost:5173
```

Sign in with the seeded super admin (`SUPER_ADMIN_EMAIL` / `SUPER_ADMIN_PASSWORD`
from the backend `.env`).

```bash
npm run typecheck    # tsc --noEmit
npm run build        # type-checked production bundle → dist/
```

The built `dist/` can be embedded and served by the Go binary (the plan's
single-binary deployment pattern) or hosted separately behind the same domain.

## Structure

```
src/
  api/client.ts        axios + envelope unwrap + refresh interceptor
  stores/auth.ts       zustand auth store (login/bootstrap/roles)
  i18n/                react-i18next init, backend hydration, RTL direction
  routes/              router + RequireAuth / RequireRole guards
  components/
    layout/            role-filtered sidebar shell, locale switcher
    ui/                DataTable, primitives, toast
  pages/               one file per screen
  types/api.ts         response envelope + domain types
```

## Screens

**Super Admin** — Dashboard (analytics), Agencies (approve/suspend), Users
(block), Reviews (moderate), Boosts (confirm offline payment → activate),
Banners (CRUD), Languages (add without redeploy), Translations (UI string
editor), Settings (feature-flag toggles).

**Agency** — My Tours (publish/unpublish), Bookings. The sidebar renders the
agency menu when the signed-in user holds an agency role; the backend resolves
agency scope from membership (or the `X-Agency-ID` header for multi-agency
staff).

## Notes / next

- Tour authoring forms (create/edit with translations, departures, prices,
  image manager with drag-reorder + past-gallery) are stubbed at the list level;
  the create/edit detail screens are the next build-out.
- Image uploads go through `POST /agency/uploads/image` (multipart) — wire the
  drag-drop manager into the tour editor when added.
