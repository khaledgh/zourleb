import { Stack } from "expo-router";

// Agency portal — a modal stack pushed from the Profile tab. No bottom tabs;
// tourists see the tourist experience, agency users can drill into this stack.
export default function AgencyLayout() {
  return (
    <Stack
      screenOptions={{
        headerTintColor: "#1f8a5b",
        headerStyle: { backgroundColor: "#fff" },
        headerShadowVisible: false,
      }}
    >
      <Stack.Screen
        name="index"
        options={{ title: "Agency Portal", headerLargeTitle: true }}
      />
      <Stack.Screen name="tours" options={{ title: "My Tours" }} />
      <Stack.Screen name="bookings" options={{ title: "Bookings" }} />
    </Stack>
  );
}
