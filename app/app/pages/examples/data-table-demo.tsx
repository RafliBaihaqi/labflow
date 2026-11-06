import * as React from "react";
import type { ColumnDef } from "@tanstack/react-table";
import {
  Copy,
  Download,
  MoreHorizontal,
  Pencil,
  Plus,
  Printer,
  Trash2,
} from "lucide-react";

import { MainLayout } from "~/components/main/main-layout";
import {
  DataTable,
  selectColumn,
} from "~/components/data-table/data-table";
import { DataTableColumnHeader } from "~/components/data-table/data-table-column-header";
import { Avatar, AvatarFallback } from "~/components/ui/avatar";
import { Badge } from "~/components/ui/badge";
import { Button } from "~/components/ui/button";
import { cn } from "~/lib/utils";

type OrderStatus = "Pending" | "Completed" | "Refunded";
type PaymentStatus = "Paid" | "Unpaid";

interface Order {
  id: string;
  number: string;
  customer: string;
  date: string;
  status: OrderStatus;
  amount: number;
  payment: PaymentStatus;
}

const ORDERS: Order[] = [
  { id: "1008", number: "#ORD1008", customer: "Esther Kiehn", date: "17 Dec 2024", status: "Pending", amount: 10.5, payment: "Unpaid" },
  { id: "1007", number: "#ORD1007", customer: "Denise Kuhn", date: "16 Dec 2024", status: "Pending", amount: 100.5, payment: "Unpaid" },
  { id: "1006", number: "#ORD1006", customer: "Clint Hoppe", date: "16 Dec 2024", status: "Completed", amount: 60.56, payment: "Paid" },
  { id: "1005", number: "#ORD1005", customer: "Darin Deckow", date: "16 Dec 2024", status: "Refunded", amount: 640.5, payment: "Paid" },
  { id: "1004", number: "#ORD1004", customer: "Jacquelyn Robel", date: "15 Dec 2024", status: "Completed", amount: 39.5, payment: "Paid" },
  { id: "1003", number: "#ORD1003", customer: "Clint Hoppe", date: "16 Dec 2024", status: "Completed", amount: 29.5, payment: "Paid" },
  { id: "1002", number: "#ORD1002", customer: "Erin Bins", date: "16 Dec 2024", status: "Completed", amount: 120.35, payment: "Paid" },
  { id: "1001", number: "#ORD1001", customer: "Gretchen Quitzon", date: "14 Dec 2024", status: "Refunded", amount: 123.5, payment: "Paid" },
  { id: "1000", number: "#ORD1000", customer: "Stewart Kulas", date: "14 Dec 2024", status: "Completed", amount: 80.0, payment: "Paid" },
];

const STATUS_VARIANT: Record<OrderStatus, "warning" | "success" | "destructive"> = {
  Pending: "warning",
  Completed: "success",
  Refunded: "destructive",
};

const TABS = ["All", "Incomplete", "Overdue", "Ongoing", "Finished"] as const;

function initials(name: string): string {
  return name
    .split(" ")
    .map((n) => n[0])
    .slice(0, 2)
    .join("")
    .toUpperCase();
}

const columns: ColumnDef<Order>[] = [
  selectColumn<Order>(),
  {
    accessorKey: "number",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Order Number" />
    ),
    cell: ({ row }) => (
      <span className="font-medium text-foreground">{row.original.number}</span>
    ),
  },
  {
    accessorKey: "customer",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Customer Name" />
    ),
    cell: ({ row }) => (
      <div className="flex items-center gap-2">
        <Avatar className="size-6">
          <AvatarFallback className="text-[10px]">
            {initials(row.original.customer)}
          </AvatarFallback>
        </Avatar>
        <span>{row.original.customer}</span>
      </div>
    ),
  },
  {
    accessorKey: "date",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Order Date" />
    ),
  },
  {
    accessorKey: "status",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Status" />
    ),
    cell: ({ row }) => (
      <Badge variant={STATUS_VARIANT[row.original.status]}>
        {row.original.status}
      </Badge>
    ),
  },
  {
    accessorKey: "amount",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Total Amount" />
    ),
    cell: ({ row }) => (
      <span className="tabular-nums">${row.original.amount.toFixed(2)}</span>
    ),
  },
  {
    accessorKey: "payment",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Payment Status" />
    ),
    cell: ({ row }) => (
      <span
        className={cn(
          "text-sm",
          row.original.payment === "Paid"
            ? "text-foreground"
            : "text-muted-foreground"
        )}
      >
        {row.original.payment}
      </span>
    ),
  },
  {
    id: "actions",
    header: "Action",
    enableSorting: false,
    cell: () => (
      <div className="flex items-center gap-1">
        <Button variant="ghost" size="icon-sm" aria-label="Edit">
          <Pencil className="size-4" />
        </Button>
        <Button variant="ghost" size="icon-sm" aria-label="Delete">
          <Trash2 className="size-4 text-destructive" />
        </Button>
        <Button variant="ghost" size="icon-sm" aria-label="More">
          <MoreHorizontal className="size-4" />
        </Button>
      </div>
    ),
  },
];

export default function DataTableDemoPage() {
  const [activeTab, setActiveTab] = React.useState<(typeof TABS)[number]>("All");

  return (
    <MainLayout breadcrumbs={[{ label: "Orders" }]}>
      <DataTable
        columns={columns}
        data={ORDERS}
        getRowId={(row) => row.id}
        toolbar={
          <div className="space-y-4">
            <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
              <h2 className="font-heading text-xl font-semibold">All Orders</h2>
              <div className="flex flex-wrap items-center gap-2">
                <Button variant="outline" size="sm">
                  <Download className="size-4" />
                  Export Orders
                </Button>
                <Button size="sm">
                  <Plus className="size-4" />
                  Add Orders
                </Button>
              </div>
            </div>

            <div className="flex items-center gap-1 border-b border-border">
              {TABS.map((tab) => (
                <button
                  key={tab}
                  type="button"
                  onClick={() => setActiveTab(tab)}
                  className={cn(
                    "cursor-pointer border-b-2 px-3 py-2 text-sm font-medium transition-colors",
                    activeTab === tab
                      ? "border-primary text-foreground"
                      : "border-transparent text-muted-foreground hover:text-foreground"
                  )}
                >
                  {tab}
                </button>
              ))}
            </div>
          </div>
        }
        bulkActions={(selected, clear) => (
          <>
            <Button variant="ghost" size="sm">
              <Copy className="size-4" />
              Duplicate
            </Button>
            <Button variant="ghost" size="sm">
              <Printer className="size-4" />
              Print
            </Button>
            <Button
              variant="ghost"
              size="sm"
              className="text-destructive hover:text-destructive"
              onClick={clear}
            >
              <Trash2 className="size-4" />
              Delete ({selected.length})
            </Button>
          </>
        )}
      />
    </MainLayout>
  );
}
