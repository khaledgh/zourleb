import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiList, apiPost } from "@/api/client";
import { PageHeader, Spinner } from "@/components/ui/primitives";
import { toast } from "@/components/ui/toast";
import type { Setting } from "@/types/api";

export function SettingsPage() {
  const qc = useQueryClient();
  const { data, isLoading } = useQuery({
    queryKey: ["settings"],
    queryFn: () => apiList<Setting>("/admin/settings"),
  });

  const save = useMutation({
    mutationFn: (s: Setting) =>
      apiPost("/admin/settings", { key: s.key, value: s.value, type: s.type }),
    onSuccess: () => {
      toast.success("Setting saved");
      qc.invalidateQueries({ queryKey: ["settings"] });
    },
    onError: () => toast.error("Save failed"),
  });

  if (isLoading) return <Spinner />;

  return (
    <div>
      <PageHeader title="Settings" subtitle="Feature flags & platform defaults" />
      <div className="card divide-y divide-gray-100">
        {(data?.items ?? []).map((s) => (
          <div key={s.key} className="flex items-center justify-between gap-4 p-4">
            <div>
              <div className="font-medium text-gray-900">{s.key}</div>
              <div className="text-xs text-gray-400">type: {s.type}</div>
            </div>
            {s.type === "bool" ? (
              <button
                role="switch"
                aria-checked={s.value === "true"}
                onClick={() =>
                  save.mutate({ ...s, value: s.value === "true" ? "false" : "true" })
                }
                className={
                  "relative h-6 w-11 rounded-full transition " +
                  (s.value === "true" ? "bg-brand-500" : "bg-gray-300")
                }
              >
                <span
                  className={
                    "absolute top-0.5 h-5 w-5 rounded-full bg-white transition " +
                    (s.value === "true" ? "start-5" : "start-0.5")
                  }
                />
              </button>
            ) : (
              <div className="flex gap-2">
                <input
                  className="input w-48"
                  defaultValue={s.value}
                  onBlur={(e) =>
                    e.target.value !== s.value &&
                    save.mutate({ ...s, value: e.target.value })
                  }
                />
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
