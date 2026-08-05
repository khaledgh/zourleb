import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { apiList } from "@/api/client";
import { PageHeader, StatusBadge } from "@/components/ui/primitives";
import { DataTable, Pager, type Column } from "@/components/ui/DataTable";

interface AdminTour {
  id: number;
  slug: string;
  status: string;
  type: string;
  price_from: number;
  base_currency: string;
  duration_days: number;
  difficulty: string;
  translations?: Array<{ locale: string; title: string }>;
}

const STATUSES = ["", "draft", "pending", "published", "archived"] as const;

export function AdminToursPage() {
  const [status, setStatus] = useState<string>("");
  const [q, setQ] = useState<string>("");
  const [page, setPage] = useState(1);

  const { data, isLoading } = useQuery({
    queryKey: ["admin-tours", status, q, page],
    queryFn: () =>
      apiList<AdminTour>("/admin/tours", {
        status: status || undefined,
        q: q || undefined,
        page,
        per_page: 20,
      }),
  });

  const columns: Column<AdminTour>[] = [
    {
      header: "Title",
      cell: (t) => <span className="font-medium">{t.translations?.[0]?.title || t.slug}</span>,
    },
    { header: "Type", cell: (t) => t.type.replace(/_/g, " ") },
    { header: "Duration", cell: (t) => `${t.duration_days} days` },
    { header: "Price", cell: (t) => `${t.price_from} ${t.base_currency}` },
    { header: "Difficulty", cell: (t) => t.difficulty || "—" },
    { header: "Status", cell: (t) => <StatusBadge status={t.status} /> },
  ];

  return (
    <div>
      <PageHeader
        title="Trips (All Tours)"
        subtitle="View and moderate all agency trips on the platform"
        actions={
          <div className="flex gap-3">
            <input
              className="input w-64"
              placeholder="Search tours…"
              value={q}
              onChange={(e) => {
                setQ(e.target.value);
                setPage(1);
              }}
            />
            <select
              className="input w-44"
              value={status}
              onChange={(e) => {
                setStatus(e.target.value);
                setPage(1);
              }}
            >
              {STATUSES.map((s) => (
                <option key={s} value={s}>
                  {s === "" ? "All statuses" : s}
                </option>
              ))}
            </select>
          </div>
        }
      />
      <DataTable
        columns={columns}
        rows={data?.items ?? []}
        loading={isLoading}
        rowKey={(t) => t.id}
        empty="No tours found."
      />
      <Pager
        page={page}
        totalPages={data?.meta?.total_pages ?? 1}
        onPage={setPage}
      />
    </div>
  );
}
