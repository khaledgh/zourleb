import { useQuery } from "@tanstack/react-query";
import { apiData } from "@/api/client";
import { useAuth } from "@/stores/auth";

interface FavoritesPayload {
  tour_ids: number[];
}

/** Fetches the signed-in user's favorite tour IDs. Returns an empty set for guests. */
export function useFavoriteIds() {
  const user = useAuth((s) => s.user);
  const { data, isLoading } = useQuery({
    queryKey: ["favorite-ids"],
    queryFn: () => apiData<FavoritesPayload>("/favorites"),
    enabled: !!user,
    staleTime: 60_000,
  });

  const ids = new Set(data?.tour_ids ?? []);
  return { ids, isLoading };
}
