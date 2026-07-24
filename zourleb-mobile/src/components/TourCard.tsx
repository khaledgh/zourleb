import { Image, Pressable, Text, View } from "react-native";
import { useRouter } from "expo-router";
import { useTranslation } from "react-i18next";
import { money } from "@/lib/format";
import type { TourCard as TourCardType } from "@/types/api";

// A tour card with a fixed aspect-ratio cover (cover-crop) to avoid the card
// clipping issues that come from variable image sizes.
export function TourCard({ tour }: { tour: TourCardType }) {
  const router = useRouter();
  const { t } = useTranslation();

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
      </View>
      <View className="p-3">
        <Text numberOfLines={1} className="text-base font-semibold text-gray-900">
          {tour.title}
        </Text>
        <Text numberOfLines={2} className="mt-0.5 text-sm text-gray-500">
          {tour.summary}
        </Text>
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
