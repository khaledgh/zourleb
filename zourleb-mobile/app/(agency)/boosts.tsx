import { FlatList, Text, View } from "react-native";
import { useQuery } from "@tanstack/react-query";
import { apiData, apiList } from "@/api/client";
import { Card, EmptyState, Loading } from "@/components/ui";
import { money } from "@/lib/format";
import type { AgencyBoost, BoostPackage } from "@/types/api";

export default function AgencyBoosts() {
  const { data: boosts, isLoading: boostsLoading } = useQuery({
    queryKey: ["agency-boosts"],
    queryFn: () => apiList<AgencyBoost>("/agency/boosts", { per_page: 50 }),
  });

  const { data: packages, isLoading: packagesLoading } = useQuery({
    queryKey: ["boost-packages"],
    queryFn: () => apiData<BoostPackage[]>("/agency/boosts/packages"),
  });

  if (boostsLoading || packagesLoading) return <Loading />;

  return (
    <FlatList
      className="flex-1 bg-gray-50"
      data={boosts?.items ?? []}
      keyExtractor={(b) => String(b.id)}
      contentContainerClassName="p-4 pb-12"
      ListHeaderComponent={
        <View className="mb-4">
          <Text className="mb-2 text-base font-semibold text-gray-900">Available packages</Text>
          {(packages ?? []).length === 0 ? (
            <Text className="text-sm text-gray-500">No boost packages available.</Text>
          ) : (
            packages!.map((p) => (
              <Card key={p.id} className="mb-2 p-3">
                <View className="flex-row items-center justify-between">
                  <View>
                    <Text className="font-semibold text-gray-900">{p.name}</Text>
                    <Text className="text-xs text-gray-500">{p.placement} · {p.duration_days} day(s)</Text>
                  </View>
                  <Text className="font-bold text-brand-600">{money(p.price, p.currency)}</Text>
                </View>
              </Card>
            ))
          )}
          <Text className="mt-4 mb-2 text-base font-semibold text-gray-900">My boosts</Text>
        </View>
      }
      renderItem={({ item }) => (
        <Card className="mb-3 p-4">
          <View className="flex-row items-center justify-between">
            <View>
              <Text className="font-semibold text-gray-900">Tour #{item.tour_id ?? "—"}</Text>
              <Text className="text-xs text-gray-500">{item.placement} · {item.status}</Text>
            </View>
            <Text className={`text-xs ${item.status === "active" ? "text-brand-600" : "text-gray-500"}`}>
              {item.status}
            </Text>
          </View>
        </Card>
      )}
      ListEmptyComponent={
        <EmptyState icon="🚀" message="No active boosts. Purchase a package to promote a tour." />
      }
    />
  );
}
