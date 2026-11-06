import type React from "react";
import { useState, useEffect } from "react";
import { cn } from "~/lib/utils";
import { Button } from "~/components/ui/button";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "~/components/ui/field";
import { Input } from "~/components/ui/input";
import * as z from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { useMutation } from "@tanstack/react-query";
import { resetPassword } from "~/lib/api-client";
import { useLocation, useNavigate } from "react-router";
import { Eye, EyeOff, Lock, CheckCircle2 } from "lucide-react";
import { useT } from "~/i18n/use-t";

const resetSchema = z
  .object({
    password: z.string().min(8, "Password must be at least 8 characters"),
    confirmPassword: z.string().min(1, "Confirm your password"),
  })
  .refine((d) => d.password === d.confirmPassword, {
    message: "Passwords do not match",
    path: ["confirmPassword"],
  });

type ResetInput = z.infer<typeof resetSchema>;

const INPUT_CLASSNAME =
  "border-brand-baby-blue/40 focus:border-brand-ocean-blue focus:ring-brand-ocean-blue/30 pl-4 py-3 text-sm 2xl:text-base bg-white/80 backdrop-blur-sm transition-all duration-300 h-10";

export function UpdatePasswordForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const t = useT();
  const navigate = useNavigate();
  const location = useLocation();
  const resetToken =
    (location.state as { resetToken?: string } | null)?.resetToken ?? "";

  useEffect(() => {
    if (!resetToken) navigate("/forgot-password", { replace: true });
  }, [resetToken, navigate]);

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<ResetInput>({
    resolver: zodResolver(resetSchema),
    defaultValues: { password: "", confirmPassword: "" },
  });

  const passwordValue = watch("password", "");
  const confirmPasswordValue = watch("confirmPassword", "");

  const calculatePasswordStrength = (val: string): number => {
    let score = 0;
    if (val.length >= 8) score++;
    if (/[A-Z]/.test(val)) score++;
    if (/[a-z]/.test(val)) score++;
    if (/[0-9]/.test(val)) score++;
    if (/[^A-Za-z0-9]/.test(val)) score++;
    return score;
  };

  const passwordStrength = calculatePasswordStrength(passwordValue);

  const getPasswordStrengthColor = (score: number) => {
    return (
      [
        "bg-red-500",
        "bg-orange-500",
        "bg-yellow-500",
        "bg-blue-500",
        "bg-green-500",
      ][score - 1] || "bg-gray-300"
    );
  };

  const getPasswordStrengthText = (score: number) => {
    return (
      [
        t("auth.signup.passLevel1"),
        t("auth.signup.passLevel2"),
        t("auth.signup.passLevel3"),
        t("auth.signup.passLevel4"),
        t("auth.signup.passLevel5"),
      ][score - 1] || ""
    );
  };

  const mutation = useMutation({
    mutationFn: resetPassword,
    onSuccess: () => {
      navigate("/login", { state: { passwordReset: true } });
    },
  });

  const onSubmit = (data: ResetInput) => {
    mutation.mutate({ token: resetToken, password: data.password });
  };

  return (
    <div
      className={cn("space-y-2 lg:space-y-6 pb-4 2xl:pb-6 px-4 2xl:px-6", className)}
      {...props}
    >
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <FieldGroup className="space-y-4">
          <Field>
            <FieldLabel
              htmlFor="password"
              className="text-brand-dark-grey font-semibold flex items-center text-sm 2xl:text-base"
            >
              <Lock className="h-4 w-4 mr-2 text-brand-ocean-blue" />
              {t("auth.reset.newPassword")}{" "}
              <span className="text-red-500 ml-1">*</span>
            </FieldLabel>
            <div className="relative">
              <Input
                id="password"
                type={showPassword ? "text" : "password"}
                placeholder="••••••••"
                className={cn(INPUT_CLASSNAME, "pr-12")}
                aria-invalid={!!errors.password}
                {...register("password")}
              />
              <button
                type="button"
                onClick={() => setShowPassword((v) => !v)}
                className="absolute inset-y-0 right-0 flex items-center pr-3 text-muted-foreground transition-colors hover:text-foreground"
                aria-label={showPassword ? "Hide password" : "Show password"}
              >
                {showPassword ? (
                  <EyeOff className="size-4" />
                ) : (
                  <Eye className="size-4" />
                )}
              </button>
            </div>

            {passwordValue && (
              <div className="space-y-2 mt-2">
                <div className="flex items-center justify-between">
                  <span className="text-xs text-brand-dark-grey/70">
                    {t("auth.signup.passwordStrength")}
                  </span>
                  <span
                    className={`text-xs font-semibold ${
                      passwordStrength <= 1
                        ? "text-red-600"
                        : passwordStrength === 2
                          ? "text-orange-600"
                          : passwordStrength === 3
                            ? "text-yellow-600"
                            : passwordStrength === 4
                              ? "text-blue-600"
                              : "text-green-600"
                    }`}
                  >
                    {getPasswordStrengthText(passwordStrength)}
                  </span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div
                    className={cn(
                      "h-2 rounded-full transition-all duration-300",
                      getPasswordStrengthColor(passwordStrength)
                    )}
                    style={{
                      width: `${(passwordStrength / 5) * 100}%`,
                    }}
                  ></div>
                </div>
              </div>
            )}
            <p className="text-xs text-brand-dark-grey/60 mt-1">
              {t("auth.signup.passwordMessage")}
            </p>
            {errors.password && (
              <FieldError className="mt-1">
                {errors.password.message}
              </FieldError>
            )}
          </Field>

          <Field>
            <FieldLabel
              htmlFor="confirmPassword"
              className="text-brand-dark-grey font-semibold flex items-center text-sm 2xl:text-base"
            >
              <Lock className="h-4 w-4 mr-2 text-brand-ocean-blue" />
              {t("auth.reset.confirmNewPassword")}{" "}
              <span className="text-red-500 ml-1">*</span>
            </FieldLabel>
            <div className="relative">
              <Input
                id="confirmPassword"
                type={showConfirmPassword ? "text" : "password"}
                placeholder="••••••••"
                className={cn(INPUT_CLASSNAME, "pr-12")}
                aria-invalid={!!errors.confirmPassword}
                {...register("confirmPassword")}
              />
              <button
                type="button"
                onClick={() => setShowConfirmPassword((v) => !v)}
                className="absolute inset-y-0 right-0 flex items-center pr-3 text-muted-foreground transition-colors hover:text-foreground"
                aria-label={
                  showConfirmPassword ? "Hide password" : "Show password"
                }
              >
                {showConfirmPassword ? (
                  <EyeOff className="size-4" />
                ) : (
                  <Eye className="size-4" />
                )}
              </button>
            </div>

            {confirmPasswordValue &&
              passwordValue !== confirmPasswordValue && (
                <p className="text-xs text-red-600 flex items-center mt-2">
                  <span className="w-1.5 h-1.5 bg-red-500 rounded-full mr-2"></span>
                  {t("auth.signup.passwordMismatch")}
                </p>
              )}
            {confirmPasswordValue &&
              passwordValue === confirmPasswordValue && (
                <p className="text-xs text-green-600 flex items-center mt-2">
                  <CheckCircle2 className="w-3.5 h-3.5 mr-1 text-green-600" />
                  {t("auth.signup.passwordMatch")}
                </p>
              )}
            {errors.confirmPassword && (
              <FieldError className="mt-1">
                {errors.confirmPassword.message}
              </FieldError>
            )}
          </Field>

          {mutation.isError && (
            <div
              role="alert"
              className="rounded-lg border border-destructive/30 bg-destructive/5 px-3 py-2 text-xs text-destructive"
            >
              {t("auth.reset.genericError")}
            </div>
          )}

          <Button
            type="submit"
            className="w-full bg-gradient-to-r from-brand-ocean-blue to-brand-ocean-blue/80 hover:from-brand-ocean-blue/90 hover:to-brand-ocean-blue/70 text-white font-semibold py-1.5 lg:py-3 shadow-lg hover:shadow-xl transition-all duration-300 disabled:opacity-70 text-sm lg:text-base"
            disabled={
              mutation.isPending ||
              (confirmPasswordValue !== "" &&
                passwordValue !== confirmPasswordValue)
            }
          >
            {mutation.isPending ? (
              <div className="flex items-center space-x-2">
                <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                <span>{t("auth.reset.resetting")}</span>
              </div>
            ) : (
              <span>{t("auth.reset.resetButton")}</span>
            )}
          </Button>
        </FieldGroup>
      </form>
    </div>
  );
}
