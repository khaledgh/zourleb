import { create } from "zustand";
import { useEffect } from "react";
import clsx from "clsx";

type ToastKind = "success" | "error" | "info";
interface Toast {
  id: number;
  kind: ToastKind;
  message: string;
}

interface ToastState {
  toasts: Toast[];
  push: (kind: ToastKind, message: string) => void;
  dismiss: (id: number) => void;
}

let seq = 0;

export const useToasts = create<ToastState>((set) => ({
  toasts: [],
  push: (kind, message) =>
    set((s) => ({ toasts: [...s.toasts, { id: ++seq, kind, message }] })),
  dismiss: (id) => set((s) => ({ toasts: s.toasts.filter((t) => t.id !== id) })),
}));

// Convenience helpers usable outside React components.
export const toast = {
  success: (m: string) => useToasts.getState().push("success", m),
  error: (m: string) => useToasts.getState().push("error", m),
  info: (m: string) => useToasts.getState().push("info", m),
};

export function ToastViewport() {
  const { toasts, dismiss } = useToasts();
  return (
    <div className="fixed bottom-4 end-4 z-50 flex w-80 flex-col gap-2">
      {toasts.map((t) => (
        <ToastItem key={t.id} toast={t} onDone={() => dismiss(t.id)} />
      ))}
    </div>
  );
}

function ToastItem({ toast, onDone }: { toast: Toast; onDone: () => void }) {
  useEffect(() => {
    const h = setTimeout(onDone, 4000);
    return () => clearTimeout(h);
  }, [onDone]);
  return (
    <div
      className={clsx(
        "card flex items-start gap-2 p-3 text-sm shadow-lg",
        toast.kind === "success" && "border-l-4 border-l-brand-500",
        toast.kind === "error" && "border-l-4 border-l-red-500",
        toast.kind === "info" && "border-l-4 border-l-blue-500",
      )}
    >
      <span className="flex-1">{toast.message}</span>
      <button onClick={onDone} className="text-gray-400 hover:text-gray-600">
        ✕
      </button>
    </div>
  );
}
