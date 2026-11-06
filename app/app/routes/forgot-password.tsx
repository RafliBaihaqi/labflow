import { GuestRoute } from "~/components/auth/guest-route";
import SendOTPToEmailPage from "~/pages/auth/send-otp-to-email";
import type { Route } from "./+types/forgot-password";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Forgot Password — LabFlow" },
    { name: "description", content: "Reset your LabFlow password." },
  ];
}

export default function ForgotPassword() {
  return (
    <GuestRoute>
      <SendOTPToEmailPage />
    </GuestRoute>
  );
}
