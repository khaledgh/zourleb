import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiList, apiPost } from "@/api/client";
import { PageHeader, StatusBadge } from "@/components/ui/primitives";
import { DataTable, Pager, type Column } from "@/components/ui/DataTable";
import { toast } from "@/components/ui/toast";

interface AdminUser {
  id: number;
  name: string;
  email: string;
  status: string;
  locale: string;
}

export function UsersPage() {
  const qc = useQueryClient();
  const [q, setQ] = useState("");
  const [page, setPage] = useState(1);

  const { data, isLoading } = useQuery({
    queryKey: ["users", q, page],
    queryFn: () => apiList<AdminUser>("/admin/users", { q: q || undefined, page, per_page: 20 }),
  });

  const setStatus = useMutation({
    mutationFn: ({ id, status }: { id: number; status: string }) =>
      apiPost(`/admin/users/${id}/status`, { status }),
    onSuccess: () => {
      toast.success("User updated");
      qc.invalidateQueries({ queryKey: ["users"] });
    },
    onError: () => toast.error("Action failed"),
  });

  const columns: Column<AdminUser>[] = [
    { header: "Name", cell: (u) => <span className="font-medium">{u.name}</span> },
    { header: "Email", cell: (u) => u.email },
    { header: "Locale", cell: (u) => u.locale },
    { header: "Status", cell: (u) => <StatusBadge status={u.status} /> },
    {
      header: "",
      className: "text-end",
      cell: (u) => (
        <button
          className={u.status === "blocked" ? "btn-primary px-3 py-1 text-xs" : "btn-danger px-3 py-1 text-xs"}
          onClick={() =>
            setStatus.mutate({
              id: u.id,
              status: u.status === "blocked" ? "active" : "blocked",
            })
          }
        >
          {u.status === "blocked" ? "Unblock" : "Block"}
        </button>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title="Users"
        actions={
          <input
            className="input w-64"
            placeholder="Search name or email…"
            value={q}
            onChange={(e) => {
              setQ(e.target.value);
              setPage(1);
            }}
          />
        }
      />
      <DataTable columns={columns} rows={data?.items ?? []} loading={isLoading} rowKey={(u) => u.id} />
      <Pager page={page} totalPages={data?.meta?.total_pages ?? 1} onPage={setPage} />
    </div>
  );
}
