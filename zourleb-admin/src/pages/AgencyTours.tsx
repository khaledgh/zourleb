import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiList, apiPost, apiPut } from "@/api/client";
import { PageHeader, StatusBadge } from "@/components/ui/primitives";
import { DataTable, type Column } from "@/components/ui/DataTable";
import { toast } from "@/components/ui/toast";
import type { Category, Region } from "@/types/api";

interface TourTranslationInput {
  locale: string;
  title: string;
  summary: string;
  description: string;
  itinerary: string;
  included: string;
  excluded: string;
}

interface SaveTourRequest {
  category_id: number | null;
  region_id: number | null;
  type: "day_trip" | "multi_day";
  duration_days: number;
  difficulty: string;
  min_age: number;
  max_capacity: number;
  base_currency: "USD" | "LBP";
  price_from: number;
  translations: TourTranslationInput[];
}

interface TourTranslation {
  locale: string;
  title: string;
  summary: string;
  description: string;
  itinerary: string;
  included: string;
  excluded: string;
}

interface AgencyTour {
  id: number;
  slug: string;
  status: string;
  type: "day_trip" | "multi_day";
  price_from: number;
  base_currency: "USD" | "LBP";
  featured: boolean;
  category_id: number | null;
  region_id: number | null;
  duration_days: number;
  difficulty: string;
  min_age: number;
  max_capacity: number;
  translations: TourTranslation[];
}

function getTranslation(t: AgencyTour, locale: string): TourTranslation | undefined {
  return t.translations?.find((tr) => tr.locale === locale) || t.translations?.[0];
}

