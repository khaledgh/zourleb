# Zourleb — Full‑Stack Architecture & Implementation Plan

> **Zourleb** (زورليبان / *Visit Lebanon*) — a multi‑sided tourism marketplace connecting Lebanese travel agencies with local & international tourists. Agencies publish tours, build a portfolio from past trips, sell optional products, and pay to boost their visibility. Tourists discover, book, and pay for tours with verified phone numbers.

**Stack:** Go (Echo + GORM + MySQL) · React + TS admin (React Query + Tailwind v3) · Expo (React Native + NativeWind + React Query) · OneSignal · WhatsApp/SMS OTP · backend‑driven i18n (ar/fr/en + extensible).

---

## 1. Actors & Roles

| Actor | Surface | Can do |
|---|---|---|
| **Super Admin** (platform owner) | Web + Mobile | Approve/suspend agencies, manage languages & translations, manage boost packages & banners, toggle modules (e.g. Shop), see platform analytics & revenue, manage commission rates, moderate reviews. |
| **Agency Owner** | Web + Mobile | Manage agency profile, invite staff, create/publish tours, manage departures & pricing, upload past‑tour galleries, buy boosts, manage shop products, view bookings & payouts. |
| **Agency Staff** (manager / editor) | Web + Mobile | Scoped subset of owner actions via RBAC (e.g. editor can draft tours but not buy boosts). |
| **Tourist** | Mobile (primary) + Web (browse) | Discover tours, favorite, register/book, verify phone, pay, review, get push notifications, and apply to register an agency. |
| **Guest** | Mobile + Web | Browse public catalog before auth. |

RBAC is **role + permission** based (not hard‑coded roles) so new permission sets can be added without code changes.

---

## 2. High‑Level Architecture

```
                       ┌─────────────────────────────────────┐
   Tourist (Expo)──────┤                                     │
   Agency  (Expo)──────┤        Echo API (Go)                │
   Admin   (React)─────┤   handler → service → repository    │──► MySQL (GORM)
                       │   middleware: auth/RBAC/i18n/rate   │
                       │                                     │──► Redis (cache, OTP, rate‑limit)
                       └───────────────┬─────────────────────┘──► Object storage (images)
                                       │
        ┌──────────────┬──────────────┼───────────────┬──────────────┐
   OneSignal      WhatsApp/SMS     Payment GW       Google OAuth    Cron worker
   (push)         (OTP)            (boosts/orders)  (sign‑in)       (boost expiry, etc.)
```

Single Go binary can also **embed and serve the compiled admin build** (your preferred deployment pattern), or admin can be deployed separately behind the same domain.

---

## 3. Backend — Go Architecture (Echo + GORM + MySQL)

### 3.1 Folder structure
Your required core (`handler`, `service`, `models`) is kept front‑and‑center; a thin `repository` layer is added so business logic never touches GORM directly (testability + swap‑ability).

```
zourleb-api/
├── cmd/api/main.go                # entrypoint, DI wiring, graceful shutdown
├── cmd/worker/main.go             # cron/queue worker (boost expiry, notifications)
├── config/                        # env loading, typed Config struct
├── internal/
│   ├── handler/                   # Echo HTTP handlers (request/response only)
│   │   ├── auth_handler.go
│   │   ├── tour_handler.go
│   │   ├── booking_handler.go
│   │   ├── boost_handler.go
│   │   ├── shop_handler.go
│   │   ├── i18n_handler.go
│   │   └── ...
│   ├── service/                   # business logic (the brain)
│   │   ├── auth_service.go
│   │   ├── otp_service.go
│   │   ├── tour_service.go
│   │   ├── booking_service.go
│   │   ├── boost_service.go
│   │   ├── payment_service.go
│   │   ├── notification_service.go
│   │   └── ...
│   ├── repository/                # GORM data access, one per aggregate
│   ├── models/                    # GORM models + DTOs + request/response structs
│   ├── middleware/                # JWTAuth, RequirePermission, Locale, RateLimit, RequestID
│   ├── router/                    # route registration grouped by version & role
│   ├── dto/                       # request/response (if you prefer separate from models)
│   └── validator/                 # go-playground/validator setup + custom rules (LB phone)
├── pkg/                           # reusable, framework-agnostic packages
│   ├── otp/                       # code gen + hashing
│   ├── whatsapp/                  # WhatsApp Cloud API client
│   ├── sms/                       # SMS gateway client
│   ├── onesignal/                 # OneSignal REST client
│   ├── oauth/                     # Google token verification
│   ├── storage/                   # S3/local image storage abstraction
│   ├── payment/                   # payment provider abstraction
│   ├── i18n/                      # translation resolver + fallback chain
│   ├── pagination/
│   ├── response/                  # unified API envelope + error codes
│   └── phone/                     # Lebanese number normalization (+961)
├── migrations/                    # versioned SQL (golang-migrate)
├── seeds/                         # languages, default settings, super admin
└── Makefile / docker-compose.yml
```

