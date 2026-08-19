import { useState } from "react";

import { Alert, Image, ScrollView, Text, TextInput, View, Pressable } from "react-native";

import { useLocalSearchParams, useRouter, Stack } from "expo-router";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";

import { apiData, apiList, apiPost, apiDelete, ApiException } from "@/api/client";

import { useAuth } from "@/stores/auth";

import { Badge, Button, Divider, Loading, SectionHeader, StarRating } from "@/components/ui";

import { money, shortDate } from "@/lib/format";

import type { Review, TourDetail } from "@/types/api";



export default function TourDetailScreen() {

  const { slug } = useLocalSearchParams<{ slug: string }>();

  const router = useRouter();

  const user = useAuth((s) => s.user);

  const queryClient = useQueryClient();

  // Favorite + review form state
  const [faved, setFaved] = useState(false);
  const [rating, setRating] = useState(0);

  const [comment, setComment] = useState("");
  const [reviewError, setReviewError] = useState<string | null>(null);
  const [showReviewForm, setShowReviewForm] = useState(false);



  const { data: tour, isLoading } = useQuery({

    queryKey: ["tour", slug],

    queryFn: () => apiData<TourDetail>(`/tours/${slug}`),

    enabled: !!slug,

  });



  const { data: reviewsData } = useQuery({

    queryKey: ["tour-reviews", tour?.id],

    queryFn: () => apiList<Review>("/reviews", { tour_id: tour!.id, per_page: 20 }),

    enabled: !!tour?.id,

  });



  const favorite = useMutation({

    mutationFn: () =>

      faved

        ? apiDelete(`/favorites/${tour!.id}`)

        : apiPost(`/favorites/${tour!.id}`, {}),

    onSuccess: () => {

      setFaved((v) => !v);

      queryClient.invalidateQueries({ queryKey: ["favorite-ids"] });

    },

  });



  const submitReview = useMutation({

    mutationFn: () =>

      apiPost("/reviews", { tour_id: tour!.id, rating, comment }),

    onSuccess: () => {

      setShowReviewForm(false);

      setRating(0);

      setComment("");

      queryClient.invalidateQueries({ queryKey: ["tour-reviews", tour?.id] });

      Alert.alert("Thanks!", "Your review has been submitted.");

    },

    onError: (e) => {

      setReviewError(e instanceof ApiException ? e.message : "Could not submit review.");

    },

  });



  if (isLoading || !tour) return <Loading />;



  const reviews = reviewsData?.items ?? [];



  return (

    <View className="flex-1 bg-white">

      <Stack.Screen options={{ title: tour.title }} />

      <ScrollView contentContainerClassName="pb-28">

        {/* Cover image */}

        <View className="aspect-[16/10] w-full bg-gray-100">

          {tour.cover ? (

            <Image source={{ uri: tour.cover }} className="h-full w-full" resizeMode="cover" />

          ) : null}

        </View>



        <View className="p-4">

          {/* Title + favorite */}

          <View className="flex-row items-start justify-between">

            <Text className="flex-1 text-xl font-bold text-gray-900">{tour.title}</Text>

            {user ? (

              <Pressable

                onPress={() => favorite.mutate()}

                className="ms-2 h-9 w-9 items-center justify-center rounded-full border border-gray-100"

              >

                <Text className={`text-xl ${faved ? "text-red-500" : "text-gray-400"}`}>

                  {faved ? "♥" : "♡"}

                </Text>

              </Pressable>

            ) : null}

          </View>



          {/* Badges */}

          <View className="mt-2 flex-row flex-wrap gap-2">

            <Badge label={tour.type.replace(/_/g, " ")} />

            <Badge label={`${tour.duration_days} day${tour.duration_days !== 1 ? "s" : ""}`} />

            {tour.difficulty ? <Badge label={tour.difficulty} /> : null}

          </View>



          {/* Rating */}

          {tour.rating_avg > 0 ? (

            <View className="mt-2">

              <StarRating rating={tour.rating_avg} showNumber size={16} />

            </View>

          ) : null}



          <Text className="mt-3 text-sm leading-5 text-gray-600">{tour.summary}</Text>



          {/* About */}

          {tour.description ? (

            <>

              <Divider />

              <SectionHeader title="About" />

              <Text className="text-sm leading-5 text-gray-600">{tour.description}</Text>

            </>

          ) : null}



          {/* What's included / excluded */}

          {(tour.included || tour.excluded) ? (

            <>

              <Divider />

              {tour.included ? (

                <>

                  <Text className="mb-1 text-sm font-semibold text-gray-900">✅ Included</Text>

                  <Text className="text-sm text-gray-600">{tour.included}</Text>

                </>

              ) : null}

              {tour.excluded ? (

                <View className="mt-2">

                  <Text className="mb-1 text-sm font-semibold text-gray-900">❌ Excluded</Text>

                  <Text className="text-sm text-gray-600">{tour.excluded}</Text>

                </View>

              ) : null}

            </>

          ) : null}



          {/* Departures */}

          {tour.departures.length > 0 ? (

            <>

              <Divider />

              <SectionHeader title="Upcoming departures" />

              {tour.departures.map((d) => (

                <View

                  key={d.id}

                  className="mt-2 flex-row items-center justify-between rounded-xl border border-gray-100 p-3"

                >

                  <Text className="text-sm text-gray-700">

                    {shortDate(d.start_date)} → {shortDate(d.end_date)}

                  </Text>

                  <Badge

                    label={d.seats_left > 0 ? `${d.seats_left} seats` : "Full"}

                    tone={d.seats_left > 0 ? "success" : "danger"}

                  />

                </View>

              ))}

            </>

          ) : null}



          {/* Past gallery */}

          {tour.past_gallery.length > 0 ? (

            <>

              <Divider />

              <SectionHeader title="From past trips" />

              <ScrollView horizontal showsHorizontalScrollIndicator={false} className="mt-2 -mx-1">

                {tour.past_gallery.map((img, i) => (

                  <Image

                    key={i}

                    source={{ uri: img.url }}

                    className="mx-1 h-28 w-36 rounded-xl bg-gray-100"

                    resizeMode="cover"

                  />

                ))}

              </ScrollView>

            </>

          ) : null}



          {/* Reviews section */}

          <Divider />

          <SectionHeader

            title={`Reviews (${reviews.length})`}

            action={user && !showReviewForm ? "Write a review" : undefined}

            onAction={() => setShowReviewForm(true)}

          />



          {/* Write review form */}

          {showReviewForm && user ? (

            <View className="mb-4 rounded-2xl border border-brand-200 bg-brand-50 p-4">

              <Text className="mb-2 text-sm font-semibold text-gray-900">Your rating</Text>

              <View className="mb-3 flex-row gap-1">

                {[1, 2, 3, 4, 5].map((s) => (

                  <Pressable key={s} onPress={() => setRating(s)}>

                    <Text className={`text-3xl ${s <= rating ? "text-amber-400" : "text-gray-300"}`}>

                      ★

                    </Text>

                  </Pressable>

                ))}

              </View>

              <TextInput

                className="mb-3 min-h-[80px] rounded-xl border border-gray-300 bg-white p-3 text-sm"

                placeholder="Share your experience…"

                multiline

                textAlignVertical="top"

                value={comment}

                onChangeText={setComment}

              />

              {reviewError ? (

                <Text className="mb-2 text-sm text-red-600">{reviewError}</Text>

              ) : null}

              <View className="flex-row gap-2">

                <View className="flex-1">

                  <Button

                    title="Submit"

                    onPress={() => submitReview.mutate()}

                    loading={submitReview.isPending}

                    disabled={rating === 0}

                  />

                </View>

                <View className="flex-1">

                  <Button

                    title="Cancel"

                    variant="ghost"

                    onPress={() => setShowReviewForm(false)}

                  />

                </View>

              </View>

            </View>

          ) : null}



          {/* Review list */}

          {reviews.length === 0 ? (

            <Text className="text-sm text-gray-400">No reviews yet. Be the first!</Text>

          ) : (

            reviews.map((r) => (

              <View key={r.id} className="mb-3 rounded-xl border border-gray-100 p-3">

                <View className="flex-row items-center justify-between">

                  <Text className="text-sm font-semibold text-gray-800">{r.user_name}</Text>

                  <StarRating rating={r.rating} size={12} />

                </View>

                {r.comment ? (

                  <Text className="mt-1 text-sm text-gray-600">{r.comment}</Text>

                ) : null}

                <Text className="mt-1 text-xs text-gray-400">{shortDate(r.created_at)}</Text>

              </View>

            ))

          )}

        </View>

      </ScrollView>



      {/* Sticky bottom bar */}

      <View className="absolute bottom-0 w-full flex-row items-center justify-between border-t border-gray-100 bg-white p-4">

        <View>

          <Text className="text-xs text-gray-400">From</Text>

          <Text className="text-lg font-bold text-brand-600">

            {money(tour.price_from, tour.currency)}

          </Text>

        </View>

        <View className="w-44">

          <Button

            title="Book now"

            onPress={() =>

              user

                ? router.push(`/booking/${tour.slug}`)

                : router.push("/(auth)/login")

            }

          />

        </View>

      </View>

    </View>

  );

}

