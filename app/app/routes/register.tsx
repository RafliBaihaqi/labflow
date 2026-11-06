import { GuestRoute } from "~/components/auth/guest-route";
import RegisterPage from "~/pages/auth/register";
import type { Route } from "./+types/register";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Register — LabFlow" },
    { name: "description", content: "Create your LabFlow account." },
  ];
}

export default function Register() {
  return (
    <GuestRoute>
      <RegisterPage />
    </GuestRoute>
  );
}
