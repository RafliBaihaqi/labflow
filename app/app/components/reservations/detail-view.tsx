import { ArrowRightLeft, Edit } from "lucide-react";
import type { ColumnDef } from "@tanstack/react-table";
import type {ReservationResponse, TestItemResponse} from "~/lib/types/reservation";
import {StatusBadge} from "~/components/reservations/status-badge";
import { useAuthStore } from "~/lib/store/auth-store";
import { useUIStore } from "~/lib/store/ui-store";
import { Button } from "~/components/ui/button";
import { Dialog, DialogContent, DialogTrigger } from "~/components/ui/dialog";
import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
} from "~/components/ui/card";
import { DataTable } from "~/components/data-table/data-table";
import {formatDate} from "~/lib/date-formatter";

interface DetailViewProps {
    reservation: ReservationResponse;
}

function renderValue(value: string | number | Date | null | undefined): string {
    if (value === null || value === undefined || value === "") return "—";
    return String(value);
}

function DetailField({
                         label,
                         children,
                     }: {
    label: string;
    children: React.ReactNode;
}) {
    return (
        <div className="space-y-1">
            <p className="text-xs font-medium text-muted-foreground">{label}</p>
            <div className="text-sm">{children}</div>
        </div>
    );
}

export function DetailView({reservation}: DetailViewProps) {
    const user = useAuthStore((s) => s.user);
    const hasPermission = useAuthStore((s) => s.hasPermission);
    const openModal = useUIStore((s) => s.openModal);
    
    const canUpdateStatus = hasPermission("update_reservation_status") &&
        reservation.status !== "test_completed" &&
        reservation.status !== "rejected";
    
    const canEditTest = (test: TestItemResponse) => {
        if (reservation.status !== "test_on_progress") return false;
        if (!hasPermission("update_test")) return false;
        
        const roles = user?.roles || [];
        if (roles.includes("admin")) return true;
        
        const isMicroTest = test.category_name === "Microbiology";
        const isChemTest = test.category_name === "Chemistry";
        
        if (isMicroTest && (roles.includes("micro-analyst") || roles.includes("micro-qc"))) return true;
        if (isChemTest && (roles.includes("chemical-analyst") || roles.includes("chemical-qc"))) return true;
        
        return false;
    };

    const testColumns: ColumnDef<TestItemResponse>[] = [
        {
            accessorKey: "test_list_name",
            header: "Test Type",
            cell: ({ row }) => (
                <span className="font-medium">{renderValue(row.original.test_list_name)}</span>
            ),
        },
        {
            accessorKey: "method",
            header: "Method",
            cell: ({ row }) => renderValue(row.original.method),
        },
        {
            accessorKey: "test_unit",
            header: "Unit",
            cell: ({ row }) => (
                <span className="tabular-nums">{renderValue(row.original.test_unit)}</span>
            ),
        },
        {
            accessorKey: "criteria",
            header: "Criteria",
            cell: ({ row }) => (
                <span className="tabular-nums">{renderValue(row.original.criteria)}</span>
            ),
        },
        {
            accessorKey: "test_result",
            header: "Result",
            cell: ({ row }) => (
                <span className="tabular-nums">{renderValue(row.original.test_result)}</span>
            ),
        },
        {
            accessorKey: "test_summary",
            header: "Comply",
            cell: ({ row }) => renderValue(row.original.test_summary),
        },
        {
            accessorKey: "notes",
            header: "Notes",
            cell: ({ row }) => (
                <span className="text-pretty">{renderValue(row.original.notes)}</span>
            ),
        },
        {
            id: "actions",
            cell: ({ row }) => {
                if (!canEditTest(row.original)) return null;
                return (
                    <Button
                        variant="ghost"
                        size="sm"
                        className="h-8 w-8 p-0"
                        onClick={() =>
                            openModal("edit-test", {
                                test: row.original,
                                reservationXid: reservation.xid,
                            })
                        }
                    >
                        <span className="sr-only">Edit</span>
                        <Edit className="h-4 w-4" />
                    </Button>
                );
            },
        },
    ];

    return (
        <div className="space-y-4">
            <div className="flex items-center justify-between">
                <div>
                    <h2 className="text-balance text-lg font-semibold">Reservation Details</h2>
                    <p className="text-pretty text-sm text-muted-foreground">View and manage this reservation.</p>
                </div>
                {canUpdateStatus && (
                    <Button onClick={() => openModal("status-update", { xid: reservation.xid, currentStatus: reservation.status })}>
                        <ArrowRightLeft className="mr-2 size-4" />
                        Update Status
                    </Button>
                )}
            </div>
            <Card>
                <CardHeader>
                    <CardTitle className="text-lg font-semibold">Sample Information</CardTitle>
                </CardHeader>
                <CardContent className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                    <DetailField label="Sample Name">
                        {renderValue(reservation.sample_name)}
                    </DetailField>
                    <DetailField label="Sample Amount">
                        <span className="tabular-nums">{renderValue(reservation.sample_amount)}</span>
                    </DetailField>
                    <DetailField label="Sample Type">
                        {renderValue(reservation.sample_type)}
                    </DetailField>
                    <DetailField label="Sample Residue">
                        {renderValue(reservation.sample_residue)}
                    </DetailField>
                    <DetailField label="Batch Number">
                        <span className="tabular-nums">{renderValue(reservation.batch_number)}</span>
                    </DetailField>
                    <DetailField label="Status">
                        <StatusBadge status={reservation.status}/>
                    </DetailField>
                    {reservation.sample_photo && (
                        <DetailField label="Sample Photo">
                            <Dialog>
                                <DialogTrigger asChild>
                                    <button
                                        type="button"
                                        aria-label="View large sample photo"
                                        className="mt-2 block w-full max-w-[200px] aspect-video relative overflow-hidden rounded-md border cursor-pointer hover:opacity-80 transition-opacity focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                                    >
                                        <img
                                            src={reservation.sample_photo}
                                            alt="Sample thumbnail"
                                            className="object-cover w-full h-full"
                                        />
                                    </button>
                                </DialogTrigger>
                                <DialogContent className="max-w-4xl p-0 overflow-hidden bg-transparent border-none shadow-none ring-0">
                                    <img
                                        src={reservation.sample_photo}
                                        alt="Sample expanded view"
                                        className="w-full h-auto max-h-[85vh] object-contain rounded-md"
                                    />
                                </DialogContent>
                            </Dialog>
                        </DetailField>
                    )}
                </CardContent>
            </Card>

            <Card>
                <CardHeader>
                    <CardTitle>Logistics &amp; Documents</CardTitle>
                </CardHeader>
                <CardContent className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                    <DetailField label="Test Request Number">
                        <span className="tabular-nums">{renderValue(reservation.test_request_number)}</span>
                    </DetailField>
                    <DetailField label="Reservation Date">
                        <span className="tabular-nums">{formatDate(reservation.reservation_date)}</span>
                    </DetailField>
                    <DetailField label="Sample Taken Date">
                        <span className="tabular-nums">{formatDate(reservation.sample_taken_date)}</span>
                    </DetailField>
                    <DetailField label="Sample Taken By">
                        {renderValue(reservation.sample_taken_by)}
                    </DetailField>
                    <DetailField label="Lead Time">
                        <span className="tabular-nums">{renderValue(reservation.lead_time)}</span>
                    </DetailField>
                    <DetailField label="PO Number">
                        <span className="tabular-nums">{renderValue(reservation.po_number)}</span>
                    </DetailField>
                    <DetailField label="PO URL">
                        {reservation.po_url ? (
                            <a
                                href={reservation.po_url}
                                target="_blank"
                                rel="noopener noreferrer"
                                className="text-primary underline-offset-4 hover:underline break-all"
                            >
                                {reservation.po_url}
                            </a>
                        ) : (
                            "—"
                        )}
                    </DetailField>
                </CardContent>
            </Card>

            <Card>
                <CardHeader>
                    <CardTitle>Notes</CardTitle>
                </CardHeader>
                <CardContent className="grid gap-4 sm:grid-cols-2">
                    <DetailField label="Special Instructions">
                        <p className="text-pretty">{renderValue(reservation.special_instructions)}</p>
                    </DetailField>
                    <DetailField label="Note">
                        <p className="text-pretty">{renderValue(reservation.note)}</p>
                    </DetailField>
                </CardContent>
            </Card>

            <Card>
                <CardHeader>
                    <CardTitle>Test Items</CardTitle>
                </CardHeader>
                <CardContent>
                    <DataTable
                        columns={testColumns}
                        data={reservation.tests}
                        getRowId={(row) => row.xid}
                        hidePagination
                        emptyMessage="No test items"
                    />
                </CardContent>
            </Card>
        </div>
    );
}
