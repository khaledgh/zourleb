import * as SecureStore from "expo-secure-store";
import { Platform } from "react-native";

// Token persistence. expo-secure-store isn't available on web, so we fall back
// to localStorage there (the app's primary targets are iOS/Android).

const ACCESS = "zb_access";
const REFRESH = "zb_refresh";
const LOCALE = "zb_locale";

const isWeb = Platform.OS === "web";

async function setItem(key: string, value: string) {
  if (isWeb) {
    globalThis.localStorage?.setItem(key, value);
    return;
  }
  await SecureStore.setItemAsync(key, value);
}

async function getItem(key: string): Promise<string | null> {
  if (isWeb) {
    return globalThis.localStorage?.getItem(key) ?? null;
  }
  return SecureStore.getItemAsync(key);
}

async function deleteItem(key: string) {
  if (isWeb) {
    globalThis.localStorage?.removeItem(key);
    return;
  }
  await SecureStore.deleteItemAsync(key);
}

export const tokenStore = {
  async get() {
    const [access, refresh] = await Promise.all([
      getItem(ACCESS),
      getItem(REFRESH),
    ]);
    return { access, refresh };
  },
  async set(access: string, refresh: string) {
    await Promise.all([setItem(ACCESS, access), setItem(REFRESH, refresh)]);
  },
  async clear() {
    await Promise.all([deleteItem(ACCESS), deleteItem(REFRESH)]);
  },
};

export const localeStore = {
  get: () => getItem(LOCALE),
  set: (code: string) => setItem(LOCALE, code),
};
