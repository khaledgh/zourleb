import { FlatList, Pressable, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiList, apiPost } from "@/api/client";
import { useAuth } from "@/stores/auth";
import { SignInPrompt } from "@/components/SignInPrompt";
import { EmptyState, Loading } from "@/components/ui";
import { shortDate } from "@/lib/format";
import type { Notification } from "@/types/api";

export default function Notifications() {
  const user = useAuth((s) => s.user);
  const queryClient = useQueryClient();

  const { data, isLoading, refetch, isRefetching } = useQuery({
    queryKey: ["notifications"],
    queryFn: () => apiList<Notification>("/notifications", { per_page: 50 }),
    enabled: !!user,
    refetchInterval: 60_000, // poll every 60 s while screen is open
  });

  const markRead = useMutation({
    mutationFn: (id: number) => apiPost(`/notifications/${id}/read`, {}),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] });
    },
  });

  if (!user) {
    return (
      <SafeAreaView className="flex-1 bg-gray-50">
        <SignInPrompt message="Sign in to see your notifications." />
      </SafeAreaView>
    );
  }

  const notifications = data?.items ?? [];
  const unreadCount = notifications.filter((n) => !n.read).length;

  return (
    <SafeAreaView className="flex-1 bg-gray-50" edges={["top"]}>
      <View className="px-4 pt-4 pb-2 flex-row items-center justify-between">
        <Text className="text-2xl font-bold text-gray-900">Notifications</Text>
        {unreadCount > 0 ? (
          <View className="rounded-full bg-brand-500 px-2.5 py-0.5">
            <Text className="text-xs font-bold text-white">{unreadCount}</Text>
          </View>
        ) : null}
      </View>

      {isLoading ? (
        <Loading />
      ) : (
        <FlatList
          data={notifications}
          keyExtractor={(n) => String(n.id)}
          contentContainerClassName="px-4 pb-8 pt-2"
          onRefresh={refetch}
          refreshing={isRefetching}
          renderItem={({ item }) => (
            <Pressable
              onPress={() => {
                if (!item.read) markRead.mutate(item.id);
              }}
              className={`mb-2 overflow-hidden rounded-2xl border bg-white active:opacity-80 ${
                item.read ? "border-gray-100" : "border-brand-200"
              }`}
            >
              {/* Unread indicator strip */}
              {!item.read ? (
                <View className="h-1 w-full bg-brand-400" />
              ) : null}
              <View className="p-4">
                <View className="flex-row items-start gap-3">
                  <View
                    className={`mt-0.5 h-8 w-8 items-center justify-center rounded-full ${
                      item.read ? "bg-gray-100" : "bg-brand-50"
                    }`}
                  >
                    <Text className="text-base">{notifIcon(item.type)}</Text>
                  </View>
                  <View className="flex-1">
                    <Text className={`text-sm font-semibold ${item.read ? "text-gray-700" : "text-gray-900"}`}>
                      {item.title}
                    </Text>
                    <Text className="mt-0.5 text-sm text-gray-500">{item.body}</Text>
                    <Text className="mt-1 text-xs text-gray-400">{shortDate(item.created_at)}</Text>
                  </View>
                  {!item.read ? (
                    <View className="mt-1.5 h-2 w-2 rounded-full bg-brand-500" />
                  ) : null}
                </View>
              </View>
            </Pressable>
          )}
          ListEmptyComponent={
            <EmptyState icon="🔔" message="No notifications yet." />
          }
        />
      )}
    </SafeAreaView>
  );
}

function notifIcon(type: string): string {
  const map: Record<string, string> = {
    booking_confirmed: "✅",
    booking_cancelled: "❌",
    booking_reminder: "⏰",
    review_reply: "💬",
    boost_activated: "⚡",
    boost_expired: "📉",
    general: "📢",
  };
  return map[type] ?? "🔔";
}
