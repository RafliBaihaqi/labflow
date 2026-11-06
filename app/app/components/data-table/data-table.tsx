import * as React from "react";
import {
  type ColumnDef,
  type OnChangeFn,
  type RowSelectionState,
  type SortingState,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { ChevronLeft, ChevronRight, X } from "lucide-react";

import { cn } from "~/lib/utils";
import { Button } from "../ui/button";
import { Checkbox } from "../ui/checkbox";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../ui/table";

/** Prebuilt selection column (checkbox header + per-row checkbox). */
export function selectColumn<TData>(): ColumnDef<TData> {
  return {
    id: "select",
    header: ({ table }) => (
      <Checkbox
        checked={
          table.getIsAllPageRowsSelected()
            ? true
            : table.getIsSomePageRowsSelected()
              ? "indeterminate"
              : false
        }
        onCheckedChange={(value) =>
          table.toggleAllPageRowsSelected(!!value)
        }
        aria-label="Select all"
      />
    ),
    cell: ({ row }) => (
      <span onClick={(e) => e.stopPropagation()}>
        <Checkbox
          checked={row.getIsSelected()}
          onCheckedChange={(value) => row.toggleSelected(!!value)}
          aria-label="Select row"
        />
      </span>
    ),
    enableSorting: false,
    enableHiding: false,
    size: 40,
  };
}

interface ServerPagination {
  /** Current page index (0-based). */
  pageIndex: number;
  /** Total number of pages from the server. */
  pageCount: number;
  /** Rows per page. */
  pageSize: number;
  /** Total number of rows across all pages (for the entries count). */
  totalRows: number;
  /** Called with the next 0-based page index. */
  onPageChange: (pageIndex: number) => void;
}

interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
  /** Title + action buttons rendered above the table. */
  toolbar?: React.ReactNode;
  /** Rows per page for client-side pagination. Ignored when serverPagination is set. Defaults to 10. */
  pageSize?: number;
  /** Stable row id accessor (recommended for selection). */
  getRowId?: (row: TData, index: number) => string;
  /** Floating bar content shown when rows are selected. */
  bulkActions?: (selectedRows: TData[], clearSelection: () => void) => React.ReactNode;
  /** Message shown when there are no rows. */
  emptyMessage?: string;
  /** Makes rows clickable. The checkbox/action cells stop propagation. */
  onRowClick?: (row: TData) => void;
  /** Enables server-side (manual) pagination. When set, `data` should be the current page only. */
  serverPagination?: ServerPagination;
  /** Enables server-side (manual) sorting with controlled sorting state. */
  serverSorting?: {
    sorting: SortingState;
    onSortingChange: OnChangeFn<SortingState>;
  };
  /** Hides the entries count + pagination footer (useful for small, embedded tables). */
  hidePagination?: boolean;
}

function pageRange(current: number, total: number): (number | "ellipsis")[] {
  if (total <= 7) {
    return Array.from({ length: total }, (_, i) => i + 1);
  }
  const pages: (number | "ellipsis")[] = [1];
  const start = Math.max(2, current - 1);
  const end = Math.min(total - 1, current + 1);
  if (start > 2) pages.push("ellipsis");
  for (let i = start; i <= end; i++) pages.push(i);
  if (end < total - 1) pages.push("ellipsis");
  pages.push(total);
  return pages;
}

