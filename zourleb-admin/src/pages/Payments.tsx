import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiList, apiPost } from "@/api/client";
import { PageHeader, StatusBadge } from "@/components/ui/primitives";
import { DataTable, Pager, type Column } from "@/components/ui/DataTable";
import { toast } from "@/components/ui/toast";
import type { Payment } from "@/types/api";

export function PaymentsPage() {
  const qc = useQueryClient();
  const [page, setPage] = useState(1);

  const { data, isLoading } = useQuery({
    queryKey: ["payments", page],
    queryFn: () => apiList<Payment>("/admin/payments", { page, per_page: 20 }),
  });

  const confirm = useMutation({
    mutationFn: (id: number) => apiPost("/admin/payments/confirm", { payment_id: id }),
    onSuccess: () => {
      toast.success("Payment confirmed");
      qc.invalidateQueries({ queryKey: ["payments"] });
    },
    onError: () => toast.error("Confirm failed"),
  });

  const columns: Column<Payment>[] = [
    { header: "ID", cell: (p) => p.id },
    { header: "Payable", cell: (p) => `${p.payable_type} #${p.payable_id}` },
    { header: "Payer", cell: (p) => `${p.payer_type} #${p.payer_id}` },
    { header: "Amount", cell: (p) => `${p.amount.toFixed(2)} ${p.currency}` },
    { header: "Provider", cell: (p) => p.provider },
    { header: "Status", cell: (p) => <StatusBadge status={p.status} /> },
    {
      header: "",
      className: "text-end",
      cell: (p) =>
        p.status !== "paid" ? (
          <button
            className="btn-primary px-3 py-1 text-xs"
            onClick={() => confirm.mutate(p.id)}
            disabled={confirm.isPending}
          >
            Confirm
          </button>
        ) : null,
    },
  ];

  return (
    <div>
      <PageHeader title="Payments" subtitle="Confirm offline / manual payments" />
      <DataTable
        columns={columns}
        rows={data?.items ?? []}
        loading={isLoading}
        rowKey={(p) => p.id}
        empty="No payments found."
      />
      <Pager
        page={page}
        totalPages={data?.meta?.total_pages ?? 1}
        onPage={setPage}
      />
    </div>
  );
}
