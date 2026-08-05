import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiList, apiPost, apiPut, apiDelete } from "@/api/client";
import { PageHeader, StatusBadge } from "@/components/ui/primitives";
import { DataTable, Pager, type Column } from "@/components/ui/DataTable";
import { toast } from "@/components/ui/toast";
import type { Agency } from "@/types/api";

interface Region {
  id: number;
  name: string;
}

const STATUSES = ["", "pending", "approved", "suspended"] as const;

export function AgenciesPage() {
  const qc = useQueryClient();
  const [status, setStatus] = useState<string>("pending");
  const [page, setPage] = useState(1);

  // Form / Modal states
  const [modalOpen, setModalOpen] = useState(false);
  const [editingAgency, setEditingAgency] = useState<Agency | null>(null);
  
  // Fields state
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [website, setWebsite] = useState("");
  const [regionId, setRegionId] = useState<string>("");
  const [agencyStatus, setAgencyStatus] = useState<string>("pending");
  const [commissionRate, setCommissionRate] = useState<number>(10);
  const [subscriptionTier, setSubscriptionTier] = useState<string>("basic");

  // Fetch regions
  const { data: regions } = useQuery({
    queryKey: ["regions"],
    queryFn: () => apiList<Region>("/regions"),
  });

  // Fetch agencies
  const { data, isLoading } = useQuery({
    queryKey: ["agencies", status, page],
    queryFn: () =>
      apiList<Agency>("/admin/agencies", { status: status || undefined, page, per_page: 20 }),
  });

  const createAgency = useMutation({
    mutationFn: (body: any) => apiPost("/admin/agencies", body),
    onSuccess: () => {
      toast.success("Agency created successfully");
      closeModal();
      qc.invalidateQueries({ queryKey: ["agencies"] });
    },
    onError: () => toast.error("Failed to create agency"),
  });

  const updateAgency = useMutation({
    mutationFn: ({ id, body }: { id: number; body: any }) =>
      apiPut(`/admin/agencies/${id}`, body), // Using PUT backend endpoint via post helper (it is mapped to PUT in the router but post client can handle it or we can use generic client)
    onSuccess: () => {
      toast.success("Agency updated successfully");
      closeModal();
      qc.invalidateQueries({ queryKey: ["agencies"] });
    },
    onError: () => toast.error("Failed to update agency"),
  });

  const deleteAgency = useMutation({
    mutationFn: (id: number) => apiDelete(`/admin/agencies/${id}`),
    onSuccess: () => {
      toast.success("Agency deleted successfully");
      qc.invalidateQueries({ queryKey: ["agencies"] });
    },
    onError: () => toast.error("Failed to delete agency"),
  });

  const decide = useMutation({
    mutationFn: ({ id, approve }: { id: number; approve: boolean }) =>
      apiPost(`/admin/agencies/${id}/approve`, { approve }),
    onSuccess: (_d, v) => {
      toast.success(v.approve ? "Agency approved" : "Agency suspended");
      qc.invalidateQueries({ queryKey: ["agencies"] });
    },
    onError: () => toast.error("Action failed"),
  });

  function openAddModal() {
    setEditingAgency(null);
    setName("");
    setEmail("");
    setPhone("");
    setWebsite("");
    setRegionId("");
    setAgencyStatus("pending");
    setCommissionRate(10);
    setSubscriptionTier("basic");
    setModalOpen(true);
  }

  function openEditModal(a: Agency) {
    // Cast to any to access additional details from raw API payload
    const raw = a as any;
    setEditingAgency(a);
    setName(a.name);
    setEmail(a.email);
    setPhone(a.phone);
    setWebsite(raw.website || "");
    setRegionId(raw.region_id ? String(raw.region_id) : "");
    setAgencyStatus(a.status);
    setCommissionRate(raw.commission_rate ?? 10);
    setSubscriptionTier(raw.subscription_tier ?? "basic");
    setModalOpen(true);
  }

  function closeModal() {
    setModalOpen(false);
    setEditingAgency(null);
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!name || !email || !phone) {
      toast.error("Please fill in required fields (Name, Email, Phone)");
      return;
    }

    const payload = {
      name,
      email,
      phone,
      website: website || undefined,
      region_id: regionId ? Number(regionId) : undefined,
      status: agencyStatus,
      commission_rate: Number(commissionRate),
      subscription_tier: subscriptionTier,
    };

    if (editingAgency) {
      updateAgency.mutate({ id: editingAgency.id, body: payload });
    } else {
      createAgency.mutate(payload);
    }
  }

  function handleDelete(id: number) {
    if (confirm("Are you sure you want to delete this agency? This is a soft-delete.")) {
      deleteAgency.mutate(id);
    }
  }

  const columns: Column<Agency>[] = [
    { header: "Name", cell: (a) => <span className="font-medium">{a.name}</span> },
    { header: "Email", cell: (a) => a.email },
    { header: "Phone", cell: (a) => a.phone },
    { header: "Status", cell: (a) => <StatusBadge status={a.status} /> },
    { header: "Rating", cell: (a) => `${a.rating_avg.toFixed(1)} (${a.rating_count})` },
    {
      header: "",
      className: "text-end",
      cell: (a) => (
        <div className="flex justify-end gap-2">
          <button
            className="btn-ghost px-2 py-1 text-xs"
            onClick={() => openEditModal(a)}
          >
            Edit
          </button>
          {a.status !== "approved" && (
            <button
              className="btn-primary px-3 py-1 text-xs"
              onClick={() => decide.mutate({ id: a.id, approve: true })}
            >
              Approve
            </button>
          )}
          {a.status !== "suspended" && (
            <button
              className="btn-danger px-3 py-1 text-xs"
              onClick={() => decide.mutate({ id: a.id, approve: false })}
            >
              Suspend
            </button>
          )}
          <button
            className="btn-danger bg-red-50 text-red-600 border border-red-100 hover:bg-red-100 px-3 py-1 text-xs"
            onClick={() => handleDelete(a.id)}
          >
            Delete
          </button>
        </div>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title="Agencies"
        subtitle="Manage travel agencies"
        actions={
          <div className="flex gap-3">
            <select
              className="input w-44"
              value={status}
              onChange={(e) => {
                setStatus(e.target.value);
                setPage(1);
              }}
            >
              {STATUSES.map((s) => (
                <option key={s} value={s}>
                  {s === "" ? "All statuses" : s}
                </option>
              ))}
            </select>
            <button className="btn-primary" onClick={openAddModal}>
              Add Agency
            </button>
          </div>
        }
      />
      <DataTable
        columns={columns}
        rows={data?.items ?? []}
        loading={isLoading}
        rowKey={(a) => a.id}
        empty="No agencies match this filter."
      />
      <Pager
        page={page}
        totalPages={data?.meta?.total_pages ?? 1}
        onPage={setPage}
      />

      {/* Modal dialog */}
      {modalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
          <div className="card w-full max-w-lg bg-white p-6 shadow-2xl animate-fade-in max-h-[90vh] overflow-y-auto">
            <h3 className="text-lg font-bold text-slate-800 mb-4">
              {editingAgency ? "Edit Agency" : "Add New Agency"}
            </h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="label">Agency Name *</label>
                <input
                  className="input"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  required
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="label">Email *</label>
                  <input
                    type="email"
                    className="input"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                  />
                </div>
                <div>
                  <label className="label">Phone *</label>
                  <input
                    className="input"
                    placeholder="e.g. +961 3 123 456"
                    value={phone}
                    onChange={(e) => setPhone(e.target.value)}
                    required
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="label">Website</label>
                  <input
                    className="input"
                    placeholder="e.g. www.agency.com"
                    value={website}
                    onChange={(e) => setWebsite(e.target.value)}
                  />
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
                  <label className="label">Status</label>
                  <select
                    className="input"
                    value={agencyStatus}
                    onChange={(e) => setAgencyStatus(e.target.value)}
                  >
                    <option value="pending">Pending</option>
                    <option value="approved">Approved</option>
                    <option value="suspended">Suspended</option>
                  </select>
                </div>
                <div>
                  <label className="label">Tier</label>
                  <select
                    className="input"
                    value={subscriptionTier}
                    onChange={(e) => setSubscriptionTier(e.target.value)}
                  >
                    <option value="basic">Basic</option>
                    <option value="pro">Pro</option>
                    <option value="premium">Premium</option>
                  </select>
                </div>
                <div>
                  <label className="label">Commission %</label>
                  <input
                    type="number"
                    className="input"
                    min="0"
                    max="100"
                    value={commissionRate}
                    onChange={(e) => setCommissionRate(Number(e.target.value))}
                  />
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
                  disabled={createAgency.isPending || updateAgency.isPending}
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
