import { useState, useEffect, useMemo } from "react";
import type { ReactNode } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "~/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogFooter,
} from "~/components/ui/dialog";
import { useAuthStore } from "~/lib/store/auth-store";
import { updateReservationStatus, uploadReservationPhoto } from "~/lib/api-client";
import type { ReservationStatus } from "~/components/reservations/status-badge";
import { useLifecycleTransitions } from "~/lib/hooks/use-lifecycle-transitions";

interface StatusUpdateModalProps {
  xid: string;
  currentStatus: ReservationStatus | null;
  trigger?: ReactNode;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
}

const statusLabels: Record<string, string> = {
  pending: "Pending",
  sample_received: "Sample Received",
  test_on_progress: "Test On Progress",
  test_completed: "Test Completed",
  coa_submitted: "CoA Submitted",
  head_approved: "Head Approved",
  order_completed: "Order Completed",
  rejected: "Rejected",
};

export function StatusUpdateModal({
  xid,
  currentStatus,
  trigger,
  open: controlledOpen,
  onOpenChange,
}: StatusUpdateModalProps) {
  const hasPermission = useAuthStore((s) => s.hasPermission);
  const queryClient = useQueryClient();
  const { data: transitionMap = {} } = useLifecycleTransitions();
  const status = currentStatus ?? "pending";

  const [uncontrolledOpen, setUncontrolledOpen] = useState(false);
  const open = controlledOpen !== undefined ? controlledOpen : uncontrolledOpen;

  const handleOpenChange = (newOpen: boolean) => {
    if (onOpenChange) {
      onOpenChange(newOpen);
    } else {
      setUncontrolledOpen(newOpen);
    }
    if (!newOpen) {
      setSelectedStatus(status);
      setFile(null);
      setError(null);
    } else {
      setSelectedStatus(status);
    }
  };

  const setOpen = (newOpen: boolean) => {
    handleOpenChange(newOpen);
  };

  const [selectedStatus, setSelectedStatus] = useState<string>(status);
  const [file, setFile] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  // Sync state when modal is opened or status changes
  useEffect(() => {
    if (open) {
      setSelectedStatus(status);
      setFile(null);
      setError(null);
    }
  }, [open, status]);

  // Manage preview URL lifecycle to avoid memory leaks
  useEffect(() => {
    if (!file) {
      setPreviewUrl(null);
      return;
    }
    const url = URL.createObjectURL(file);
    setPreviewUrl(url);
    return () => {
      URL.revokeObjectURL(url);
    };
  }, [file]);

  const canUpdate = hasPermission("update_reservation_status");
  const transitions = useMemo(() => transitionMap[status] || [], [status]);
  const isTerminal = status === "order_completed" || status === "rejected";
  const showPhotoUpload = selectedStatus === "sample_received" && status !== "sample_received";

  const uploadMutation = useMutation({
    mutationFn: uploadReservationPhoto,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["reservation", xid] });
    },
  });

  const statusMutation = useMutation({
    mutationFn: (payload: { status: string; sample_photo?: string }) =>
      updateReservationStatus(xid, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["reservation", xid] });
      queryClient.invalidateQueries({ queryKey: ["reservations"] });
      queryClient.invalidateQueries({ queryKey: ["reservation-stats"] });
      setOpen(false);
      setFile(null);
      setError(null);
    },
    onError: (err: Error) => {
      setError(err.message || "Failed to update status");
    },
  });

  if (!canUpdate || isTerminal) {
    return null;
  }

  const handleSubmit = async () => {
    setError(null);
    if (!selectedStatus || selectedStatus === status) {
      setError("Please select a new status");
      return;
    }

    let samplePhoto: string | undefined;
    if (showPhotoUpload) {
      if (!file) {
        setError("Sample photo is required for this status");
        return;
      }
      try {
        const uploadRes = await uploadMutation.mutateAsync(file);
        samplePhoto = uploadRes.data.object_name;
      } catch {
        setError("Failed to upload sample photo");
        return;
      }
    }

    statusMutation.mutate({
      status: selectedStatus,
      ...(samplePhoto ? { sample_photo: samplePhoto } : {}),
    });
  };

  const isLoading = uploadMutation.isPending || statusMutation.isPending;

  // Make sure the current status is in the list so Select can display its label
  const displayOptions = useMemo(
    () => Array.from(new Set([status, ...transitions])),
    [status, transitions]
  );

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      {controlledOpen === undefined && (
        <DialogTrigger asChild>
          {trigger ? trigger : (
            <Button variant="outline" size="sm">
              Update Status
            </Button>
          )}
        </DialogTrigger>
      )}
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Update Reservation Status</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="space-y-2">
            <Label>Status</Label>
            <Select
              value={selectedStatus}
              onValueChange={(value) => {
                setSelectedStatus(value);
                setFile(null);
                setError(null);
              }}
            >
              <SelectTrigger>
                <SelectValue placeholder="Select a status" />
              </SelectTrigger>
              <SelectContent>
                {displayOptions.map((s) => (
                  <SelectItem key={s} value={s} disabled={s === status}>
                    {statusLabels[s] || s}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          {showPhotoUpload && (
            <div className="space-y-2">
              <Label>Sample Photo</Label>
              <Input
                type="file"
                accept="image/*"
                onChange={(e) => setFile(e.target.files?.[0] || null)}
              />
              {previewUrl && (
                <div className="mt-2 w-full max-w-[200px] aspect-video relative overflow-hidden rounded-md border">
                  <img
                    src={previewUrl}
                    alt="Preview"
                    className="object-cover w-full h-full"
                  />
                </div>
              )}
              <p className="text-xs text-muted-foreground">
                Required when changing status to Sample Received
              </p>
            </div>
          )}

          {error && (
            <p className="text-sm text-red-600">{error}</p>
          )}
        </div>
        <DialogFooter>
          <Button
            onClick={handleSubmit}
            disabled={isLoading || !selectedStatus || selectedStatus === status}
            className="w-full sm:w-auto"
          >
            {isLoading ? "Updating..." : "Update Status"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
