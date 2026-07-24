import { useForm } from "react-hook-form";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiList, apiPost, apiDelete } from "@/api/client";
import { PageHeader, StatusBadge } from "@/components/ui/primitives";
import { DataTable, type Column } from "@/components/ui/DataTable";
import { toast } from "@/components/ui/toast";
import type { Banner } from "@/types/api";

interface BannerForm {
  image: string;
  title: string;
  link: string;
  sort_order: number;
  active: boolean;
}

export function BannersPage() {
  const qc = useQueryClient();
  const { data, isLoading } = useQuery({
    queryKey: ["banners"],
    queryFn: () => apiList<Banner>("/admin/banners"),
  });
  const { register, handleSubmit, reset } = useForm<BannerForm>({
    defaultValues: { active: true, sort_order: 0 },
  });

  const create = useMutation({
    mutationFn: (v: BannerForm) =>
      apiPost("/admin/banners", { ...v, sort_order: Number(v.sort_order), type: "manual" }),
    onSuccess: () => {
      toast.success("Banner created");
      qc.invalidateQueries({ queryKey: ["banners"] });
      reset();
    },
    onError: () => toast.error("Create failed"),
  });

  const remove = useMutation({
    mutationFn: (id: number) => apiDelete(`/admin/banners/${id}`),
    onSuccess: () => {
      toast.success("Banner removed");
      qc.invalidateQueries({ queryKey: ["banners"] });
    },
  });

  const columns: Column<Banner>[] = [
    {
      header: "Image",
      cell: (b) =>
        b.image ? (
          <img src={b.image} alt="" className="h-10 w-16 rounded object-cover" />
        ) : (
          "—"
        ),
    },
    { header: "Title", cell: (b) => b.title },
    { header: "Link", cell: (b) => <span className="text-xs text-gray-500">{b.link}</span> },
    { header: "Order", cell: (b) => b.sort_order },
    { header: "Status", cell: (b) => <StatusBadge status={b.active ? "active" : "suspended"} /> },
    {
      header: "",
      className: "text-end",
      cell: (b) => (
        <button className="btn-danger px-3 py-1 text-xs" onClick={() => remove.mutate(b.id)}>
          Delete
        </button>
      ),
    },
  ];

  return (
    <div>
      <PageHeader title="Banners" subtitle="Home carousel (manual banners)" />
      <form
        onSubmit={handleSubmit((v) => create.mutate(v))}
        className="card mb-6 grid grid-cols-2 gap-4 p-4 md:grid-cols-4"
      >
        <div className="md:col-span-2">
          <label className="label">Image URL</label>
          <input className="input" {...register("image", { required: true })} />
        </div>
        <div>
          <label className="label">Title</label>
          <input className="input" {...register("title")} />
        </div>
        <div>
          <label className="label">Link</label>
          <input className="input" {...register("link")} />
        </div>
        <div>
          <label className="label">Sort order</label>
          <input type="number" className="input" {...register("sort_order")} />
        </div>
        <label className="flex items-center gap-2 pt-6 text-sm">
          <input type="checkbox" {...register("active")} /> Active
        </label>
        <div className="col-span-full">
          <button type="submit" className="btn-primary">Add banner</button>
        </div>
      </form>
      <DataTable columns={columns} rows={data?.items ?? []} loading={isLoading} rowKey={(b) => b.id} />
    </div>
  );
}
