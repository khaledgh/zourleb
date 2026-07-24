import { Text, View } from "react-native";
import { useRouter } from "expo-router";
import { Button } from "./ui";

// Shown on guarded tourist screens when the user is a guest.
export function SignInPrompt({ message }: { message: string }) {
  const router = useRouter();
  return (
    <View className="flex-1 items-center justify-center gap-4 px-8">
      <Text className="text-center text-base text-gray-600">{message}</Text>
      <View className="w-full">
        <Button title="Sign in" onPress={() => router.push("/(auth)/login")} />
      </View>
    </View>
  );
}
