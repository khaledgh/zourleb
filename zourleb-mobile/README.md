# Zourleb Mobile

Expo (React Native) app for Zourleb — the tourist-facing experience with an
agency mode unlocked for agency accounts. Built with Expo Router, NativeWind,
React Query, react-hook-form, and i18next.

## Stack

- **Expo SDK 52** + **Expo Router** (file-based routing, typed groups)
- **NativeWind v4** (Tailwind for RN) with the brand theme
- **React Query v5** for server state
- **zustand** auth store; **expo-secure-store** for tokens (localStorage on web)
- **react-hook-form** for forms
- **i18next + react-i18next + expo-localization** with RTL via `I18nManager`,
  hydrated from the backend `GET /i18n/:locale` bundle (ar/fr/en, extensible)
- **axios** client mirroring the admin: bearer + `Accept-Language`, in-memory
  token cache, single-flight **refresh-on-401**

## Running

Point the app at the API via `app.json` → `expo.extra.apiBaseUrl`. On a physical
device use your machine's LAN IP (not `localhost`). Start the backend
(`cd ../zourleb-api && make run`), then:

```bash
npm install
npx expo start          # press i / a, or scan the QR with Expo Go
npm run typecheck       # tsc --noEmit
```

## Structure

```
app/                         Expo Router routes
  _layout.tsx                providers, i18n init, session bootstrap
  index.tsx                  redirect → tabs (browse is public)
  (auth)/login | register
  (tabs)/                    bottom tabs (tourist)
    index      Home          banners + categories + featured tours
    search     Search        query + category filters
    bookings   My bookings   (sign-in gated)
    favorites  Favorites     (sign-in gated)
    profile    Profile       language switch, agency mode, sign out
  tour/[slug]                tour detail (gallery, departures, past-gallery)
  booking/[slug]             booking flow with phone OTP (3 steps)
src/
  api/client.ts              axios + envelope + refresh
  stores/auth.ts             zustand auth (login/register/google/bootstrap)
  i18n/                      i18next init + backend hydration + RTL
  lib/                       secureStore (tokens), query client, formatters
  components/                ui primitives, TourCard, SignInPrompt
  types/api.ts               envelope + domain DTOs
```

## Two faces in one app

The tourist experience is the default. When the signed-in user holds an agency
role, the **Profile** screen surfaces an "Agency mode" entry point (the agency
tour/booking management screens are the next build-out — the backend endpoints
already exist under `/agency/**`).

## Booking flow (phone OTP)

`booking/[slug]` is a 3-step flow matching the backend contract:

1. **Details** — pick a departure, enter lead-traveler name + Lebanese phone
2. **OTP** — `POST /otp/request` (WhatsApp) → enter code → `POST /otp/verify`
3. **Confirm** — `POST /bookings` (the backend enforces phone verification at
   booking time)

Idempotency and offline-safe submission (queued retries with idempotency keys)
are the natural hardening step for flaky connectivity.

## Notes / next

- **OneSignal push** + device registration (`POST /me/devices`) — add the
  `onesignal-expo-plugin` and wire player-id registration on login.
- **Google sign-in** — `auth.googleSignIn(idToken)` is wired in the store; add
  `expo-auth-session`/Google provider to obtain the ID token.
- Full RTL re-layout requires an app reload after switching to Arabic
  (`I18nManager.forceRTL`); the language switch sets it for next launch.
- Agency-mode screens, favorites toggle persistence, e-voucher/QR, and the shop
  tab (when `modules.shop` is enabled) are follow-ups.
```
