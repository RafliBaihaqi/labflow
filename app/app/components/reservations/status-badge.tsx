import { cn } from "~/lib/utils";

export type ReservationStatus =
  | "pending"
  | "sample_received"
  | "test_on_progress"
  | "test_completed"
  | "coa_submitted"
  | "head_approved"
  | "order_completed"
  | "rejected";

const statusConfig: Record<ReservationStatus, { bg: string; dot: string }> = {
  pending: { bg: "bg-gray-100 text-gray-700", dot: "bg-gray-500" },
  sample_received: { bg: "bg-amber-50 text-amber-700", dot: "bg-amber-500" },
  test_on_progress: {
    bg: "bg-purple-50 text-purple-700",
    dot: "bg-purple-500",
  },
  test_completed: {
    bg: "bg-emerald-50 text-emerald-700",
    dot: "bg-emerald-500",
  },
  coa_submitted: {
    bg: "bg-sky-50 text-sky-700",
    dot: "bg-sky-500",
  },
  head_approved: {
    bg: "bg-indigo-50 text-indigo-700",
    dot: "bg-indigo-500",
  },
  order_completed: {
    bg: "bg-emerald-50 text-emerald-700",
    dot: "bg-emerald-600",
  },
  rejected: {
    bg: "bg-red-50 text-red-700",
    dot: "bg-red-500",
  },
};

interface StatusBadgeProps {
  status: ReservationStatus;
  className?: string;
}

export function StatusBadge({ status, className }: StatusBadgeProps) {
  const config = statusConfig[status] || {
    bg: "bg-gray-100 text-gray-700",
    dot: "bg-gray-500",
  };
  return (
    <span
      className={cn(
        "inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium",
        config.bg,
        className,
      )}
    >
      <span
        aria-hidden="true"
        className={cn("size-1.5 rounded-full", config.dot)}
      />
      {status}
    </span>
  );
}
