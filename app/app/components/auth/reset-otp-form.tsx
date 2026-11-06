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
import { verifyResetOtp, forgotPassword } from "~/lib/api-client";
import { useLocation, useNavigate } from "react-router";
import { Loader2 } from "lucide-react";
import { useT } from "~/i18n/use-t";

const otpSchema = z.object({
  otp: z.string().regex(/^\d{6}$/, "Enter a 6-digit code"),
});

type OtpInput = z.infer<typeof otpSchema>;

const RESEND_COOLDOWN_SECONDS = 60;

export function ResetOtpForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const location = useLocation();
  const navigate = useNavigate();
  const t = useT();
  const email = (location.state as { email?: string } | null)?.email ?? "";

  useEffect(() => {
    if (!email) navigate("/forgot-password", { replace: true });
  }, [email, navigate]);

  const [cooldown, setCooldown] = useState(0);
  useEffect(() => {
    if (cooldown <= 0) return;
    const timer = setTimeout(() => setCooldown((c) => c - 1), 1000);
    return () => clearTimeout(timer);
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
    mutationFn: verifyResetOtp,
    onSuccess: (res) => {
      navigate("/reset-password", {
        state: { resetToken: res.data.reset_token },
      });
    },
  });

  const resendMutation = useMutation({
    mutationFn: forgotPassword,
    onSuccess: () => setCooldown(RESEND_COOLDOWN_SECONDS),
  });

  const onSubmit = (data: OtpInput) => {
    verifyMutation.mutate({ email, otp: data.otp });
  };

  const verifyMessage =
    axios.isAxiosError(verifyMutation.error) &&
    verifyMutation.error.response?.status === 401
      ? t("auth.reset.invalidOtp")
      : t("auth.reset.genericError");

  const resendMessage =
    axios.isAxiosError(resendMutation.error) &&
    resendMutation.error.response?.status === 429
      ? t("auth.reset.genericError")
      : t("auth.reset.genericError");

  return (
    <div className={cn("flex flex-col gap-4 px-4 2xl:px-6 pb-4 2xl:pb-6", className)} {...props}>
      <form onSubmit={handleSubmit(onSubmit)}>
        <FieldGroup>
          <Field>
            <Controller
              control={control}
              name="otp"
              render={({ field }) => (
                <InputOTP
                  maxLength={6}
                  value={field.value}
                  onChange={field.onChange}
                >
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

          <Button
            type="submit"
            className="w-full bg-gradient-to-r from-brand-ocean-blue to-brand-ocean-blue/80 hover:from-brand-ocean-blue/90 hover:to-brand-ocean-blue/70 text-white font-semibold py-1.5 lg:py-3 shadow-lg hover:shadow-xl transition-all duration-300 disabled:opacity-70 text-sm lg:text-base"
            disabled={verifyMutation.isPending}
          >
            {verifyMutation.isPending && (
              <Loader2 className="size-4 animate-spin mr-2" />
            )}
            {verifyMutation.isPending
              ? t("auth.reset.verifying")
              : t("auth.reset.verifyButton")}
          </Button>

          <div className="text-center text-sm text-muted-foreground">
            {t("auth.reset.noCode")}{" "}
            <button
              type="button"
              disabled={cooldown > 0 || resendMutation.isPending}
              onClick={() => resendMutation.mutate({ email })}
              className="text-primary underline-offset-4 hover:underline disabled:opacity-50 disabled:no-underline"
            >
              {cooldown > 0
                ? t("auth.reset.resendIn").replace("{seconds}", String(cooldown))
                : t("auth.reset.resendCode")}
            </button>
          </div>

          {resendMutation.isError && (
            <p className="text-center text-xs text-destructive">
              {resendMessage}
            </p>
          )}
        </FieldGroup>
      </form>
    </div>
  );
}
