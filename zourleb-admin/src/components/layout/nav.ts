// Navigation config. Each item declares which roles may see it; the sidebar
// filters by the signed-in user's roles, giving the two menus (admin/agency).

export interface NavItem {
  to: string;
  labelKey: string;
  icon: string;
  roles: ("super_admin" | "agency_owner" | "agency_staff")[];
}

export const NAV: NavItem[] = [
  { to: "/", labelKey: "nav.dashboard", icon: "▦", roles: ["super_admin", "agency_owner", "agency_staff"] },
  { to: "/agencies", labelKey: "nav.agencies", icon: "🏢", roles: ["super_admin"] },
  { to: "/users", labelKey: "nav.users", icon: "👤", roles: ["super_admin"] },
  { to: "/reviews", labelKey: "nav.reviews", icon: "★", roles: ["super_admin"] },
  { to: "/boosts", labelKey: "nav.boosts", icon: "📣", roles: ["super_admin"] },
  { to: "/payments", labelKey: "nav.payments", icon: "💳", roles: ["super_admin"] },
  { to: "/banners", labelKey: "nav.banners", icon: "🖼", roles: ["super_admin"] },
  { to: "/audit-logs", labelKey: "nav.audit_logs", icon: "🔍", roles: ["super_admin"] },
  { to: "/languages", labelKey: "nav.languages", icon: "🌐", roles: ["super_admin"] },
  { to: "/translations", labelKey: "nav.translations", icon: "🔤", roles: ["super_admin"] },
  { to: "/settings", labelKey: "nav.settings", icon: "⚙", roles: ["super_admin"] },
  { to: "/tours", labelKey: "nav.all_tours", icon: "🧭", roles: ["super_admin"] },
  // Agency portal items (rendered when the user has an agency role)
  { to: "/agency/tours", labelKey: "nav.tours", icon: "🧭", roles: ["agency_owner", "agency_staff"] },
  { to: "/agency/bookings", labelKey: "nav.bookings", icon: "🎫", roles: ["agency_owner", "agency_staff"] },
];
