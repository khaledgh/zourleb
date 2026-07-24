import { create } from "zustand";
import {
  apiData,
  apiPost,
  clearAuthTokens,
  hydrateTokens,
  setAuthTokens,
  setUnauthorizedHandler,
} from "@/api/client";
import type { AuthResponse, User } from "@/types/api";

interface AuthState {
  user: User | null;
  ready: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (name: string, email: string, password: string) => Promise<void>;
  googleSignIn: (idToken: string) => Promise<void>;
  logout: () => Promise<void>;
  bootstrap: () => Promise<void>;
  refreshMe: () => Promise<void>;
  isAgency: () => boolean;
}

export const useAuth = create<AuthState>((set, get) => ({
  user: null,
  ready: false,

  async login(email, password) {
    const data = await apiPost<AuthResponse>("/auth/login", { email, password });
    await setAuthTokens(data.access_token, data.refresh_token);
    set({ user: data.user });
  },

  async register(name, email, password) {
    const data = await apiPost<AuthResponse>("/auth/register", {
      name,
      email,
      password,
    });
    await setAuthTokens(data.access_token, data.refresh_token);
    set({ user: data.user });
  },

  async googleSignIn(idToken) {
    const data = await apiPost<AuthResponse>("/auth/google", { id_token: idToken });
    await setAuthTokens(data.access_token, data.refresh_token);
    set({ user: data.user });
  },

  async logout() {
    const { refresh } = await hydrateTokens();
    try {
      if (refresh) await apiPost("/auth/logout", { refresh_token: refresh });
    } catch {
      // ignore
    }
    await clearAuthTokens();
    set({ user: null });
  },

  async bootstrap() {
    setUnauthorizedHandler(() => set({ user: null }));
    const { access } = await hydrateTokens();
    if (!access) {
      set({ ready: true });
      return;
    }
    try {
      const user = await apiData<User>("/me");
      set({ user, ready: true });
    } catch {
      await clearAuthTokens();
      set({ user: null, ready: true });
    }
  },

  async refreshMe() {
    try {
      const user = await apiData<User>("/me");
      set({ user });
    } catch {
      // keep existing
    }
  },

  isAgency() {
    const u = get().user;
    return (
      !!u &&
      (u.roles.includes("agency_owner") || u.roles.includes("agency_staff"))
    );
  },
}));