export function DataTable<TData, TValue>({
  columns,
  data,
  toolbar,
  pageSize = 10,
  getRowId,
  bulkActions,
  emptyMessage = "No results.",
  onRowClick,
  serverPagination,
  serverSorting,
  hidePagination,
}: DataTableProps<TData, TValue>) {
  const [internalSorting, setInternalSorting] = React.useState<SortingState>([]);
  const [rowSelection, setRowSelection] = React.useState<RowSelectionState>({});

  const isServer = !!serverPagination;
  const sorting = serverSorting ? serverSorting.sorting : internalSorting;
  const onSortingChange = serverSorting
    ? serverSorting.onSortingChange
    : setInternalSorting;

  const table = useReactTable({
    data,
    columns,
    state: {
      sorting,
      rowSelection,
      ...(serverPagination
        ? {
            pagination: {
              pageIndex: serverPagination.pageIndex,
              pageSize: serverPagination.pageSize,
            },
          }
        : {}),
    },
    getRowId,
    onSortingChange,
    onRowSelectionChange: setRowSelection,
    enableRowSelection: true,
    manualSorting: !!serverSorting,
    manualPagination: isServer,
    pageCount: serverPagination?.pageCount,
    onPaginationChange: serverPagination
      ? (updater) => {
          const prev = {
            pageIndex: serverPagination.pageIndex,
            pageSize: serverPagination.pageSize,
          };
          const next = typeof updater === "function" ? updater(prev) : updater;
          serverPagination.onPageChange(next.pageIndex);
        }
      : undefined,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    ...(serverSorting ? {} : { getSortedRowModel: getSortedRowModel() }),
    ...(isServer
      ? {}
      : {
          getPaginationRowModel: getPaginationRowModel(),
          initialState: {
            pagination: {
              pageSize: hidePagination ? Math.max(data.length, 1) : pageSize,
            },
          },
        }),
  });

  const currentPageRows = table.getRowModel().rows.length;
  const totalRows = isServer
    ? serverPagination.totalRows
    : table.getFilteredRowModel().rows.length;
  const pageIndex = isServer
    ? serverPagination.pageIndex
    : table.getState().pagination.pageIndex;
  const currentPageSize = isServer
    ? serverPagination.pageSize
    : table.getState().pagination.pageSize;
  const firstRow = totalRows === 0 ? 0 : pageIndex * currentPageSize + 1;
  const lastRow = pageIndex * currentPageSize + currentPageRows;
  const pageCount = isServer ? serverPagination.pageCount : table.getPageCount();

  const selectedRows = table
    .getFilteredSelectedRowModel()
    .rows.map((r) => r.original);
  const selectedCount = selectedRows.length;

  return (
    <div className="space-y-4">
      {toolbar}

      <div className="overflow-hidden bg-card">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow
                key={headerGroup.id}
                className="border-border bg-muted/40 hover:bg-muted/40"
              >
                {headerGroup.headers.map((header) => (
                  <TableHead
                    key={header.id}
                    className="h-11 text-xs font-medium text-muted-foreground"
                  >
                    {header.isPlaceholder
                      ? null
                      : flexRender(
                          header.column.columnDef.header,
                          header.getContext()
                        )}
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && "selected"}
                  onClick={onRowClick ? () => onRowClick(row.original) : undefined}
                  className={cn(
                    "border-border data-[state=selected]:bg-primary/5",
                    onRowClick && "cursor-pointer"
                  )}
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id} className="py-3">
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center text-sm text-muted-foreground"
                >
                  {emptyMessage}
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      {/* Footer: entries count + pagination */}
      {!hidePagination && (
        <div className="flex flex-col items-center justify-between gap-3 sm:flex-row">
        <p className="text-sm text-muted-foreground">
          Showing <span className="font-medium text-foreground">{firstRow}</span>
          {"-"}
          <span className="font-medium text-foreground">{lastRow}</span> of{" "}
          <span className="font-medium text-foreground">{totalRows}</span>{" "}
          entries
        </p>

        <div className="flex items-center gap-1">
          <Button
            variant="outline"
            size="sm"
            onClick={() => table.previousPage()}
            disabled={!table.getCanPreviousPage()}
          >
            <ChevronLeft className="size-4" />
            Previous
          </Button>

          {pageRange(pageIndex + 1, pageCount).map((p, i) =>
            p === "ellipsis" ? (
              <span
                key={`ellipsis-${i}`}
                className="px-2 text-sm text-muted-foreground"
              >
                …
              </span>
            ) : (
              <Button
                key={p}
                variant={p === pageIndex + 1 ? "default" : "ghost"}
                size="icon-sm"
                className="size-8"
                onClick={() => table.setPageIndex(p - 1)}
              >
                {p}
              </Button>
            )
          )}

          <Button
            variant="outline"
            size="sm"
            onClick={() => table.nextPage()}
            disabled={!table.getCanNextPage()}
          >
            Next
            <ChevronRight className="size-4" />
          </Button>
        </div>
      </div>
      )}

      {/* Floating bulk-action bar */}
      {bulkActions && selectedCount > 0 && (
        <div
          className={cn(
            "fixed inset-x-0 bottom-6 z-50 mx-auto flex w-fit items-center gap-2 rounded-xl border border-border bg-card px-3 py-2 shadow-lg",
            "animate-in fade-in slide-in-from-bottom-4"
          )}
        >
          <span className="flex items-center gap-2 rounded-lg bg-muted px-2.5 py-1 text-sm font-medium">
            <span className="flex size-5 items-center justify-center rounded-full bg-primary text-xs text-primary-foreground">
              {selectedCount}
            </span>
            Selected
          </span>
          {bulkActions(selectedRows, () => table.resetRowSelection())}
          <Button
            variant="ghost"
            size="icon-sm"
            onClick={() => table.resetRowSelection()}
            aria-label="Clear selection"
          >
            <X className="size-4" />
          </Button>
        </div>
      )}
    </div>
  );
}
