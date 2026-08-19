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
      <Stack.Screen name="members" options={{ title: "Members" }} />
      <Stack.Screen name="new-tour" options={{ title: "Create Tour" }} />
      <Stack.Screen name="boosts" options={{ title: "Boosts" }} />
      <Stack.Screen name="products" options={{ title: "Products" }} />
      <Stack.Screen name="new-product" options={{ title: "Add Product" }} />
      <Stack.Screen name="edit-agency" options={{ title: "Edit Agency" }} />
      <Stack.Screen name="tour/[id]" options={{ title: "Tour" }} />
    </Stack>
  );
}
