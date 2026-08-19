import { useState } from "react";
import { Alert, ScrollView, Text, TextInput, View } from "react-native";
import { useRouter, Stack } from "expo-router";
import { useMutation } from "@tanstack/react-query";
import { apiPost, ApiException } from "@/api/client";
import { Button, Card } from "@/components/ui";

export default function NewProduct() {
  const router = useRouter();
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [price, setPrice] = useState("");
  const [currency, setCurrency] = useState("USD");
  const [stock, setStock] = useState("");

  const create = useMutation({
    mutationFn: () =>
      apiPost("/agency/products", {
        price: Number(price) || 0,
        currency,
        stock: Number(stock) || 0,
        translations: [
          { locale: "en", name: name.trim(), description: description.trim() },
        ],
      }),
    onSuccess: () => {
      Alert.alert("Product added");
      router.back();
    },
    onError: (e) => {
      Alert.alert("Failed", e instanceof ApiException ? e.message : "Could not add product.");
    },
  });

  return (
    <ScrollView className="flex-1 bg-gray-50" contentContainerClassName="p-4 pb-12">
      <Stack.Screen options={{ title: "Add Product" }} />
      <Text className="mb-4 text-xl font-bold text-gray-900">Add a product</Text>

      <Card className="p-4 mb-4">
        <Field label="Name (English)" value={name} onChange={setName} />
        <Field label="Description" value={description} onChange={setDescription} multiline />
        <View className="mb-3 flex-row gap-3">
          <View className="flex-1">
            <Text className="mb-1 text-sm font-medium text-gray-700">Price</Text>
            <TextInput
              className="h-12 rounded-xl border border-gray-300 bg-white px-3 text-base"
              keyboardType="numeric"
              value={price}
              onChangeText={setPrice}
            />
          </View>
          <View className="w-28">
            <Text className="mb-1 text-sm font-medium text-gray-700">Currency</Text>
            <TextInput
              className="h-12 rounded-xl border border-gray-300 bg-white px-3 text-base"
              value={currency}
              onChangeText={setCurrency}
              maxLength={3}
              autoCapitalize="characters"
            />
          </View>
        </View>
        <View className="mb-3">
          <Text className="mb-1 text-sm font-medium text-gray-700">Stock</Text>
          <TextInput
            className="h-12 rounded-xl border border-gray-300 bg-white px-3 text-base"
            keyboardType="numeric"
            value={stock}
            onChangeText={setStock}
          />
        </View>
      </Card>

      <Button
        title="Add product"
        onPress={() => create.mutate()}
        loading={create.isPending}
        disabled={!name.trim()}
      />
    </ScrollView>
  );
}

function Field({
  label,
  value,
  onChange,
  multiline,
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
  multiline?: boolean;
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
      />
    </View>
  );
}
