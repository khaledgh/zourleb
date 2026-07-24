import clsx from "clsx";
import type { ReactNode } from "react";

export function PageHeader({
  title,
  subtitle,
  actions,
}: {
  title: string;
  subtitle?: string;
  actions?: ReactNode;
}) {
  return (
    <div className="mb-8 flex flex-wrap items-center justify-between gap-4 border-b border-slate-100 pb-5">
      <div>
        <h1 className="font-display text-2xl font-extrabold tracking-tight text-slate-800">{title}</h1>
        {subtitle && <p className="mt-1.5 text-xs font-semibold text-slate-400">{subtitle}</p>}
      </div>
      {actions && <div className="flex items-center gap-3">{actions}</div>}
    </div>
  );
}

const statusStyles: Record<string, string> = {
  approved: "bg-emerald-50 text-emerald-600 border border-emerald-100/50",
  active: "bg-emerald-50 text-emerald-600 border border-emerald-100/50",
  pending: "bg-amber-50 text-amber-600 border border-amber-100/50",
  pending_payment: "bg-amber-50 text-amber-600 border border-amber-100/50",
  suspended: "bg-rose-50 text-rose-600 border border-rose-100/50",
  rejected: "bg-rose-50 text-rose-600 border border-rose-100/50",
  expired: "bg-slate-50 text-slate-500 border border-slate-200/50",
  blocked: "bg-rose-50 text-rose-600 border border-rose-100/50",
};

export function StatusBadge({ status }: { status: string }) {
  return (
    <span
      className={clsx(
        "inline-flex items-center rounded-xl px-2.5 py-1 text-[10px] font-bold uppercase tracking-wider",
        statusStyles[status] ?? "bg-slate-50 text-slate-500 border border-slate-200/50"
      )}
    >
      {status.replace(/_/g, " ")}
    </span>
  );
}

export function Spinner({ label }: { label?: string }) {
  return (
    <div className="flex flex-col items-center justify-center gap-3 py-16 text-xs font-semibold text-slate-400">
      <span className="h-6 w-6 animate-spin rounded-full border-2 border-slate-200 border-t-[#2F80ED]" />
      {label && <span>{label}</span>}
    </div>
  );
}

export function EmptyState({ message }: { message: string }) {
  return (
    <div className="card flex flex-col items-center justify-center py-16 text-center border border-dashed border-slate-200 bg-slate-50/20 text-xs font-semibold text-slate-400">
      <svg className="w-10 h-10 text-slate-300 mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
        <path strokeLinecap="round" strokeLinejoin="round" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0a2 2 0 01-2 2H6a2 2 0 01-2-2m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5" />
      </svg>
      {message}
    </div>
  );
}

export function StatCard({ label, value }: { label: string; value: ReactNode }) {
  return (
    <div className="card border border-slate-100 p-6 flex flex-col justify-between hover:scale-[1.02] transition-all">
      <div className="text-[10px] font-bold uppercase tracking-wider text-slate-400">{label}</div>
      <div className="mt-4 font-display text-2xl font-extrabold text-slate-800 leading-none">{value}</div>
    </div>
  );
}
