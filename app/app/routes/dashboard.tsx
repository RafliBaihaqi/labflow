import {ProtectedRoute} from "~/components/auth/protected-route";
import DashboardPage from "~/pages/dashboard/page";

export default function Dashboard() {
    return (
        <ProtectedRoute>
            <DashboardPage/>
        </ProtectedRoute>
    );
}