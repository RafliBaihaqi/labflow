import { useState } from "react";
import { useForm, useFieldArray, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "react-router";
import { Loader2 } from "lucide-react";
import axios from "axios";
import { toast } from "sonner";

import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "~/components/ui/field";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import { TestSelector } from "~/components/reservations/test-selector";
import { createReservation } from "~/lib/api-client";
import {
  reservationSchema,
  type ReservationFormData,
} from "~/lib/schemas/reservation";
import { DatePickerSimple } from "./date-picker";

export function CreationForm() {
  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
  } = useForm<ReservationFormData>({
    resolver: zodResolver(reservationSchema),
    defaultValues: { tests: [] },
  });

  const { fields, append, remove } = useFieldArray({
    control,
    name: "tests",
  });

  const navigate = useNavigate();
  const [backendErrors, setBackendErrors] = useState<string[]>([]);

  const mutation = useMutation({
    mutationFn: createReservation,
    onSuccess: (res) => {
      toast.success("Reservation created successfully");
      navigate(`/reservations/${res.data.xid}`);
    },
    onError: (error) => {
      if (axios.isAxiosError(error)) {
        const status = error.response?.status;
        if (status === 400) {
          const data = error.response?.data;
          if (data?.errors && Array.isArray(data.errors)) {
            setBackendErrors(data.errors);
          } else if (data?.message) {
            setBackendErrors([data.message]);
          } else {
            setBackendErrors(["Validation failed. Please check your input."]);
          }
        } else {
          toast.error("An unexpected error occurred. Please try again later.");
        }
      } else {
        toast.error("Network error. Please check your connection and try again.");
      }
    },
  });

  const onSubmit = (data: ReservationFormData) => {
    setBackendErrors([]);
    const payload = {
      ...data,
      sample_taken_date: data.sample_taken_date.toISOString(),
    };
    mutation.mutate(payload);
  };

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      className="grid gap-4 lg:grid-cols-[minmax(0,1fr)_minmax(0,1.5fr)]"
    >
      {/* Left column: Sample fields */}
      <div className="flex flex-col gap-4">
        <Card>
          <CardHeader>
            <CardTitle className="text-balance">Sample Information</CardTitle>
            <CardDescription className="text-pretty">
              Provide sample details for this reservation.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <FieldGroup>
              <Field>
                <FieldLabel htmlFor="sample_name">
                  Sample Name <span className="text-destructive">*</span>
                </FieldLabel>
                <Input
                  id="sample_name"
                  type="text"
                  maxLength={150}
                  aria-invalid={!!errors.sample_name}
                  {...register("sample_name")}
                />
                {errors.sample_name && (
                  <FieldError>{errors.sample_name.message}</FieldError>
                )}
              </Field>

              <Field>
                <FieldLabel htmlFor="sample_amount">
                  Sample Amount <span className="text-destructive">*</span>
                </FieldLabel>
                <Input
                  id="sample_amount"
                  type="number"
                  step="0.01"
                  min="0.01"
                  max="999999.99"
                  aria-invalid={!!errors.sample_amount}
                  {...register("sample_amount", { valueAsNumber: true })}
                />
                {errors.sample_amount && (
                  <FieldError>{errors.sample_amount.message}</FieldError>
                )}
              </Field>

              <Field>
                <FieldLabel htmlFor="sample_type">
                  Sample Type <span className="text-destructive">*</span>
                </FieldLabel>
                <select
                  id="sample_type"
                  aria-invalid={!!errors.sample_type}
                  className="h-9 w-full min-w-0 rounded-lg border border-input bg-transparent px-3 py-1.5 text-sm transition-colors outline-none focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/30"
                  {...register("sample_type")}
                >
                  <option value="">Select sample type...</option>
                  <option value="Solid">Solid</option>
                  <option value="Liquid">Liquid</option>
                  <option value="Powder">Powder</option>
                </select>
                {errors.sample_type && (
                  <FieldError>{errors.sample_type.message}</FieldError>
                )}
              </Field>

              <Field>
                <FieldLabel htmlFor="sample_residue">
                  Sample Residue <span className="text-destructive">*</span>
                </FieldLabel>
                <Input
                  id="sample_residue"
                  type="text"
                  maxLength={150}
                  aria-invalid={!!errors.sample_residue}
                  {...register("sample_residue")}
                />
                {errors.sample_residue && (
                  <FieldError>{errors.sample_residue.message}</FieldError>
                )}
              </Field>

              <Field>
                <FieldLabel htmlFor="batch_number">Batch Number</FieldLabel>
                <Input
                  id="batch_number"
                  type="text"
                  maxLength={50}
                  aria-invalid={!!errors.batch_number}
                  {...register("batch_number")}
                />
                {errors.batch_number && (
                  <FieldError>{errors.batch_number.message}</FieldError>
                )}
              </Field>

              <Field>
                <FieldLabel htmlFor="sample_taken_date">
                  Sample Taken Date <span className="text-destructive">*</span>
                </FieldLabel>
                <Controller
                  name="sample_taken_date"
                  control={control}
                  render={({ field }) => (
                    <DatePickerSimple
                      value={field.value}
                      onChange={field.onChange}
                    />
                  )}
                />
                {errors.sample_taken_date && (
                  <FieldError>{errors.sample_taken_date.message}</FieldError>
                )}
              </Field>

              <Field>
                <FieldLabel htmlFor="sample_taken_by">
                  Sample Taken By <span className="text-destructive">*</span>
                </FieldLabel>
                <Input
                  id="sample_taken_by"
                  type="text"
                  maxLength={100}
                  aria-invalid={!!errors.sample_taken_by}
                  {...register("sample_taken_by")}
                />
                {errors.sample_taken_by && (
                  <FieldError>{errors.sample_taken_by.message}</FieldError>
                )}
              </Field>

              <Field>
                <FieldLabel htmlFor="special_instructions">
                  Special Instructions{" "}
                  <span className="text-destructive">*</span>
                </FieldLabel>
                <textarea
                  id="special_instructions"
                  maxLength={500}
                  rows={3}
                  aria-invalid={!!errors.special_instructions}
                  className="w-full min-w-0 rounded-lg border border-input bg-transparent px-3 py-2 text-sm transition-colors outline-none focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/30 disabled:pointer-events-none disabled:cursor-not-allowed disabled:bg-input/50 disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-2 aria-invalid:ring-destructive/20"
                  {...register("special_instructions")}
                />
                {errors.special_instructions && (
                  <FieldError>{errors.special_instructions.message}</FieldError>
                )}
              </Field>

              <Field>
                <FieldLabel htmlFor="note">Note</FieldLabel>
                <textarea
                  id="note"
                  maxLength={1000}
                  rows={3}
                  aria-invalid={!!errors.note}
                  className="w-full min-w-0 rounded-lg border border-input bg-transparent px-3 py-2 text-sm transition-colors outline-none focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/30 disabled:pointer-events-none disabled:cursor-not-allowed disabled:bg-input/50 disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-2 aria-invalid:ring-destructive/20"
                  {...register("note")}
                />
                {errors.note && <FieldError>{errors.note.message}</FieldError>}
              </Field>
            </FieldGroup>
          </CardContent>
        </Card>

        {/* Backend error summary (HTTP 400) */}
        {backendErrors.length > 0 && (
          <div
            role="alert"
            className="rounded-lg border border-destructive/30 bg-destructive/5 px-4 py-3 text-sm text-destructive"
          >
            <p className="font-medium">Please fix the following errors:</p>
            <ul className="mt-1 ml-4 list-disc">
              {backendErrors.map((err, i) => (
                <li key={i}>{err}</li>
              ))}
            </ul>
          </div>
        )}

        {/* Submit button */}
        <Button
          type="submit"
          disabled={mutation.isPending}
          className="w-full cursor-pointer"
        >
          {mutation.isPending && (
            <Loader2 className="size-4 animate-spin" data-icon="inline-start" />
          )}
          {mutation.isPending ? "Submitting..." : "Create Reservation"}
        </Button>
      </div>

      {/* Right column: Test selection & line items */}
      <div className="lg:sticky lg:top-16 lg:self-start lg:max-h-[calc(100vh-5rem)] lg:overflow-y-auto">
        <Card>
          <CardHeader>
            <CardTitle>Test Items</CardTitle>
            <CardDescription>
              Select a test group then add individual tests.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <TestSelector
              fields={fields}
              append={append}
              remove={remove}
              register={register}
              errors={errors.tests}
            />
          </CardContent>
        </Card>
      </div>

    </form>
  );
}
