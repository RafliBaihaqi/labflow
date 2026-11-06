import axios from "axios";
import type {
    CreateReservationPayload,
    PaginatedResponse,
    ReservationListItem,
    ReservationListParams,
    ReservationResponse,
    ReservationStats,
    TestCatalogItem,
    UpdateStatusPayload,
} from "~/lib/types/reservation";
import type {
    User,
    UpdateProfilePayload,
    UpdateProfileResponse,
    RegisterPayload,
    RegisterResponse,
    VerifyOtpPayload,
    ResendOtpPayload,
    MessageResponse,
} from "~/lib/types/auth";
import type { QuotationInput, QuotationItem, QuotationStatus } from "~/lib/types/quotation";
import type { TestListInput, TestListItem } from "~/lib/types/test-list";

interface loginRequest {
    identifier: string;
    password: string;
    turnstileToken: string;
}

const API_BASE_URL = (import.meta as any).env?.VITE_API_BASE_URL || "";

const instance = axios.create({
    baseURL: API_BASE_URL,
    timeout: 10_000,
    withCredentials: true,
});

// Centralized error normalization + 401 handling
instance.interceptors.response.use(
    (response) => response,
    (error) => {
        if (axios.isAxiosError(error) && error.response?.status === 401) {
            // Clear auth state and redirect to login (avoid loop if already on login)
            if (typeof window !== "undefined") {
                const path = window.location.pathname;
                if (!path.includes("/login") && !path.includes("/register")) {
                    // Dynamic import to avoid circular deps
                    import("~/lib/store/auth-store").then(({ useAuthStore }) => {
                        try {
                            useAuthStore.getState().clearUser();
                        } catch {
                            // ignore
                        }
                    });
                }
            }
        }
        return Promise.reject(error);
    }
);

function normalizeErrorMessage(error: unknown): string {
    if (axios.isAxiosError(error)) {
        return (
            (error.response?.data as any)?.error ||
            (error.response?.data as any)?.message ||
            error.message ||
            "An unexpected error occurred"
        );
    }
    return error instanceof Error ? error.message : "An unexpected error occurred";
}

export { normalizeErrorMessage };

export const login = async (data: loginRequest) => {
    return await instance.post("/auth/login", data);
};

export const logout = async () => {
    return await instance.post("/auth/logout");
};

export const me = async () => {
    const res = await instance.get<User>("/api/auth/me");
    return res.data;

}

export const register = async (data: RegisterPayload) => {
    return await instance.post<RegisterResponse>("/auth/register", data);
};

export const verifyOtp = async (data: VerifyOtpPayload) => {
    return await instance.post<MessageResponse>("/auth/verify-otp", data);
};

export const resendOtp = async (data: ResendOtpPayload) => {
    return await instance.post<MessageResponse>("/auth/resend-otp", data);
};

export const getReservations = async (params: ReservationListParams = {}) => {
    const searchParams = new URLSearchParams();
    if (params.page) searchParams.set("page", String(params.page));
    if (params.limit) searchParams.set("limit", String(params.limit));
    if (params.status?.length) {
        params.status.forEach((s) => searchParams.append("status", s));
    }
    if (params.search) searchParams.set("search", params.search);
    if (params.sort_by) searchParams.set("sort_by", params.sort_by);
    if (params.sort_order) searchParams.set("sort_order", params.sort_order);

    return await instance.get<PaginatedResponse<ReservationListItem>>(
        `/api/reservations?${searchParams.toString()}`,
        {timeout: 10_000, withCredentials: true},
    );
};

export const getReservationStats = async () => {
    return await instance.get<ReservationStats>("/api/reservations/stats", {
        timeout: 10_000,
        withCredentials: true,
    });
};

export const createReservation = async (data: CreateReservationPayload) => {
    return await instance.post<{ xid: string }>("/api/reservations", data, {
        timeout: 30_000,
    });
};

export const getReservation = async (xid: string) => {
    return await instance.get<ReservationResponse>(`/api/reservations/${xid}`, {
        timeout: 10_000,
    });
};

export const getTestCatalog = async (testGroupId: number) => {
    return await instance.get<TestCatalogItem[]>(
        `/api/tests/catalog/${testGroupId}`,
    );
};

