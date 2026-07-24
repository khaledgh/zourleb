import { useState } from "react";
import { Text, TextInput, View } from "react-native";
import { Link, useRouter } from "expo-router";
import { SafeAreaView } from "react-native-safe-area-context";
import { useForm, Controller } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { useAuth } from "@/stores/auth";
import { ApiException } from "@/api/client";
import { Button, Field } from "@/components/ui";

interface Form {
  email: string;
  password: string;
}

export default function Login() {
  const { t } = useTranslation();
  const router = useRouter();
  const login = useAuth((s) => s.login);
  const [submitting, setSubmitting] = useState(false);
  const [serverError, setServerError] = useState<string | null>(null);
  const { control, handleSubmit, formState: { errors } } = useForm<Form>();

  async function onSubmit(v: Form) {
    setSubmitting(true);
    setServerError(null);
    try {
      await login(v.email, v.password);
      router.replace("/(tabs)");
    } catch (e) {
      setServerError(e instanceof ApiException ? e.message : "Sign-in failed");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <SafeAreaView className="flex-1 bg-white">
      <View className="flex-1 justify-center px-6">
        <Text className="mb-1 text-center text-3xl font-bold text-brand-600">Zourleb</Text>
        <Text className="mb-8 text-center text-sm text-gray-500">
          {t("action.signin")}
        </Text>

        <Field label="Email" error={errors.email?.message}>
          <Controller
            control={control}
            name="email"
            rules={{ required: "Email is required" }}
            render={({ field: { onChange, value } }) => (
              <TextInput
                className="h-12 rounded-xl border border-gray-300 px-3 text-base"
                autoCapitalize="none"
                keyboardType="email-address"
                value={value}
                onChangeText={onChange}
              />
            )}
          />
        </Field>

        <Field label="Password" error={errors.password?.message}>
          <Controller
            control={control}
            name="password"
            rules={{ required: "Password is required" }}
            render={({ field: { onChange, value } }) => (
              <TextInput
                className="h-12 rounded-xl border border-gray-300 px-3 text-base"
                secureTextEntry
                value={value}
                onChangeText={onChange}
              />
            )}
          />
        </Field>

        {serverError ? (
          <Text className="mb-3 text-center text-sm text-red-600">{serverError}</Text>
        ) : null}

        <Button title={t("action.signin")} onPress={handleSubmit(onSubmit)} loading={submitting} />

        <View className="mt-6 flex-row justify-center gap-1">
          <Text className="text-sm text-gray-500">No account?</Text>
          <Link href="/(auth)/register" className="text-sm font-semibold text-brand-600">
            Register
          </Link>
        </View>
      </View>
    </SafeAreaView>
  );
}