export function AgencyToursPage() {
  const qc = useQueryClient();
  const [modalOpen, setModalOpen] = useState(false);
  const [editingTour, setEditingTour] = useState<AgencyTour | null>(null);

  // Form states
  const [categoryId, setCategoryId] = useState<string>("");
  const [regionId, setRegionId] = useState<string>("");
  const [tourType, setTourType] = useState<"day_trip" | "multi_day">("day_trip");
  const [durationDays, setDurationDays] = useState<number>(1);
  const [difficulty, setDifficulty] = useState<string>("medium");
  const [minAge, setMinAge] = useState<number>(6);
  const [maxCapacity, setMaxCapacity] = useState<number>(20);
  const [baseCurrency, setBaseCurrency] = useState<"USD" | "LBP">("USD");
  const [priceFrom, setPriceFrom] = useState<number>(0);

  // Translations form states (we allow editing English and Arabic)
  const [titleAr, setTitleAr] = useState("");
  const [summaryAr, setSummaryAr] = useState("");
  const [descriptionAr, setDescriptionAr] = useState("");
  const [titleEn, setTitleEn] = useState("");
  const [summaryEn, setSummaryEn] = useState("");
  const [descriptionEn, setDescriptionEn] = useState("");


  const { data: categories } = useQuery({
    queryKey: ["categories"],
    queryFn: () => apiList<Category>("/categories"),
  });

  const { data: regions } = useQuery({
    queryKey: ["regions"],
    queryFn: () => apiList<Region>("/regions"),
  });

  const { data, isLoading } = useQuery({
    queryKey: ["agency-tours"],
    queryFn: () => apiList<AgencyTour>("/agency/tours", { per_page: 50 }),
  });

  const createTour = useMutation({
    mutationFn: (body: SaveTourRequest) => apiPost("/agency/tours", body),
    onSuccess: () => {
      toast.success("Tour created successfully");
      closeModal();
      qc.invalidateQueries({ queryKey: ["agency-tours"] });
    },
    onError: () => toast.error("Failed to create tour"),
  });

  const updateTour = useMutation({
    mutationFn: ({ id, body }: { id: number; body: SaveTourRequest }) =>
      apiPut(`/agency/tours/${id}`, body),
    onSuccess: () => {
      toast.success("Tour updated successfully");
      closeModal();
      qc.invalidateQueries({ queryKey: ["agency-tours"] });
    },
    onError: () => toast.error("Failed to update tour"),
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

  function openAddModal() {
    setEditingTour(null);
    setCategoryId("");
    setRegionId("");
    setTourType("day_trip");
    setDurationDays(1);
    setDifficulty("medium");
    setMinAge(6);
    setMaxCapacity(20);
    setBaseCurrency("USD");
    setPriceFrom(0);
    setTitleAr("");
    setSummaryAr("");
    setDescriptionAr("");
    setTitleEn("");
    setSummaryEn("");
    setDescriptionEn("");
    setModalOpen(true);
  }

  function openEditModal(t: AgencyTour) {
    setEditingTour(t);
    setCategoryId(t.category_id ? String(t.category_id) : "");
    setRegionId(t.region_id ? String(t.region_id) : "");
    setTourType(t.type);
    setDurationDays(t.duration_days);
    setDifficulty(t.difficulty || "medium");
    setMinAge(t.min_age);
    setMaxCapacity(t.max_capacity);
    setBaseCurrency(t.base_currency || "USD");
    setPriceFrom(t.price_from || 0);

    const transAr = t.translations?.find((tr) => tr.locale === "ar");
    const transEn = t.translations?.find((tr) => tr.locale === "en");

    setTitleAr(transAr?.title || "");
    setSummaryAr(transAr?.summary || "");
    setDescriptionAr(transAr?.description || "");

    setTitleEn(transEn?.title || "");
    setSummaryEn(transEn?.summary || "");
    setDescriptionEn(transEn?.description || "");

    setModalOpen(true);
  }

  function closeModal() {
    setModalOpen(false);
    setEditingTour(null);
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!titleAr && !titleEn) {
      toast.error("Please fill in at least one title (Arabic or English)");
      return;
    }

    const translations: TourTranslationInput[] = [];
    if (titleAr) {
      translations.push({
        locale: "ar",
        title: titleAr,
        summary: summaryAr,
        description: descriptionAr,
        itinerary: "",
        included: "",
        excluded: "",
      });
    }
    if (titleEn) {
      translations.push({
        locale: "en",
        title: titleEn,
        summary: summaryEn,
        description: descriptionEn,
        itinerary: "",
        included: "",
        excluded: "",
      });
    }

    const payload: SaveTourRequest = {
      category_id: categoryId ? Number(categoryId) : null,
      region_id: regionId ? Number(regionId) : null,
      type: tourType,
      duration_days: Number(durationDays),
      difficulty,
      min_age: Number(minAge),
      max_capacity: Number(maxCapacity),
      base_currency: baseCurrency,
      price_from: Number(priceFrom),
      translations,
    };

    if (editingTour) {
      updateTour.mutate({ id: editingTour.id, body: payload });
    } else {
      createTour.mutate(payload);
    }
  }

  const columns: Column<AgencyTour>[] = [
    {
      header: "Title",
      cell: (t) => <span className="font-medium">{getTranslation(t, "en")?.title || t.slug}</span>,
    },
    { header: "Type", cell: (t) => t.type.replace(/_/g, " ") },
    { header: "From", cell: (t) => `${t.price_from} ${t.base_currency}` },
    { header: "Status", cell: (t) => <StatusBadge status={t.status} /> },
    {
      header: "",
      className: "text-end",
      cell: (t) => (
        <div className="flex justify-end gap-2">
          <button
            className="btn-ghost px-3 py-1 text-xs"
            onClick={() => openEditModal(t)}
          >
            Edit
          </button>
          <button
            className="btn-ghost px-3 py-1 text-xs"
            onClick={() => publish.mutate({ id: t.id, publish: t.status !== "published" })}
          >
            {t.status === "published" ? "Unpublish" : "Publish"}
          </button>
        </div>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title="My tours"
        subtitle="Manage and publish your tours"
        actions={
          <button className="btn-primary" onClick={openAddModal}>
            Create Tour
          </button>
        }
      />
      <DataTable
        columns={columns}
        rows={data?.items ?? []}
        loading={isLoading}
        rowKey={(t) => t.id}
        empty="You haven't created any tours yet."
      />

      {/* Modal Dialog */}
      {modalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
          <div className="card w-full max-w-2xl bg-white p-6 shadow-2xl animate-fade-in max-h-[90vh] overflow-y-auto">
            <h3 className="text-lg font-bold text-slate-800 mb-4">
              {editingTour ? "Edit Tour" : "Create New Tour"}
            </h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="label">Category</label>
                  <select
                    className="input"
                    value={categoryId}
                    onChange={(e) => setCategoryId(e.target.value)}
                  >
                    <option value="">Select Category</option>
                    {(categories?.items ?? []).map((c) => (
                      <option key={c.id} value={c.id}>
                        {c.name}
                      </option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="label">Region</label>
                  <select
                    className="input"
                    value={regionId}
                    onChange={(e) => setRegionId(e.target.value)}
                  >
                    <option value="">Select Region</option>
                    {(regions?.items ?? []).map((r) => (
                      <option key={r.id} value={r.id}>
                        {r.name}
                      </option>
                    ))}
                  </select>
                </div>
              </div>

              <div className="grid grid-cols-3 gap-4">
                <div>
                  <label className="label">Type</label>
                  <select
                    className="input"
                    value={tourType}
                    onChange={(e) => setTourType(e.target.value as any)}
                  >
                    <option value="day_trip">Day Trip</option>
                    <option value="multi_day">Multi Day</option>
                  </select>
                </div>
                <div>
                  <label className="label">Duration (Days)</label>
                  <input
                    type="number"
                    min="1"
                    className="input"
                    value={durationDays}
                    onChange={(e) => setDurationDays(Number(e.target.value))}
                  />
                </div>
                <div>
                  <label className="label">Difficulty</label>
                  <select
                    className="input"
                    value={difficulty}
                    onChange={(e) => setDifficulty(e.target.value)}
                  >
                    <option value="easy">Easy</option>
                    <option value="medium">Medium</option>
                    <option value="hard">Hard</option>
                  </select>
                </div>
              </div>

              <div className="grid grid-cols-4 gap-4">
                <div>
                  <label className="label">Min Age</label>
                  <input
                    type="number"
                    className="input"
                    value={minAge}
                    onChange={(e) => setMinAge(Number(e.target.value))}
                  />
                </div>
                <div>
                  <label className="label">Max Capacity</label>
                  <input
                    type="number"
                    className="input"
                    value={maxCapacity}
                    onChange={(e) => setMaxCapacity(Number(e.target.value))}
                  />
                </div>
                <div>
                  <label className="label">Currency</label>
                  <select
                    className="input"
                    value={baseCurrency}
                    onChange={(e) => setBaseCurrency(e.target.value as any)}
                  >
                    <option value="USD">USD</option>
                    <option value="LBP">LBP</option>
                  </select>
                </div>
                <div>
                  <label className="label">Price From *</label>
                  <input
                    type="number"
                    className="input"
                    value={priceFrom}
                    onChange={(e) => setPriceFrom(Number(e.target.value))}
                    required
                  />
                </div>
              </div>

              <div className="border-t border-slate-100 pt-4">
                <h4 className="font-semibold text-sm text-slate-700 mb-2">Translations (Arabic)</h4>
                <div className="space-y-3">
                  <div>
                    <label className="label">Title (AR)</label>
                    <input
                      className="input text-right"
                      dir="rtl"
                      value={titleAr}
                      onChange={(e) => setTitleAr(e.target.value)}
                    />
                  </div>
                  <div>
                    <label className="label">Summary (AR)</label>
                    <input
                      className="input text-right"
                      dir="rtl"
                      value={summaryAr}
                      onChange={(e) => setSummaryAr(e.target.value)}
                    />
                  </div>
                  <div>
                    <label className="label">Description (AR)</label>
                    <textarea
                      className="input h-20 text-right"
                      dir="rtl"
                      value={descriptionAr}
                      onChange={(e) => setDescriptionAr(e.target.value)}
                    />
                  </div>
                </div>
              </div>

              <div className="border-t border-slate-100 pt-4">
                <h4 className="font-semibold text-sm text-slate-700 mb-2">Translations (English)</h4>
                <div className="space-y-3">
                  <div>
                    <label className="label">Title (EN)</label>
                    <input
                      className="input"
                      value={titleEn}
                      onChange={(e) => setTitleEn(e.target.value)}
                    />
                  </div>
                  <div>
                    <label className="label">Summary (EN)</label>
                    <input
                      className="input"
                      value={summaryEn}
                      onChange={(e) => setSummaryEn(e.target.value)}
                    />
                  </div>
                  <div>
                    <label className="label">Description (EN)</label>
                    <textarea
                      className="input h-20"
                      value={descriptionEn}
                      onChange={(e) => setDescriptionEn(e.target.value)}
                    />
                  </div>
                </div>
              </div>

              <div className="flex justify-end gap-3 pt-4 border-t border-slate-100">
                <button
                  type="button"
                  className="btn-ghost"
                  onClick={closeModal}
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="btn-primary"
                  disabled={createTour.isPending || updateTour.isPending}
                >
                  Save
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
