import { FlatList, Pressable, Text, View } from "react-native";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiList, apiPost } from "@/api/client";
import { Badge, Card, EmptyState, Loading } from "@/components/ui";
import { money } from "@/lib/format";
import type { AgencyTour } from "@/types/api";

export default function AgencyTours() {
  const queryClient = useQueryClient();

  const { data, isLoading, refetch, isRefetching } = useQuery({
    queryKey: ["agency-tours"],
    queryFn: () => apiList<AgencyTour>("/agency/tours", { per_page: 50 }),
  });

  const togglePublish = useMutation({
    mutationFn: ({ id, publish }: { id: number; publish: boolean }) =>
      apiPost(`/agency/tours/${id}/publish`, { publish }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["agency-tours"] });
    },
  });

  if (isLoading) return <Loading />;

  return (
    <FlatList
      className="flex-1 bg-gray-50"
      data={data?.items ?? []}
      keyExtractor={(t) => String(t.id)}
      contentContainerClassName="p-4 pb-12"
      onRefresh={refetch}
      refreshing={isRefetching}
      renderItem={({ item }) => {
        const isPublished = item.status === "published";
        return (
          <Card className="mb-3 overflow-hidden">
            {/* Status strip */}
            <View className={`h-1 w-full ${isPublished ? "bg-green-400" : "bg-gray-200"}`} />
            <View className="p-4">
              <View className="flex-row items-start justify-between gap-2">
                <View className="flex-1">
                  <Text className="text-base font-semibold text-gray-900" numberOfLines={2}>
                    {item.title}
                  </Text>
                  <Text className="mt-0.5 text-xs text-gray-400">
                    {item.type?.replace(/_/g, " ")} · {item.duration_days} day{item.duration_days !== 1 ? "s" : ""}
                  </Text>
                </View>
                <Badge
                  label={isPublished ? "Published" : "Draft"}
                  tone={isPublished ? "success" : "neutral"}
                />
              </View>

              <View className="mt-3 flex-row items-center justify-between">
                <Text className="text-sm font-bold text-brand-600">
                  {money(item.price_from, item.currency)}
                </Text>
                <Pressable
                  onPress={() =>
                    togglePublish.mutate({ id: item.id, publish: !isPublished })
                  }
                  disabled={togglePublish.isPending}
                  className={`rounded-xl px-4 py-2 active:opacity-70 ${
                    isPublished ? "bg-gray-100" : "bg-brand-500"
                  }`}
                >
                  <Text
                    className={`text-sm font-semibold ${
                      isPublished ? "text-gray-700" : "text-white"
                    }`}
                  >
                    {isPublished ? "Unpublish" : "Publish"}
                  </Text>
                </Pressable>
              </View>
            </View>
          </Card>
        );
      }}
      ListEmptyComponent={
        <EmptyState icon="🗺️" message="No tours yet. Create your first tour from the admin panel." />
      }
    />
  );
}
