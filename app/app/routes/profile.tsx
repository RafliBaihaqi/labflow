import ProfilePage from "~/pages/profile/page";
import { ProtectedRoute } from "~/components/auth/protected-route";

export default function Profile() {
  return (
    <ProtectedRoute>
      <ProfilePage />
    </ProtectedRoute>
  );
}
