"use client";

import type React from "react";
import { useState, useCallback } from "react";
import { useTurnstile } from "~/hooks/use-turnstile";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { Mail, Lock, Eye, EyeOff } from "lucide-react";
import z from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { NavLink, useLocation, useNavigate } from "react-router";
import { login } from "~/lib/api-client";
import { Field, FieldError, FieldGroup, FieldLabel } from "../ui/field";
import { useT } from "~/i18n/use-t";
import { cn } from "~/lib/utils";

const TURNSTILE_SITE_KEY = (import.meta.env.VITE_TURNSTILE_SITE_KEY as string) || "";
const DUMMY_TOKEN = "1x00000000000000000000AA";
// In dev, we allow bypass to stay true: if DEV mode or sitekey missing/test key, prefill dummy token
const IS_DEV = import.meta.env.DEV;
const IS_TEST_KEY = TURNSTILE_SITE_KEY === DUMMY_TOKEN;
const IS_DEV_BYPASS = IS_DEV || !TURNSTILE_SITE_KEY || IS_TEST_KEY;

function createLoginSchema(t: (id: string) => string) {
  return z.object({
    identifier: z.string().min(1, t("auth.validation.identifierRequired")),
    password: z.string().min(8, t("auth.validation.passwordMin")),
  });
}

type LoginInput = z.infer<ReturnType<typeof createLoginSchema>>;

