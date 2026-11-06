import type React from "react";
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
import axios from "axios";
import { forgotPassword } from "~/lib/api-client";
import { useNavigate } from "react-router";
import { NavLink } from "react-router";
import { Mail } from "lucide-react";
import { useT } from "~/i18n/use-t";

function createForgotSchema(t: (id: string) => string) {
  return z.object({
    email: z.string().email(t("auth.reset.genericError")),
  });
}

type ForgotInput = z.infer<ReturnType<typeof createForgotSchema>>;

const INPUT_CLASSNAME =
  "border-brand-baby-blue/40 focus:border-brand-ocean-blue focus:ring-brand-ocean-blue/30 pl-4 py-3 text-sm 2xl:text-base bg-white/80 backdrop-blur-sm transition-all duration-300 h-10";

export function ResetForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const t = useT();
  const navigate = useNavigate();
  const forgotSchema = createForgotSchema(t);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ForgotInput>({
    resolver: zodResolver(forgotSchema),
    defaultValues: { email: "" },
  });

  const mutation = useMutation({
    mutationFn: forgotPassword,
    onSuccess: (_res, variables) => {
      navigate("/verify-reset-otp", { state: { email: variables.email } });
    },
  });

  const onSubmit = (data: ForgotInput) => {
    mutation.mutate({ email: data.email });
  };

  const serverMessage =
    axios.isAxiosError(mutation.error) &&
    mutation.error.response?.status === 404
      ? t("auth.reset.userNotFound")
      : t("auth.reset.genericError");

  return (
    <div
      className={cn("space-y-2 lg:space-y-6 pb-4 2xl:pb-6 px-4 2xl:px-6", className)}
      {...props}
    >
      <form onSubmit={handleSubmit(onSubmit)}>
        <FieldGroup>
          <Field>
            <FieldLabel
              htmlFor="email"
              className="text-brand-dark-grey font-semibold flex items-center text-sm 2xl:text-base"
            >
              <Mail className="h-4 w-4 mr-2 text-brand-ocean-blue" />
              {t("auth.reset.emailLabel")}
            </FieldLabel>
            <Input
              id="email"
              type="email"
              placeholder={t("auth.reset.emailPlaceholder")}
              aria-invalid={!!errors.email}
              {...register("email")}
              required
              className={INPUT_CLASSNAME}
            />
            {errors.email && (
              <FieldError>{errors.email.message}</FieldError>
            )}
          </Field>

          {mutation.isError && (
            <div
              role="alert"
              className="rounded-lg border border-destructive/30 bg-destructive/5 px-3 py-2 text-xs text-destructive"
            >
              {serverMessage}
            </div>
          )}

          <Button
            type="submit"
            className="w-full bg-gradient-to-r from-brand-ocean-blue to-brand-ocean-blue/80 hover:from-brand-ocean-blue/90 hover:to-brand-ocean-blue/70 text-white font-semibold py-1.5 lg:py-3 shadow-lg hover:shadow-xl transition-all duration-300 disabled:opacity-70 text-sm lg:text-base"
            disabled={mutation.isPending}
          >
            {mutation.isPending ? (
              <div className="flex items-center space-x-2">
                <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                <span>{t("auth.reset.sendingOtp")}</span>
              </div>
            ) : (
              <span>{t("auth.reset.sendOtp")}</span>
            )}
          </Button>
        </FieldGroup>
      </form>

      <div className="w-full space-y-2 2xl:space-y-4 pt-4 border-t border-brand-baby-blue/30">
        <div className="flex justify-center items-center bg-gradient-to-r from-brand-baby-blue/10 to-brand-light-orange/10 rounded-lg p-4 border border-brand-baby-blue/20 gap-1">
          <p className="text-brand-dark-grey/70 text-sm 2xl:text-base leading-tight 2xl:leading-normal">
            {t("auth.reset.rememberPassword")}
          </p>
          <NavLink
            to="/login"
            className="inline-flex items-center space-x-2 text-brand-ocean-blue hover:text-brand-ocean-blue/80 font-semibold hover:underline transition-colors text-sm 2xl:text-base"
          >
            <span className="leading-tight">
              {t("auth.reset.backToLogin")}
            </span>
          </NavLink>
        </div>
      </div>
    </div>
  );
}
