import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiList, apiPost } from "@/api/client";
import { PageHeader, StatusBadge } from "@/components/ui/primitives";
import { DataTable, Pager, type Column } from "@/components/ui/DataTable";
import { toast } from "@/components/ui/toast";
import type { Agency } from "@/types/api";

const STATUSES = ["", "pending", "approved", "suspended"] as const;

export function AgenciesPage() {
  const qc = useQueryClient();
  const [status, setStatus] = useState<string>("pending");
  const [page, setPage] = useState(1);

  const { data, isLoading } = useQuery({
    queryKey: ["agencies", status, page],
    queryFn: () =>
      apiList<Agency>("/admin/agencies", { status: status || undefined, page, per_page: 20 }),
  });

  const decide = useMutation({
    mutationFn: ({ id, approve }: { id: number; approve: boolean }) =>
      apiPost(`/admin/agencies/${id}/approve`, { approve }),
    onSuccess: (_d, v) => {
      toast.success(v.approve ? "Agency approved" : "Agency suspended");
      qc.invalidateQueries({ queryKey: ["agencies"] });
    },
    onError: () => toast.error("Action failed"),
  });

  const columns: Column<Agency>[] = [
    { header: "Name", cell: (a) => <span className="font-medium">{a.name}</span> },
    { header: "Email", cell: (a) => a.email },
    { header: "Phone", cell: (a) => a.phone },
    { header: "Status", cell: (a) => <StatusBadge status={a.status} /> },
    { header: "Rating", cell: (a) => `${a.rating_avg.toFixed(1)} (${a.rating_count})` },
    {
      header: "",
      className: "text-end",
      cell: (a) => (
        <div className="flex justify-end gap-2">
          {a.status !== "approved" && (
            <button
              className="btn-primary px-3 py-1 text-xs"
              onClick={() => decide.mutate({ id: a.id, approve: true })}
            >
              Approve
            </button>
          )}
          {a.status !== "suspended" && (
            <button
              className="btn-danger px-3 py-1 text-xs"
              onClick={() => decide.mutate({ id: a.id, approve: false })}
            >
              Suspend
            </button>
          )}
        </div>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title="Agencies"
        subtitle="Review and moderate travel agencies"
        actions={
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
        }
      />
      <DataTable
        columns={columns}
        rows={data?.items ?? []}
        loading={isLoading}
        rowKey={(a) => a.id}
        empty="No agencies match this filter."
      />
      <Pager
        page={page}
        totalPages={data?.meta?.total_pages ?? 1}
        onPage={setPage}
      />
    </div>
  );
}
