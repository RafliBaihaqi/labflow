export interface User {
    id: number;
    xid: string;
    username: string;
    email: string;
    name: string;
    phone_number?: string;
    company?: string;
    address?: string;
    avatar?: string;
    roles: string[];
    permissions: string[];
}

export interface UpdateProfilePayload {
    company?: string;
    phone_number?: string;
    address?: string;
    avatar?: string;
}

export interface UpdateProfileResponse {
    xid: string;
    phone_number?: string;
    company?: string;
    address?: string;
    avatar?: string | null;
    message: string;
}

export interface RegisterPayload {
    name: string;
    username: string;
    email: string;
    password: string;
    phone_number?: string;
    company?: string;
    address?: string;
}

export interface RegisterResponse {
    id: number;
    xid: string;
    username: string;
    email: string;
    name: string;
    phone_number?: string;
    company?: string;
    address?: string;
    created_at: string;
    message: string;
}

export interface VerifyOtpPayload {
    email: string;
    otp: string;
}

export interface ResendOtpPayload {
    email: string;
}

export interface MessageResponse {
    message: string;
}
