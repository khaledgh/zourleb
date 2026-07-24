import { useState } from "react";
import { useForm } from "react-hook-form";
import { useMutation } from "@tanstack/react-query";
import { apiGet, apiPost } from "@/api/client";
import { PageHeader } from "@/components/ui/primitives";
import { toast } from "@/components/ui/toast";

interface TForm {
  locale: string;
  namespace: string;
  key: string;
  value: string;
}

interface Bundle {
  locale: string;
  translations: Record<string, Record<string, string>>;
}

// Translations editor: load a locale's bundle and upsert individual strings.
export function TranslationsPage() {
  const [locale, setLocale] = useState("en");
  const [rows, setRows] = useState<TForm[]>([]);
  const { register, handleSubmit, reset } = useForm<TForm>({
    defaultValues: { locale: "en", namespace: "common" },
  });

  async function load(code: string) {
    setLocale(code);
    try {
      const env = await apiGet<Bundle>(`/i18n/${code}`);
      const b = env.data;
      const flat: TForm[] = [];
      for (const [ns, kv] of Object.entries(b?.translations ?? {})) {
        for (const [k, v] of Object.entries(kv)) {
          flat.push({ locale: code, namespace: ns, key: k, value: v });
        }
      }
      setRows(flat);
    } catch {
      setRows([]);
    }
  }

  const save = useMutation({
    mutationFn: (v: TForm) => apiPost("/admin/translations", v),
    onSuccess: () => {
      toast.success("Translation saved");
      void load(locale);
    },
    onError: () => toast.error("Save failed"),
  });

  return (
    <div>
      <PageHeader
        title="Translations"
        subtitle="UI strings served to the apps as a JSON bundle"
        actions={
          <div className="flex items-center gap-2">
            <input
              className="input w-24"
              value={locale}
              onChange={(e) => setLocale(e.target.value)}
            />
            <button className="btn-ghost" onClick={() => load(locale)}>
              Load
            </button>
          </div>
        }
      />

      <form
        onSubmit={handleSubmit((v) => save.mutate({ ...v, locale }))}
        className="card mb-6 grid grid-cols-2 gap-4 p-4 md:grid-cols-4"
      >
        <div>
          <label className="label">Namespace</label>
          <input className="input" {...register("namespace", { required: true })} />
        </div>
        <div>
          <label className="label">Key</label>
          <input className="input" {...register("key", { required: true })} />
        </div>
        <div className="md:col-span-2">
          <label className="label">Value</label>
          <input className="input" {...register("value", { required: true })} />
        </div>
        <div className="col-span-full">
          <button type="submit" className="btn-primary">
            Save string
          </button>
          <button
            type="button"
            className="btn-ghost ms-2"
            onClick={() => reset()}
          >
            Clear
          </button>
        </div>
      </form>

      <div className="card overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 text-xs uppercase text-gray-500">
            <tr>
              <th className="px-4 py-3 text-start">Namespace</th>
              <th className="px-4 py-3 text-start">Key</th>
              <th className="px-4 py-3 text-start">Value</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {rows.length === 0 && (
              <tr>
                <td colSpan={3} className="px-4 py-8 text-center text-gray-400">
                  Load a locale to view its strings.
                </td>
              </tr>
            )}
            {rows.map((r) => (
              <tr key={`${r.namespace}.${r.key}`} className="hover:bg-gray-50">
                <td className="px-4 py-2 font-mono text-xs">{r.namespace}</td>
                <td className="px-4 py-2 font-mono text-xs">{r.key}</td>
                <td className="px-4 py-2">
                  <input
                    className="input"
                    defaultValue={r.value}
                    onBlur={(e) =>
                      e.target.value !== r.value &&
                      save.mutate({ ...r, value: e.target.value })
                    }
                  />
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
