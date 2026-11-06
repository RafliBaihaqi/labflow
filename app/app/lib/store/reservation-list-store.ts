import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";
import type { SortingState, Updater } from "@tanstack/react-table";

interface ReservationListState {
  // page/search are now URL primary, but kept for backward compat/session cache
  page: number;
  search: string;
  activeStatus: string;
  sorting: SortingState;
  setPage: (page: number) => void;
  setSearch: (search: string) => void;
  setActiveStatus: (status: string) => void;
  setSorting: (sorting: Updater<SortingState>) => void;
  resetFilters: () => void;
}

function safeSessionStorage() {
  if (typeof window === "undefined") return undefined as any;
  try {
    return sessionStorage;
  } catch {
    return undefined as any;
  }
}

export const useReservationListStore = create<ReservationListState>()(
  persist(
    (set) => ({
      page: 1,
      search: "",
      activeStatus: "",
      sorting: [],
      setPage: (page) => set({ page }),
      setSearch: (search) => set({ search, page: 1 }),
      setActiveStatus: (activeStatus) => set({ activeStatus, page: 1 }),
      setSorting: (updaterOrValue) =>
        set((state) => {
          const sorting =
            typeof updaterOrValue === "function"
              ? (updaterOrValue as (old: SortingState) => SortingState)(state.sorting)
              : updaterOrValue;
          return { sorting };
        }),
      resetFilters: () => set({ page: 1, search: "", activeStatus: "", sorting: [] }),
    }),
    {
      name: "reservation-list-state",
      version: 1,
      storage: createJSONStorage(() => safeSessionStorage() as any),
      // Only persist UI prefs, not per-keystroke search
      partialize: (state) => ({
        activeStatus: state.activeStatus,
        sorting: state.sorting,
      }),
    }
  )
);
