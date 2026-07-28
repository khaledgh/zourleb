import { Pressable, ScrollView, Text, View } from "react-native";
import { useRouter } from "expo-router";
import { useQuery } from "@tanstack/react-query";
import { apiData, apiList } from "@/api/client";
import { Card, Loading, SectionHeader } from "@/components/ui";
import type { AgencyProfile, AgencyTour, Booking } from "@/types/api";

export default function AgencyDashboard() {
  const router = useRouter();

  const { data: profile, isLoading: profileLoading } = useQuery({
    queryKey: ["agency-profile"],
    queryFn: () => apiData<AgencyProfile>("/agency/profile"),
  });

  const { data: toursData } = useQuery({
    queryKey: ["agency-tours"],
    queryFn: () => apiList<AgencyTour>("/agency/tours", { per_page: 3 }),
  });

  const { data: bookingsData } = useQuery({
    queryKey: ["agency-bookings"],
    queryFn: () => apiList<Booking>("/agency/bookings", { per_page: 3 }),
  });

  if (profileLoading) return <Loading />;

  const totalTours = toursData?.meta?.total ?? 0;
  const pendingBookings = (bookingsData?.items ?? []).filter(
    (b) => b.status === "pending",
  ).length;

  return (
    <ScrollView className="flex-1 bg-gray-50" contentContainerClassName="p-4 pb-12">
      {/* Agency identity card */}
      <Card className="mb-5 p-4">
        <Text className="text-lg font-bold text-gray-900">{profile?.name ?? "My Agency"}</Text>
        <View className="mt-1 flex-row items-center gap-2">
          <Text className="text-sm text-gray-500">{profile?.email}</Text>
          {profile?.verified ? (
            <Text className="text-xs text-brand-600 font-semibold">✓ Verified</Text>
          ) : (
            <Text className="text-xs text-amber-600">Pending approval</Text>
          )}
        </View>
      </Card>

      {/* Summary tiles */}
      <View className="mb-5 flex-row gap-3">
        <StatTile icon="🗺️" label="Total tours" value={totalTours} />
        <StatTile icon="⏳" label="Pending bookings" value={pendingBookings} urgent={pendingBookings > 0} />
      </View>

      {/* Tours shortcut */}
      <SectionHeader
        title="Tours"
        action="View all"
        onAction={() => router.push("/(agency)/tours")}
      />
      <Pressable
        onPress={() => router.push("/(agency)/tours")}
        className="mb-5 rounded-2xl border border-gray-100 bg-white p-4 active:opacity-80"
      >
        <Text className="text-sm text-gray-600">
          Manage your tours, publish or unpublish them, and track departures.
        </Text>
        <Text className="mt-2 text-sm font-semibold text-brand-600">
          {totalTours} tour{totalTours !== 1 ? "s" : ""} →
        </Text>
      </Pressable>

      {/* Bookings shortcut */}
      <SectionHeader
        title="Bookings"
        action="View all"
        onAction={() => router.push("/(agency)/bookings")}
      />
      <Pressable
        onPress={() => router.push("/(agency)/bookings")}
        className="rounded-2xl border border-gray-100 bg-white p-4 active:opacity-80"
      >
        <Text className="text-sm text-gray-600">
          Track customer bookings across all your tours.
        </Text>
        {pendingBookings > 0 ? (
          <Text className="mt-2 text-sm font-semibold text-amber-600">
            {pendingBookings} pending booking{pendingBookings !== 1 ? "s" : ""} →
          </Text>
        ) : (
          <Text className="mt-2 text-sm font-semibold text-brand-600">View all →</Text>
        )}
      </Pressable>
    </ScrollView>
  );
}

function StatTile({
  icon,
  label,
  value,
  urgent = false,
}: {
  icon: string;
  label: string;
  value: number;
  urgent?: boolean;
}) {
  return (
    <View className={`flex-1 rounded-2xl border p-4 ${urgent ? "border-amber-200 bg-amber-50" : "border-gray-100 bg-white"}`}>
      <Text className="text-2xl">{icon}</Text>
      <Text className={`mt-2 text-2xl font-bold ${urgent ? "text-amber-700" : "text-gray-900"}`}>
        {value}
      </Text>
      <Text className="mt-0.5 text-xs text-gray-500">{label}</Text>
    </View>
  );
}
