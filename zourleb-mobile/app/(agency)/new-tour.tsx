import { useState } from "react";
import { Alert, ScrollView, Text, TextInput, View } from "react-native";
import { useRouter, Stack } from "expo-router";
import { useQuery, useMutation } from "@tanstack/react-query";
import { apiData, apiPost, ApiException } from "@/api/client";
import { Button, Card, Loading } from "@/components/ui";
import type { Category, Region } from "@/types/api";

export default function NewTour() {
  const router = useRouter();
  const [title, setTitle] = useState("");
  const [summary, setSummary] = useState("");
  const [description, setDescription] = useState("");
  const [duration, setDuration] = useState("1");
  const [difficulty, setDifficulty] = useState("");
  const [minAge, setMinAge] = useState("");
  const [capacity, setCapacity] = useState("");
  const [price, setPrice] = useState("");
  const [categoryId, setCategoryId] = useState<number | null>(null);
  const [regionId, setRegionId] = useState<number | null>(null);

  const { data: categories, isLoading: catsLoading } = useQuery({
    queryKey: ["categories"],
    queryFn: () => apiData<Category[]>("/categories"),
  });

  const { data: regions, isLoading: regionsLoading } = useQuery({
    queryKey: ["regions"],
    queryFn: () => apiData<Region[]>("/regions"),
  });

  const create = useMutation({
    mutationFn: () =>
      apiPost("/agency/tours", {
        title: title.trim(),
        summary: summary.trim(),
        description: description.trim() || undefined,
        type: "guided",
        duration_days: Number(duration) || 1,
        difficulty: difficulty.trim() || undefined,
        min_age: minAge ? Number(minAge) : undefined,
        max_capacity: capacity ? Number(capacity) : undefined,
        price_from: price ? Number(price) : 0,
        currency: "USD",
        category_id: categoryId ?? undefined,
        region_id: regionId ?? undefined,
      }),
    onSuccess: () => {
      Alert.alert("Tour created", "It will appear in your tours list.");
      router.back();
    },
    onError: (e) => {
      Alert.alert("Failed", e instanceof ApiException ? e.message : "Could not create tour.");
    },
  });

  if (catsLoading || regionsLoading) return <Loading />;

  return (
    <ScrollView className="flex-1 bg-gray-50" contentContainerClassName="p-4 pb-12">
      <Stack.Screen options={{ title: "Create Tour" }} />
      <Text className="mb-4 text-xl font-bold text-gray-900">Create a new tour</Text>

      <Card className="p-4 mb-4">
        <Field label="Title" value={title} onChange={setTitle} />
        <Field label="Summary" value={summary} onChange={setSummary} />
        <Field label="Description" value={description} onChange={setDescription} multiline />
        <View className="mb-3 flex-row gap-3">
          <View className="flex-1">
            <Text className="mb-1 text-sm font-medium text-gray-700">Duration (days)</Text>
            <TextInput
              className="h-12 rounded-xl border border-gray-300 bg-white px-3 text-base"
              keyboardType="numeric"
              value={duration}
              onChangeText={setDuration}
            />
          </View>
          <View className="flex-1">
            <Text className="mb-1 text-sm font-medium text-gray-700">Price from</Text>
            <TextInput
              className="h-12 rounded-xl border border-gray-300 bg-white px-3 text-base"
              keyboardType="numeric"
              value={price}
              onChangeText={setPrice}
              placeholder="0"
            />
          </View>
        </View>
        <Field label="Difficulty" value={difficulty} onChange={setDifficulty} placeholder="e.g. easy" />
        <View className="mb-3 flex-row gap-3">
          <View className="flex-1">
            <Text className="mb-1 text-sm font-medium text-gray-700">Min age</Text>
            <TextInput
              className="h-12 rounded-xl border border-gray-300 bg-white px-3 text-base"
              keyboardType="numeric"
              value={minAge}
              onChangeText={setMinAge}
            />
          </View>
          <View className="flex-1">
            <Text className="mb-1 text-sm font-medium text-gray-700">Max capacity</Text>
            <TextInput
              className="h-12 rounded-xl border border-gray-300 bg-white px-3 text-base"
              keyboardType="numeric"
              value={capacity}
              onChangeText={setCapacity}
            />
          </View>
        </View>
      </Card>

      <Button
        title="Create tour"
        onPress={() => create.mutate()}
        loading={create.isPending}
        disabled={!title.trim()}
      />
    </ScrollView>
  );
}

function Field({
  label,
  value,
  onChange,
  multiline,
  placeholder,
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
  multiline?: boolean;
  placeholder?: string;
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
        placeholder={placeholder}
      />
    </View>
  );
}
