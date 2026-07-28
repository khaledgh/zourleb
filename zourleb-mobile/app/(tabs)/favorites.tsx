import { FlatList, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiData, apiDelete, apiList } from "@/api/client";
import { useAuth } from "@/stores/auth";
import { SignInPrompt } from "@/components/SignInPrompt";
import { TourCard } from "@/components/TourCard";
import { EmptyState, Loading } from "@/components/ui";
import type { TourCard as TourCardType } from "@/types/api";

interface FavoritesPayload {
  tour_ids: number[];
}

export default function Favorites() {
  const user = useAuth((s) => s.user);
  const queryClient = useQueryClient();

  const { data: favPayload } = useQuery({
    queryKey: ["favorite-ids"],
    queryFn: () => apiData<FavoritesPayload>("/favorites"),
    enabled: !!user,
  });

  const favIds = new Set(favPayload?.tour_ids ?? []);

  const { data: allTours, isLoading } = useQuery({
    queryKey: ["favorite-tours", favPayload?.tour_ids],
    queryFn: async () => {
      const result = await apiList<TourCardType>("/tours", { per_page: 100 });
      return result.items;
    },
    enabled: !!user && !!favPayload,
  });

  const removeFav = useMutation({
    mutationFn: (tourId: number) => apiDelete(`/favorites/${tourId}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["favorite-ids"] });
      queryClient.invalidateQueries({ queryKey: ["favorite-tours"] });
    },
  });

  const favoriteTours = (allTours ?? []).filter((t) => favIds.has(t.id));

  if (!user) {
    return (
      <SafeAreaView className="flex-1 bg-gray-50">
        <SignInPrompt message="Sign in to save and view your favorite tours." />
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView className="flex-1 bg-gray-50" edges={["top"]}>
      <Text className="px-4 py-4 text-2xl font-bold text-gray-900">Favorites</Text>
      {isLoading ? (
        <Loading />
      ) : (
        <FlatList
          data={favoriteTours}
          keyExtractor={(t) => String(t.id)}
          renderItem={({ item }) => (
            <View className="px-4">
              <TourCard
                tour={item}
                isFavorited={true}
                onFavChange={(tourId, nowFaved) => {
                  if (!nowFaved) removeFav.mutate(tourId);
                }}
              />
            </View>
          )}
          contentContainerClassName="pb-8"
          ListEmptyComponent={
            <EmptyState
              icon="♡"
              message="No favorites yet. Tap the heart on a tour to save it."
            />
          }
        />
      )}
    </SafeAreaView>
  );
}
