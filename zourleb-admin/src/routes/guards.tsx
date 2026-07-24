import { Navigate, useLocation } from "react-router-dom";
import type { ReactNode } from "react";
import { useAuth } from "@/stores/auth";
import { Spinner } from "@/components/ui/primitives";

// RequireAuth gates the authenticated app; RequireRole gates by role.

export function RequireAuth({ children }: { children: ReactNode }) {
  const { user, ready } = useAuth();
  const location = useLocation();
  if (!ready) return <Spinner label="…" />;
  if (!user) return <Navigate to="/login" state={{ from: location }} replace />;
  return <>{children}</>;
}

export function RequireRole({
  roles,
  children,
}: {
  roles: string[];
  children: ReactNode;
}) {
  const { user } = useAuth();
  const ok = user && roles.some((r) => user.roles.includes(r));
  if (!ok) return <Navigate to="/" replace />;
  return <>{children}</>;
}
