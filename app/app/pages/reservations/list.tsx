import { useState, useEffect } from "react";
import type { ColumnDef, SortingState } from "@tanstack/react-table";
import { useQuery, keepPreviousData } from "@tanstack/react-query";
import { useNavigate, Link } from "react-router";
import { getReservations, getReservationStats } from "~/lib/api-client";
import { MainLayout } from "~/components/main/main-layout";
import { useReservationListStore } from "~/lib/store/reservation-list-store";
import { DataTable, selectColumn } from "~/components/data-table/data-table";
import { DataTableColumnHeader } from "~/components/data-table/data-table-column-header";
import { StatusBadge } from "~/components/reservations/status-badge";
import type { ReservationListItem } from "~/lib/types/reservation";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { Skeleton } from "~/components/ui/skeleton";
import { cn } from "~/lib/utils";
import { Plus, Search, Inbox, TriangleAlert } from "lucide-react";
import { queryKeys } from "~/lib/query-keys";

const STATUS_TABS = [
  { label: "All", value: "" },
  { label: "Pending", value: "pending" },
  { label: "Sample Received", value: "sample_received" },
  { label: "Test On Progress", value: "test_on_progress" },
  { label: "Test Completed", value: "test_completed" },
  { label: "CoA Submitted", value: "coa_submitted" },
  { label: "Head Approved", value: "head_approved" },
  { label: "Order Completed", value: "order_completed" },
  { label: "Rejected", value: "rejected" },
] as const;

const STAT_CARDS = [
  { label: "Total Reservations", key: "total", dot: "bg-primary" },
  { label: "Pending", key: "pending", dot: "bg-gray-400" },
  { label: "Test On Progress", key: "test_on_progress", dot: "bg-purple-500" },
  { label: "Test Completed", key: "test_completed", dot: "bg-emerald-500" },
] as const;

function formatDate(isoDate: string): string {
  return new Date(isoDate).toLocaleDateString(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

const reservationColumns: ColumnDef<ReservationListItem>[] = [
  selectColumn<ReservationListItem>(),
  {
    accessorKey: "test_request_number",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Test Request Number" />
    ),
    cell: ({ row }) => (
      <span className="font-medium text-foreground">
        {row.original.test_request_number}
      </span>
    ),
  },
  {
    accessorKey: "sample_name",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Sample Name" />
    ),
  },
  {
    accessorKey: "status",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Status" />
    ),
    cell: ({ row }) => <StatusBadge status={row.original.status} />,
  },
  {
    accessorKey: "reservation_date",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Reservation Date" />
    ),
    cell: ({ row }) => formatDate(row.original.reservation_date),
  },
];

function useDebounce(value: string, delay: number): string {
  const [debouncedValue, setDebouncedValue] = useState(value);

  useEffect(() => {
    const timer = setTimeout(() => setDebouncedValue(value), delay);
    return () => clearTimeout(timer);
  }, [value, delay]);

  return debouncedValue;
}

function TableSkeleton() {
  return (
    <div className="space-y-2 rounded-xl border border-border bg-card p-4 shadow-sm">
      {Array.from({ length: 6 }).map((_, i) => (
        <Skeleton key={i} className="h-10 w-full" />
      ))}
    </div>
  );
}

function ErrorState({ onRetry }: { onRetry: () => void }) {
  return (
    <div className="flex flex-col items-center justify-center gap-3 rounded-xl border border-border bg-card py-16 shadow-sm">
      <TriangleAlert className="size-8 text-muted-foreground" />
      <p className="text-sm text-muted-foreground">
        Could not load reservations
      </p>
      <Button onClick={onRetry} variant="outline" size="sm">
        Retry
      </Button>
    </div>
  );
}

function EmptyState() {
  return (
    <div className="flex flex-col items-center justify-center gap-3 rounded-xl border border-border bg-card py-16 shadow-sm">
      <Inbox className="size-8 text-muted-foreground" />
      <p className="text-sm text-muted-foreground">No reservations found</p>
      <Button asChild variant="outline" size="sm">
        <Link to="/reservations/new">
          <Plus data-icon="inline-start" />
          Create a reservation
        </Link>
      </Button>
    </div>
  );
}

