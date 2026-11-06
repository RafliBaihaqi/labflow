export interface QuotationInput {
  email: string;
  company_name: string;
  phone_number?: string;
  message?: string;
}

export type QuotationStatus = "pending" | "contacted" | "closed";

export interface QuotationItem extends QuotationInput {
  xid: string;
  status?: QuotationStatus;
}
