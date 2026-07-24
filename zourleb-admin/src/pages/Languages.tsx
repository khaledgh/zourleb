import { useState } from "react";
import { useForm } from "react-hook-form";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiList, apiPost } from "@/api/client";
import { PageHeader, StatusBadge } from "@/components/ui/primitives";
import { DataTable, type Column } from "@/components/ui/DataTable";
import { toast } from "@/components/ui/toast";
import type { Language } from "@/types/api";

interface LangForm {
  code: string;
  name: string;
  native_name: string;
  is_rtl: boolean;
  is_active: boolean;
  sort_order: number;
}

export function LanguagesPage() {
  const qc = useQueryClient();
  const [open, setOpen] = useState(false);
  const { data, isLoading } = useQuery({
    queryKey: ["admin-languages"],
    queryFn: () => apiList<Language>("/admin/languages"),
  });

  const { register, handleSubmit, reset } = useForm<LangForm>({
    defaultValues: { is_active: true, is_rtl: false, sort_order: 0 },
  });

  const create = useMutation({
    mutationFn: (v: LangForm) => apiPost("/admin/languages", v),
    onSuccess: () => {
      toast.success("Language added");
      qc.invalidateQueries({ queryKey: ["admin-languages"] });
      setOpen(false);
      reset();
    },
    onError: () => toast.error("Could not add language"),
  });

  const toggleActive = useMutation({
    mutationFn: (l: Language) =>
      apiPost(`/admin/languages/${l.id}`, {
        code: l.code,
        name: l.name,
        native_name: l.native_name,
        is_rtl: l.is_rtl,
        is_active: !l.is_active,
        sort_order: l.sort_order,
      }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["admin-languages"] }),
  });

  const columns: Column<Language>[] = [
    { header: "Code", cell: (l) => <span className="font-mono">{l.code}</span> },
    { header: "Name", cell: (l) => l.name },
    { header: "Native", cell: (l) => l.native_name },
    { header: "RTL", cell: (l) => (l.is_rtl ? "✓" : "—") },
    { header: "Default", cell: (l) => (l.is_default ? "✓" : "—") },
    { header: "Status", cell: (l) => <StatusBadge status={l.is_active ? "active" : "suspended"} /> },
    {
      header: "",
      className: "text-end",
      cell: (l) => (
        <button className="btn-ghost px-3 py-1 text-xs" onClick={() => toggleActive.mutate(l)}>
          {l.is_active ? "Disable" : "Enable"}
        </button>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title="Languages"
        subtitle="Add languages without a redeploy"
        actions={
          <button className="btn-primary" onClick={() => setOpen((o) => !o)}>
            {open ? "Cancel" : "Add language"}
          </button>
        }
      />
      {open && (
        <form
          onSubmit={handleSubmit((v) => create.mutate({ ...v, sort_order: Number(v.sort_order) }))}
          className="card mb-4 grid grid-cols-2 gap-4 p-4 md:grid-cols-3"
        >
          <div>
            <label className="label">Code</label>
            <input className="input" placeholder="de" {...register("code", { required: true })} />
          </div>
          <div>
            <label className="label">Name</label>
            <input className="input" placeholder="German" {...register("name", { required: true })} />
          </div>
          <div>
            <label className="label">Native name</label>
            <input className="input" placeholder="Deutsch" {...register("native_name", { required: true })} />
          </div>
          <div>
            <label className="label">Sort order</label>
            <input type="number" className="input" {...register("sort_order")} />
          </div>
          <label className="flex items-center gap-2 pt-6 text-sm">
            <input type="checkbox" {...register("is_rtl")} /> RTL
          </label>
          <label className="flex items-center gap-2 pt-6 text-sm">
            <input type="checkbox" {...register("is_active")} /> Active
          </label>
          <div className="col-span-full">
            <button type="submit" className="btn-primary">Save language</button>
          </div>
        </form>
      )}
      <DataTable columns={columns} rows={data?.items ?? []} loading={isLoading} rowKey={(l) => l.id} />
    </div>
  );
}
