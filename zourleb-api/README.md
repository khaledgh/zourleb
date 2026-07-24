# Zourleb API

Go (Echo + GORM + MySQL) backend for the Zourleb tourism marketplace.

## Architecture

```
cmd/api        — HTTP API entrypoint (graceful shutdown)
cmd/worker     — background job runner (boost expiry, OTP cleanup, …)
config         — typed env configuration
internal/
  handler      — Echo handlers (parse/validate → service → response)
  service      — business logic
  repository   — GORM data access
  models       — GORM entities + DTOs
  middleware   — JWT auth, RBAC, locale, request id
  router       — route registration
  validator    — go-playground/validator + lb_phone rule
  database     — connection + AutoMigrate
  app          — dependency wiring (container)
pkg/           — framework-agnostic packages (token, hash, phone, oauth, …)
seeds          — baseline data (languages, RBAC, settings, super admin)
```

Layering: **handler → service → repository → models**. Handlers never touch
GORM; services own all rules and orchestration.

## Running locally

1. Start MySQL (and Redis):

   ```bash
   make docker-up        # or bring your own MySQL on :3306
   ```

2. Configure env:

   ```bash
   cp .env.example .env  # adjust DB_*, JWT_SECRET, GOOGLE_CLIENT_ID, …
   ```

   Env vars are read directly from the process environment. On Windows
   PowerShell, set them with `$env:DB_PASSWORD = "..."` before `make run`,
   or use a tool like `direnv`/`godotenv` wrapper.

3. Run the API (auto-migrates schema and seeds baseline data on boot):

   ```bash
   make run
   ```

   The server listens on `:8080`. A super admin is seeded from
   `SUPER_ADMIN_EMAIL` / `SUPER_ADMIN_PASSWORD` — change the password after
   first login.

## Implemented

**Foundation** — unified response envelope with stable error codes; env config;
slog logging; graceful shutdown; full GORM schema (AutoMigrate) across all
domains; seeds (ar/fr/en languages, permission catalog, roles, settings, super
admin, sample translations).

**Auth & account** — register, login, Google OAuth verify, JWT access + rotating
hashed refresh tokens, logout; `/me`, device registration; RBAC middleware
(`RequirePermission`, agency-scope), locale + request-id middleware.

**Phone OTP** — `/otp/request` + `/otp/verify` over WhatsApp Cloud API / SMS
(pluggable, dev-mode logs the code), codes hashed in Redis/in-memory cache with
TTL, resend cooldown, attempt cap, daily cap, DB audit trail.

**Public catalog** — `/home` aggregation, `/tours` (filter/sort/paginate),
`/tours/:slug`, `/categories`, `/regions`, `/agencies/:slug`, `/reviews`; all
content localized via the `*_translations` fallback chain (requested → default).

**Booking** — create (phone-verification gated, atomic seat hold, multi-currency
pricing), list, detail, pay (provider abstraction), cancel (seat release).

**Engagement** — favorites, verified-booking reviews, in-app notifications
(persisted + OneSignal push).

**Agency portal** (RBAC + agency-scoped) — apply, profile, members, tours CRUD +
publish, departures, prices, images incl. **past-gallery**, bookings, image
upload, boosts (packages, create, pay).

**Monetization** — boost lifecycle (create → pay → activate, priority-ordered
active injection, impressions/clicks), payment ledger, provider registry with a
built-in **manual/offline** adapter for the Lebanese market.

**Shop** (toggleable behind `shop.enabled`) — public products, orders with atomic
stock decrement; agency product management.

**Super admin** — agency approval, languages CRUD, UI translations (+ bulk),
settings/feature-flags, banners, categories/regions, review moderation, users,
pending boosts, manual payment confirmation, platform analytics.

**Worker** (`cmd/worker`) — boost expiry, OTP audit cleanup, departure
reconciliation on a ticker.

## Endpoint map (v1, under `/api/v1`)

- **Auth**: `POST /auth/{register,login,google,refresh,logout}`, `POST /otp/{request,verify}`
- **Public**: `GET /home`, `GET /tours`, `GET /tours/:slug`, `GET /categories`, `GET /regions`, `GET /agencies/:slug`, `GET /languages`, `GET /i18n/:locale`, `GET /reviews`
- **Me**: `GET/PATCH /me`, `POST /me/devices`, `GET /notifications`, `POST /notifications/:id/read`
- **Bookings**: `POST /bookings`, `GET /bookings`, `GET /bookings/:code`, `POST /bookings/:code/{pay,cancel}`
- **Favorites/Reviews**: `GET/POST/DELETE /favorites[/:tourId]`, `POST /reviews`
- **Agency** (`/agency/**`): profile, members, tours+departures+prices+images, bookings, uploads, boosts
- **Shop** (`/shop/**`, gated): products, orders; `/agency/products` management
- **Admin** (`/admin/**`): agencies, languages, translations, settings, banners, categories, regions, reviews, users, boosts, payments/confirm, analytics
- **Ops**: `GET /health`, `GET /ready`

## Next (not yet built)

In-app chat, loyalty/referral, map view, blog/events, subscription tiers, S3
storage adapter, Stripe/local payment adapters, e-voucher QR, the React admin
dashboard and Expo mobile app (separate packages).
```
