export interface CreateReservationPayload {
  sample_name: string;
  sample_amount: number;
  sample_type: "Solid" | "Liquid" | "Powder";
  sample_residue: string;
  batch_number?: string;
  sample_taken_date: string;
  sample_taken_by: string;
  special_instructions: string;
  note?: string;
  tests: TestItemPayload[];
}

export interface TestItemPayload {
  test_list_id: number;
  method: string;
  test_unit: string;
  criteria: string;
}

export interface ReservationResponse {
  xid: string;
  sample_name: string;
  sample_amount: number;
  sample_type: string;
  sample_residue: string;
  batch_number: string | null;
  reservation_date: string;
  status:
    | "pending"
    | "sample_received"
    | "test_on_progress"
    | "test_completed"
    | "rejected";
  test_request_number: string;
  special_instructions: string;
  note: string | null;
  sample_taken_date: string;
  sample_taken_by: string;
  po_number: string | null;
  po_url: string | null;
  sample_photo: string | null;
  lead_time: string | null;
  tests: TestItemResponse[];
}

export interface TestItemResponse {
  id: number;
  xid: string;
  reservation_id: number;
  test_list_id: number;
  method: string;
  criteria: string;
  status: string | null;
  test_unit: string;
  test_result: string | null;
  test_summary: string | null;
  approved_date: string | null;
  notes: string | null;
  created_at: string | null;
  updated_at: string | null;
  test_list_name: string;
  category_name: string;
}

export interface TestCatalogItem {
  id: number;
  name: string;
  category: string;
}

export interface ReservationListItem {
  xid: string;
  test_request_number: string;
  sample_name: string;
  status:
    | "sample_received"
    | "test_on_progress"
    | "test_completed"
    | "rejected";
  reservation_date: string;
}

export interface ReservationListParams {
  page?: number;
  limit?: number;
  status?: string[];
  search?: string;
  sort_by?: string;
  sort_order?: "asc" | "desc";
}

export interface UpdateStatusPayload {
  status: string;
  sample_photo?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface ReservationStats {
  total: number;
  pending: number;
  sample_received: number;
  test_on_progress: number;
  test_completed: number;
  rejected: number;
}
