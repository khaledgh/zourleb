import type { ReactNode } from "react";
import { Spinner, EmptyState } from "./primitives";

export interface Column<T> {
  header: string;
  cell: (row: T) => ReactNode;
  className?: string;
}

interface DataTableProps<T> {
  columns: Column<T>[];
  rows: T[];
  loading?: boolean;
  rowKey: (row: T) => string | number;
  empty?: string;
}

export function DataTable<T>({
  columns,
  rows,
  loading,
  rowKey,
  empty = "No records found.",
}: DataTableProps<T>) {
  if (loading) return <Spinner />;
  if (rows.length === 0) return <EmptyState message={empty} />;

  return (
    <div className="card overflow-hidden border border-slate-100/80 p-0 shadow-[0_8px_30px_rgb(0,0,0,0.015)]">
      <div className="overflow-x-auto">
        <table className="w-full text-sm border-collapse text-left">
          <thead className="border-b border-slate-100/80 text-[10px] font-bold text-slate-400 uppercase tracking-wider bg-slate-50/30">
            <tr>
              {columns.map((c, i) => (
                <th key={i} className="px-6 py-3.5 text-left font-bold">
                  {c.header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-50">
            {rows.map((row) => (
              <tr key={rowKey(row)} className="hover:bg-slate-50/40 transition-colors duration-150">
                {columns.map((c, i) => (
                  <td key={i} className={`px-6 py-4 text-slate-600 font-semibold text-xs leading-normal ${c.className ?? ""}`}>
                    {c.cell(row)}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export function Pager({
  page,
  totalPages,
  onPage,
}: {
  page: number;
  totalPages: number;
  onPage: (p: number) => void;
}) {
  if (totalPages <= 1) return null;
  return (
    <div className="mt-6 flex items-center justify-end gap-3 text-xs font-bold text-slate-400">
      <button
        className="btn bg-white hover:bg-slate-50 text-slate-600 border border-slate-200/60 rounded-xl p-2 h-9 w-9 disabled:opacity-30 disabled:pointer-events-none transition-all"
        disabled={page <= 1}
        onClick={() => onPage(page - 1)}
      >
        <svg className="w-4 h-4 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={3}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M15 19l-7-7 7-7" />
        </svg>
      </button>
      <span className="bg-slate-100/80 text-slate-600 px-3.5 py-1.5 rounded-xl text-center min-w-[3rem]">
        {page} / {totalPages}
      </span>
      <button
        className="btn bg-white hover:bg-slate-50 text-slate-600 border border-slate-200/60 rounded-xl p-2 h-9 w-9 disabled:opacity-30 disabled:pointer-events-none transition-all"
        disabled={page >= totalPages}
        onClick={() => onPage(page + 1)}
      >
        <svg className="w-4 h-4 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={3}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M9 5l7 7-7 7" />
        </svg>
      </button>
    </div>
  );
}
