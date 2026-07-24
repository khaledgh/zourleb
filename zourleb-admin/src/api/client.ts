import axios, {
  AxiosError,
  type AxiosInstance,
  type InternalAxiosRequestConfig,
} from "axios";
import { tokenStore } from "@/lib/tokens";
import type { ApiError, AuthResponse, Envelope } from "@/types/api";

// Locale is read lazily so language switches affect subsequent requests.
let currentLocale = localStorage.getItem("zb_locale") || "ar";
export function setApiLocale(locale: string) {
  currentLocale = locale;
  localStorage.setItem("zb_locale", locale);
}
export function getApiLocale() {
  return currentLocale;
}

// A normalized error our UI layer can switch on by `code`.
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
  baseURL: "/api/v1",
  headers: { "Content-Type": "application/json" },
});

// Request: attach bearer + Accept-Language.
raw.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = tokenStore.access();
  if (token) config.headers.Authorization = `Bearer ${token}`;
  config.headers["Accept-Language"] = currentLocale;
  return config;
});

// --- Single-flight refresh ---
let refreshing: Promise<string | null> | null = null;

async function doRefresh(): Promise<string | null> {
  const refresh = tokenStore.refresh();
  if (!refresh) return null;
  try {
    const res = await axios.post<Envelope<AuthResponse>>(
      "/api/v1/auth/refresh",
      { refresh_token: refresh },
    );
    const data = res.data.data;
    if (!data) return null;
    tokenStore.set(data.access_token, data.refresh_token);
    return data.access_token;
  } catch {
    tokenStore.clear();
    return null;
  }
}

// Response: on 401 (once), refresh and retry; otherwise normalize the error.
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
      const newToken = await refreshing;
      refreshing = null;
      if (newToken) {
        original.headers.Authorization = `Bearer ${newToken}`;
        return raw(original);
      }
      // refresh failed → bounce to login
      tokenStore.clear();
      if (location.pathname !== "/login") location.href = "/login";
    }

    const env = error.response?.data;
    if (env?.error) throw new ApiException(env.error, status);
    throw new ApiException(
      { code: "NETWORK", message: error.message || "Network error" },
      status,
    );
  },
);

// --- Typed helpers that unwrap the envelope ---

export async function apiGet<T>(
  url: string,
  params?: Record<string, unknown>,
): Promise<Envelope<T>> {
  const res = await raw.get<Envelope<T>>(url, { params });
  return res.data;
}

export async function apiPost<T>(url: string, body?: unknown): Promise<T> {
  const res = await raw.post<Envelope<T>>(url, body);
  return res.data.data as T;
}

export async function apiPatch<T>(url: string, body?: unknown): Promise<T> {
  const res = await raw.patch<Envelope<T>>(url, body);
  return res.data.data as T;
}

export async function apiPut<T>(url: string, body?: unknown): Promise<T> {
  const res = await raw.put<Envelope<T>>(url, body);
  return res.data.data as T;
}

export async function apiDelete<T>(url: string): Promise<T> {
  const res = await raw.delete<Envelope<T>>(url);
  return res.data.data as T;
}

// apiList returns both items and pagination meta from a list envelope.
export async function apiList<T>(
  url: string,
  params?: Record<string, unknown>,
): Promise<{ items: T[]; meta?: Envelope<T[]>["meta"] }> {
  const env = await apiGet<T[]>(url, params);
  return { items: env.data ?? [], meta: env.meta };
}

export { raw as axiosClient };
