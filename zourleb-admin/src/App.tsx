import { useEffect, useState } from "react";
import { RouterProvider } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { router } from "@/routes/AppRouter";
import { useAuth } from "@/stores/auth";
import { ToastViewport } from "@/components/ui/toast";
import { Spinner } from "@/components/ui/primitives";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: { retry: 1, refetchOnWindowFocus: false, staleTime: 30_000 },
  },
});

export default function App() {
  const bootstrap = useAuth((s) => s.bootstrap);
  const ready = useAuth((s) => s.ready);
  const [booted, setBooted] = useState(false);

  useEffect(() => {
    bootstrap().finally(() => setBooted(true));
  }, [bootstrap]);

  if (!booted && !ready) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Spinner label="Loading…" />
      </div>
    );
  }

  return (
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
      <ToastViewport />
    </QueryClientProvider>
  );
}
