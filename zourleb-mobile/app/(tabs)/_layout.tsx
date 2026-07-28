import { Tabs } from "expo-router";
import { Text, View } from "react-native";
import { useTranslation } from "react-i18next";
import { useQuery } from "@tanstack/react-query";
import { apiList } from "@/api/client";
import { useAuth } from "@/stores/auth";
import type { Notification } from "@/types/api";

function TabIcon({ icon, color }: { icon: string; color: string }) {
  return <Text style={{ fontSize: 20, color }}>{icon}</Text>;
}

function BellIcon({ color }: { color: string }) {
  const user = useAuth((s) => s.user);
  const { data } = useQuery({
    queryKey: ["notifications"],
    queryFn: () => apiList<Notification>("/notifications", { per_page: 50 }),
    enabled: !!user,
    refetchInterval: 60_000,
  });
  const unread = (data?.items ?? []).filter((n) => !n.read).length;

  return (
    <View className="relative">
      <Text style={{ fontSize: 20, color }}>🔔</Text>
      {unread > 0 ? (
        <View className="absolute -right-1.5 -top-1 h-4 w-4 items-center justify-center rounded-full bg-red-500">
          <Text className="text-[9px] font-bold text-white">{unread > 9 ? "9+" : unread}</Text>
        </View>
      ) : null}
    </View>
  );
}

// Bottom tabs for the tourist experience. Browse is public; Bookings/Favorites/
// Notifications prompt sign-in inside the screen when the user is a guest.
export default function TabsLayout() {
  const { t } = useTranslation();
  return (
    <Tabs
      screenOptions={{
        headerShown: false,
        tabBarActiveTintColor: "#1f8a5b",
        tabBarInactiveTintColor: "#9ca3af",
        tabBarStyle: { height: 58, paddingBottom: 6, paddingTop: 6 },
      }}
    >
      <Tabs.Screen
        name="index"
        options={{
          title: t("tab.home"),
          tabBarIcon: ({ color }: { color: string }) => <TabIcon icon="🏠" color={color} />,
        }}
      />
      <Tabs.Screen
        name="search"
        options={{
          title: t("tab.search"),
          tabBarIcon: ({ color }: { color: string }) => <TabIcon icon="🔎" color={color} />,
        }}
      />
      <Tabs.Screen
        name="bookings"
        options={{
          title: t("tab.bookings"),
          tabBarIcon: ({ color }: { color: string }) => <TabIcon icon="🎫" color={color} />,
        }}
      />
      <Tabs.Screen
        name="favorites"
        options={{
          title: t("tab.favorites"),
          tabBarIcon: ({ color }: { color: string }) => <TabIcon icon="♥" color={color} />,
        }}
      />
      <Tabs.Screen
        name="notifications"
        options={{
          title: "Alerts",
          tabBarIcon: ({ color }: { color: string }) => <BellIcon color={color} />,
        }}
      />
      <Tabs.Screen
        name="profile"
        options={{
          title: t("tab.profile"),
          tabBarIcon: ({ color }: { color: string }) => <TabIcon icon="👤" color={color} />,
        }}
      />
    </Tabs>
  );
}
