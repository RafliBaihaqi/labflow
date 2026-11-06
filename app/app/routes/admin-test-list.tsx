import { ProtectedRoute } from "~/components/auth/protected-route";
import AdminTestListPage from "~/pages/admin/test-list/page";

export default function AdminTestList() {
  return (
    <ProtectedRoute>
      <AdminTestListPage />
    </ProtectedRoute>
  );
}
