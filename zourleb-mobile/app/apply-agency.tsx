import { useState } from "react";
import { KeyboardAvoidingView, Platform, ScrollView, Text, TextInput, View, Pressable } from "react-native";
import { useRouter, Stack } from "expo-router";
import { useMutation, useQuery } from "@tanstack/react-query";
import { apiPost, apiData, ApiException } from "@/api/client";
import { Button, Field, Card } from "@/components/ui";
import type { Region } from "@/types/api";

export default function ApplyAgency() {
  const router = useRouter();

  const [name, setName] = useState("");
  const [phone, setPhone] = useState("");
  const [email, setEmail] = useState("");
  const [regionId, setRegionId] = useState<number | null>(null);

  const [nameError, setNameError] = useState<string | null>(null);
  const [phoneError, setPhoneError] = useState<string | null>(null);
  const [emailError, setEmailError] = useState<string | null>(null);
  const [serverError, setServerError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  // Fetch regions
  const { data: regions } = useQuery({
    queryKey: ["regions"],
    queryFn: () => apiData<Region[]>("/regions"),
  });

  function validate(): boolean {
    let valid = true;

    if (name.trim().length < 2) {
      setNameError("Agency name must be at least 2 characters.");
      valid = false;
    } else {
      setNameError(null);
    }

    if (!email.trim() || !email.includes("@")) {
      setEmailError("Please enter a valid email address.");
      valid = false;
    } else {
      setEmailError(null);
    }

    if (!phone.trim()) {
      setPhoneError("Phone number is required.");
      valid = false;
    } else {
      setPhoneError(null);
    }

    return valid;
  }

  const applyMutation = useMutation({
    mutationFn: () =>
      apiPost("/agency/apply", {
        name: name.trim(),
        phone: phone.trim(),
        email: email.trim().toLowerCase(),
        region_id: regionId ?? undefined,
      }),
    onSuccess: () => {
      setSuccess(true);
    },
    onError: (e: any) => {
      setServerError(e instanceof ApiException ? e.message : "Application failed. Try again.");
    },
  });

  function onSubmit() {
    if (!validate()) return;
    setServerError(null);
    applyMutation.mutate();
  }

  if (success) {
    return (
      <View className="flex-1 bg-white items-center justify-center p-6">
        <Stack.Screen options={{ title: "Application Success" }} />
        <Text className="text-6xl mb-4">🎉</Text>
        <Text className="text-2xl font-bold text-center text-gray-900 mb-2">Application Submitted!</Text>
        <Text className="text-sm text-center text-gray-500 mb-8">
          Thank you for applying. Our admin team will review your application and get back to you shortly.
        </Text>
        <View className="w-full">
          <Button title="Back to Profile" onPress={() => router.back()} />
        </View>
      </View>
    );
  }

  return (
    <KeyboardAvoidingView
      className="flex-1"
      behavior={Platform.OS === "ios" ? "padding" : "height"}
    >
      <Stack.Screen options={{ title: "Apply as Agency" }} />
      <ScrollView className="flex-1 bg-gray-50" contentContainerClassName="p-4 pb-12">
        <View className="mb-5">
          <Text className="text-xl font-bold text-gray-900 mb-1">Register your Agency</Text>
          <Text className="text-sm text-gray-500">
            Submit your application to become an agency on Zourleb and start listing your tours.
          </Text>
        </View>

        <Card className="p-4 mb-5">
          <Field label="Agency Name" error={nameError ?? undefined}>
            <TextInput
              className="h-12 rounded-xl border border-gray-300 px-3 text-base bg-white"
              autoCapitalize="words"
              placeholder="e.g. Lebanon Travel Co."
              value={name}
              onChangeText={setName}
            />
          </Field>

          <Field label="Agency Email" error={emailError ?? undefined}>
            <TextInput
              className="h-12 rounded-xl border border-gray-300 px-3 text-base bg-white"
              autoCapitalize="none"
              keyboardType="email-address"
              placeholder="e.g. contact@agency.com"
              value={email}
              onChangeText={setEmail}
            />
          </Field>

          <Field label="Agency Phone (Lebanese)" error={phoneError ?? undefined}>
            <TextInput
              className="h-12 rounded-xl border border-gray-300 px-3 text-base bg-white"
              keyboardType="phone-pad"
              placeholder="e.g. +961 3 123 456"
              value={phone}
              onChangeText={setPhone}
            />
          </Field>

          <Text className="text-sm font-medium text-gray-700 mb-2">Base Region (Lebanon)</Text>
          <View className="flex-row flex-wrap gap-2 mb-4">
            {(regions ?? []).map((r: any) => {
              const selected = regionId === r.id;
              return (
                <Pressable
                  key={r.id}
                  onPress={() => setRegionId(r.id)}
                  className={`rounded-full px-4 py-2 border ${
                    selected ? "bg-brand-500 border-brand-500" : "bg-white border-gray-200"
                  }`}
                >
                  <Text className={`text-sm ${selected ? "text-white font-semibold" : "text-gray-700"}`}>
                    {r.name}
                  </Text>
                </Pressable>
              );
            })}
          </View>
        </Card>

        {serverError ? (
          <Text className="mb-4 text-center text-sm text-red-600">{serverError}</Text>
        ) : null}

        <Button title="Submit Application" onPress={onSubmit} loading={applyMutation.isPending} />
        <View className="mt-3">
          <Button title="Cancel" variant="ghost" onPress={() => router.back()} />
        </View>
      </ScrollView>
    </KeyboardAvoidingView>
  );
}
