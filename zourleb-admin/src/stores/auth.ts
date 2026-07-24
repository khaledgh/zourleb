import { create } from "zustand";
import { apiGet, apiPost } from "@/api/client";
import { tokenStore } from "@/lib/tokens";
import type { AuthResponse, User } from "@/types/api";

interface AuthState {
  user: User | null;
  ready: boolean; // initial session check completed
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  bootstrap: () => Promise<void>;
  hasRole: (...roles: string[]) => boolean;
}

export const useAuth = create<AuthState>((set, get) => ({
  user: null,
  ready: false,

  async login(email, password) {
    const data = await apiPost<AuthResponse>("/auth/login", { email, password });
    tokenStore.set(data.access_token, data.refresh_token);
    set({ user: data.user });
  },

  async logout() {
    const refresh = tokenStore.refresh();
    try {
      if (refresh) await apiPost("/auth/logout", { refresh_token: refresh });
    } catch {
      // ignore — clear locally regardless
    }
    tokenStore.clear();
    set({ user: null });
  },

  // bootstrap restores the session from a stored token on app load.
  async bootstrap() {
    if (!tokenStore.access()) {
      set({ ready: true });
      return;
    }
    try {
      const env = await apiGet<User>("/me");
      set({ user: env.data ?? null, ready: true });
    } catch {
      tokenStore.clear();
      set({ user: null, ready: true });
    }
  },

  hasRole(...roles) {
    const u = get().user;
    if (!u) return false;
    return roles.some((r) => u.roles.includes(r));
  },
}));

export const isSuperAdmin = (u: User | null) =>
  !!u && u.roles.includes("super_admin");
export const isAgency = (u: User | null) =>
  !!u && (u.roles.includes("agency_owner") || u.roles.includes("agency_staff"));
