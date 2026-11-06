import { Link } from "react-router";
import { ArrowLeft } from "lucide-react";

import { MainLayout } from "~/components/main/main-layout";
import { CreationForm } from "~/components/reservations/creation-form";
import { Button } from "~/components/ui/button";

export default function ReservationNewPage() {
  return (
    <MainLayout
      breadcrumbs={[
        { label: "Reservations", to: "/reservations" },
        { label: "New" },
      ]}
    >
      <div className="flex items-center gap-2">
        <Button asChild variant="ghost" size="icon-sm" aria-label="Back to reservations">
          <Link to="/reservations">
            <ArrowLeft />
          </Link>
        </Button>
        <h1 className="font-heading text-lg font-semibold">
          Create Reservation
        </h1>
      </div>

      <CreationForm />
    </MainLayout>
  );
}
