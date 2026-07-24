import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiList, apiPost } from "@/api/client";
import { PageHeader, StatusBadge } from "@/components/ui/primitives";
import { DataTable, type Column } from "@/components/ui/DataTable";
import { toast } from "@/components/ui/toast";
import type { Boost } from "@/types/api";

// Pending boosts await offline/manual payment confirmation. Confirming the
// payment activates the boost server-side.
export function BoostsPage() {
  const qc = useQueryClient();
  const { data, isLoading } = useQuery({
    queryKey: ["pending-boosts"],
    queryFn: () => apiList<Boost>("/admin/boosts/pending"),
  });

  const confirm = useMutation({
    mutationFn: (paymentId: number) =>
      apiPost("/admin/payments/confirm", { payment_id: paymentId }),
    onSuccess: () => {
      toast.success("Payment confirmed — boost activated");
      qc.invalidateQueries({ queryKey: ["pending-boosts"] });
    },
    onError: () => toast.error("Confirmation failed"),
  });

  const columns: Column<Boost>[] = [
    { header: "Agency", cell: (b) => `#${b.agency_id}` },
    { header: "Tour", cell: (b) => (b.tour_id ? `#${b.tour_id}` : "—") },
    { header: "Placement", cell: (b) => b.placement.replace(/_/g, " ") },
    { header: "Status", cell: (b) => <StatusBadge status={b.status} /> },
    {
      header: "",
      className: "text-end",
      cell: (b) => (
        <button
          className="btn-primary px-3 py-1 text-xs"
          // The backend stores the payment id on the boost once a charge is
          // created; here we confirm by boost-derived payment. For pending
          // boosts an admin confirms the linked payment id (entered inline).
          onClick={() => {
            const pid = window.prompt(`Payment ID for boost #${b.id}?`);
            if (pid) confirm.mutate(Number(pid));
          }}
        >
          Confirm payment
        </button>
      ),
    },
  ];

  return (
    <div>
      <PageHeader title="Boosts" subtitle="Pending payment confirmation" />
      <DataTable
        columns={columns}
        rows={data?.items ?? []}
        loading={isLoading}
        rowKey={(b) => b.id}
        empty="No boosts awaiting confirmation."
      />
    </div>
  );
}
