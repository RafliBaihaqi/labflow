import { AuthLayout } from "~/components/auth/auth-layout";
import { useT } from "~/i18n/use-t";
import { ResetForm } from "~/components/auth/reset-form";

export default function SendOTPToEmailPage() {
  const t = useT();

  return (
    <AuthLayout
      title={t("auth.reset.forgotTitle")}
      description={t("auth.reset.forgotDescription")}
    >
      <ResetForm />
    </AuthLayout>
  );
}
