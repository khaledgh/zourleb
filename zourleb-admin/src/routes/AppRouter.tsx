import { createBrowserRouter, Navigate } from "react-router-dom";
import { AppLayout } from "@/components/layout/AppLayout";
import { RequireAuth, RequireRole } from "./guards";
import { LoginPage } from "@/pages/Login";
import { DashboardPage } from "@/pages/Dashboard";
import { AgenciesPage } from "@/pages/Agencies";
import { UsersPage } from "@/pages/Users";
import { SettingsPage } from "@/pages/Settings";
import { LanguagesPage } from "@/pages/Languages";
import { TranslationsPage } from "@/pages/Translations";
import { BannersPage } from "@/pages/Banners";
import { ReviewsPage } from "@/pages/Reviews";
import { BoostsPage } from "@/pages/Boosts";
import { AdminToursPage } from "@/pages/AdminTours";
import { AgencyToursPage } from "@/pages/AgencyTours";
import { AgencyBookingsPage } from "@/pages/AgencyBookings";

const admin = (el: JSX.Element) => (
  <RequireRole roles={["super_admin"]}>{el}</RequireRole>
);
const agency = (el: JSX.Element) => (
  <RequireRole roles={["agency_owner", "agency_staff"]}>{el}</RequireRole>
);

export const router = createBrowserRouter([
  { path: "/login", element: <LoginPage /> },
  {
    path: "/",
    element: (
      <RequireAuth>
        <AppLayout />
      </RequireAuth>
    ),
    children: [
      { index: true, element: <DashboardPage /> },
      // super admin
      { path: "agencies", element: admin(<AgenciesPage />) },
      { path: "users", element: admin(<UsersPage />) },
      { path: "settings", element: admin(<SettingsPage />) },
      { path: "languages", element: admin(<LanguagesPage />) },
      { path: "translations", element: admin(<TranslationsPage />) },
      { path: "banners", element: admin(<BannersPage />) },
      { path: "reviews", element: admin(<ReviewsPage />) },
      { path: "boosts", element: admin(<BoostsPage />) },
      { path: "tours", element: admin(<AdminToursPage />) },
      // agency portal
      { path: "agency/tours", element: agency(<AgencyToursPage />) },
      { path: "agency/bookings", element: agency(<AgencyBookingsPage />) },
      { path: "*", element: <Navigate to="/" replace /> },
    ],
  },
]);