### 3.2 Layering rules
- **handler** → parse/validate input, call service, format response. No business logic, no GORM.
- **service** → all rules, transactions, orchestration across repositories & external clients. Returns domain errors.
- **repository** → GORM queries only, returns models. No business rules.
- **models** → GORM structs + relationships; DTOs separated so you never leak DB internals to the API.

### 3.3 Cross‑cutting
- **Unified response envelope:** `{ "success": bool, "data": ..., "error": { "code": "TOUR_NOT_FOUND", "message": "...localized..." }, "meta": { pagination } }`.
- **Error codes are stable strings** (clients map them to localized UI; messages also localized server‑side via `Accept-Language`).
- **Validation:** `go-playground/validator` with a custom `lb_phone` rule. Every write endpoint validates a request DTO.
- **Auth:** JWT access token (short‑lived) + refresh token (rotating, stored hashed). `Authorization: Bearer`.
- **Soft deletes** (`gorm.DeletedAt`) on tours/agencies/products so nothing is hard‑lost.
- **Audit log** table for sensitive actions (agency approval, boost approval, refunds).
- **Idempotency keys** on payment & booking creation (you've used this pattern before — reuse it here).
- **DB transactions** kept short; long external calls (WhatsApp/payment) happen **outside** the transaction.

---

## 4. Database Schema (MySQL / GORM)

> Translatable content uses sibling `*_translations` tables keyed by `(parent_id, locale)`. UI strings live in a separate `translations` table managed by Super Admin.

### Identity & access
- **users** — `id, name, email, email_verified_at, password_hash(null for OAuth-only), google_id, avatar, phone, phone_verified_at, locale, status(active/blocked), last_login_at`
- **roles** — `id, key(super_admin/agency_owner/agency_staff/tourist), name`
- **permissions** — `id, key(tour.create, boost.purchase, …)`
- **role_permissions** — `role_id, permission_id`
- **user_roles** — `user_id, role_id, agency_id(null for global roles)`  ← scopes a role to an agency

### Agencies
- **agencies** — `id, slug, name, logo, cover, phone, email, region_id, website, verified(bool), status(pending/approved/suspended), subscription_tier, commission_rate, rating_avg, rating_count, created_by`
- **agency_translations** — `agency_id, locale, description, about`
- **agency_members** — `agency_id, user_id, role_id, invited_at, joined_at`  (staff invites)

### Catalog
- **regions** — `id, slug, parent_id(null)` (+ **region_translations**: name) — Lebanese governorates/areas
- **categories** — `id, slug, icon, sort_order` (+ **category_translations**: name) — adventure, religious, historical, culinary, ski, beach, hiking…
- **tours** — `id, agency_id, slug, category_id, region_id, type(day_trip/multi_day), duration_days, difficulty, min_age, max_capacity, base_currency, price_from, status(draft/pending/published/archived), is_showcase(bool ← "old/past tours" portfolio), featured(bool), starts_from_date, created_by`
- **tour_translations** — `tour_id, locale, title, summary, description, itinerary(json/long), included(text), excluded(text), meta_title, meta_description`
- **tour_images** — `tour_id, url, sort_order, is_cover, source(current/past_gallery)`  ← old‑tour images flagged as `past_gallery`
- **tour_departures** — `id, tour_id, start_date, end_date, capacity, seats_booked, status(open/full/cancelled)`
- **tour_prices** — `id, tour_id, departure_id(null=applies to all), traveler_type(adult/child/infant/student), currency(USD/LBP), amount`  ← dual‑currency ready

### Bookings
- **bookings** — `id, code, user_id, tour_id, departure_id, status(pending/awaiting_payment/confirmed/cancelled/completed), travelers_count, subtotal, currency, contact_phone, contact_phone_verified(bool), payment_status, notes`
- **booking_travelers** — `id, booking_id, full_name, traveler_type, phone, phone_verified_at`
- **booking_payments** — `id, booking_id, provider, provider_ref, amount, currency, status, idempotency_key`

### Phone verification
- **phone_verifications** — `id, phone, channel(whatsapp/sms), purpose(register/booking), code_hash, attempts, max_attempts, expires_at, verified_at, ip` (live OTPs cached in Redis; this table is the audit trail)

### Shop (toggleable module)
- **products** — `id, agency_id, slug, category_id, price, currency, stock, status(active/hidden/out_of_stock)`
- **product_translations** — `product_id, locale, name, description`
- **product_images** — `product_id, url, sort_order, is_cover`
- **orders** / **order_items** — standard commerce tables, gated by the `shop.enabled` setting

### Monetization
- **boost_packages** — `id, name, placement(home_banner/featured_list/category_top/search_top), duration_days, price, currency, max_active_slots, priority`
- **boosts** — `id, agency_id, tour_id(or product_id), package_id, placement, status(pending_payment/active/expired/rejected), starts_at, ends_at, impressions, clicks, payment_id`
- **banners** — `id, type(boosted/manual), tour_id(null), image, title, link, sort_order, active, starts_at, ends_at` ← home carousel; boosted ones auto‑injected
- **payments** — polymorphic ledger: `id, payable_type(boost/order/booking/subscription), payable_id, payer_type(agency/user), payer_id, amount, currency, provider, provider_ref, status, idempotency_key`

### Platform
- **languages** — `id, code, name, native_name, is_rtl, is_active, is_default, sort_order` ← Super Admin adds languages here
- **translations** — `id, locale, namespace, key, value` ← UI strings, served as a JSON bundle
- **settings** — `key, value, type` ← feature flags (`shop.enabled`, `reviews.enabled`, `boost.enabled`, `default_currency`, etc.)
- **reviews** — `id, user_id, tour_id, booking_id, rating(1‑5), comment, status(pending/approved/rejected)`
- **favorites** — `user_id, tour_id`
- **device_tokens** — `id, user_id, onesignal_player_id, platform(ios/android/web)`
- **notifications** — `id, user_id, type, title, body, data(json), read_at`
- **audit_logs** — `id, actor_id, action, target_type, target_id, payload(json), ip, created_at`

---

## 5. API Design (v1, grouped)

All under `/api/v1`. Locale resolved from `Accept-Language` or `?locale=`. Pagination via `?page&per_page`, filtering via query params.

**Auth & account**
- `POST /auth/register` · `POST /auth/login` · `POST /auth/google` · `POST /auth/refresh` · `POST /auth/logout`
- `POST /auth/email/verify` · `POST /auth/email/resend`
- `POST /otp/request` `{ phone, channel, purpose }` · `POST /otp/verify` `{ phone, code }`
- `GET /me` · `PATCH /me` · `POST /me/devices` (register OneSignal player id)

**Public catalog**
- `GET /tours` (filters: region, category, date, price, type, currency, sort, `featured`) · `GET /tours/:slug`
- `GET /home` → bundles banners + boosted/featured + categories + regions in one call
- `GET /agencies/:slug` · `GET /categories` · `GET /regions`
- `GET /i18n/:locale` → UI translation bundle · `GET /languages`

**Tourist**
- `POST /bookings` (requires verified contact phone) · `GET /bookings` · `GET /bookings/:code` · `POST /bookings/:code/cancel`
- `POST /bookings/:code/pay` · favorites · `POST /reviews`

**Agency portal** (`/agency/...`, RBAC‑scoped to caller's agency)
- agency profile, members/invites
- tours CRUD + publish, departures CRUD, prices CRUD, image upload/reorder, **past‑gallery upload**
- `GET /agency/bookings`, `GET /agency/analytics`
- boosts: `GET /agency/boost-packages` · `POST /agency/boosts` · `POST /agency/boosts/:id/pay`
- shop: products CRUD (only if `shop.enabled`)

**Super admin** (`/admin/...`)
- agencies approve/suspend, users, categories/regions
- languages CRUD, translations CRUD (+ bulk import/export)
- boost packages CRUD, boost approvals, manual banners
- settings (feature flags), reviews moderation, platform analytics & revenue

---

## 6. Authentication & Authorization

1. **Email/password** — bcrypt hashing, email verification link/code before booking is allowed.
2. **Google OAuth** — mobile/web obtains Google ID token → `POST /auth/google` → backend verifies token signature & audience → links or creates user by `google_id`/email.
3. **JWT** — access (15 min) + rotating refresh (30 days, hashed in DB, revocable on logout/breach).
4. **RBAC middleware** — `RequirePermission("tour.create")`; permissions resolved from the caller's roles, **scoped by `agency_id`** so staff only touch their own agency.
5. **Phone is a separate verified attribute**, required specifically at booking time (see §7).

---

## 7. Phone Verification (WhatsApp + SMS OTP)

**When:** required before a booking is confirmed; optionally at registration. Each traveler can be individually verified, or just the contact phone — configurable via `settings`.

**Flow**
1. `POST /otp/request { phone, channel }` → normalize to E.164 (`+961…` via `pkg/phone`), rate‑limit per phone+IP (Redis), generate 6‑digit code, store **hash** in Redis with TTL (e.g. 5 min) + audit row.
2. Send via chosen channel:
   - **WhatsApp** → WhatsApp Cloud API (Meta) template message, or Twilio WhatsApp.
   - **SMS** → Twilio / Vonage / a local Lebanese SMS gateway.
   - Channel abstraction (`pkg/otp` dispatches to `pkg/whatsapp` or `pkg/sms`) so providers are swappable.
3. `POST /otp/verify { phone, code }` → compare hash, enforce max attempts & expiry → set `phone_verified_at`.

**Hardening:** resend cooldown, max 3 attempts then lock, per‑phone daily cap, code never returned in API responses, only hashes stored.

---

## 8. Internationalization (backend‑driven)

- **Languages managed in DB** (`languages` table). Super Admin can add a new language (e.g. Armenian, German) without a deploy. Base set: **Arabic (RTL, default), French, English**.
- **UI strings** in `translations` (`locale, namespace, key, value`); `GET /i18n/:locale` returns a cached JSON bundle (versioned/ETag) that web admin & Expo app load on boot and cache offline.
- **Content translations** (tours, agencies, categories, products) via `*_translations` tables with a **fallback chain**: requested locale → default locale → key. The `pkg/i18n` resolver handles this server‑side; lists return already‑localized fields based on `Accept-Language`.
- **RTL** handled on clients: web `dir="rtl"` + Tailwind logical utilities; Expo `I18nManager` + RTL‑aware NativeWind classes.
- **Currency**, given Lebanon's reality: dual **USD/LBP** pricing supported at the `tour_prices` level; a display currency setting + optional live/admin‑set exchange rate.

---

## 9. Notifications (OneSignal)

- Mobile/web register their OneSignal **player id** → `POST /me/devices`.
- `notification_service` sends via OneSignal REST API and also persists an in‑app `notifications` row.
- **Triggers:** booking confirmed/cancelled, payment received, OTP fallback, new tour from a favorited agency, price drop / last‑minute deal, **boost about to expire** (to the agency), review reply, admin broadcasts.
- Localized payloads using the recipient's `locale`. Quiet‑hours & per‑type opt‑out stored per user.

---

## 10. Monetization & Ad‑Boost System

**Boost = paid placement of a tour (or product) in a high‑visibility slot for a date range.**

1. Super Admin defines **boost_packages** (placement + duration + price + priority + max slots).
2. Agency picks a tour + package → creates a `boost` (`pending_payment`).
3. Agency pays → on success the boost becomes **active** for `[starts_at, ends_at]`.
4. Home/featured/search queries inject **active boosts** ordered by package `priority` (then recency), capped at `max_active_slots`. If full, agency can queue for the next window.
5. **Impressions & clicks tracked** per boost → analytics + invoices.
6. Cron worker auto‑expires boosts and notifies the agency to renew.

**Other revenue levers (pick what fits):**
- **Booking commission** — platform takes `commission_rate %` per confirmed booking (already on `agencies`).
- **Agency subscription tiers** (Basic / Pro / Premium) — more tours, lower commission, included boost credits.
- **Featured agency** badge, **shop listing fees**, **sponsored blog posts**.

**Payments in Lebanon — be realistic:** international card rails (Stripe) are limited. Architect `pkg/payment` as a **provider interface** with adapters for:
- **Whish Money / OMT / Areeba** (local) — common for in‑country payments.
- **Stripe/Checkout** — for international cards / diaspora bookings.
- **Manual bank transfer / cash‑on‑arrival** — a built‑in "offline payment + admin confirmation" provider for both bookings and boosts (very practical for the Lebanese market).

---

## 11. The Shop Module (toggleable)

- Entire module gated behind `settings['shop.enabled']`. When off, all `/shop` & `/agency/products` routes return `404/feature_disabled`, and clients hide the tab/menu (the home call reports enabled modules).
- Agencies list souvenirs / travel gear / local products; standard cart → order → payment flow reusing `pkg/payment`.
- You can ship the platform **with shop hidden** and flip it on later with zero redeploy.

---

## 12. Home Page & Banners

`GET /home` returns one optimized payload:
- **Banner carousel** = active manual banners + injected active home‑banner boosts.
- **Featured tours** (boosted `featured_list` + editorially `featured`).
- Categories, regions, "Upcoming departures", "Last‑minute deals", optional "Shop highlights" (if enabled), agencies spotlight.
- Everything pre‑localized & currency‑aware; heavily cacheable (Redis, short TTL) with cache‑busting on boost/banner changes.

---

## 13. Admin Dashboard (Web — React + TS)

- **Vite + React + TypeScript**, React Router, **React Query** (server state), **Tailwind v3**.
- **Forms:** `react-hook-form + zod` everywhere; every form validated, errors localized.
- **Axios** instance with interceptors (bearer token, refresh‑on‑401, `Accept-Language`), unified error → toast mapping by `error.code`.
- **Two layouts/menus** driven by role: Super Admin vs Agency. Route guards from permissions.
- **Modules:** Auth, Agencies, Tours (+ departures, prices, image manager with drag‑reorder & past‑gallery), Bookings, Boosts & Packages, Banners, Shop, Languages & Translations editor, Users & Roles, Settings/Feature‑flags, Analytics.
- **i18n:** `react-i18next` hydrated from `GET /i18n/:locale`; RTL toggle for Arabic.
- Data tables with server pagination/sorting/filtering; image uploads to storage via presigned URLs.

---

## 14. Mobile App (Expo — React Native)

- **Expo Router**, **NativeWind** (Tailwind for RN), **React Query**, `react-hook-form + zod`.
- **i18next + expo-localization**, RTL via `I18nManager`, bundle from `GET /i18n/:locale` cached locally.
- **OneSignal** SDK + secure token storage (`expo-secure-store`).
- Two faces in one app: *Tourist* experience by default; *Agency* portal unlocked when the signed‑in user has an agency role (mode switch) — so agencies manage tours/bookings/boosts on mobile too. Tourists can also submit an application to register a new agency directly from the Profile screen.
- Screens: Home (banners/featured), Search & filters, Tour detail (gallery, itinerary, dates, prices, reviews), Booking flow with **phone OTP**, My bookings (with e‑voucher/QR), Favorites, Agency profile, Agency registration/application form, Notifications, Profile/Language/Currency.
- Offline‑friendly: cache catalog & i18n; queue‑safe booking submission with idempotency keys.

---

## 15. Storage, Media & Jobs

- **Images:** object storage (S3‑compatible / your VPS via MinIO) behind a `pkg/storage` interface; presigned uploads; server stores only URLs. Generate thumbnail variants for cards (you've fought card layout/clipping before — fixed aspect ratios + cover crop help here).
- **Background worker** (`cmd/worker`): boost expiry, departure seat reconciliation, scheduled notifications, OTP cleanup, exchange‑rate refresh, analytics rollups. Use a DB‑backed queue or Redis + a scheduler.
- **Caching:** Redis for `/home`, i18n bundles, catalog lists, OTP, and rate limiting.

---

## 16. Feature Ideas (menu — pick & prioritize)

**Discovery & UX**
- Lebanon **map view** of tours (Leaflet/Mapbox) with region pins.
- Rich **filters**: price range, duration, region, difficulty, date, language of guide, family‑friendly.
- **Wishlist/Favorites**, "tours like this", recently viewed.
- **Seasonal collections** (Ski season, Summer beaches, Religious tours, Cedars & nature, Food & wine).
- **Last‑minute deals** & **early‑bird** pricing.
- **Tour guide profiles** with ratings.

**Trust & social**
- **Verified‑booking reviews** & photos.
- Agency **verification badges** & response‑time stats.
- **Referral program** + **loyalty points** (very effective for diaspora word‑of‑mouth).
- **In‑app chat** tourist ↔ agency (pre‑booking questions).
- Social sharing & deep links to tours.

**Booking power features**
- **Group / private tour** requests & custom quotes.
- **Waitlist** when a departure is full.
- **Saved travelers** (repeat bookings skip re‑entry; phone stays verified).
- **E‑vouchers / QR tickets**, check‑in scan for the agency.
- **Cancellation/refund policy** display + automated refund states.
- **Multi‑currency USD/LBP** with admin‑set or live rate.

**Content & growth**
- **Blog / travel guides** (SEO; "best hikes in Lebanon", "Byblos in a day").
- **Events/festivals calendar** (Baalbeck Festival, Beiteddine, etc.).
- **Newsletter** + segmented push campaigns.

**Agency tooling & monetization**
- **Analytics dashboard** (views, conversion, top tours, boost ROI).
- **Subscription tiers** with included boost credits & lower commission.
- **Boost auction/bidding** (v2 upgrade over fixed packages).
- **Payouts/statements** & commission ledger.

**Safety & ops**
- **Emergency contacts & safety info** per tour.
- **Weather** for departure dates.
- **Offline itinerary** access.
- Multi‑pickup‑point selection & transport details.

---

## 17. Phased Roadmap

**Phase 0 — Foundations**
Repo + CI, Docker, config, migrations & seeds (languages ar/fr/en, super admin, default settings), response/error envelope, auth skeleton (email + Google + JWT), RBAC, i18n bundle endpoint, storage & image upload, OneSignal wiring.

**Phase 1 — MVP (bookable marketplace)**
Agencies (apply from web/mobile → approve), tours + translations + images (incl. **past‑gallery**), departures & prices, public catalog + `/home` + banners (manual), search/filters, **booking with phone OTP** (WhatsApp/SMS), email/Google auth, tourist app + agency portal (web & mobile), OneSignal events, base multi‑language.

**Phase 2 — Monetization & trust**
Boost packages + boost purchase + payment provider(s) incl. manual/offline, boosted home/featured injection, impressions/clicks, agency analytics, reviews & ratings, agency verification, commission ledger.

**Phase 3 — Growth modules**
Shop module (behind flag), loyalty + referral, in‑app chat, advanced filters & map, blog/guides, events calendar, waitlist & saved travelers.

**Phase 4 — Scale & optimize**
Subscription tiers, boost bidding, deeper analytics/rollups, performance & caching hardening, payout automation, A/B on home placements.

---

## 18. Non‑Functional Checklist (senior concerns)

- **Security:** bcrypt, JWT rotation, OTP hashing + rate limits, RBAC scoping, idempotency on money paths, audit logs, input validation on every write, signed image URLs, OWASP basics, secrets via env, least‑privilege DB user.
- **Reliability:** short DB transactions, external calls outside transactions, queued/async sends, graceful shutdown, health/readiness endpoints, structured logging + request IDs.
- **Performance:** Redis caching for hot reads, indexed filter columns (region/category/status/dates), pagination everywhere, image thumbnails/CDN, N+1 avoidance via GORM preloads.
- **Scalability:** stateless API (scale horizontally), separate worker, object storage off the app server.
- **Observability:** metrics (Prometheus), error tracking (Sentry), boost/booking funnels.
- **i18n/RTL correctness** verified across web + mobile; currency formatting per locale.
- **Lebanon‑specific:** offline/manual payment path, USD/LBP duality, WhatsApp‑first OTP (higher deliverability than SMS locally), resilient to flaky connectivity (client caching, retries with idempotency).

---

*Next step suggestion: I can turn any single section into a build‑ready spec for an AI coding assistant — e.g. the full MySQL migration set + GORM models, the boost engine, or the OTP service — in the same hardened, no‑silent‑failure style you use.*
