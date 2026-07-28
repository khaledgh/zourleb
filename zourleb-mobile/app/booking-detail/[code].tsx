import { Alert, ScrollView, Text, View } from "react-native";
import { useLocalSearchParams, useRouter, Stack } from "expo-router";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiData, apiPost, ApiException } from "@/api/client";
import { Badge, Button, Divider, Loading, RowItem } from "@/components/ui";
import { money, shortDate } from "@/lib/format";
import type { BookingDetail } from "@/types/api";

const CANCELLABLE = new Set(["pending", "confirmed"]);

const statusTone = (s: string) =>
  s === "confirmed" || s === "completed"
    ? "success"
    : s === "cancelled"
      ? "danger"
      : "warn";

export default function BookingDetailScreen() {
  const { code } = useLocalSearchParams<{ code: string }>();
  const router = useRouter();
  const queryClient = useQueryClient();

  const { data: booking, isLoading } = useQuery({
    queryKey: ["booking", code],
    queryFn: () => apiData<BookingDetail>(`/bookings/${code}`),
    enabled: !!code,
  });

  const cancel = useMutation({
    mutationFn: () => apiPost(`/bookings/${code}/cancel`, {}),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["my-bookings"] });
      queryClient.invalidateQueries({ queryKey: ["booking", code] });
      Alert.alert("Cancelled", "Your booking has been cancelled.");
      router.back();
    },
    onError: (e) => {
      Alert.alert("Error", e instanceof ApiException ? e.message : "Could not cancel.");
    },
  });

  function confirmCancel() {
    Alert.alert(
      "Cancel booking?",
      "This action cannot be undone. Are you sure you want to cancel?",
      [
        { text: "Keep booking", style: "cancel" },
        { text: "Cancel booking", style: "destructive", onPress: () => cancel.mutate() },
      ],
    );
  }

  if (isLoading || !booking) return <Loading />;

  return (
    <ScrollView className="flex-1 bg-gray-50" contentContainerClassName="p-4 pb-12">
      <Stack.Screen options={{ title: `Booking ${booking.code}` }} />

      {/* Status banner */}
      <View className={`mb-4 rounded-2xl p-4 ${booking.status === "cancelled" ? "bg-red-50" : "bg-brand-50"}`}>
        <View className="flex-row items-center justify-between">
          <Text className="text-xs text-gray-500">Booking code</Text>
          <Badge label={booking.status.toUpperCase()} tone={statusTone(booking.status)} />
        </View>
        <Text className="mt-1 text-xl font-bold tracking-widest text-gray-800">{booking.code}</Text>
      </View>

      {/* Tour info */}
      <View className="mb-4 rounded-2xl border border-gray-100 bg-white p-4">
        <Text className="mb-2 text-sm font-semibold text-gray-500 uppercase tracking-wide">Tour</Text>
        <Text className="text-base font-semibold text-gray-900">
          {booking.tour_title || `Tour #${booking.tour_id}`}
        </Text>
        {booking.departure_start ? (
          <Text className="mt-1 text-sm text-gray-500">
            {shortDate(booking.departure_start)}
            {booking.departure_end ? ` → ${shortDate(booking.departure_end)}` : ""}
          </Text>
        ) : null}
      </View>

      {/* Details */}
      <View className="mb-4 rounded-2xl border border-gray-100 bg-white p-4">
        <Text className="mb-2 text-sm font-semibold text-gray-500 uppercase tracking-wide">Details</Text>
        <RowItem label="Contact phone" value={booking.contact_phone} />
        <Divider />
        <RowItem label="Travelers" value={String(booking.travelers_count)} />
        <Divider />
        <RowItem label="Payment status" value={booking.payment_status} />
        <Divider />
        <RowItem label="Booked on" value={shortDate(booking.created_at)} />
      </View>

      {/* Travelers list */}
      {booking.travelers && booking.travelers.length > 0 ? (
        <View className="mb-4 rounded-2xl border border-gray-100 bg-white p-4">
          <Text className="mb-3 text-sm font-semibold text-gray-500 uppercase tracking-wide">Travelers</Text>
          {booking.travelers.map((tr, i) => (
            <View key={tr.id ?? i} className="flex-row items-center justify-between py-1.5">
              <Text className="text-sm text-gray-800">{tr.full_name}</Text>
              <Badge label={tr.traveler_type} />
            </View>
          ))}
        </View>
      ) : null}

      {/* Price */}
      <View className="mb-6 rounded-2xl border border-gray-100 bg-white p-4">
        <Text className="mb-2 text-sm font-semibold text-gray-500 uppercase tracking-wide">Total</Text>
        <Text className="text-2xl font-bold text-brand-600">
          {money(booking.subtotal, booking.currency)}
        </Text>
      </View>

      {/* Cancel button — only for cancellable statuses */}
      {CANCELLABLE.has(booking.status) ? (
        <Button
          title="Cancel booking"
          variant="danger"
          onPress={confirmCancel}
          loading={cancel.isPending}
        />
      ) : null}
    </ScrollView>
  );
}
