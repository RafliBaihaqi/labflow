import { GuestRoute } from "~/components/auth/guest-route";
import VerifyOtpPage from "~/pages/auth/verify-otp";
import type { Route } from "./+types/verify-otp";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Verify Email — LabFlow" },
    { name: "description", content: "Verify your LabFlow account." },
  ];
}

export default function VerifyOtp() {
  return (
    <GuestRoute>
      <VerifyOtpPage />
    </GuestRoute>
  );
}
