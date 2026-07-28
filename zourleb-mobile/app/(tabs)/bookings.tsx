import { FlatList, Pressable, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useQuery } from "@tanstack/react-query";
import { useRouter } from "expo-router";
import { apiList } from "@/api/client";
import { useAuth } from "@/stores/auth";
import { SignInPrompt } from "@/components/SignInPrompt";
import { Badge, Card, EmptyState, Loading } from "@/components/ui";
import { money, shortDate } from "@/lib/format";
import type { Booking } from "@/types/api";

const tone = (status: string) =>
  status === "confirmed" || status === "completed"
    ? "success"
    : status === "cancelled"
      ? "danger"
      : "warn";

export default function Bookings() {
  const user = useAuth((s) => s.user);
  const router = useRouter();

  const { data, isLoading, refetch, isRefetching } = useQuery({
    queryKey: ["my-bookings"],
    queryFn: () => apiList<Booking>("/bookings", { per_page: 30 }),
    enabled: !!user,
  });

  if (!user) {
    return (
      <SafeAreaView className="flex-1 bg-gray-50">
        <SignInPrompt message="Sign in to see your bookings." />
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView className="flex-1 bg-gray-50" edges={["top"]}>
      <Text className="px-4 py-4 text-2xl font-bold text-gray-900">My bookings</Text>
      {isLoading ? (
        <Loading />
      ) : (
        <FlatList
          data={data?.items ?? []}
          keyExtractor={(b) => b.code}
          contentContainerClassName="px-4 pb-8"
          onRefresh={refetch}
          refreshing={isRefetching}
          renderItem={({ item }) => (
            <Pressable onPress={() => router.push(`/booking-detail/${item.code}`)}>
              <Card className="mb-3 overflow-hidden active:opacity-80">
                {/* Colored status strip */}
                <View
                  className={`h-1 w-full ${
                    item.status === "confirmed" || item.status === "completed"
                      ? "bg-green-400"
                      : item.status === "cancelled"
                        ? "bg-red-400"
                        : "bg-amber-400"
                  }`}
                />
                <View className="p-4">
                  <View className="flex-row items-center justify-between">
                    <Text className="flex-1 text-base font-semibold text-gray-900" numberOfLines={1}>
                      {item.tour_title || `Tour #${item.tour_id}`}
                    </Text>
                    <Badge label={item.status} tone={tone(item.status)} />
                  </View>
                  <Text className="mt-1 text-xs text-gray-400">
                    {item.code} · {shortDate(item.created_at)}
                  </Text>
                  <View className="mt-2 flex-row items-center justify-between">
                    <Text className="text-sm text-gray-600">
                      {item.travelers_count} traveler{item.travelers_count !== 1 ? "s" : ""}
                    </Text>
                    <Text className="text-sm font-bold text-brand-600">
                      {money(item.subtotal, item.currency)}
                    </Text>
                  </View>
                  <Text className="mt-2 text-xs text-brand-600">Tap to view details →</Text>
                </View>
              </Card>
            </Pressable>
          )}
          ListEmptyComponent={
            <EmptyState
              icon="🎫"
              message="No bookings yet — explore tours to get started."
            />
          }
        />
      )}
    </SafeAreaView>
  );
}
