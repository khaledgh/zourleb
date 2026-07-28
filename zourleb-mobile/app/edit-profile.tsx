import { useState } from "react";
import { KeyboardAvoidingView, Platform, ScrollView, Text, TextInput, View } from "react-native";
import { useRouter, Stack } from "expo-router";
import { useMutation } from "@tanstack/react-query";
import { apiPatch, ApiException } from "@/api/client";
import { useAuth } from "@/stores/auth";
import { Button, Field } from "@/components/ui";
import type { User } from "@/types/api";

export default function EditProfile() {
  const router = useRouter();
  const user = useAuth((s) => s.user);
  const refreshMe = useAuth((s) => s.refreshMe);

  const [name, setName] = useState(user?.name ?? "");
  const [phone, setPhone] = useState(user?.phone ?? "");
  const [nameError, setNameError] = useState<string | null>(null);
  const [serverError, setServerError] = useState<string | null>(null);

  function validate(): boolean {
    if (name.trim().length < 2) {
      setNameError("Name must be at least 2 characters.");
      return false;
    }
    setNameError(null);
    return true;
  }

  const update = useMutation({
    mutationFn: () =>
      apiPatch<User>("/me", { name: name.trim(), phone: phone.trim() || undefined }),
    onSuccess: async () => {
      await refreshMe();
      router.back();
    },
    onError: (e) => {
      setServerError(e instanceof ApiException ? e.message : "Update failed. Try again.");
    },
  });

  function onSave() {
    if (!validate()) return;
    setServerError(null);
    update.mutate();
  }

  return (
    <KeyboardAvoidingView
      className="flex-1"
      behavior={Platform.OS === "ios" ? "padding" : "height"}
    >
      <Stack.Screen options={{ title: "Edit Profile" }} />
      <ScrollView className="flex-1 bg-gray-50" contentContainerClassName="p-4 pb-12">
        <View className="mb-6 rounded-2xl border border-gray-100 bg-white p-4">
          <Field label="Display name" error={nameError ?? undefined}>
            <TextInput
              className="h-12 rounded-xl border border-gray-300 px-3 text-base"
              autoCapitalize="words"
              value={name}
              onChangeText={setName}
            />
          </Field>

          <Field label="Phone number (Lebanese)">
            <TextInput
              className="h-12 rounded-xl border border-gray-300 px-3 text-base"
              keyboardType="phone-pad"
              placeholder="03 123 456"
              value={phone}
              onChangeText={setPhone}
            />
          </Field>
        </View>

        {serverError ? (
          <Text className="mb-4 text-center text-sm text-red-600">{serverError}</Text>
        ) : null}

        <Button title="Save changes" onPress={onSave} loading={update.isPending} />
        <View className="mt-3">
          <Button title="Cancel" variant="ghost" onPress={() => router.back()} />
        </View>
      </ScrollView>
    </KeyboardAvoidingView>
  );
}