export const updateReservationStatus = async (xid: string, data: UpdateStatusPayload) => {
    return await instance.put<{ message: string }>(`/api/reservations/${xid}/status`, data);
};

export const uploadReservationPhoto = async (file: File) => {
    const formData = new FormData();
    formData.append("file", file);
    return await instance.post<{ message: string; object_name: string }>(
        "/api/uploads/reservation",
        formData,
        {
            headers: {
                "Content-Type": "multipart/form-data",
            },
            timeout: 30_000,
        }
    );
};

export const uploadUserAvatar = async (file: File) => {
    const formData = new FormData();
    formData.append("file", file);
    return await instance.post<{ message: string; object_name: string }>(
        "/api/uploads/user",
        formData,
        {
            headers: {
                "Content-Type": "multipart/form-data",
            },
            timeout: 30_000,
        }
    );
};

export const updateTest = async (xid: string, data: { test_result: string | null; test_summary: string | null; notes: string | null }) => {
    return await instance.patch<{ message: string }>(`/api/tests/${xid}`, data);
};

export const updateProfile = async (xid: string, data: UpdateProfilePayload) => {
    return await instance.patch<UpdateProfileResponse>(`/api/users/${xid}`, data);
};

export const submitQuotation = async (data: QuotationInput) => {
    return await instance.post("/api/quotations", data);
};

export const listQuotations = async (params: { page?: number; limit?: number; search?: string } = {}) => {
    const searchParams = new URLSearchParams();
    if (params.page) searchParams.set("page", String(params.page));
    if (params.limit) searchParams.set("limit", String(params.limit));
    if (params.search) searchParams.set("search", params.search);
    return await instance.get<PaginatedResponse<QuotationItem>>(`/api/quotations?${searchParams.toString()}`);
};

export const updateQuotationStatus = async (xid: string, status: QuotationStatus) => {
    return await instance.patch<QuotationItem>(`/api/quotations/${xid}`, { status });
};

export const listTestLists = async (params: { page?: number; limit?: number; search?: string; category_id?: number } = {}) => {
    const searchParams = new URLSearchParams();
    if (params.page) searchParams.set("page", String(params.page));
    if (params.limit) searchParams.set("limit", String(params.limit));
    if (params.search) searchParams.set("search", params.search);
    if (params.category_id) searchParams.set("category_id", String(params.category_id));
    return await instance.get<PaginatedResponse<TestListItem>>(`/api/test-lists?${searchParams.toString()}`);
};

export const createTestList = async (data: TestListInput) => {
    return await instance.post<TestListItem>("/api/test-lists", data);
};

export const updateTestList = async (xid: string, data: TestListInput) => {
    return await instance.put<TestListItem>(`/api/test-lists/${xid}`, data);
};

export const deactivateTestList = async (xid: string) => {
    return await instance.delete<TestListItem>(`/api/test-lists/${xid}`);
};

export const forgotPassword = async (data: { email: string }) => {
    return await instance.post<MessageResponse>("/auth/forgot-password", data);
};

export const verifyResetOtp = async (data: { email: string; otp: string }) => {
    return await instance.post<{ reset_token: string }>("/auth/verify-reset-otp", data);
};

export const resetPassword = async (data: { token: string; password: string }) => {
    return await instance.post<MessageResponse>("/auth/reset-password", data);
};

type LifecycleTransitionsResponse = {
    transitions: Record<string, string[]>;
};

export const getLifecycleTransitions = async () => {
    try {
        const res = await instance.get<LifecycleTransitionsResponse>("/api/lifecycle/transitions");
        return res.data.transitions;
    } catch {
        // Fallback to extended map matching backend lifecycle.go until endpoint exists
        return {
            pending: ["sample_received", "rejected"],
            sample_received: ["test_on_progress", "rejected"],
            test_on_progress: ["test_completed", "rejected"],
            test_completed: ["coa_submitted", "rejected"],
            coa_submitted: ["head_approved", "rejected"],
            head_approved: ["order_completed", "rejected"],
            order_completed: [],
            rejected: [],
        } as Record<string, string[]>;
    }
};
