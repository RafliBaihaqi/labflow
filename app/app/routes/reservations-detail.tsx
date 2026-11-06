import type {Route} from "./+types/reservations-detail";
import ReservationDetailPage from "~/pages/reservations/detail";
import {ProtectedRoute} from "~/components/auth/protected-route";

export default function ReservationsDetail({params}: Route.ComponentProps) {
    return (
        <ProtectedRoute>
            <ReservationDetailPage xid={params.xid}/>
        </ProtectedRoute>);
}
