import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import { apiGet, setApiLocale } from "@/api/client";

// Minimal built-in strings so the UI renders before the backend bundle loads
// (and as an offline fallback).
const fallback: Record<string, Record<string, string>> = {
  en: {
    "nav.dashboard": "Dashboard",
    "nav.agencies": "Agencies",
    "nav.tours": "Tours",
    "nav.bookings": "Bookings",
    "nav.boosts": "Boosts",
    "nav.banners": "Banners",
    "nav.languages": "Languages",
    "nav.translations": "Translations",
    "nav.users": "Users",
    "nav.settings": "Settings",
    "nav.reviews": "Reviews",
    "action.logout": "Sign out",
    "action.save": "Save",
    "action.approve": "Approve",
    "action.suspend": "Suspend",
    "common.loading": "Loading…",
  },
  ar: {
    "nav.dashboard": "لوحة التحكم",
    "nav.agencies": "الوكالات",
    "nav.tours": "الجولات",
    "nav.bookings": "الحجوزات",
    "nav.boosts": "الإعلانات",
    "nav.banners": "اللافتات",
    "nav.languages": "اللغات",
    "nav.translations": "الترجمات",
    "nav.users": "المستخدمون",
    "nav.settings": "الإعدادات",
    "nav.reviews": "المراجعات",
    "action.logout": "تسجيل الخروج",
    "action.save": "حفظ",
    "action.approve": "موافقة",
    "action.suspend": "تعليق",
    "common.loading": "جارٍ التحميل…",
  },
};

export const RTL_LOCALES = new Set(["ar", "he", "fa", "ur"]);

export function applyDirection(locale: string) {
  const dir = RTL_LOCALES.has(locale) ? "rtl" : "ltr";
  document.documentElement.dir = dir;
  document.documentElement.lang = locale;
}

export async function initI18n() {
  const stored = localStorage.getItem("zb_locale") || "ar";
  await i18n.use(initReactI18next).init({
    lng: stored,
    fallbackLng: "en",
    resources: {
      en: { translation: fallback.en },
      ar: { translation: fallback.ar },
    },
    interpolation: { escapeValue: false },
  });
  applyDirection(stored);
  // hydrate from the backend bundle (admin namespace flattened)
  void hydrateFromBackend(stored);
  return i18n;
}

interface Bundle {
  locale: string;
  is_rtl: boolean;
  translations: Record<string, Record<string, string>>;
}

// hydrateFromBackend merges server-managed UI strings over the fallback.
async function hydrateFromBackend(locale: string) {
  try {
    const env = await apiGet<Bundle>(`/i18n/${locale}`);
    const b = env.data;
    if (!b) return;
    const flat: Record<string, string> = {};
    for (const [ns, kv] of Object.entries(b.translations)) {
      for (const [k, v] of Object.entries(kv)) flat[`${ns}.${k}`] = v;
    }
    i18n.addResourceBundle(locale, "translation", flat, true, true);
  } catch {
    // keep fallback strings
  }
}

export async function changeLocale(locale: string) {
  setApiLocale(locale);
  localStorage.setItem("zb_locale", locale);
  if (!i18n.hasResourceBundle(locale, "translation") && fallback[locale]) {
    i18n.addResourceBundle(locale, "translation", fallback[locale]);
  }
  await i18n.changeLanguage(locale);
  applyDirection(locale);
  await hydrateFromBackend(locale);
}

export default i18n;
