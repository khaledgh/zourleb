import { FlatList, Text, View } from "react-native";
import { useQuery } from "@tanstack/react-query";
import { apiData } from "@/api/client";
import { Card, EmptyState, Loading } from "@/components/ui";
import type { AgencyMember } from "@/types/api";

export default function AgencyMembers() {
  const { data, isLoading, refetch, isRefetching } = useQuery({
    queryKey: ["agency-members"],
    queryFn: () => apiData<AgencyMember[]>("/agency/members"),
  });

  if (isLoading) return <Loading />;

  return (
    <FlatList
      className="flex-1 bg-gray-50"
      data={data ?? []}
      keyExtractor={(m) => `${m.agency_id}-${m.user_id}`}
      contentContainerClassName="p-4 pb-12"
      onRefresh={refetch}
      refreshing={isRefetching}
      renderItem={({ item }) => (
        <Card className="mb-3 p-4">
          <View className="flex-row items-center justify-between">
            <View>
              <Text className="text-base font-semibold text-gray-900">User #{item.user_id}</Text>
              <Text className="text-xs text-gray-500">Role {item.role_id}</Text>
            </View>
            <Text className="text-xs text-brand-600">
              {item.joined_at ? "Joined" : item.invited_at ? "Invited" : "Pending"}
            </Text>
          </View>
        </Card>
      )}
      ListEmptyComponent={
        <EmptyState icon="👥" message="No members yet. Invite staff from the admin panel." />
      }
    />
  );
}
