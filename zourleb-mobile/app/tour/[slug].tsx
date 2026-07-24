import { Image, ScrollView, Text, View, Pressable } from "react-native";
import { useLocalSearchParams, useRouter, Stack } from "expo-router";
import { useQuery, useMutation } from "@tanstack/react-query";
import { useState } from "react";
import { apiData, apiPost } from "@/api/client";
import { useAuth } from "@/stores/auth";
import { Button, Loading, Badge } from "@/components/ui";
import { money, shortDate } from "@/lib/format";
import type { TourDetail } from "@/types/api";

export default function TourDetailScreen() {
  const { slug } = useLocalSearchParams<{ slug: string }>();
  const router = useRouter();
  const user = useAuth((s) => s.user);
  const [faved, setFaved] = useState(false);

  const { data: tour, isLoading } = useQuery({
    queryKey: ["tour", slug],
    queryFn: () => apiData<TourDetail>(`/tours/${slug}`),
    enabled: !!slug,
  });

  const favorite = useMutation({
    mutationFn: (tourId: number) => apiPost(`/favorites/${tourId}`),
    onSuccess: () => setFaved(true),
  });

  if (isLoading || !tour) return <Loading />;

  return (
    <View className="flex-1 bg-white">
      <Stack.Screen options={{ title: tour.title }} />
      <ScrollView contentContainerClassName="pb-28">
        <View className="aspect-[16/10] w-full bg-gray-100">
          {tour.cover ? (
            <Image source={{ uri: tour.cover }} className="h-full w-full" resizeMode="cover" />
          ) : null}
        </View>

        <View className="p-4">
          <View className="flex-row items-start justify-between">
            <Text className="flex-1 text-xl font-bold text-gray-900">{tour.title}</Text>
            {user ? (
              <Pressable onPress={() => favorite.mutate(tour.id)} className="ms-2 p-1">
                <Text className="text-2xl">{faved ? "♥" : "♡"}</Text>
              </Pressable>
            ) : null}
          </View>

          <View className="mt-2 flex-row gap-2">
            <Badge label={tour.type.replace(/_/g, " ")} />
            <Badge label={`${tour.duration_days} day(s)`} />
            {tour.difficulty ? <Badge label={tour.difficulty} /> : null}
          </View>

          <Text className="mt-3 text-sm leading-5 text-gray-600">{tour.summary}</Text>

          {tour.description ? (
            <>
              <Text className="mt-5 text-base font-semibold text-gray-900">About</Text>
              <Text className="mt-1 text-sm leading-5 text-gray-600">{tour.description}</Text>
            </>
          ) : null}

          {tour.departures.length > 0 ? (
            <>
              <Text className="mt-5 text-base font-semibold text-gray-900">Upcoming departures</Text>
              {tour.departures.map((d) => (
                <View key={d.id} className="mt-2 flex-row items-center justify-between rounded-xl border border-gray-100 p-3">
                  <Text className="text-sm text-gray-700">
                    {shortDate(d.start_date)} → {shortDate(d.end_date)}
                  </Text>
                  <Badge
                    label={d.seats_left > 0 ? `${d.seats_left} seats` : "Full"}
                    tone={d.seats_left > 0 ? "success" : "danger"}
                  />
                </View>
              ))}
            </>
          ) : null}

          {tour.past_gallery.length > 0 ? (
            <>
              <Text className="mt-5 text-base font-semibold text-gray-900">From past trips</Text>
              <ScrollView horizontal showsHorizontalScrollIndicator={false} className="mt-2">
                {tour.past_gallery.map((img, i) => (
                  <Image
                    key={i}
                    source={{ uri: img.url }}
                    className="me-2 h-24 w-32 rounded-xl bg-gray-100"
                    resizeMode="cover"
                  />
                ))}
              </ScrollView>
            </>
          ) : null}
        </View>
      </ScrollView>

      <View className="absolute bottom-0 w-full flex-row items-center justify-between border-t border-gray-100 bg-white p-4">
        <View>
          <Text className="text-xs text-gray-400">From</Text>
          <Text className="text-lg font-bold text-brand-600">
            {money(tour.price_from, tour.currency)}
          </Text>
        </View>
        <View className="w-44">
          <Button
            title="Book now"
            onPress={() =>
              user
                ? router.push(`/booking/${tour.slug}`)
                : router.push("/(auth)/login")
            }
          />
        </View>
      </View>
    </View>
  );
}
