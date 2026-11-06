import { useEffect, useState } from "react";
import { cn } from "~/lib/utils";
import { Button } from "~/components/ui/button";
import { Field, FieldError, FieldGroup } from "~/components/ui/field";
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
} from "~/components/ui/input-otp";
import * as z from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { useMutation } from "@tanstack/react-query";
import axios from "axios";
import { verifyOtp, resendOtp } from "~/lib/api-client";
import { useLocation, useNavigate } from "react-router";
import { Loader2 } from "lucide-react";

const otpSchema = z.object({
  otp: z.string().regex(/^\d{6}$/, "Masukkan 6 digit kode"),
});

type OtpInput = z.infer<typeof otpSchema>;

const RESEND_COOLDOWN_SECONDS = 60;

export function OtpForm({ className, ...props }: React.ComponentProps<"div">) {
  const location = useLocation();
  const navigate = useNavigate();
  const email = (location.state as { email?: string } | null)?.email ?? "";

  useEffect(() => {
    if (!email) navigate("/register", { replace: true });
  }, [email, navigate]);

  const [cooldown, setCooldown] = useState(0);
  useEffect(() => {
    if (cooldown <= 0) return;
    const t = setTimeout(() => setCooldown((c) => c - 1), 1000);
    return () => clearTimeout(t);
  }, [cooldown]);

  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<OtpInput>({
    resolver: zodResolver(otpSchema),
    defaultValues: { otp: "" },
  });

  const verifyMutation = useMutation({
    mutationFn: verifyOtp,
    onSuccess: () => {
      navigate("/", { state: { verified: true } });
    },
  });

  const resendMutation = useMutation({
    mutationFn: resendOtp,
    onSuccess: () => setCooldown(RESEND_COOLDOWN_SECONDS),
  });

  const onSubmit = (data: OtpInput) => {
    verifyMutation.mutate({ email, otp: data.otp });
  };

  const verifyMessage =
    axios.isAxiosError(verifyMutation.error) &&
    verifyMutation.error.response?.status === 400
      ? "Invalid or expired code. Please try again."
      : "Something went wrong. Please try again.";

  const resendMessage =
    axios.isAxiosError(resendMutation.error) &&
    resendMutation.error.response?.status === 429
      ? "Please wait before resending."
      : "Could not resend the code. Please try again.";

  return (
    <div className={cn("flex flex-col gap-4", className)} {...props}>
      <form onSubmit={handleSubmit(onSubmit)}>
        <FieldGroup>
          <Field>
            <Controller
              control={control}
              name="otp"
              render={({ field }) => (
                <InputOTP maxLength={6} value={field.value} onChange={field.onChange}>
                  <InputOTPGroup>
                    <InputOTPSlot index={0} />
                    <InputOTPSlot index={1} />
                    <InputOTPSlot index={2} />
                    <InputOTPSlot index={3} />
                    <InputOTPSlot index={4} />
                    <InputOTPSlot index={5} />
                  </InputOTPGroup>
                </InputOTP>
              )}
            />
            {errors.otp && <FieldError>{errors.otp.message}</FieldError>}
          </Field>

          {verifyMutation.isError && (
            <div
              role="alert"
              className="rounded-lg border border-destructive/30 bg-destructive/5 px-3 py-2 text-xs text-destructive"
            >
              {verifyMessage}
            </div>
          )}

          <Button type="submit" size="lg" className="w-full" disabled={verifyMutation.isPending}>
            {verifyMutation.isPending && (
              <Loader2 className="size-4 animate-spin" data-icon="inline-start" />
            )}
            {verifyMutation.isPending ? "Verifying..." : "Verify"}
          </Button>

          <div className="text-center text-sm text-muted-foreground">
            Didn&apos;t get the code?{" "}
            <button
              type="button"
              disabled={cooldown > 0 || resendMutation.isPending}
              onClick={() => resendMutation.mutate({ email })}
              className="text-primary underline-offset-4 hover:underline disabled:opacity-50 disabled:no-underline"
            >
              {cooldown > 0 ? `Resend in ${cooldown}s` : "Resend code"}
            </button>
          </div>

          {resendMutation.isError && (
            <p className="text-center text-xs text-destructive">{resendMessage}</p>
          )}
        </FieldGroup>
      </form>
    </div>
  );
}
