import { ScrollView, Text, View, Pressable } from "react-native";
import { useLocalSearchParams, useRouter, Stack } from "expo-router";
import { useQuery } from "@tanstack/react-query";
import { apiData } from "@/api/client";
import { Badge, Card, Divider, Loading, SectionHeader } from "@/components/ui";
import { money, shortDate } from "@/lib/format";
import type { TourDetail } from "@/types/api";

export default function AgencyTourDetail() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();

  const { data: tour, isLoading } = useQuery({
    queryKey: ["agency-tour", id],
    queryFn: () => apiData<TourDetail>(`/agency/tours/${id}`),
    enabled: !!id,
  });

  if (isLoading || !tour) return <Loading />;

  return (
    <ScrollView className="flex-1 bg-gray-50" contentContainerClassName="p-4 pb-12">
      <Stack.Screen options={{ title: tour.title }} />
      <Text className="text-xl font-bold text-gray-900">{tour.title}</Text>
      <View className="mt-2 flex-row flex-wrap gap-2">
        <Badge label={tour.status} />
        <Badge label={`${tour.duration_days} day${tour.duration_days !== 1 ? "s" : ""}`} />
        {tour.difficulty ? <Badge label={tour.difficulty} /> : null}
      </View>

      <Text className="mt-3 text-sm text-gray-600">{tour.summary}</Text>

      <Divider />
      <SectionHeader
        title="Departures"
        action="+ Add"
        onAction={() => router.push(`/agency/new-departure/${id}`)}
      />
      {tour.departures.length === 0 ? (
        <Text className="text-sm text-gray-500">No departures yet.</Text>
      ) : (
        tour.departures.map((d) => (
          <Card key={d.id} className="mb-3 p-3">
            <View className="flex-row items-center justify-between">
              <Text className="text-sm text-gray-900">
                {shortDate(d.start_date)} → {shortDate(d.end_date)}
              </Text>
              <Badge
                label={d.seats_left > 0 ? `${d.seats_left} seats` : "Full"}
                tone={d.seats_left > 0 ? "success" : "danger"}
              />
            </View>
            <Text className="text-xs text-gray-500">Capacity {d.capacity}</Text>
          </Card>
        ))
      )}

      <Divider />
      <SectionHeader
        title="Prices"
        action="+ Add"
        onAction={() => router.push(`/agency/new-price/${id}`)}
      />
      {tour.prices.length === 0 ? (
        <Text className="text-sm text-gray-500">No prices yet.</Text>
      ) : (
        tour.prices.map((p, i) => (
          <Card key={i} className="mb-3 p-3">
            <View className="flex-row items-center justify-between">
              <Text className="text-sm text-gray-900">{p.traveler_type}</Text>
              <Text className="font-bold text-brand-600">{money(p.amount, p.currency)}</Text>
            </View>
          </Card>
        ))
      )}

      <Divider />
      <SectionHeader
        title="Images"
        action="+ Add"
        onAction={() => router.push(`/agency/new-image/${id}`)}
      />
      <Text className="text-sm text-gray-500">{tour.images.length} image(s)</Text>
    </ScrollView>
  );
}
