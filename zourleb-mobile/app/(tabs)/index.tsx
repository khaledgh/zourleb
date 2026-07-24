import { FlatList, Image, ScrollView, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useQuery } from "@tanstack/react-query";
import { apiData } from "@/api/client";
import { TourCard } from "@/components/TourCard";
import { Loading, EmptyState } from "@/components/ui";
import type { HomePayload } from "@/types/api";

export default function Home() {
  const { data, isLoading, refetch, isRefetching } = useQuery({
    queryKey: ["home"],
    queryFn: () => apiData<HomePayload>("/home"),
  });

  if (isLoading) return <Loading />;

  return (
    <SafeAreaView className="flex-1 bg-gray-50" edges={["top"]}>
      <FlatList
        data={data?.featured ?? []}
        keyExtractor={(t) => String(t.id)}
        renderItem={({ item }) => <TourCard tour={item} />}
        contentContainerClassName="px-4 pb-8"
        onRefresh={refetch}
        refreshing={isRefetching}
        ListHeaderComponent={
          <View>
            <Text className="py-4 text-2xl font-bold text-gray-900">
              Discover Lebanon
            </Text>

            {(data?.banners?.length ?? 0) > 0 ? (
              <ScrollView
                horizontal
                showsHorizontalScrollIndicator={false}
                className="mb-5"
              >
                {data!.banners.map((b) => (
                  <View
                    key={b.id}
                    className="me-3 h-40 w-72 overflow-hidden rounded-2xl bg-gray-200"
                  >
                    {b.image ? (
                      <Image source={{ uri: b.image }} className="h-full w-full" resizeMode="cover" />
                    ) : null}
                    {b.title ? (
                      <View className="absolute bottom-0 w-full bg-black/30 p-3">
                        <Text className="font-semibold text-white">{b.title}</Text>
                      </View>
                    ) : null}
                  </View>
                ))}
              </ScrollView>
            ) : null}

            {(data?.categories?.length ?? 0) > 0 ? (
              <ScrollView horizontal showsHorizontalScrollIndicator={false} className="mb-5">
                {data!.categories.map((c) => (
                  <View key={c.id} className="me-2 rounded-full bg-white px-4 py-2 border border-gray-100">
                    <Text className="text-sm text-gray-700">{c.name}</Text>
                  </View>
                ))}
              </ScrollView>
            ) : null}

            <Text className="mb-3 text-lg font-semibold text-gray-900">Featured tours</Text>
          </View>
        }
        ListEmptyComponent={<EmptyState message="No featured tours yet." />}
      />
    </SafeAreaView>
  );
}
