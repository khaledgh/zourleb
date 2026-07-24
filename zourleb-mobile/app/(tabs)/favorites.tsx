import { FlatList, Text } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useQuery } from "@tanstack/react-query";
import { apiData } from "@/api/client";
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

  // The backend returns favorite tour ids; we resolve them to cards via the
  // public catalog (filtered client-side here for simplicity).
  const { data: favIds } = useQuery({
    queryKey: ["favorite-ids"],
    queryFn: () => apiData<FavoritesPayload>("/favorites"),
    enabled: !!user,
  });

  const { data: tours, isLoading } = useQuery({
    queryKey: ["favorite-tours", favIds?.tour_ids],
    queryFn: async () => {
      const env = await apiData<TourCardType[]>("/tours", { per_page: 100 });
      const ids = new Set(favIds?.tour_ids ?? []);
      return env.filter((t) => ids.has(t.id));
    },
    enabled: !!user && !!favIds,
  });

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
          data={tours ?? []}
          keyExtractor={(t) => String(t.id)}
          renderItem={({ item }) => <TourCard tour={item} />}
          contentContainerClassName="px-4 pb-8"
          ListEmptyComponent={<EmptyState message="No favorites yet. Tap ♥ on a tour to save it." />}
        />
      )}
    </SafeAreaView>
  );
}
