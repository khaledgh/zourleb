import { Image, Pressable, Text, View } from "react-native";
import { useRouter } from "expo-router";
import { useTranslation } from "react-i18next";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiPost, apiDelete } from "@/api/client";
import { StarRating } from "@/components/ui";
import { money } from "@/lib/format";
import type { TourCard as TourCardType } from "@/types/api";

interface Props {
  tour: TourCardType;
  /** Whether this tour is already favorited (controlled from parent). */
  isFavorited?: boolean;
  /** Called after a successful favorite toggle so the parent can update its list. */
  onFavChange?: (tourId: number, nowFaved: boolean) => void;
}

// A tour card with a fixed aspect-ratio cover and an inline favorite toggle.
export function TourCard({ tour, isFavorited = false, onFavChange }: Props) {
  const router = useRouter();
  const { t } = useTranslation();
  const queryClient = useQueryClient();

  const toggleFav = useMutation({
    mutationFn: () =>
      isFavorited
        ? apiDelete(`/favorites/${tour.id}`)
        : apiPost(`/favorites/${tour.id}`, {}),
    onSuccess: () => {
      const next = !isFavorited;
      onFavChange?.(tour.id, next);
      queryClient.invalidateQueries({ queryKey: ["favorite-ids"] });
      queryClient.invalidateQueries({ queryKey: ["favorite-tours"] });
    },
  });

  return (
    <Pressable
      onPress={() => router.push(`/tour/${tour.slug}`)}
      className="mb-4 overflow-hidden rounded-2xl border border-gray-100 bg-white active:opacity-90"
    >
      <View className="aspect-[16/10] w-full bg-gray-100">
        {tour.cover ? (
          <Image
            source={{ uri: tour.cover }}
            className="h-full w-full"
            resizeMode="cover"
          />
        ) : null}
        {tour.featured ? (
          <View className="absolute left-3 top-3 rounded-full bg-brand-500 px-2 py-0.5">
            <Text className="text-xs font-semibold text-white">★ Featured</Text>
          </View>
        ) : null}
        {/* Favorite heart toggle */}
        <Pressable
          onPress={(e) => {
            e.stopPropagation?.();
            toggleFav.mutate();
          }}
          className="absolute right-3 top-3 h-8 w-8 items-center justify-center rounded-full bg-white/80"
          hitSlop={8}
        >
          <Text className={`text-lg ${isFavorited ? "text-red-500" : "text-gray-400"}`}>
            {isFavorited ? "♥" : "♡"}
          </Text>
        </Pressable>
      </View>

      <View className="p-3">
        <Text numberOfLines={1} className="text-base font-semibold text-gray-900">
          {tour.title}
        </Text>
        <Text numberOfLines={2} className="mt-0.5 text-sm text-gray-500">
          {tour.summary}
        </Text>

        {tour.rating_avg > 0 ? (
          <View className="mt-1.5">
            <StarRating rating={tour.rating_avg} showNumber />
          </View>
        ) : null}

        <View className="mt-2 flex-row items-center justify-between">
          <View className="flex-row items-center gap-1">
            <Text className="text-xs text-gray-400">{tour.agency.name}</Text>
            {tour.agency.verified ? (
              <Text className="text-xs text-brand-600">✓</Text>
            ) : null}
          </View>
          <Text className="text-sm font-bold text-brand-600">
            {t("common.from")} {money(tour.price_from, tour.currency)}
          </Text>
        </View>
      </View>
    </Pressable>
  );
}
