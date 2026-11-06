import ReservationListPage from "~/pages/reservations/list";
import {ProtectedRoute} from "~/components/auth/protected-route";

export default function ReservationsList() {
    return (
        <ProtectedRoute>
            <ReservationListPage/>
        </ProtectedRoute>);
}
