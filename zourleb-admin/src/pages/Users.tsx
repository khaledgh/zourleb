import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiList, apiPost, apiPut, apiDelete } from "@/api/client";
import { PageHeader, StatusBadge } from "@/components/ui/primitives";
import { DataTable, Pager, type Column } from "@/components/ui/DataTable";
import { toast } from "@/components/ui/toast";

interface AdminUser {
  id: number;
  name: string;
  email: string;
  phone: string;
  status: string;
  locale: string;
  roles?: string[];
}

const ROLES = ["super_admin", "agency_owner", "agency_staff", "tourist"] as const;

export function UsersPage() {
  const qc = useQueryClient();
  const [q, setQ] = useState("");
  const [page, setPage] = useState(1);

  // Modal / Form states
  const [modalOpen, setModalOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<AdminUser | null>(null);

  // Field states
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [password, setPassword] = useState("");
  const [userStatus, setUserStatus] = useState<string>("active");
  const [selectedRoles, setSelectedRoles] = useState<string[]>(["tourist"]);

  const { data, isLoading } = useQuery({
    queryKey: ["users", q, page],
    queryFn: () => apiList<AdminUser>("/admin/users", { q: q || undefined, page, per_page: 20 }),
  });

  const createUser = useMutation({
    mutationFn: (body: any) => apiPost("/admin/users", body),
    onSuccess: () => {
      toast.success("User created successfully");
      closeModal();
      qc.invalidateQueries({ queryKey: ["users"] });
    },
    onError: () => toast.error("Failed to create user"),
  });

  const updateUser = useMutation({
    mutationFn: ({ id, body }: { id: number; body: any }) =>
      apiPut(`/admin/users/${id}`, body),
    onSuccess: () => {
      toast.success("User updated successfully");
      closeModal();
      qc.invalidateQueries({ queryKey: ["users"] });
    },
    onError: () => toast.error("Failed to update user"),
  });

  const deleteUser = useMutation({
    mutationFn: (id: number) => apiDelete(`/admin/users/${id}`),
    onSuccess: () => {
      toast.success("User deleted successfully");
      qc.invalidateQueries({ queryKey: ["users"] });
    },
    onError: () => toast.error("Failed to delete user"),
  });

  const setStatus = useMutation({
    mutationFn: ({ id, status }: { id: number; status: string }) =>
      apiPost(`/admin/users/${id}/status`, { status }),
    onSuccess: () => {
      toast.success("User status updated");
      qc.invalidateQueries({ queryKey: ["users"] });
    },
    onError: () => toast.error("Action failed"),
  });

  function openAddModal() {
    setEditingUser(null);
    setName("");
    setEmail("");
    setPhone("");
    setPassword("");
    setUserStatus("active");
    setSelectedRoles(["tourist"]);
    setModalOpen(true);
  }

  function openEditModal(u: AdminUser) {
    setEditingUser(u);
    setName(u.name);
    setEmail(u.email);
    setPhone(u.phone || "");
    setPassword(""); // Keep blank to not change password
    setUserStatus(u.status);
    setSelectedRoles(u.roles ?? []);
    setModalOpen(true);
  }

  function closeModal() {
    setModalOpen(false);
    setEditingUser(null);
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!name || !email) {
      toast.error("Please fill in required fields (Name, Email)");
      return;
    }

    if (!editingUser && !password) {
      toast.error("Password is required for new users");
      return;
    }

    const payload: any = {
      name,
      email,
      phone: phone || undefined,
      status: userStatus,
      roles: selectedRoles,
    };

    if (password) {
      payload.password = password;
    }

    if (editingUser) {
      updateUser.mutate({ id: editingUser.id, body: payload });
    } else {
      createUser.mutate(payload);
    }
  }

  function handleDelete(id: number) {
    if (confirm("Are you sure you want to delete this user? This is a soft-delete.")) {
      deleteUser.mutate(id);
    }
  }

  function toggleRole(role: string) {
    setSelectedRoles((prev) =>
      prev.includes(role) ? prev.filter((r) => r !== role) : [...prev, role]
    );
  }

  const columns: Column<AdminUser>[] = [
    { header: "Name", cell: (u) => <span className="font-medium">{u.name}</span> },
    { header: "Email", cell: (u) => u.email },
    { header: "Phone", cell: (u) => u.phone || "—" },
    { header: "Status", cell: (u) => <StatusBadge status={u.status} /> },
    {
      header: "",
      className: "text-end",
      cell: (u) => (
        <div className="flex justify-end gap-2">
          <button
            className="btn-ghost px-2 py-1 text-xs"
            onClick={() => openEditModal(u)}
          >
            Edit
          </button>
          <button
            className={u.status === "blocked" ? "btn-primary px-3 py-1 text-xs" : "btn-danger px-3 py-1 text-xs"}
            onClick={() =>
              setStatus.mutate({
                id: u.id,
                status: u.status === "blocked" ? "active" : "blocked",
              })
            }
          >
            {u.status === "blocked" ? "Unblock" : "Block"}
          </button>
          <button
            className="btn-danger bg-red-50 text-red-600 border border-red-100 hover:bg-red-100 px-3 py-1 text-xs"
            onClick={() => handleDelete(u.id)}
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
        title="Users"
        actions={
          <div className="flex gap-3">
            <input
              className="input w-64"
              placeholder="Search name or email…"
              value={q}
              onChange={(e) => {
                setQ(e.target.value);
                setPage(1);
              }}
            />
            <button className="btn-primary" onClick={openAddModal}>
              Add User
            </button>
          </div>
        }
      />
      <DataTable columns={columns} rows={data?.items ?? []} loading={isLoading} rowKey={(u) => u.id} />
      <Pager page={page} totalPages={data?.meta?.total_pages ?? 1} onPage={setPage} />

      {/* Modal Dialog */}
      {modalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
          <div className="card w-full max-w-lg bg-white p-6 shadow-2xl animate-fade-in max-h-[90vh] overflow-y-auto">
            <h3 className="text-lg font-bold text-slate-800 mb-4">
              {editingUser ? "Edit User" : "Add New User"}
            </h3>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="label">Full Name *</label>
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
                  <label className="label">Phone (Lebanese)</label>
                  <input
                    className="input"
                    placeholder="e.g. +961 3 123 456"
                    value={phone}
                    onChange={(e) => setPhone(e.target.value)}
                  />
                </div>
              </div>

              <div>
                <label className="label">Password {editingUser && "(Leave blank to keep current)"} *</label>
                <input
                  type="password"
                  className="input"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required={!editingUser}
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="label">Status</label>
                  <select
                    className="input"
                    value={userStatus}
                    onChange={(e) => setUserStatus(e.target.value)}
                  >
                    <option value="active">Active</option>
                    <option value="blocked">Blocked</option>
                  </select>
                </div>
              </div>

              <div>
                <label className="label mb-2 block">Roles</label>
                <div className="flex flex-wrap gap-2">
                  {ROLES.map((role) => {
                    const active = selectedRoles.includes(role);
                    return (
                      <button
                        type="button"
                        key={role}
                        onClick={() => toggleRole(role)}
                        className={`rounded-full px-3 py-1.5 text-xs font-semibold border transition-all ${
                          active
                            ? "bg-brand-500 text-white border-brand-500 shadow-sm"
                            : "bg-slate-50 text-slate-600 border-slate-200 hover:bg-slate-100"
                        }`}
                      >
                        {role.replace(/_/g, " ")}
                      </button>
                    );
                  })}
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
                  disabled={createUser.isPending || updateUser.isPending}
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
