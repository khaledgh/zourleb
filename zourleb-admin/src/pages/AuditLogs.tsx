import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { apiList } from "@/api/client";
import { PageHeader } from "@/components/ui/primitives";
import { DataTable, Pager, type Column } from "@/components/ui/DataTable";
import type { AuditLog } from "@/types/api";

export function AuditLogsPage() {
  const [page, setPage] = useState(1);

  const { data, isLoading } = useQuery({
    queryKey: ["audit-logs", page],
    queryFn: () => apiList<AuditLog>("/admin/audit-logs", { page, per_page: 20 }),
  });

  const columns: Column<AuditLog>[] = [
    { header: "ID", cell: (a) => a.id },
    { header: "Action", cell: (a) => a.action },
    { header: "Target", cell: (a) => `${a.target_type} #${a.target_id}` },
    { header: "Actor", cell: (a) => a.actor_id ?? "system" },
    { header: "IP", cell: (a) => a.ip || "-" },
    {
      header: "Payload",
      cell: (a) =>
        a.payload ? (
          <pre className="text-[10px] text-slate-500 max-w-xs truncate">
            {JSON.stringify(a.payload)}
          </pre>
        ) : (
          "-"
        ),
    },
    { header: "Created", cell: (a) => new Date(a.created_at).toLocaleString() },
  ];

  return (
    <div>
      <PageHeader title="Audit Logs" subtitle="Sensitive action history" />
      <DataTable
        columns={columns}
        rows={data?.items ?? []}
        loading={isLoading}
        rowKey={(a) => a.id}
        empty="No audit logs found."
      />
      <Pager
        page={page}
        totalPages={data?.meta?.total_pages ?? 1}
        onPage={setPage}
      />
    </div>
  );
}