export default function LoginForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const [showPassword, setShowPassword] = useState(false);
  // Dev bypass stays true: prefill dummy token in dev or when sitekey is test/missing
  const [turnstileToken, setTurnstileToken] = useState<string | null>(
    IS_DEV_BYPASS ? DUMMY_TOKEN : null
  );
  const t = useT();
  const loginSchema = createLoginSchema(t);

  const onTurnstileSuccess = useCallback((token: string) => {
    setTurnstileToken(token);
  }, []);

  const onTurnstileExpire = useCallback(() => {
    // In dev bypass, keep dummy token instead of null to keep button enabled
    setTurnstileToken(IS_DEV_BYPASS ? DUMMY_TOKEN : null);
  }, []);

  const { containerRef, reset: resetTurnstile } = useTurnstile({
    siteKey: TURNSTILE_SITE_KEY || DUMMY_TOKEN, // fallback to dummy to avoid undefined error in hook
    onSuccess: onTurnstileSuccess,
    onExpire: onTurnstileExpire,
  });

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginInput>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      identifier: "",
      password: "",
    },
  });

  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const location = useLocation();
  const justVerified =
    (location.state as { verified?: boolean } | null)?.verified ?? false;
  const passwordReset =
    (location.state as { passwordReset?: boolean } | null)?.passwordReset ?? false;

  const mutation = useMutation({
    mutationFn: login,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["me"] });
      navigate("/dashboard");
    },
    onError: () => {
      // Expire the current token so the user must solve a fresh challenge.
      setTurnstileToken(null);
      resetTurnstile();
    },
  });

  const onSubmit = (data: LoginInput) => {
    const token = turnstileToken || (IS_DEV_BYPASS ? DUMMY_TOKEN : null);
    if (!token) return;
    mutation.mutate({ ...data, turnstileToken: token });
  };

  return (
    <div className={cn("space-y-2 lg:space-y-6 pb-4 2xl:pb-6 px-4 2xl:px-6", className)} {...props}>
      <form onSubmit={handleSubmit(onSubmit)}>
        <FieldGroup>
          {justVerified && (
            <div
              role="status"
              className="rounded-lg border border-primary/30 bg-primary/5 px-3 py-2 text-xs text-primary"
            >
              {t("auth.login.emailVerified")}
            </div>
          )}

          {passwordReset && (
            <div
              role="status"
              className="rounded-lg border border-primary/30 bg-primary/5 px-3 py-2 text-xs text-primary"
            >
              {t("auth.login.passwordResetSuccess")}
            </div>
          )}

          <Field>
            <FieldLabel
              htmlFor="identifier"
              className="text-brand-dark-grey font-semibold flex items-center text-sm 2xl:text-base"
            >
              <Mail className="h-4 w-4 mr-2 text-brand-ocean-blue" />
              {t("auth.login.identifier")}
            </FieldLabel>
            <Input
              id="identifier"
              type="text"
              placeholder={t("auth.login.identifierPlaceholder")}
              aria-invalid={!!errors.identifier}
              {...register("identifier")}
              required
              className="border-brand-baby-blue/40 focus:border-brand-ocean-blue focus:ring-brand-ocean-blue/30 pl-4 py-3 text-sm 2xl:text-base bg-white/80 backdrop-blur-sm transition-all duration-300 h-10"
            />
            {errors.identifier && (
              <FieldError>{errors.identifier.message}</FieldError>
            )}
          </Field>

          <Field>
            <FieldLabel
              htmlFor="password"
              className="text-brand-dark-grey font-semibold flex items-center text-sm 2xl:text-base"
            >
              <Lock className="h-4 w-4 mr-2 text-brand-ocean-blue" />
              {t("auth.login.password")}
            </FieldLabel>
            <div className="relative">
              <Input
                id="password"
                type={showPassword ? "text" : "password"}
                placeholder={t("auth.login.passwordPlaceholder")}
                aria-invalid={!!errors.password}
                {...register("password")}
                required
                className="border-brand-baby-blue/40 focus:border-brand-ocean-blue focus:ring-brand-ocean-blue/30 pl-4 pr-10 py-3 text-sm 2xl:text-base bg-white/80 backdrop-blur-sm transition-all duration-300 h-10"
              />
              <button
                type="button"
                onClick={() => setShowPassword((v) => !v)}
                className="absolute inset-y-0 right-0 flex items-center pr-3 text-muted-foreground transition-colors hover:text-foreground"
                aria-label={
                  showPassword
                    ? t("auth.login.hidePassword")
                    : t("auth.login.showPassword")
                }
              >
                {showPassword ? (
                  <EyeOff className="size-4" />
                ) : (
                  <Eye className="size-4" />
                )}
              </button>
            </div>
            {errors.password && (
              <FieldError>{errors.password.message}</FieldError>
            )}
          </Field>

          <div className="flex items-center justify-between">
            <label className="flex cursor-pointer items-center gap-2 text-sm text-muted-foreground">
              <input
                type="checkbox"
                className="size-4 rounded border-input accent-secondary"
              />
              {t("auth.login.rememberMe")}
            </label>
            <NavLink
              to="/forgot-password"
              className="text-sm text-primary underline-offset-4 hover:underline"
            >
              {t("auth.login.forgotPassword")}
            </NavLink>
          </div>

          {/* Server error */}
          {mutation.isError && (
            <div
              role="alert"
              className="rounded-lg border border-destructive/30 bg-destructive/5 px-3 py-2 text-xs text-destructive"
            >
              {(mutation.error as any)?.response?.data?.message ||
                (mutation.error as any)?.response?.data?.error ||
                (mutation.error as any)?.message ||
                t("auth.login.invalidCredentials")}
            </div>
          )}

          {/* Cloudflare Turnstile widget - suppressHydrationWarning to fix React #418 */}
          <div
            ref={containerRef}
            suppressHydrationWarning
            className="flex justify-center"
          />
          {IS_DEV_BYPASS && !TURNSTILE_SITE_KEY && (
            <div className="text-xs text-amber-600 text-center">
              Dev mode: Turnstile bypass active (dummy token)
            </div>
          )}

          {/* Login Button */}
          <Button
            type="submit"
            className="w-full bg-gradient-to-r from-brand-ocean-blue to-brand-ocean-blue/80 hover:from-brand-ocean-blue/90 hover:to-brand-ocean-blue/70 text-white font-semibold py-1.5 lg:py-3 shadow-lg hover:shadow-xl transition-all duration-300 disabled:opacity-70 text-sm lg:text-base"
            disabled={mutation.isPending || (!IS_DEV_BYPASS && !turnstileToken)}
          >
            {mutation.isPending ? (
              <div className="flex items-center space-x-2">
                <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                <span>{t("auth.login.submitting")}</span>
              </div>
            ) : (
              <div className="flex items-center space-x-2">
                <span>{t("auth.login.submit")}</span>
              </div>
            )}
          </Button>
        </FieldGroup>
      </form>

      {/* Links Section */}
      <div className="w-full space-y-2 2xl:space-y-4 pt-4 border-t border-brand-baby-blue/30">
        <div className="flex justify-center items-center bg-gradient-to-r from-brand-baby-blue/10 to-brand-light-orange/10 rounded-lg p-4 border border-brand-baby-blue/20 gap-1">
          <p className="text-brand-dark-grey/70 text-sm 2xl:text-base leading-tight 2xl:leading-normal">
            {t("auth.login.noAccount")}
          </p>
          <NavLink
            to="/register"
            className="inline-flex items-center space-x-2 text-brand-ocean-blue hover:text-brand-ocean-blue/80 font-semibold hover:underline transition-colors text-sm 2xl:text-base"
          >
            <span className="leading-tight">
              {t("auth.login.register")}
            </span>
          </NavLink>
        </div>
      </div>
    </div>
  );
}
