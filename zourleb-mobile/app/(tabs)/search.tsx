import { useState } from "react";

import { FlatList, ScrollView, Text, TextInput, View, Pressable } from "react-native";

import { SafeAreaView } from "react-native-safe-area-context";

import { useQuery } from "@tanstack/react-query";

import { apiList, apiData } from "@/api/client";

import { TourCard } from "@/components/TourCard";

import { Loading, EmptyState } from "@/components/ui";

import type { Category, TourCard as TourCardType } from "@/types/api";

export default function Search() {
  const [q, setQ] = useState("");
  const [categoryId, setCategoryId] = useState<number | null>(null);

  const { data: categories } = useQuery({

    queryKey: ["categories"],

    queryFn: () => apiData<Category[]>("/categories"),

  });



  const { data, isLoading } = useQuery({

    queryKey: ["search-tours", q, categoryId],

    queryFn: () =>

      apiList<TourCardType>("/tours", {

        q: q || undefined,

        category_id: categoryId ?? undefined,

        per_page: 20,

      }),

  });



  return (

    <SafeAreaView className="flex-1 bg-gray-50" edges={["top"]}>

      <View className="px-4 pt-3">

        <TextInput

          className="mb-3 h-12 rounded-xl border border-gray-200 bg-white px-4 text-base"

          placeholder="Search tours…"

          value={q}

          onChangeText={setQ}

        />

        <ScrollView horizontal showsHorizontalScrollIndicator={false} className="mb-2">

          <FilterChip label="All" active={categoryId === null} onPress={() => setCategoryId(null)} />

          {(categories ?? []).map((c) => (

            <FilterChip

              key={c.id}

              label={c.name}

              active={categoryId === c.id}

              onPress={() => setCategoryId(c.id)}

            />

          ))}

        </ScrollView>

      </View>



      {isLoading ? (

        <Loading />

      ) : (

        <FlatList

          data={data?.items ?? []}

          keyExtractor={(t) => String(t.id)}

          renderItem={({ item }) => <TourCard tour={item} />}

          contentContainerClassName="px-4 pb-8 pt-2"

          ListEmptyComponent={<EmptyState message="No tours match your search." />}

        />

      )}

    </SafeAreaView>

  );

}



function FilterChip({

  label,

  active,

  onPress,

}: {

  label: string;

  active: boolean;

  onPress: () => void;

}) {

  return (

    <Pressable

      onPress={onPress}

      className={`me-2 rounded-full px-4 py-2 ${active ? "bg-brand-500" : "bg-white border border-gray-200"}`}

    >

      <Text className={`text-sm ${active ? "text-white" : "text-gray-700"}`}>{label}</Text>

    </Pressable>

  );

}

