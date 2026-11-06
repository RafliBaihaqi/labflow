import { AuthLayout } from "~/components/auth/auth-layout";
import { OtpForm } from "~/components/auth/otp-form";
import { useLocation } from "react-router";

export default function VerifyOtpPage() {
  const location = useLocation();
  const email = (location.state as { email?: string } | null)?.email ?? "";

  return (
    <AuthLayout
      title="Verify your email"
      description={
        <span>
          We sent a 6-digit code to{" "}
          <span className="font-semibold text-brand-dark-grey">{email}</span>.
          Enter it below to activate your account.
        </span>
      }
    >
      <OtpForm />
    </AuthLayout>
  );
}
