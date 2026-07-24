import axios, {
  AxiosError,
  type AxiosInstance,
  type InternalAxiosRequestConfig,
} from "axios";
import Constants from "expo-constants";
import { tokenStore } from "@/lib/secureStore";
import type { ApiError, AuthResponse, Envelope } from "@/types/api";

// Base URL comes from app.json `extra.apiBaseUrl`. On a physical device this
// must be the dev machine's LAN IP, not localhost.
const baseURL =
  (Constants.expoConfig?.extra?.apiBaseUrl as string) ??
  "http://localhost:8080/api/v1";

// In-memory token cache (secure-store is async; interceptors are sync).
let accessToken: string | null = null;
let refreshToken: string | null = null;
let currentLocale = "ar";

// onUnauthorized is set by the auth store so a failed refresh can reset state.
let onUnauthorized: (() => void) | null = null;
export function setUnauthorizedHandler(fn: () => void) {
  onUnauthorized = fn;
}

export function setLocale(locale: string) {
  currentLocale = locale;
}

export async function hydrateTokens() {
  const { access, refresh } = await tokenStore.get();
  accessToken = access;
  refreshToken = refresh;
  return { access, refresh };
}

export async function setAuthTokens(access: string, refresh: string) {
  accessToken = access;
  refreshToken = refresh;
  await tokenStore.set(access, refresh);
}

export async function clearAuthTokens() {
  accessToken = null;
  refreshToken = null;
  await tokenStore.clear();
}

export class ApiException extends Error {
  code: string;
  fields?: Record<string, string>;
  status: number;
  constructor(err: ApiError, status: number) {
    super(err.message);
    this.code = err.code;
    this.fields = err.fields;
    this.status = status;
  }
}

const raw: AxiosInstance = axios.create({
  baseURL,
  headers: { "Content-Type": "application/json" },
});

raw.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  if (accessToken) config.headers.Authorization = `Bearer ${accessToken}`;
  config.headers["Accept-Language"] = currentLocale;
  return config;
});

let refreshing: Promise<string | null> | null = null;

async function doRefresh(): Promise<string | null> {
  if (!refreshToken) return null;
  try {
    const res = await axios.post<Envelope<AuthResponse>>(
      `${baseURL}/auth/refresh`,
      { refresh_token: refreshToken },
    );
    const data = res.data.data;
    if (!data) return null;
    await setAuthTokens(data.access_token, data.refresh_token);
    return data.access_token;
  } catch {
    await clearAuthTokens();
    return null;
  }
}

raw.interceptors.response.use(
  (r) => r,
  async (error: AxiosError<Envelope<unknown>>) => {
    const original = error.config as InternalAxiosRequestConfig & {
      _retried?: boolean;
    };
    const status = error.response?.status ?? 0;

    if (status === 401 && original && !original._retried) {
      original._retried = true;
      refreshing = refreshing ?? doRefresh();
      const token = await refreshing;
      refreshing = null;
      if (token) {
        original.headers.Authorization = `Bearer ${token}`;
        return raw(original);
      }
      onUnauthorized?.();
    }

    const env = error.response?.data;
    if (env?.error) throw new ApiException(env.error, status);
    throw new ApiException(
      { code: "NETWORK", message: error.message || "Network error" },
      status,
    );
  },
);

// --- Typed helpers ---

export async function apiGet<T>(
  url: string,
  params?: Record<string, unknown>,
): Promise<Envelope<T>> {
  const res = await raw.get<Envelope<T>>(url, { params });
  return res.data;
}

export async function apiData<T>(
  url: string,
  params?: Record<string, unknown>,
): Promise<T> {
  return (await apiGet<T>(url, params)).data as T;
}

export async function apiPost<T>(url: string, body?: unknown): Promise<T> {
  const res = await raw.post<Envelope<T>>(url, body);
  return res.data.data as T;
}

export async function apiPatch<T>(url: string, body?: unknown): Promise<T> {
  const res = await raw.patch<Envelope<T>>(url, body);
  return res.data.data as T;
}

export async function apiList<T>(
  url: string,
  params?: Record<string, unknown>,
): Promise<{ items: T[]; meta?: Envelope<T[]>["meta"] }> {
  const env = await apiGet<T[]>(url, params);
  return { items: env.data ?? [], meta: env.meta };
}
