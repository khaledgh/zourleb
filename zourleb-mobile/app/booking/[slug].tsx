import { useState } from "react";
import { ScrollView, Text, TextInput, View, Pressable, Alert } from "react-native";
import { useLocalSearchParams, useRouter, Stack } from "expo-router";
import { useQuery } from "@tanstack/react-query";
import { apiData, apiPost, ApiException } from "@/api/client";
import { Button, Loading, Field, Badge } from "@/components/ui";
import { shortDate } from "@/lib/format";
import type { Booking, TourDetail } from "@/types/api";

type Step = "details" | "otp" | "confirm";

// Booking flow: collect contact phone + traveler → verify phone via OTP →
// create the booking. The backend enforces phone verification at booking time.
export default function BookingScreen() {
  const { slug } = useLocalSearchParams<{ slug: string }>();
  const router = useRouter();

  const { data: tour, isLoading } = useQuery({
    queryKey: ["tour", slug],
    queryFn: () => apiData<TourDetail>(`/tours/${slug}`),
    enabled: !!slug,
  });

  const [step, setStep] = useState<Step>("details");
  const [phone, setPhone] = useState("");
  const [travelerName, setTravelerName] = useState("");
  const [departureId, setDepartureId] = useState<number | null>(null);
  const [code, setCode] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  if (isLoading || !tour) return <Loading />;

  async function requestOtp() {
    setError(null);
    if (!phone || !travelerName) {
      setError("Enter your name and phone number.");
      return;
    }
    setBusy(true);
    try {
      await apiPost("/otp/request", { phone, channel: "whatsapp", purpose: "booking" });
      setStep("otp");
    } catch (e) {
      setError(e instanceof ApiException ? e.message : "Could not send code.");
    } finally {
      setBusy(false);
    }
  }

  async function verifyOtp() {
    setError(null);
    setBusy(true);
    try {
      await apiPost("/otp/verify", { phone, code });
      setStep("confirm");
    } catch (e) {
      setError(e instanceof ApiException ? e.message : "Invalid code.");
    } finally {
      setBusy(false);
    }
  }

  async function confirmBooking() {
    if (!tour) return;
    setError(null);
    setBusy(true);
    try {
      const booking = await apiPost<Booking>("/bookings", {
        tour_id: tour.id,
        departure_id: departureId,
        contact_phone: phone,
        travelers: [{ full_name: travelerName, traveler_type: "adult" }],
      });
      Alert.alert("Booking confirmed", `Your booking code is ${booking.code}.`);
      router.replace("/(tabs)/bookings");
    } catch (e) {
      setError(e instanceof ApiException ? e.message : "Booking failed.");
    } finally {
      setBusy(false);
    }
  }

  return (
    <ScrollView className="flex-1 bg-white" contentContainerClassName="p-4">
      <Stack.Screen options={{ title: "Book" }} />
      <Text className="text-lg font-bold text-gray-900">{tour.title}</Text>
      <StepIndicator step={step} />

      {step === "details" && (
        <View className="mt-4">
          {tour.departures.length > 0 ? (
            <View className="mb-4">
              <Text className="mb-1 text-sm font-medium text-gray-700">Choose a departure</Text>
              {tour.departures.map((d) => (
                <Pressable
                  key={d.id}
                  onPress={() => setDepartureId(d.id)}
                  className={`mt-2 flex-row items-center justify-between rounded-xl border p-3 ${
                    departureId === d.id ? "border-brand-500 bg-brand-50" : "border-gray-200"
                  }`}
                >
                  <Text className="text-sm text-gray-700">
                    {shortDate(d.start_date)} → {shortDate(d.end_date)}
                  </Text>
                  <Badge
                    label={d.seats_left > 0 ? `${d.seats_left} left` : "Full"}
                    tone={d.seats_left > 0 ? "success" : "danger"}
                  />
                </Pressable>
              ))}
            </View>
          ) : null}

          <Field label="Lead traveler name">
            <TextInput
              className="h-12 rounded-xl border border-gray-300 px-3 text-base"
              value={travelerName}
              onChangeText={setTravelerName}
            />
          </Field>
          <Field label="Contact phone (Lebanese)">
            <TextInput
              className="h-12 rounded-xl border border-gray-300 px-3 text-base"
              keyboardType="phone-pad"
              placeholder="03 123 456"
              value={phone}
              onChangeText={setPhone}
            />
          </Field>

          {error ? <Text className="mb-3 text-sm text-red-600">{error}</Text> : null}
          <Button title="Send verification code" onPress={requestOtp} loading={busy} />
          <Text className="mt-2 text-center text-xs text-gray-400">
            We'll send a WhatsApp code to verify your number.
          </Text>
        </View>
      )}

      {step === "otp" && (
        <View className="mt-4">
          <Text className="mb-3 text-sm text-gray-600">
            Enter the code sent to {phone}.
          </Text>
          <Field label="Verification code">
            <TextInput
              className="h-12 rounded-xl border border-gray-300 px-3 text-center text-lg tracking-[8px]"
              keyboardType="number-pad"
              maxLength={6}
              value={code}
              onChangeText={setCode}
            />
          </Field>
          {error ? <Text className="mb-3 text-sm text-red-600">{error}</Text> : null}
          <Button title="Verify" onPress={verifyOtp} loading={busy} />
          <Pressable onPress={requestOtp} className="mt-3">
            <Text className="text-center text-sm text-brand-600">Resend code</Text>
          </Pressable>
        </View>
      )}

      {step === "confirm" && (
        <View className="mt-4">
          <View className="rounded-xl bg-brand-50 p-4">
            <Text className="text-sm text-brand-700">✓ Phone verified</Text>
          </View>
          <View className="mt-4 rounded-xl border border-gray-100 p-4">
            <Row label="Tour" value={tour.title} />
            <Row label="Traveler" value={travelerName} />
            <Row label="Phone" value={phone} />
            {departureId ? (
              <Row
                label="Departure"
                value={
                  tour.departures.find((d) => d.id === departureId)
                    ? shortDate(tour.departures.find((d) => d.id === departureId)!.start_date)
                    : "—"
                }
              />
            ) : null}
          </View>
          {error ? <Text className="my-3 text-sm text-red-600">{error}</Text> : null}
          <View className="mt-4">
            <Button title="Confirm booking" onPress={confirmBooking} loading={busy} />
          </View>
        </View>
      )}
    </ScrollView>
  );
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <View className="flex-row justify-between py-1">
      <Text className="text-sm text-gray-400">{label}</Text>
      <Text className="text-sm font-medium text-gray-800">{value}</Text>
    </View>
  );
}

function StepIndicator({ step }: { step: Step }) {
  const steps: Step[] = ["details", "otp", "confirm"];
  const idx = steps.indexOf(step);
  return (
    <View className="mt-3 flex-row gap-2">
      {steps.map((s, i) => (
        <View
          key={s}
          className={`h-1.5 flex-1 rounded-full ${i <= idx ? "bg-brand-500" : "bg-gray-200"}`}
        />
      ))}
    </View>
  );
}
