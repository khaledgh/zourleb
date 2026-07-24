import "../global.css";
import { useEffect, useState } from "react";
import { Stack } from "expo-router";
import { StatusBar } from "expo-status-bar";
import { GestureHandlerRootView } from "react-native-gesture-handler";
import { SafeAreaProvider } from "react-native-safe-area-context";
import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "@/lib/query";
import { initI18n } from "@/i18n";
import { useAuth } from "@/stores/auth";
import { Loading } from "@/components/ui";

// Root layout: initializes i18n + restores the session, then renders the stack.
// Auth gating is handled per-group (see app/(tabs)/_layout.tsx and (auth)).
export default function RootLayout() {
  const bootstrap = useAuth((s) => s.bootstrap);
  const [ready, setReady] = useState(false);

  useEffect(() => {
    (async () => {
      await initI18n();
      await bootstrap();
      setReady(true);
    })();
  }, [bootstrap]);

  if (!ready) {
    return (
      <SafeAreaProvider>
        <Loading label="…" />
      </SafeAreaProvider>
    );
  }

  return (
    <GestureHandlerRootView style={{ flex: 1 }}>
      <SafeAreaProvider>
        <QueryClientProvider client={queryClient}>
          <StatusBar style="dark" />
          <Stack screenOptions={{ headerShown: false }}>
            <Stack.Screen name="(tabs)" />
            <Stack.Screen name="(auth)" />
            <Stack.Screen
              name="tour/[slug]"
              options={{ headerShown: true, title: "" }}
            />
            <Stack.Screen
              name="booking/[slug]"
              options={{ headerShown: true, title: "Booking", presentation: "modal" }}
            />
          </Stack>
        </QueryClientProvider>
      </SafeAreaProvider>
    </GestureHandlerRootView>
  );
}
