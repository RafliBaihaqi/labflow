import LoginForm from "~/components/auth/login-form";
import { AuthLayout } from "~/components/auth/auth-layout";
import { useT } from "~/i18n/use-t";

export default function LoginPage() {
  const t = useT();

  return (
    <AuthLayout
      title={t("auth.login.title")}
      description={t("auth.login.subtitle")}
    >
      <LoginForm />
    </AuthLayout>
  );
}
