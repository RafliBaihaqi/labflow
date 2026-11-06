import { useQuery } from "@tanstack/react-query";
import { AxiosError } from "axios";
import { Link } from "react-router";
import { ArrowLeft, TriangleAlert, FileQuestion } from "lucide-react";

import { getReservation } from "~/lib/api-client";
import { MainLayout } from "~/components/main/main-layout";
import { DetailView } from "~/components/reservations/detail-view";
import { Button } from "~/components/ui/button";
import { Skeleton } from "~/components/ui/skeleton";
import { queryKeys } from "~/lib/query-keys";

function DetailSkeleton() {
  return (
    <div className="space-y-4">
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {Array.from({ length: 6 }).map((_, i) => (
          <Skeleton key={i} className="h-16 w-full" />
        ))}
      </div>
      <Skeleton className="h-48 w-full" />
    </div>
  );
}

function NotFound() {
  return (
    <div className="flex flex-col items-center justify-center gap-3 py-16">
      <FileQuestion className="size-8 text-muted-foreground" />
      <p className="text-sm text-muted-foreground">Reservation not found</p>
      <Button asChild variant="outline" size="sm">
        <Link to="/reservations">Back to reservations</Link>
      </Button>
    </div>
  );
}

function ErrorState() {
  return (
    <div className="flex flex-col items-center justify-center gap-3 py-16">
      <TriangleAlert className="size-8 text-muted-foreground" />
      <p className="text-sm text-muted-foreground">
        Something went wrong. Please try again later.
      </p>
    </div>
  );
}

export default function ReservationDetailPage({ xid }: { xid: string }) {
  const { data, isPending, isError, error } = useQuery({
    queryKey: queryKeys.reservation(xid),
    queryFn: () => getReservation(xid),
    retry: (failureCount, err) => {
      if (err instanceof AxiosError) {
        if (err.response && err.response.status >= 400 && err.response.status < 500) {
          return false;
        }
        if (err.code === "ECONNABORTED" || (err.response && err.response.status >= 500)) {
          return failureCount < 1;
        }
      }
      return failureCount < 1;
    },
  });

  const reservation = data?.data;
  const title = reservation?.sample_name || "Reservation";

  let content: React.ReactNode;
  if (isPending) {
    content = <DetailSkeleton />;
  } else if (isError) {
    content =
      error instanceof AxiosError && error.response?.status === 404 ? (
        <NotFound />
      ) : (
        <ErrorState />
      );
  } else {
    content = <DetailView reservation={reservation!} />;
  }

  return (
    <MainLayout
      breadcrumbs={[
        { label: "Reservations", to: "/reservations" },
        { label: isPending ? "Detail" : title },
      ]}
    >
      <div className="flex items-center gap-2">
        <Button asChild variant="ghost" size="icon-sm" aria-label="Back to reservations">
          <Link to="/reservations">
            <ArrowLeft />
          </Link>
        </Button>
        <h1 className="font-heading text-lg font-semibold">
          {isPending ? "Reservation" : title}
        </h1>
      </div>

      {content}
    </MainLayout>
  );
}
