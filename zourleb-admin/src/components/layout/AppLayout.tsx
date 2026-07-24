import { useState } from "react";
import { NavLink, Outlet, useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import clsx from "clsx";
import { useAuth } from "@/stores/auth";
import { NAV } from "./nav";
import { LocaleSwitcher } from "./LocaleSwitcher";

function getNavIcon(labelKey: string) {
  switch (labelKey) {
    case "nav.dashboard":
      return (
        <svg className="w-5 h-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v4a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v4a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v4a2 2 0 01-2 2H6a2 2 0 01-2-2v-4zM14 16a2 2 0 012-2h2a2 2 0 012 2v4a2 2 0 01-2 2h-2a2 2 0 01-2-2v-4z" />
        </svg>
      );
    case "nav.agencies":
      return (
        <svg className="w-5 h-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
        </svg>
      );
    case "nav.users":
      return (
        <svg className="w-5 h-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
        </svg>
      );
    case "nav.reviews":
      return (
        <svg className="w-5 h-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.907c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
        </svg>
      );
    case "nav.boosts":
      return (
        <svg className="w-5 h-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
        </svg>
      );
    case "nav.banners":
      return (
        <svg className="w-5 h-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
      );
    case "nav.languages":
      return (
        <svg className="w-5 h-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
        </svg>
      );
    case "nav.translations":
      return (
        <svg className="w-5 h-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M3 5h12M9 3v2m1.048 9.5A18.022 18.022 0 016.412 9m6.088 9h7M11 21l5-10 5 10M12.751 5c-.313 1.56-.98 3.08-1.993 4.488a18.022 18.022 0 01-3.376-2.738" />
        </svg>
      );
    case "nav.settings":
      return (
        <svg className="w-5 h-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
          <path strokeLinecap="round" strokeLinejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
        </svg>
      );
    case "nav.tours":
      return (
        <svg className="w-5 h-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7" />
        </svg>
      );
    case "nav.bookings":
      return (
        <svg className="w-5 h-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z" />
        </svg>
      );
    default:
      return (
        <svg className="w-5 h-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
        </svg>
      );
  }
}

export function AppLayout() {
  const { t } = useTranslation();
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);

  const roles = user?.roles ?? [];
  const items = NAV.filter((n) => n.roles.some((r) => roles.includes(r)));

  async function onLogout() {
    await logout();
    navigate("/login");
  }

  // Active page title calculation
  const activeItem = items.find((item) => {
    if (item.to === "/") {
      return location.pathname === "/";
    }
    return location.pathname.startsWith(item.to);
  });
  const pageTitle = activeItem ? t(activeItem.labelKey) : "Documents";

  return (
    <div className="flex min-h-screen items-center justify-center p-0">
      <div className="admin-panel-viewport flex w-full h-screen bg-white overflow-hidden">
        
        {/* Sidebar */}
        <aside 
          className={clsx(
            "hidden md:flex shrink-0 flex-col border-e border-slate-100 bg-white p-6 justify-between transition-all duration-300",
            isSidebarCollapsed ? "w-20 px-3 items-center" : "w-64"
          )}
        >
          <div className="flex flex-col gap-8 w-full">
            {/* Logo Section */}
            <div className={clsx("flex items-center gap-3 px-2", isSidebarCollapsed ? "justify-center" : "justify-between")}>
              <div className="flex items-center gap-3">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-blue-500 text-white shadow-lg shadow-blue-500/30">
                  <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                    <path strokeLinecap="round" strokeLinejoin="round" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                </div>
                {!isSidebarCollapsed && (
                  <span className="font-display text-xl font-bold tracking-tight text-slate-800 animate-fadeIn">Zourleb</span>
                )}
              </div>
              
              {/* Toggle Trigger */}
              <button 
                onClick={() => setIsSidebarCollapsed(!isSidebarCollapsed)} 
                className="p-1.5 rounded-xl border border-slate-100 bg-slate-50 text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition-all active:scale-95 shrink-0 hidden md:block"
                title={isSidebarCollapsed ? "Expand sidebar" : "Collapse sidebar"}
              >
                <svg className={clsx("w-4 h-4 transition-transform duration-200", isSidebarCollapsed && "rotate-180")} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M15 19l-7-7 7-7" />
                </svg>
              </button>
            </div>

            {/* Navigation links */}
            <nav className="space-y-1.5 w-full">
              {items.map((item) => (
                <NavLink
                  key={item.to}
                  to={item.to}
                  end={item.to === "/"}
                  className={({ isActive }) =>
                    clsx(
                      "flex items-center gap-3 rounded-2xl text-sm font-semibold transition-all duration-200 w-full",
                      isSidebarCollapsed ? "justify-center p-3" : "px-4 py-3",
                      isActive
                        ? "bg-[#2F80ED] text-white shadow-lg shadow-blue-500/20 scale-[1.02]"
                        : "text-slate-400 hover:text-slate-700 hover:bg-slate-50",
                    )
                  }
                  title={isSidebarCollapsed ? t(item.labelKey) : undefined}
                >
                  {getNavIcon(item.labelKey)}
                  {!isSidebarCollapsed && (
                    <span className="animate-fadeIn">{t(item.labelKey)}</span>
                  )}
                </NavLink>
              ))}
            </nav>
          </div>

          {/* Bottom Area */}
          <div className="flex flex-col gap-6 w-full">
            {/* Logout button */}
            <button
              onClick={onLogout}
              className={clsx(
                "flex items-center gap-3 rounded-2xl text-sm font-semibold text-slate-400 hover:text-red-500 hover:bg-red-50/50 transition-all duration-200 w-full",
                isSidebarCollapsed ? "justify-center p-3" : "px-4 py-3"
              )}
              title={isSidebarCollapsed ? t("action.logout") : undefined}
            >
              <svg className="w-5 h-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
              {!isSidebarCollapsed && <span className="animate-fadeIn">{t("action.logout")}</span>}
            </button>
          </div>
        </aside>

        {/* Main Content Area */}
        <div className="flex flex-1 flex-col overflow-hidden bg-slate-50/30">
          <header className="flex h-20 shrink-0 items-center justify-between border-b border-slate-100 bg-white px-8">
            <div className="flex items-center gap-4">
              {/* Mobile layout toggle */}
              <div className="md:hidden flex items-center gap-2">
                <button
                  onClick={() => setIsSidebarCollapsed(!isSidebarCollapsed)}
                  className="flex h-9 w-9 items-center justify-center rounded-lg bg-blue-500 text-white"
                >
                  <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M4 6h16M4 12h16M4 18h16" />
                  </svg>
                </button>
              </div>
              
              <div className="text-lg font-bold text-slate-800 font-display">
                {pageTitle}
              </div>
            </div>

            <div className="flex items-center gap-6">
              {/* Search Bar matching mockup */}
              <div className="hidden sm:flex items-center gap-2">
                <div className="relative flex items-center bg-slate-50 border border-slate-200/50 rounded-xl px-3.5 py-2 w-64 shadow-inner">
                  <svg className="w-4 h-4 text-slate-400 mr-2 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                  <input
                    type="text"
                    placeholder="Search…"
                    className="bg-transparent text-sm w-full text-slate-700 outline-none placeholder:text-slate-400/80"
                  />
                </div>
                <button className="bg-blue-500 hover:bg-blue-600 rounded-xl p-2.5 text-white shadow-md shadow-blue-500/10 hover:shadow-lg hover:shadow-blue-500/20 active:scale-95 transition-all">
                  <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                </button>
              </div>

              {/* Locale switcher */}
              <LocaleSwitcher />

              {/* User Dropdown matching mockup */}
              <div className="flex items-center gap-3 border-l border-slate-100 pl-6 cursor-pointer group">
                <img
                  src="/user_avatar.png"
                  className="w-10 h-10 rounded-full object-cover border-2 border-slate-100 group-hover:border-blue-500 transition-all duration-200 shadow-sm"
                  alt="avatar"
                />
                <div className="hidden lg:block text-start">
                  <div className="text-sm font-bold text-slate-800 font-display leading-tight">{user?.name || "Angelina Joli"}</div>
                  <div className="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">
                    {roles.includes("super_admin") ? "Super Admin" : "Agency Owner"}
                  </div>
                </div>
                <svg className="w-4 h-4 text-slate-400 group-hover:text-slate-600 transition-all ml-1 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M19 9l-7 7-7-7" />
                </svg>
              </div>
            </div>
          </header>

          <main className="flex-1 overflow-auto p-8 relative">
            {/* Mobile collapsible overlay */}
            {isSidebarCollapsed && (
              <div className="md:hidden absolute inset-0 z-40 bg-slate-900/30 backdrop-blur-sm" onClick={() => setIsSidebarCollapsed(false)} />
            )}
            
            {/* Mobile Sidebar Flyout */}
            <div 
              className={clsx(
                "md:hidden fixed top-0 bottom-0 left-0 z-50 w-64 bg-white p-6 border-r border-slate-100 flex flex-col justify-between transition-transform duration-300 transform shadow-2xl",
                isSidebarCollapsed ? "translate-x-0" : "-translate-x-full"
              )}
            >
              <div className="flex flex-col gap-8 w-full">
                <div className="flex items-center justify-between px-2">
                  <div className="flex items-center gap-3">
                    <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-blue-500 text-white shadow-lg shadow-blue-500/30">
                      <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                        <path strokeLinecap="round" strokeLinejoin="round" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                      </svg>
                    </div>
                    <span className="font-display text-xl font-bold tracking-tight text-slate-800">Zourleb</span>
                  </div>
                  <button 
                    onClick={() => setIsSidebarCollapsed(false)}
                    className="p-1 rounded-lg hover:bg-slate-50 text-slate-400"
                  >
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
                      <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>

                <nav className="space-y-1.5 w-full">
                  {items.map((item) => (
                    <NavLink
                      key={item.to}
                      to={item.to}
                      end={item.to === "/"}
                      className={({ isActive }) =>
                        clsx(
                          "flex items-center gap-3 rounded-2xl px-4 py-3 text-sm font-semibold transition-all duration-200 w-full",
                          isActive
                            ? "bg-[#2F80ED] text-white shadow-lg shadow-blue-500/20"
                            : "text-slate-400 hover:text-slate-700 hover:bg-slate-50",
                        )
                      }
                      onClick={() => setIsSidebarCollapsed(false)}
                    >
                      {getNavIcon(item.labelKey)}
                      <span>{t(item.labelKey)}</span>
                    </NavLink>
                  ))}
                </nav>
              </div>

              <div className="flex flex-col gap-6 w-full">
                <button
                  onClick={onLogout}
                  className="flex items-center gap-3 rounded-2xl px-4 py-3 text-sm font-semibold text-slate-400 hover:text-red-500 hover:bg-red-50/50 transition-all duration-200 w-full"
                >
                  <svg className="w-5 h-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                  </svg>
                  <span>{t("action.logout")}</span>
                </button>
              </div>
            </div>

            <Outlet />
          </main>
        </div>
      </div>
    </div>
  );
}
