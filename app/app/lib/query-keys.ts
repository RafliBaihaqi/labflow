export const queryKeys = {
  me: ["me"] as const,
  reservations: (filters: Record<string, unknown> = {}) =>
    ["reservations", filters] as const,
  reservation: (xid: string) => ["reservation", xid] as const,
  reservationStats: ["reservation-stats"] as const,
  testCatalog: (categoryId: string | number) =>
    ["test-catalog", String(categoryId)] as const,
  lifecycleTransitions: ["lifecycle-transitions"] as const,
  certificates: (filters: Record<string, unknown> = {}) =>
    ["certificates", filters] as const,
  certificateCounts: ["certificate-counts"] as const,
  feedbackRequired: ["feedback-required"] as const,
  feedbackList: (filters: Record<string, unknown> = {}) =>
    ["feedback-list", filters] as const,
  quotations: (filters: Record<string, unknown> = {}) =>
    ["quotations", filters] as const,
  testLists: (filters: Record<string, unknown> = {}) =>
    ["test-lists", filters] as const,
  chatRooms: ["chat-rooms"] as const,
  chatMessages: (roomId: string) => ["chat-messages", roomId] as const,
  dashboardAdmin: ["dashboard-admin"] as const,
  dashboardCustomer: ["dashboard-customer"] as const,
  customers: (filters: Record<string, unknown> = {}) =>
    ["customers", filters] as const,
} as const;
