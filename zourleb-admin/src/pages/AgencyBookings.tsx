import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { apiList } from "@/api/client";
import { PageHeader, StatusBadge } from "@/components/ui/primitives";
import { DataTable, Pager, type Column } from "@/components/ui/DataTable";

interface AgencyBooking {
  id: number;
  code: string;
  status: string;
  payment_status: string;
  travelers_count: number;
  subtotal: number;
  currency: string;
  contact_phone: string;
  created_at: string;
}

export function AgencyBookingsPage() {
  const [page, setPage] = useState(1);
  const { data, isLoading } = useQuery({
    queryKey: ["agency-bookings", page],
    queryFn: () => apiList<AgencyBooking>("/agency/bookings", { page, per_page: 20 }),
  });

  const columns: Column<AgencyBooking>[] = [
    { header: "Code", cell: (b) => <span className="font-mono">{b.code}</span> },
    { header: "Travelers", cell: (b) => b.travelers_count },
    { header: "Total", cell: (b) => `${b.subtotal} ${b.currency}` },
    { header: "Phone", cell: (b) => b.contact_phone },
    { header: "Status", cell: (b) => <StatusBadge status={b.status} /> },
    { header: "Payment", cell: (b) => <StatusBadge status={b.payment_status} /> },
  ];

  return (
    <div>
      <PageHeader title="Bookings" subtitle="Reservations across your tours" />
      <DataTable
        columns={columns}
        rows={data?.items ?? []}
        loading={isLoading}
        rowKey={(b) => b.id}
        empty="No bookings yet."
      />
      <Pager page={page} totalPages={data?.meta?.total_pages ?? 1} onPage={setPage} />
    </div>
  );
}
