import { ProtectedRoute } from "~/components/auth/protected-route";
import AdminQuotationsPage from "~/pages/admin/quotations/page";

export default function AdminQuotations() {
  return (
    <ProtectedRoute>
      <AdminQuotationsPage />
    </ProtectedRoute>
  );
}
