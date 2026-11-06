import { AuthLayout } from "~/components/auth/auth-layout";
import { RegisterForm } from "~/components/auth/register-form";

export default function RegisterPage() {
  return (
    <AuthLayout
      title="Create your account"
      description="Set up your LabFlow account to start managing your lab."
    >
      <RegisterForm />
    </AuthLayout>
  );
}
