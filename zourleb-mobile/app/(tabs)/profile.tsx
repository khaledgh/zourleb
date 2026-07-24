import { useEffect, useState } from "react";
import { Alert, Pressable, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useRouter } from "expo-router";
import { useTranslation } from "react-i18next";
import { apiData } from "@/api/client";
import { useAuth } from "@/stores/auth";
import { changeLocale } from "@/i18n";
import { SignInPrompt } from "@/components/SignInPrompt";
import { Button, Card } from "@/components/ui";
import type { Language } from "@/types/api";

export default function Profile() {
  const { t, i18n } = useTranslation();
  const router = useRouter();
  const user = useAuth((s) => s.user);
  const logout = useAuth((s) => s.logout);
  const isAgency = useAuth((s) => s.isAgency);
  const [langs, setLangs] = useState<Language[]>([]);

  useEffect(() => {
    apiData<Language[]>("/languages").then(setLangs).catch(() => setLangs([]));
  }, []);

  async function onLogout() {
    await logout();
    router.replace("/(tabs)");
  }

  function onPickLanguage(code: string) {
    changeLocale(code).then(() => {
      Alert.alert("Language updated", "Restart the app for full RTL layout changes.");
    });
  }

  if (!user) {
    return (
      <SafeAreaView className="flex-1 bg-gray-50">
        <SignInPrompt message="Sign in to manage your profile, bookings and favorites." />
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView className="flex-1 bg-gray-50" edges={["top"]}>
      <View className="px-4 pt-4">
        <Text className="text-2xl font-bold text-gray-900">{t("tab.profile")}</Text>

        <Card className="mt-4 p-4">
          <Text className="text-lg font-semibold text-gray-900">{user.name}</Text>
          <Text className="text-sm text-gray-500">{user.email}</Text>
          <View className="mt-2 flex-row gap-2">
            {user.phone_verified ? (
              <Text className="text-xs text-brand-600">✓ Phone verified</Text>
            ) : (
              <Text className="text-xs text-amber-600">Phone not verified</Text>
            )}
          </View>
        </Card>

        {isAgency() ? (
          <Card className="mt-4 p-4">
            <Text className="text-base font-semibold text-gray-900">Agency mode</Text>
            <Text className="mb-3 text-sm text-gray-500">
              You manage an agency. Tour & booking management is available here.
            </Text>
            <Button
              title="Open agency portal"
              variant="ghost"
              onPress={() => Alert.alert("Agency portal", "Agency screens coming soon.")}
            />
          </Card>
        ) : null}

        <Card className="mt-4 p-4">
          <Text className="mb-2 text-base font-semibold text-gray-900">Language</Text>
          <View className="flex-row flex-wrap gap-2">
            {langs.map((l) => (
              <Pressable
                key={l.code}
                onPress={() => onPickLanguage(l.code)}
                className={`rounded-full px-4 py-2 ${
                  i18n.language === l.code ? "bg-brand-500" : "bg-gray-100"
                }`}
              >
                <Text className={i18n.language === l.code ? "text-white" : "text-gray-700"}>
                  {l.native_name}
                </Text>
              </Pressable>
            ))}
          </View>
        </Card>

        <View className="mt-6">
          <Button title={t("action.signout")} variant="danger" onPress={onLogout} />
        </View>
      </View>
    </SafeAreaView>
  );
}
