"use client";

import type React from "react";
import { useState } from "react";
import { cn } from "~/lib/utils";
import { Button } from "~/components/ui/button";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "~/components/ui/field";
import { Input } from "~/components/ui/input";
import { Textarea } from "~/components/ui/textarea";
import * as z from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { useMutation } from "@tanstack/react-query";
import axios from "axios";
import { register as registerApi } from "~/lib/api-client";
import { useNavigate } from "react-router";
import { NavLink } from "react-router";
import {
  Eye,
  EyeOff,
  User,
  Mail,
  Phone,
  Building,
  MapPin,
  Lock,
  Users,
  CheckCircle2,
} from "lucide-react";
import { useT } from "~/i18n/use-t";

const registerSchema = z
  .object({
    name: z.string().min(1, "Nama wajib diisi"),
    username: z.string().min(1, "Username wajib diisi"),
    email: z.string().email("Email tidak valid"),
    password: z.string().min(8, "Kata sandi minimal 8 karakter"),
    confirmPassword: z.string().min(1, "Konfirmasi kata sandi wajib diisi"),
    phone_number: z.string().optional(),
    company: z.string().optional(),
    address: z.string().optional(),
  })
  .refine((d) => d.password === d.confirmPassword, {
    message: "Kata sandi tidak cocok",
    path: ["confirmPassword"],
  });

type RegisterInput = z.infer<typeof registerSchema>;

const INPUT_CLASSNAME = "border-brand-baby-blue/40 focus:border-brand-ocean-blue focus:ring-brand-ocean-blue/30 pl-4 py-3 text-sm 2xl:text-base bg-white/80 backdrop-blur-sm transition-all duration-300 h-10";

