import { FlatList, Pressable, Text, View } from "react-native";
import { useRouter } from "expo-router";
import { useQuery } from "@tanstack/react-query";
import { apiList } from "@/api/client";
import { Card, EmptyState, Loading } from "@/components/ui";
import { money } from "@/lib/format";
import type { ProductCard } from "@/types/api";

export default function AgencyProducts() {
  const router = useRouter();

  const { data, isLoading, refetch, isRefetching } = useQuery({
    queryKey: ["agency-products"],
    queryFn: () => apiList<ProductCard>("/agency/products", { per_page: 50 }),
  });

  if (isLoading) return <Loading />;

  return (
    <FlatList
      className="flex-1 bg-gray-50"
      data={data?.items ?? []}
      keyExtractor={(p) => String(p.id)}
      contentContainerClassName="p-4 pb-12"
      onRefresh={refetch}
      refreshing={isRefetching}
      ListHeaderComponent={
        <Pressable
          onPress={() => router.push("/(agency)/new-product")}
          className="mb-4 rounded-2xl border border-brand-200 bg-brand-50 p-4 active:opacity-80"
        >
          <Text className="text-center text-base font-semibold text-brand-700">+ Add product</Text>
        </Pressable>
      }
      renderItem={({ item }) => (
        <Card className="mb-3 p-4">
          <View className="flex-row items-center justify-between">
            <View className="flex-1">
              <Text className="text-base font-semibold text-gray-900" numberOfLines={1}>
                {item.name}
              </Text>
              <Text className="text-xs text-gray-500">Stock {item.stock}</Text>
            </View>
            <Text className="font-bold text-brand-600">{money(item.price, item.currency)}</Text>
          </View>
        </Card>
      )}
      ListEmptyComponent={
        <EmptyState icon="🛍️" message="No products yet. Add your first product." />
      }
    />
  );
}
