import { format, isValid, parseISO } from "date-fns";

export function formatDate(value: string | null | undefined): string {
  if (!value) return "—";
  const date = parseISO(value);
  if (!isValid(date)) return "—";
  return format(date, "MMM d, yyyy");
}

export function formatDateTime(value: string | null | undefined): string {
  if (!value) return "—";
  const date = parseISO(value);
  if (!isValid(date)) return "—";
  return format(date, "MMM d, yyyy, h:mm a");
}
