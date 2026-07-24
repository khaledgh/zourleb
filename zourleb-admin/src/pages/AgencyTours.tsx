import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiList, apiPost } from "@/api/client";
import { PageHeader, StatusBadge } from "@/components/ui/primitives";
import { DataTable, type Column } from "@/components/ui/DataTable";
import { toast } from "@/components/ui/toast";

interface TourTranslation {
  locale: string;
  title: string;
}
interface AgencyTour {
  id: number;
  slug: string;
  status: string;
  type: string;
  price_from: number;
  base_currency: string;
  featured: boolean;
  translations: TourTranslation[];
}

function title(t: AgencyTour) {
  return t.translations?.[0]?.title || t.slug;
}

// Agency-scoped tour list. The backend resolves the caller's agency from their
// membership (or the X-Agency-ID header for multi-agency users).
export function AgencyToursPage() {
  const qc = useQueryClient();
  const { data, isLoading } = useQuery({
    queryKey: ["agency-tours"],
    queryFn: () => apiList<AgencyTour>("/agency/tours", { per_page: 50 }),
  });

  const publish = useMutation({
    mutationFn: ({ id, publish }: { id: number; publish: boolean }) =>
      apiPost(`/agency/tours/${id}/publish`, { publish }),
    onSuccess: () => {
      toast.success("Tour updated");
      qc.invalidateQueries({ queryKey: ["agency-tours"] });
    },
    onError: () => toast.error("Action failed"),
  });

  const columns: Column<AgencyTour>[] = [
    { header: "Title", cell: (t) => <span className="font-medium">{title(t)}</span> },
    { header: "Type", cell: (t) => t.type.replace(/_/g, " ") },
    { header: "From", cell: (t) => `${t.price_from} ${t.base_currency}` },
    { header: "Status", cell: (t) => <StatusBadge status={t.status} /> },
    {
      header: "",
      className: "text-end",
      cell: (t) => (
        <button
          className="btn-ghost px-3 py-1 text-xs"
          onClick={() => publish.mutate({ id: t.id, publish: t.status !== "published" })}
        >
          {t.status === "published" ? "Unpublish" : "Publish"}
        </button>
      ),
    },
  ];

  return (
    <div>
      <PageHeader title="My tours" subtitle="Manage and publish your tours" />
      <DataTable
        columns={columns}
        rows={data?.items ?? []}
        loading={isLoading}
        rowKey={(t) => t.id}
        empty="You haven't created any tours yet."
      />
    </div>
  );
}
