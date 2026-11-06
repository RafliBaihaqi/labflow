import { GuestRoute } from "~/components/auth/guest-route";
import VerifyResetOtpPage from "~/pages/auth/verify-reset-otp";
import type { Route } from "./+types/verify-reset-otp";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Verify Code — LabFlow" },
    { name: "description", content: "Verify your password reset code." },
  ];
}

export default function VerifyResetOtp() {
  return (
    <GuestRoute>
      <VerifyResetOtpPage />
    </GuestRoute>
  );
}
