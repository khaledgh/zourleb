import {
  ActivityIndicator,
  Pressable,
  Text,
  View,
  type PressableProps,
} from "react-native";
import type { ReactNode } from "react";

// ── Button ────────────────────────────────────────────────────────────────────
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
  variant?: "primary" | "outline" | "ghost" | "danger";
  disabled?: boolean;
}) {
  const base =
    "h-12 flex-row items-center justify-center rounded-xl px-5 active:opacity-80";
  const styles = {
    primary: "bg-brand-500",
    outline: "bg-white border-2 border-brand-500",
    ghost: "bg-white border border-gray-200",
    danger: "bg-red-600",
  }[variant];
  const textStyle = {
    primary: "text-white",
    outline: "text-brand-600",
    ghost: "text-gray-800",
    danger: "text-white",
  }[variant];
  return (
    <Pressable
      onPress={onPress}
      disabled={disabled || loading}
      className={`${base} ${styles} ${disabled || loading ? "opacity-50" : ""}`}
    >
      {loading ? (
        <ActivityIndicator color={variant === "ghost" || variant === "outline" ? "#176b47" : "#fff"} />
      ) : (
        <Text className={`text-base font-semibold ${textStyle}`}>{title}</Text>
      )}
    </Pressable>
  );
}

// ── Card ──────────────────────────────────────────────────────────────────────
export function Card({ children, className = "" }: { children: ReactNode; className?: string }) {
  return (
    <View className={`rounded-2xl border border-gray-100 bg-white ${className}`}>
      {children}
    </View>
  );
}

// ── Badge ─────────────────────────────────────────────────────────────────────
export function Badge({ label, tone = "neutral" }: { label: string; tone?: "neutral" | "success" | "warn" | "danger" }) {
  const bg = {
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
    <View className={`self-start rounded-full px-2.5 py-0.5 ${bg}`}>
      <Text className={`text-xs font-medium ${text}`}>{label}</Text>
    </View>
  );
}

// ── Loading ───────────────────────────────────────────────────────────────────
export function Loading({ label }: { label?: string }) {
  return (
    <View className="flex-1 items-center justify-center gap-2 py-12">
      <ActivityIndicator color="#1f8a5b" />
      {label ? <Text className="text-sm text-gray-500">{label}</Text> : null}
    </View>
  );
}

// ── EmptyState ────────────────────────────────────────────────────────────────
export function EmptyState({ icon, message }: { icon?: string; message: string }) {
  return (
    <View className="items-center justify-center py-16">
      {icon ? <Text className="mb-2 text-4xl">{icon}</Text> : null}
      <Text className="text-center text-sm text-gray-400">{message}</Text>
    </View>
  );
}

// ── Field ─────────────────────────────────────────────────────────────────────
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

// ── SectionHeader ─────────────────────────────────────────────────────────────
export function SectionHeader({ title, action, onAction }: { title: string; action?: string; onAction?: () => void }) {
  return (
    <View className="mb-3 flex-row items-center justify-between">
      <Text className="text-lg font-bold text-gray-900">{title}</Text>
      {action && onAction ? (
        <Pressable onPress={onAction}>
          <Text className="text-sm font-medium text-brand-600">{action}</Text>
        </Pressable>
      ) : null}
    </View>
  );
}

// ── Divider ───────────────────────────────────────────────────────────────────
export function Divider() {
  return <View className="my-3 h-px bg-gray-100" />;
}

// ── StarRating ────────────────────────────────────────────────────────────────
export function StarRating({
  rating,
  size = 14,
  showNumber = false,
}: {
  rating: number;
  size?: number;
  showNumber?: boolean;
}) {
  const full = Math.floor(rating);
  const stars = Array.from({ length: 5 }, (_, i) => (i < full ? "★" : "☆"));
  return (
    <View className="flex-row items-center gap-0.5">
      {stars.map((s, i) => (
        <Text key={i} style={{ fontSize: size, lineHeight: size + 4 }} className={i < full ? "text-amber-400" : "text-gray-300"}>
          {s}
        </Text>
      ))}
      {showNumber ? (
        <Text className="ms-1 text-xs text-gray-500">{rating.toFixed(1)}</Text>
      ) : null}
    </View>
  );
}

// ── RowItem — labelled row for detail screens ─────────────────────────────────
export function RowItem({ label, value }: { label: string; value: string }) {
  return (
    <View className="flex-row items-start justify-between py-2">
      <Text className="text-sm text-gray-400">{label}</Text>
      <Text className="ms-4 flex-1 text-right text-sm font-medium text-gray-800">{value}</Text>
    </View>
  );
}