export default function ReservationListPage() {
  const {
    page,
    setPage,
    activeStatus,
    setActiveStatus,
    search,
    setSearch,
    sorting,
    setSorting,
  } = useReservationListStore();
  const debouncedSearch = useDebounce(search, 300);

  const sort = sorting[0];

  const navigate = useNavigate();

  const statsQuery = useQuery({
    queryKey: queryKeys.reservationStats,
    queryFn: getReservationStats,
    staleTime: 60_000,
  });
  const stats = statsQuery.data?.data;

  const { data, isPending, isError, refetch } = useQuery({
    queryKey: queryKeys.reservations({
      page,
      status: activeStatus,
      search: debouncedSearch,
      sorting,
    }),
    queryFn: () =>
      getReservations({
        page,
        limit: 10,
        status: activeStatus ? [activeStatus] : undefined,
        search: debouncedSearch || undefined,
        sort_by: sort?.id,
        sort_order: sort ? (sort.desc ? "desc" : "asc") : undefined,
      }),
    placeholderData: keepPreviousData,
  });

  const items = data?.data.data ?? [];
  const totalPages = data?.data.total_pages ?? 1;
  const currentPage = data?.data.page ?? 1;
  const totalRows = data?.data.total ?? items.length;

  return (
    <MainLayout breadcrumbs={[{ label: "Reservations" }]}>
      <div className="flex flex-col gap-1.5 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="font-heading text-lg font-semibold">Reservations</h1>
          <p className="text-xs text-muted-foreground">
            Track and manage laboratory test reservations.
          </p>
        </div>
        <Button
          asChild
          className="bg-foreground text-background hover:bg-foreground/90"
        >
          <Link to="/reservations/new">
            <Plus data-icon="inline-start" />
            New Reservation
          </Link>
        </Button>
      </div>

      {/* Stat cards */}
      <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
        {STAT_CARDS.map((card) => (
          <div
            key={card.key}
            className="rounded-xl border border-border bg-card p-4 shadow-sm"
          >
            <div className="flex items-center gap-2 text-xs text-muted-foreground">
              <span className={cn("size-1.5 rounded-full", card.dot)} />
              {card.label}
            </div>
            {statsQuery.isPending ? (
              <Skeleton className="mt-2 h-8 w-16" />
            ) : (
              <p className="mt-2 font-heading text-2xl font-semibold tabular-nums">
                {stats?.[card.key] ?? 0}
              </p>
            )}
          </div>
        ))}
      </div>

      <div className="relative w-full sm:max-w-xs">
        <Search className="absolute left-2 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
        <Input
          placeholder="Search by request number or sample name"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="pl-7"
          aria-label="Search reservations"
        />
      </div>

      {/* Status tabs */}
      <div className="flex items-center gap-1 overflow-x-auto border-b border-border">
        {STATUS_TABS.map((tab) => (
          <button
            key={tab.label}
            type="button"
            onClick={() => setActiveStatus(tab.value)}
            aria-pressed={activeStatus === tab.value}
            className={cn(
              "cursor-pointer whitespace-nowrap border-b-2 px-3 py-2 text-sm font-medium transition-colors",
              activeStatus === tab.value
                ? "border-primary text-foreground"
                : "border-transparent text-muted-foreground hover:text-foreground",
            )}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {isPending ? (
        <TableSkeleton />
      ) : isError ? (
        <ErrorState onRetry={() => refetch()} />
      ) : items.length === 0 ? (
        <EmptyState />
      ) : (
        <DataTable
          columns={reservationColumns}
          data={items}
          getRowId={(row) => row.xid}
          onRowClick={(row) => navigate(`/reservations/${row.xid}`)}
          emptyMessage="No reservations found"
          serverSorting={{ sorting, onSortingChange: setSorting }}
          serverPagination={{
            pageIndex: currentPage - 1,
            pageCount: totalPages,
            pageSize: 10,
            totalRows,
            onPageChange: (pageIndex) => setPage(pageIndex + 1),
          }}
        />
      )}
    </MainLayout>
  );
}