export function RegisterForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const t = useT();

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<RegisterInput>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      name: "",
      username: "",
      email: "",
      password: "",
      confirmPassword: "",
      phone_number: "",
      company: "",
      address: "",
    },
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

  const navigate = useNavigate();

  const mutation = useMutation({
    mutationFn: registerApi,
    onSuccess: (_res, variables) => {
      navigate("/verify-otp", { state: { email: variables.email } });
    },
  });

  const onSubmit = (data: RegisterInput) => {
    mutation.mutate({
      name: data.name,
      username: data.username,
      email: data.email,
      password: data.password,
      phone_number: data.phone_number || undefined,
      company: data.company || undefined,
      address: data.address || undefined,
    });
  };

  const serverMessage =
    axios.isAxiosError(mutation.error) && mutation.error.response?.status === 409
      ? "Username or email already exists."
      : "Something went wrong. Please try again.";

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        <FieldGroup className="space-y-6">
          
          {/* Personal Information Section */}
          <div className="space-y-4">
            <div className="flex items-center space-x-2 pb-2 border-b border-brand-baby-blue/20">
              <div className="w-6 h-6 bg-brand-ocean-blue/10 rounded-full flex items-center justify-center">
                <User className="h-3.5 w-3.5 text-brand-ocean-blue" />
              </div>
              <h3 className="font-semibold text-brand-dark-grey text-sm 2xl:text-base">
                {t("auth.signup.privateInformation")}
              </h3>
            </div>

            <Field>
              <FieldLabel
                htmlFor="name"
                className="text-brand-dark-grey font-medium flex items-center text-sm 2xl:text-base"
              >
                <User className="h-4 w-4 mr-2 text-brand-ocean-blue" />
                {t("auth.signup.fullName")} <span className="text-red-500 ml-1">*</span>
              </FieldLabel>
              <Input
                id="name"
                type="text"
                placeholder={t("auth.signup.fullNamePlaceholder")}
                aria-invalid={!!errors.name}
                {...register("name")}
                required
                className={INPUT_CLASSNAME}
              />
              {errors.name && <FieldError>{errors.name.message}</FieldError>}
            </Field>

            <Field>
              <FieldLabel
                htmlFor="username"
                className="text-brand-dark-grey font-medium flex items-center text-sm 2xl:text-base"
              >
                <User className="h-4 w-4 mr-2 text-brand-ocean-blue" />
                {t("auth.signup.username")} <span className="text-red-500 ml-1">*</span>
              </FieldLabel>
              <Input
                id="username"
                type="text"
                placeholder={t("auth.signup.usernamePlaceholder")}
                aria-invalid={!!errors.username}
                {...register("username")}
                required
                className={INPUT_CLASSNAME}
              />
              {errors.username && (
                <FieldError>{errors.username.message}</FieldError>
              )}
            </Field>

            <Field>
              <FieldLabel
                htmlFor="email"
                className="text-brand-dark-grey font-medium flex items-center text-sm 2xl:text-base"
              >
                <Mail className="h-4 w-4 mr-2 text-brand-ocean-blue" />
                {t("auth.signup.email")} <span className="text-red-500 ml-1">*</span>
              </FieldLabel>
              <Input
                id="email"
                type="email"
                placeholder={t("auth.signup.emailPlaceholder")}
                aria-invalid={!!errors.email}
                {...register("email")}
                required
                className={INPUT_CLASSNAME}
              />
              {errors.email && <FieldError>{errors.email.message}</FieldError>}
            </Field>

            <Field>
              <FieldLabel
                htmlFor="phone_number"
                className="text-brand-dark-grey font-medium flex items-center text-sm 2xl:text-base"
              >
                <Phone className="h-4 w-4 mr-2 text-brand-ocean-blue" />
                {t("auth.signup.phoneNumber")}
              </FieldLabel>
              <Input
                id="phone_number"
                type="tel"
                placeholder={t("auth.signup.phoneNumberPlaceholder")}
                {...register("phone_number")}
                className={INPUT_CLASSNAME}
              />
            </Field>
          </div>

          {/* Company Information Section */}
          <div className="space-y-4">
            <div className="flex items-center space-x-2 pb-2 border-b border-brand-baby-blue/20">
              <div className="w-6 h-6 bg-brand-ocean-blue/10 rounded-full flex items-center justify-center">
                <Building className="h-3.5 w-3.5 text-brand-ocean-blue" />
              </div>
              <h3 className="font-semibold text-brand-dark-grey text-sm 2xl:text-base">
                {t("auth.signup.companyInfo")}
              </h3>
            </div>

            <Field>
              <FieldLabel
                htmlFor="company"
                className="text-brand-dark-grey font-medium flex items-center text-sm 2xl:text-base"
              >
                <Building className="h-4 w-4 mr-2 text-brand-ocean-blue" />
                {t("auth.signup.companyName")}
              </FieldLabel>
              <Input
                id="company"
                type="text"
                placeholder={t("auth.signup.companyNamePlaceholder")}
                {...register("company")}
                className={INPUT_CLASSNAME}
              />
            </Field>

            <Field>
              <FieldLabel
                htmlFor="address"
                className="text-brand-dark-grey font-medium flex items-center text-sm 2xl:text-base"
              >
                <MapPin className="h-4 w-4 mr-2 text-brand-ocean-blue" />
                {t("auth.signup.companyAddress")}
              </FieldLabel>
              <Textarea
                id="address"
                placeholder={t("auth.signup.addressPlaceholder")}
                {...register("address")}
                rows={3}
                className="border-brand-baby-blue/40 focus:border-brand-ocean-blue focus:ring-brand-ocean-blue/30 bg-white/80 backdrop-blur-sm transition-all duration-300 resize-none pl-4 py-3 text-sm 2xl:text-base min-h-[80px]"
              />
            </Field>
          </div>

          {/* Account Security Section */}
          <div className="space-y-4">
            <div className="flex items-center space-x-2 pb-2 border-b border-brand-baby-blue/20">
              <div className="w-6 h-6 bg-brand-ocean-blue/10 rounded-full flex items-center justify-center">
                <Lock className="h-3.5 w-3.5 text-brand-ocean-blue" />
              </div>
              <h3 className="font-semibold text-brand-dark-grey text-sm 2xl:text-base">
                {t("auth.signup.accountSecurity")}
              </h3>
            </div>

            <Field>
              <FieldLabel
                htmlFor="password"
                className="text-brand-dark-grey font-medium flex items-center text-sm 2xl:text-base"
              >
                <Lock className="h-4 w-4 mr-2 text-brand-ocean-blue" />
                {t("auth.signup.password")} <span className="text-red-500 ml-1">*</span>
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

              {/* Password Strength Indicator */}
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
                      className={cn("h-2 rounded-full transition-all duration-300", getPasswordStrengthColor(passwordStrength))}
                      style={{ width: `${(passwordStrength / 5) * 100}%` }}
                    ></div>
                  </div>
                </div>
              )}
              <p className="text-xs text-brand-dark-grey/60 mt-1">
                {t("auth.signup.passwordMessage")}
              </p>
              {errors.password && (
                <FieldError className="mt-1">{errors.password.message}</FieldError>
              )}
            </Field>

            <Field>
              <FieldLabel
                htmlFor="confirmPassword"
                className="text-brand-dark-grey font-medium flex items-center text-sm 2xl:text-base"
              >
                <Lock className="h-4 w-4 mr-2 text-brand-ocean-blue" />
                {t("auth.signup.confirmPassword")} <span className="text-red-500 ml-1">*</span>
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
                  aria-label={showConfirmPassword ? "Hide password" : "Show password"}
                >
                  {showConfirmPassword ? (
                    <EyeOff className="size-4" />
                  ) : (
                    <Eye className="size-4" />
                  )}
                </button>
              </div>

              {/* Password Match Status Messages */}
              {confirmPasswordValue && passwordValue !== confirmPasswordValue && (
                <p className="text-xs text-red-600 flex items-center mt-2">
                  <span className="w-1.5 h-1.5 bg-red-500 rounded-full mr-2"></span>
                  {t("auth.signup.passwordMismatch")}
                </p>
              )}
              {confirmPasswordValue && passwordValue === confirmPasswordValue && (
                <p className="text-xs text-green-600 flex items-center mt-2">
                  <CheckCircle2 className="w-3.5 h-3.5 mr-1 text-green-600" />
                  {t("auth.signup.passwordMatch")}
                </p>
              )}
              {errors.confirmPassword && (
                <FieldError className="mt-1">{errors.confirmPassword.message}</FieldError>
              )}
            </Field>
          </div>

          {/* Server Error Message */}
          {mutation.isError && (
            <div
              role="alert"
              className="rounded-lg border border-destructive/30 bg-destructive/5 px-3 py-2 text-xs text-destructive"
            >
              {serverMessage}
            </div>
          )}

          {/* Submit Button */}
          <Button
            type="submit"
            className="w-full bg-gradient-to-r from-brand-ocean-blue to-brand-ocean-blue/80 hover:from-brand-ocean-blue/90 hover:to-brand-ocean-blue/70 text-white font-semibold py-3 shadow-lg hover:shadow-xl transition-all duration-300 disabled:opacity-70 text-sm lg:text-base flex items-center justify-center space-x-2 h-10"
            disabled={mutation.isPending || (confirmPasswordValue !== "" && passwordValue !== confirmPasswordValue)}
          >
            {mutation.isPending ? (
              <div className="flex items-center space-x-2">
                <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                <span>{t("auth.signup.loadingRegistration")}</span>
              </div>
            ) : (
              <div className="flex items-center space-x-2">
                <span>{t("auth.signup.registerNewAccount")}</span>
                <div className="w-5 h-5 bg-white/20 rounded-full flex items-center justify-center">
                  <Users className="w-3 h-3 text-white" />
                </div>
              </div>
            )}
          </Button>
        </FieldGroup>
      </form>

      {/* Card Footer Block */}
      <div className="w-full text-center bg-gradient-to-r from-brand-baby-blue/10 to-brand-light-orange/10 rounded-lg p-4 border border-brand-baby-blue/20">
        <p className="text-brand-dark-grey/70 text-sm mb-2">
          {t("auth.signup.alreadyHaveAccount")}
        </p>
        <NavLink
          to="/login"
          className="inline-flex items-center space-x-2 text-brand-ocean-blue hover:text-brand-ocean-blue/80 font-semibold hover:underline transition-colors text-sm 2xl:text-base"
        >
          <span>{t("auth.signup.login")}</span>
          <div className="w-4 h-4 bg-brand-ocean-blue/20 rounded-full flex items-center justify-center">
            <div className="w-1.5 h-1.5 bg-brand-ocean-blue rounded-full"></div>
          </div>
        </NavLink>
      </div>
    </div>
  );
}
