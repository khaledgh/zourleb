// Token persistence in localStorage. The access token is short-lived; the
// refresh token rotates on use (the backend revokes the presented one).

const ACCESS = "zb_access";
const REFRESH = "zb_refresh";

export const tokenStore = {
  access: () => localStorage.getItem(ACCESS),
  refresh: () => localStorage.getItem(REFRESH),
  set(access: string, refresh: string) {
    localStorage.setItem(ACCESS, access);
    localStorage.setItem(REFRESH, refresh);
  },
  clear() {
    localStorage.removeItem(ACCESS);
    localStorage.removeItem(REFRESH);
  },
};
