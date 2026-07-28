import { FlatList, Text, View } from "react-native";
import { useQuery } from "@tanstack/react-query";
import { apiList } from "@/api/client";
import { Badge, Card, EmptyState, Loading } from "@/components/ui";
import { money, shortDate } from "@/lib/format";
import type { Booking } from "@/types/api";

const tone = (status: string) =>
  status === "confirmed" || status === "completed"
    ? "success"
    : status === "cancelled"
      ? "danger"
      : "warn";

export default function AgencyBookings() {
  const { data, isLoading, refetch, isRefetching } = useQuery({
    queryKey: ["agency-bookings"],
    queryFn: () => apiList<Booking>("/agency/bookings", { per_page: 50 }),
  });

  if (isLoading) return <Loading />;

  return (
    <FlatList
      className="flex-1 bg-gray-50"
      data={data?.items ?? []}
      keyExtractor={(b) => b.code}
      contentContainerClassName="p-4 pb-12"
      onRefresh={refetch}
      refreshing={isRefetching}
      renderItem={({ item }) => (
        <Card className="mb-3 overflow-hidden">
          <View className={`h-1 w-full ${
            item.status === "confirmed" || item.status === "completed"
              ? "bg-green-400"
              : item.status === "cancelled"
                ? "bg-red-400"
                : "bg-amber-400"
          }`} />
          <View className="p-4">
            <View className="flex-row items-start justify-between gap-2">
              <View className="flex-1">
                <Text className="text-sm font-semibold text-gray-900" numberOfLines={1}>
                  {item.tour_title || `Tour #${item.tour_id}`}
                </Text>
                <Text className="mt-0.5 text-xs text-gray-400">
                  {item.code} · {shortDate(item.created_at)}
                </Text>
              </View>
              <Badge label={item.status} tone={tone(item.status)} />
            </View>

            <View className="mt-3 flex-row items-center justify-between">
              <View>
                <Text className="text-xs text-gray-500">
                  📞 {item.contact_phone}
                </Text>
                <Text className="mt-0.5 text-xs text-gray-500">
                  👥 {item.travelers_count} traveler{item.travelers_count !== 1 ? "s" : ""}
                </Text>
              </View>
              <View className="items-end">
                <Text className="text-sm font-bold text-brand-600">
                  {money(item.subtotal, item.currency)}
                </Text>
                <Text className="mt-0.5 text-xs text-gray-400">{item.payment_status}</Text>
              </View>
            </View>
          </View>
        </Card>
      )}
      ListEmptyComponent={
        <EmptyState icon="🎫" message="No bookings yet for your tours." />
      }
    />
  );
}
