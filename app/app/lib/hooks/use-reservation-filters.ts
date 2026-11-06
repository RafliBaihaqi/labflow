import { useMemo, useCallback } from "react";
import { useSearchParams } from "react-router";
import type { SortingState } from "@tanstack/react-table";

/**
 * Hybrid URL + sessionStorage filters.
 * URL is primary source for page, search, status.
 * Sorting is kept in Zustand sessionStorage for UI prefs.
 */
export function useReservationFilters() {
  const [searchParams, setSearchParams] = useSearchParams();

  const page = useMemo(() => parseInt(searchParams.get("page") || "1", 10), [searchParams]);
  const search = useMemo(() => searchParams.get("search") || "", [searchParams]);
  const status = useMemo(() => searchParams.get("status") || "", [searchParams]);
  const sortBy = useMemo(() => searchParams.get("sort_by") || "", [searchParams]);
  const sortOrder = useMemo(() => searchParams.get("sort_order") || "", [searchParams]);

  const setPage = useCallback(
    (p: number) => {
      setSearchParams(
        (prev) => {
          const next = new URLSearchParams(prev);
          next.set("page", String(p));
          return next;
        },
        { replace: true }
      );
    },
    [setSearchParams]
  );

  const setSearch = useCallback(
    (s: string) => {
      setSearchParams(
        (prev) => {
          const next = new URLSearchParams(prev);
          if (s) next.set("search", s);
          else next.delete("search");
          next.set("page", "1");
          return next;
        },
        { replace: true }
      );
    },
    [setSearchParams]
  );

  const setStatus = useCallback(
    (newStatus: string) => {
      setSearchParams(
        (prev) => {
          const next = new URLSearchParams(prev);
          if (newStatus) next.set("status", newStatus);
          else next.delete("status");
          next.set("page", "1");
          return next;
        },
        { replace: true }
      );
    },
    [setSearchParams]
  );

  const setSorting = useCallback(
    (sorting: SortingState) => {
      setSearchParams(
        (prev) => {
          const next = new URLSearchParams(prev);
          if (sorting.length > 0) {
            const sort = sorting[0];
            next.set("sort_by", sort.id);
            next.set("sort_order", sort.desc ? "desc" : "asc");
          } else {
            next.delete("sort_by");
            next.delete("sort_order");
          }
          return next;
        },
        { replace: true }
      );
    },
    [setSearchParams]
  );

  return {
    page,
    search,
    status,
    sortBy,
    sortOrder,
    setPage,
    setSearch,
    setStatus,
    setSorting,
  };
}
