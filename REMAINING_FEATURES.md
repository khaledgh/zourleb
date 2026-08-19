# Remaining Features to Build

This document tracks missing features across the Zourleb platform after the current implementation pass.

## Backend

- Payment adapters
  - Stripe card payments
  - Local / offline payment confirmation flow polish
- Storage
  - S3-compatible storage adapter (currently `local` only)
- OneSignal push notification integration
  - Device registration (`POST /me/devices`) is in router; actual push dispatch needs OneSignal wiring
- E-voucher / QR code generation for confirmed bookings
- In-app chat between tourists and agencies
- Loyalty / referral program
- Map view for tours and agencies
- Blog / events module
- Subscription tiers for agencies
- Invite member endpoint for agency staff (`POST /agency/members/invite` is missing)
- Add tour sub-resources endpoints
  - `POST /agency/tours/:id/prices` and `POST /agency/tours/:id/images` exist but are not exposed in mobile
  - `POST /agency/tours/:id/departures` exists but mobile add-departure screen is not built

## Admin (`zourleb-admin`)

- Analytics dashboard full page with charts
- Category / Region full CRUD pages (list exists, create/edit may be missing)
- Shop product management for agencies
- E-voucher generation and booking QR codes
- Invite agency member / staff management
- Push-notification composer / campaign view
- Subscription and billing management

## Mobile (`zourleb-mobile`)

- Google sign-in button (`auth.googleSignIn` is wired in store; needs `expo-auth-session` UI)
- OneSignal push registration on login
- Agency member invitation (blocked by missing backend endpoint)
- Add departure / price / image screens for agency tours
- Shop tab (gated by `modules.shop`)
- Favorites persistence (Home uses favorite ids; Search and Tour detail were reverted)
- E-voucher / QR display in My Bookings
- In-app chat
- Map view
- Full RTL re-layout on Arabic without restart

## Infrastructure / Tooling

- CI/CD with `tsc`, `go test`, and build for all three packages
- End-to-end test suite
- Production Docker / deployment guide
