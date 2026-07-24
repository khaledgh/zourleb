import {
  ActivityIndicator,
  Pressable,
  Text,
  View,
  type PressableProps,
} from "react-native";
import type { ReactNode } from "react";

export function Button({
  title,
  onPress,
  loading,
  variant = "primary",
  disabled,
}: {
  title: string;
  onPress?: PressableProps["onPress"];
  loading?: boolean;
  variant?: "primary" | "ghost" | "danger";
  disabled?: boolean;
}) {
  const base =
    "h-12 flex-row items-center justify-center rounded-xl px-5 active:opacity-80";
  const styles = {
    primary: "bg-brand-500",
    ghost: "bg-white border border-gray-200",
    danger: "bg-red-600",
  }[variant];
  const textStyle = variant === "ghost" ? "text-gray-800" : "text-white";
  return (
    <Pressable
      onPress={onPress}
      disabled={disabled || loading}
      className={`${base} ${styles} ${disabled || loading ? "opacity-50" : ""}`}
    >
      {loading ? (
        <ActivityIndicator color={variant === "ghost" ? "#176b47" : "#fff"} />
      ) : (
        <Text className={`text-base font-semibold ${textStyle}`}>{title}</Text>
      )}
    </Pressable>
  );
}

export function Card({ children, className = "" }: { children: ReactNode; className?: string }) {
  return (
    <View className={`rounded-2xl border border-gray-100 bg-white ${className}`}>
      {children}
    </View>
  );
}

export function Badge({ label, tone = "neutral" }: { label: string; tone?: "neutral" | "success" | "warn" | "danger" }) {
  const styles = {
    neutral: "bg-gray-100",
    success: "bg-green-100",
    warn: "bg-amber-100",
    danger: "bg-red-100",
  }[tone];
  const text = {
    neutral: "text-gray-600",
    success: "text-green-700",
    warn: "text-amber-700",
    danger: "text-red-700",
  }[tone];
  return (
    <View className={`self-start rounded-full px-2.5 py-0.5 ${styles}`}>
      <Text className={`text-xs font-medium ${text}`}>{label}</Text>
    </View>
  );
}

export function Loading({ label }: { label?: string }) {
  return (
    <View className="flex-1 items-center justify-center gap-2 py-12">
      <ActivityIndicator color="#1f8a5b" />
      {label ? <Text className="text-sm text-gray-500">{label}</Text> : null}
    </View>
  );
}

export function EmptyState({ message }: { message: string }) {
  return (
    <View className="items-center justify-center py-16">
      <Text className="text-center text-sm text-gray-400">{message}</Text>
    </View>
  );
}

export function Field({
  label,
  children,
  error,
}: {
  label: string;
  children: ReactNode;
  error?: string;
}) {
  return (
    <View className="mb-4">
      <Text className="mb-1 text-sm font-medium text-gray-700">{label}</Text>
      {children}
      {error ? <Text className="mt-1 text-xs text-red-600">{error}</Text> : null}
    </View>
  );
}
