import { useEffect, useState } from "react";
import { Alert, ScrollView, Text, TextInput, View } from "react-native";
import { useRouter, Stack } from "expo-router";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiData, apiPatch, ApiException } from "@/api/client";
import { Button, Card, Loading } from "@/components/ui";
import type { AgencyProfile } from "@/types/api";

export default function EditAgency() {
  const router = useRouter();
  const queryClient = useQueryClient();

  const { data: profile, isLoading } = useQuery({
    queryKey: ["agency-profile"],
    queryFn: () => apiData<AgencyProfile>("/agency/profile"),
  });

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [website, setWebsite] = useState("");
  const [description, setDescription] = useState("");

  useEffect(() => {
    if (profile) {
      setName(profile.name);
      setEmail(profile.email);
      setPhone(profile.phone);
      setWebsite(profile.website);
      setDescription(profile.about ?? "");
    }
  }, [profile]);

  const save = useMutation({
    mutationFn: () =>
      apiPatch("/agency/profile", {
        name: name.trim(),
        email: email.trim(),
        phone: phone.trim(),
        website: website.trim(),
        about: description.trim(),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["agency-profile"] });
      Alert.alert("Saved");
      router.back();
    },
    onError: (e) => {
      Alert.alert("Failed", e instanceof ApiException ? e.message : "Could not update profile.");
    },
  });

  if (isLoading) return <Loading />;

  return (
    <ScrollView className="flex-1 bg-gray-50" contentContainerClassName="p-4 pb-12">
      <Stack.Screen options={{ title: "Edit Agency" }} />
      <Text className="mb-4 text-xl font-bold text-gray-900">Edit agency profile</Text>

      <Card className="p-4 mb-4">
        <Field label="Name" value={name} onChange={setName} />
        <Field label="Email" value={email} onChange={setEmail} keyboard="email-address" />
        <Field label="Phone" value={phone} onChange={setPhone} keyboard="phone-pad" />
        <Field label="Website" value={website} onChange={setWebsite} />
        <Field label="About" value={description} onChange={setDescription} multiline />
      </Card>

      <Button title="Save" onPress={() => save.mutate()} loading={save.isPending} disabled={!name.trim()} />
    </ScrollView>
  );
}

function Field({
  label,
  value,
  onChange,
  multiline,
  keyboard,
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
  multiline?: boolean;
  keyboard?: "email-address" | "phone-pad";
}) {
  return (
    <View className="mb-3">
      <Text className="mb-1 text-sm font-medium text-gray-700">{label}</Text>
      <TextInput
        className={`rounded-xl border border-gray-300 bg-white px-3 text-base ${
          multiline ? "h-24 py-3" : "h-12"
        }`}
        multiline={multiline}
        textAlignVertical={multiline ? "top" : "center"}
        value={value}
        onChangeText={onChange}
        keyboardType={keyboard}
      />
    </View>
  );
}
