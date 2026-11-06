import { AuthLayout } from "~/components/auth/auth-layout";
import { useT } from "~/i18n/use-t";
import { UpdatePasswordForm } from "~/components/auth/update-password";

export default function ResetPasswordPage() {
  const t = useT();

  return (
    <AuthLayout
      title={t("auth.reset.newPasswordTitle")}
      description={t("auth.reset.newPasswordDescription")}
    >
      <UpdatePasswordForm />
    </AuthLayout>
  );
}
