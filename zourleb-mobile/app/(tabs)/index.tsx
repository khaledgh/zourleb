import { FlatList, Image, Pressable, ScrollView, Text, TextInput, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useQuery } from "@tanstack/react-query";
import { useRouter } from "expo-router";
import { useState } from "react";
import { apiData } from "@/api/client";
import { TourCard } from "@/components/TourCard";
import { Loading, SectionHeader } from "@/components/ui";
import { useFavoriteIds } from "@/lib/favorites";
import type { HomePayload } from "@/types/api";

export default function Home() {
  const router = useRouter();
  const [selectedCategory, setSelectedCategory] = useState<number | null>(null);
  const { ids: favIds } = useFavoriteIds();

  const { data, isLoading, refetch, isRefetching } = useQuery({
    queryKey: ["home"],
    queryFn: () => apiData<HomePayload>("/home"),
  });

  if (isLoading) return <Loading />;

  // Navigate to search with category pre-selected
  function onCategoryPress(categoryId: number | null) {
    setSelectedCategory(categoryId);
    if (categoryId !== null) {
      router.push({ pathname: "/(tabs)/search", params: { category_id: String(categoryId) } });
    }
  }

  const featured = data?.featured ?? [];
  const lastMinute = data?.last_minute ?? [];
  const allTours = [...featured, ...lastMinute];

  return (
    <SafeAreaView className="flex-1 bg-gray-50" edges={["top"]}>
      <FlatList
        data={allTours}
        keyExtractor={(t) => `${t.id}-${t.slug}`}
        renderItem={({ item, index }) => {
          const isFeaturedSection = index < featured.length;
          // Show "Last Minute Deals" header before the first last-minute item
          if (!isFeaturedSection && index === featured.length) {
            return (
              <>
                <View className="px-4 pt-2">
                  <SectionHeader title="⚡ Last Minute Deals" />
                </View>
                <View className="px-4">
                  <TourCard tour={item} isFavorited={favIds.has(item.id)} />
                </View>
              </>
            );
          }
          return (
            <View className="px-4">
              <TourCard tour={item} isFavorited={favIds.has(item.id)} />
            </View>
          );
        }}
        contentContainerClassName="pb-8"
        onRefresh={refetch}
        refreshing={isRefetching}
        ListHeaderComponent={
          <View>
            {/* Header bar with title + search shortcut */}
            <View className="flex-row items-center justify-between px-4 pb-3 pt-4">
              <Text className="text-2xl font-bold text-gray-900">Discover Lebanon</Text>
            </View>

            {/* Search bar shortcut */}
            <Pressable
              onPress={() => router.push("/(tabs)/search")}
              className="mx-4 mb-4 h-12 flex-row items-center rounded-xl border border-gray-200 bg-white px-4"
            >
              <Text className="mr-2 text-gray-400">🔎</Text>
              <Text className="text-base text-gray-400">Search tours, regions…</Text>
            </Pressable>

            {/* Banners carousel */}
            {(data?.banners?.length ?? 0) > 0 ? (
              <ScrollView
                horizontal
                showsHorizontalScrollIndicator={false}
                className="mb-5 ps-4"
              >
                {data!.banners.map((b) => (
                  <View
                    key={b.id}
                    className="me-3 h-44 w-72 overflow-hidden rounded-2xl bg-gray-200"
                  >
                    {b.image ? (
                      <Image source={{ uri: b.image }} className="h-full w-full" resizeMode="cover" />
                    ) : null}
                    {b.title ? (
                      <View className="absolute bottom-0 w-full bg-black/40 p-3">
                        <Text className="font-semibold text-white">{b.title}</Text>
                      </View>
                    ) : null}
                  </View>
                ))}
                <View className="w-4" />
              </ScrollView>
            ) : null}

            {/* Category chips — tap to jump to Search filtered by category */}
            {(data?.categories?.length ?? 0) > 0 ? (
              <View className="mb-5">
                <View className="mb-2 px-4">
                  <SectionHeader title="Browse by category" />
                </View>
                <ScrollView horizontal showsHorizontalScrollIndicator={false} className="ps-4">
                  {data!.categories.map((c) => (
                    <Pressable
                      key={c.id}
                      onPress={() => onCategoryPress(c.id)}
                      className={`me-2 rounded-xl border px-4 py-3 ${
                        selectedCategory === c.id
                          ? "border-brand-500 bg-brand-50"
                          : "border-gray-100 bg-white"
                      }`}
                    >
                      {c.icon ? <Text className="mb-1 text-center text-2xl">{c.icon}</Text> : null}
                      <Text className={`text-xs ${selectedCategory === c.id ? "text-brand-600 font-semibold" : "text-gray-600"}`}>
                        {c.name}
                      </Text>
                    </Pressable>
                  ))}
                  <View className="w-4" />
                </ScrollView>
              </View>
            ) : null}

            {/* Featured tours header */}
            {featured.length > 0 ? (
              <View className="mb-3 px-4">
                <SectionHeader
                  title="Featured tours"
                  action="See all"
                  onAction={() => router.push("/(tabs)/search")}
                />
              </View>
            ) : null}
          </View>
        }
        ListEmptyComponent={
          <View className="px-4">
            <Text className="text-center text-sm text-gray-400 py-12">No tours yet — check back soon!</Text>
          </View>
        }
      />
    </SafeAreaView>
  );
}
