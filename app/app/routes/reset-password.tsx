import { GuestRoute } from "~/components/auth/guest-route";
import ResetPasswordPage from "~/pages/auth/reset-password";
import type { Route } from "./+types/reset-password";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Reset Password — LabFlow" },
    { name: "description", content: "Set a new password for your account." },
  ];
}

export default function ResetPassword() {
  return (
    <GuestRoute>
      <ResetPasswordPage />
    </GuestRoute>
  );
}
