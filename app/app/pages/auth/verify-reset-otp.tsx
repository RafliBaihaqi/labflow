import { AuthLayout } from "~/components/auth/auth-layout";
import { ResetOtpForm } from "~/components/auth/reset-otp-form";
import { useLocation } from "react-router";
import { useT } from "~/i18n/use-t";

export default function VerifyResetOtpPage() {
  const location = useLocation();
  const t = useT();
  const email = (location.state as { email?: string } | null)?.email ?? "";

  return (
    <AuthLayout
      title={t("auth.reset.verifyTitle")}
      description={
        <span>
          {t("auth.reset.verifyDescription")}{" "}
          <span className="font-semibold text-brand-dark-grey">{email}</span>.
        </span>
      }
    >
      <ResetOtpForm />
    </AuthLayout>
  );
}
