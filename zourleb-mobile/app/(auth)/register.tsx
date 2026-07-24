import { useState } from "react";
import { Text, TextInput, View } from "react-native";
import { Link, useRouter } from "expo-router";
import { SafeAreaView } from "react-native-safe-area-context";
import { useForm, Controller } from "react-hook-form";
import { useAuth } from "@/stores/auth";
import { ApiException } from "@/api/client";
import { Button, Field } from "@/components/ui";

interface Form {
  name: string;
  email: string;
  password: string;
}

export default function Register() {
  const router = useRouter();
  const register = useAuth((s) => s.register);
  const [submitting, setSubmitting] = useState(false);
  const [serverError, setServerError] = useState<string | null>(null);
  const { control, handleSubmit, formState: { errors } } = useForm<Form>();

  async function onSubmit(v: Form) {
    setSubmitting(true);
    setServerError(null);
    try {
      await register(v.name, v.email, v.password);
      router.replace("/(tabs)");
    } catch (e) {
      setServerError(e instanceof ApiException ? e.message : "Registration failed");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <SafeAreaView className="flex-1 bg-white">
      <View className="flex-1 justify-center px-6">
        <Text className="mb-8 text-center text-2xl font-bold text-gray-900">
          Create your account
        </Text>

        <Field label="Full name" error={errors.name?.message}>
          <Controller
            control={control}
            name="name"
            rules={{ required: "Name is required", minLength: { value: 2, message: "Too short" } }}
            render={({ field: { onChange, value } }) => (
              <TextInput
                className="h-12 rounded-xl border border-gray-300 px-3 text-base"
                value={value}
                onChangeText={onChange}
              />
            )}
          />
        </Field>

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
            rules={{ required: "Password is required", minLength: { value: 8, message: "Min 8 characters" } }}
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

        <Button title="Create account" onPress={handleSubmit(onSubmit)} loading={submitting} />

        <View className="mt-6 flex-row justify-center gap-1">
          <Text className="text-sm text-gray-500">Have an account?</Text>
          <Link href="/(auth)/login" className="text-sm font-semibold text-brand-600">
            Sign in
          </Link>
        </View>
      </View>
    </SafeAreaView>
  );
}
