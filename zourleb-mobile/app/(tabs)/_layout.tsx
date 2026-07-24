import { Tabs } from "expo-router";
import { Text } from "react-native";
import { useTranslation } from "react-i18next";

// Bottom tabs for the tourist experience. Browse is public; Bookings/Favorites
// prompt sign-in inside the screen when the user is a guest.
function TabIcon({ icon, color }: { icon: string; color: string }) {
  return <Text style={{ fontSize: 20, color }}>{icon}</Text>;
}

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
          tabBarIcon: ({ color }) => <TabIcon icon="🏠" color={color} />,
        }}
      />
      <Tabs.Screen
        name="search"
        options={{
          title: t("tab.search"),
          tabBarIcon: ({ color }) => <TabIcon icon="🔎" color={color} />,
        }}
      />
      <Tabs.Screen
        name="bookings"
        options={{
          title: t("tab.bookings"),
          tabBarIcon: ({ color }) => <TabIcon icon="🎫" color={color} />,
        }}
      />
      <Tabs.Screen
        name="favorites"
        options={{
          title: t("tab.favorites"),
          tabBarIcon: ({ color }) => <TabIcon icon="♥" color={color} />,
        }}
      />
      <Tabs.Screen
        name="profile"
        options={{
          title: t("tab.profile"),
          tabBarIcon: ({ color }) => <TabIcon icon="👤" color={color} />,
        }}
      />
    </Tabs>
  );
}
