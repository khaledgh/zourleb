import { Redirect } from "expo-router";

// Entry redirect. Guests land on Home (browse is public); the booking flow and
// profile actions prompt sign-in when needed.
export default function Index() {
  return <Redirect href="/(tabs)" />;
}
