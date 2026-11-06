import DataTableDemoPage from "~/pages/examples/data-table-demo";
import { ProtectedRoute } from "~/components/auth/protected-route";

export default function DataTableDemo() {
  return (
    <ProtectedRoute>
      <DataTableDemoPage />
    </ProtectedRoute>
  );
}
