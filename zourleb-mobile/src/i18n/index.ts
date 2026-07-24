import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import { I18nManager } from "react-native";
import * as Localization from "expo-localization";
import { apiGet, setLocale as setApiLocale } from "@/api/client";
import { localeStore } from "@/lib/secureStore";

// Built-in strings render before the backend bundle loads and serve as an
// offline fallback.
const fallback: Record<string, Record<string, string>> = {
  en: {
    "tab.home": "Home",
    "tab.search": "Search",
    "tab.bookings": "Bookings",
    "tab.favorites": "Favorites",
    "tab.profile": "Profile",
    "action.book": "Book now",
    "action.signin": "Sign in",
    "action.signout": "Sign out",
    "booking.verifyPhone": "Verify your phone",
    "common.loading": "Loading…",
    "common.from": "From",
  },
  ar: {
    "tab.home": "الرئيسية",
    "tab.search": "بحث",
    "tab.bookings": "حجوزاتي",
    "tab.favorites": "المفضلة",
    "tab.profile": "حسابي",
    "action.book": "احجز الآن",
    "action.signin": "تسجيل الدخول",
    "action.signout": "تسجيل الخروج",
    "booking.verifyPhone": "تحقق من رقم هاتفك",
    "common.loading": "جارٍ التحميل…",
    "common.from": "ابتداءً من",
  },
};

const RTL_LOCALES = new Set(["ar", "he", "fa", "ur"]);

function deviceLocale(): string {
  const tag = Localization.getLocales()[0]?.languageCode ?? "ar";
  return ["ar", "fr", "en"].includes(tag) ? tag : "ar";
}

export async function initI18n() {
  const stored = (await localeStore.get()) || deviceLocale();
  await i18n.use(initReactI18next).init({
    lng: stored,
    fallbackLng: "en",
    resources: {
      en: { translation: fallback.en },
      ar: { translation: fallback.ar },
    },
    interpolation: { escapeValue: false },
  });
  setApiLocale(stored);
  applyDirection(stored);
  void hydrate(stored);
  return i18n;
}

// applyDirection toggles RN's RTL. Note: a full RTL flip needs an app reload;
// expo-updates/Dev reload applies it. We set it so it takes effect next launch.
export function applyDirection(locale: string) {
  const rtl = RTL_LOCALES.has(locale);
  if (I18nManager.isRTL !== rtl) {
    I18nManager.allowRTL(rtl);
    I18nManager.forceRTL(rtl);
  }
}

interface Bundle {
  locale: string;
  is_rtl: boolean;
  translations: Record<string, Record<string, string>>;
}

async function hydrate(locale: string) {
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
    // keep fallback
  }
}

export async function changeLocale(locale: string) {
  await localeStore.set(locale);
  setApiLocale(locale);
  if (!i18n.hasResourceBundle(locale, "translation") && fallback[locale]) {
    i18n.addResourceBundle(locale, "translation", fallback[locale]);
  }
  await i18n.changeLanguage(locale);
  applyDirection(locale);
  await hydrate(locale);
}

export default i18n;
