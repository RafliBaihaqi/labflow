import { GuestRoute } from "~/components/auth/guest-route";
import LoginPage from "~/pages/auth/login";

export default function Login() {
  return (
    <GuestRoute>
      <LoginPage />
    </GuestRoute>
  );
}
