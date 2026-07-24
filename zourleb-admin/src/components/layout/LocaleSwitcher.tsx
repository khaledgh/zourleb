import { useEffect, useState } from "react";
import { apiList } from "@/api/client";
import { changeLocale } from "@/i18n";
import type { Language } from "@/types/api";

// LocaleSwitcher pulls active languages from the backend so a newly-added
// language appears without a redeploy.
export function LocaleSwitcher() {
  const [langs, setLangs] = useState<Language[]>([]);
  const [current, setCurrent] = useState(
    localStorage.getItem("zb_locale") || "ar",
  );

  useEffect(() => {
    apiList<Language>("/languages")
      .then(({ items }) => setLangs(items))
      .catch(() => setLangs([]));
  }, []);

  async function onChange(code: string) {
    setCurrent(code);
    await changeLocale(code);
    // Re-render the whole tree so localized server data refetches with the new
    // Accept-Language header.
    window.location.reload();
  }

  if (langs.length === 0) return null;
  return (
    <select
      value={current}
      onChange={(e) => onChange(e.target.value)}
      className="rounded-lg border border-gray-300 px-2 py-1 text-sm"
    >
      {langs.map((l) => (
        <option key={l.code} value={l.code}>
          {l.native_name}
        </option>
      ))}
    </select>
  );
}
