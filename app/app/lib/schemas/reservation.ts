import * as z from "zod";

export const testItemSchema = z.object({
  test_list_id: z.number().positive(),
  method: z.string().min(1, "Method is required").max(255),
  test_unit: z.string().min(1, "Unit is required").max(255),
  criteria: z.string().min(1, "Acceptance criteria is required").max(255),
});

export const reservationSchema = z.object({
  sample_name: z.string().min(1, "Sample name is required").max(150),
  sample_amount: z
    .number({ message: "Must be a number" })
    .min(0.01, "Must be at least 0.01")
    .max(999_999.99, "Must be at most 999,999.99"),
  sample_type: z.enum(["Solid", "Liquid", "Powder"], {
    message: "Sample type is required",
  }),
  sample_residue: z.string().min(1, "Sample residue is required").max(150),
  batch_number: z.string().max(50).optional().or(z.literal("")),
  sample_taken_date: z
    .date({ message: "Sample taken date is required" })
    .max(new Date(), "Date cannot be in the future"),
  sample_taken_by: z.string().min(1, "Sample taken by is required").max(100),
  special_instructions: z
    .string()
    .min(1, "Special instructions is required")
    .max(500),
  note: z.string().max(1000).optional().or(z.literal("")),
  tests: z.array(testItemSchema).min(1, "At least one test is required"),
});

export type ReservationFormData = z.infer<typeof reservationSchema>;
