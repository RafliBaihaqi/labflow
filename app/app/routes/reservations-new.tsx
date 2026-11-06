import ReservationNewPage from "~/pages/reservations/new";
import {ProtectedRoute} from "~/components/auth/protected-route";

export default function ReservationsNew() {
    return (
        <ProtectedRoute>
            <ReservationNewPage/>
        </ProtectedRoute>);
}
