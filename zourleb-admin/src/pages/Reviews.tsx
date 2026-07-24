import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiList, apiPost } from "@/api/client";
import { PageHeader } from "@/components/ui/primitives";
import { DataTable, type Column } from "@/components/ui/DataTable";
import { toast } from "@/components/ui/toast";
import type { Review } from "@/types/api";

export function ReviewsPage() {
  const qc = useQueryClient();
  const { data, isLoading } = useQuery({
    queryKey: ["pending-reviews"],
    queryFn: () => apiList<Review>("/admin/reviews/pending", { per_page: 50 }),
  });

  const moderate = useMutation({
    mutationFn: ({ id, status }: { id: number; status: string }) =>
      apiPost(`/admin/reviews/${id}/moderate`, { status }),
    onSuccess: () => {
      toast.success("Review moderated");
      qc.invalidateQueries({ queryKey: ["pending-reviews"] });
    },
    onError: () => toast.error("Action failed"),
  });

  const columns: Column<Review>[] = [
    { header: "Tour", cell: (r) => `#${r.tour_id}` },
    { header: "Rating", cell: (r) => "★".repeat(r.rating) + "☆".repeat(5 - r.rating) },
    { header: "Comment", cell: (r) => <span className="text-gray-600">{r.comment}</span> },
    {
      header: "",
      className: "text-end",
      cell: (r) => (
        <div className="flex justify-end gap-2">
          <button
            className="btn-primary px-3 py-1 text-xs"
            onClick={() => moderate.mutate({ id: r.id, status: "approved" })}
          >
            Approve
          </button>
          <button
            className="btn-danger px-3 py-1 text-xs"
            onClick={() => moderate.mutate({ id: r.id, status: "rejected" })}
          >
            Reject
          </button>
        </div>
      ),
    },
  ];

  return (
    <div>
      <PageHeader title="Reviews" subtitle="Pending moderation" />
      <DataTable
        columns={columns}
        rows={data?.items ?? []}
        loading={isLoading}
        rowKey={(r) => r.id}
        empty="No reviews awaiting moderation."
      />
    </div>
  );
}
