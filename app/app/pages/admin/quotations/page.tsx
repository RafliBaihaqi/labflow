import { useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { MainLayout } from "~/components/main/main-layout";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "~/components/ui/select";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "~/components/ui/table";
import { listQuotations, updateQuotationStatus } from "~/lib/api-client";
import { queryKeys } from "~/lib/query-keys";
import { useAuthStore } from "~/lib/store/auth-store";
import type { QuotationStatus } from "~/lib/types/quotation";

const statuses: QuotationStatus[] = ["pending", "contacted", "closed"];

export default function AdminQuotationsPage() {
  const hasPermission = useAuthStore((state) => state.hasPermission);
  const queryClient = useQueryClient();
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const filters = useMemo(() => ({ page, limit: 10, search }), [page, search]);
  const canUpdate = hasPermission("update_quotation");

  const { data, isPending, isError } = useQuery({
    queryKey: queryKeys.quotations(filters),
    queryFn: () => listQuotations(filters),
    enabled: hasPermission("read_quotation"),
  });

  const mutation = useMutation({
    mutationFn: ({ xid, status }: { xid: string; status: QuotationStatus }) =>
      updateQuotationStatus(xid, status),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["quotations"] }),
  });

  if (!hasPermission("read_quotation")) {
    return <MainLayout breadcrumbs={[{ label: "Admin Quotations" }]}>Forbidden</MainLayout>;
  }

  return (
    <MainLayout breadcrumbs={[{ label: "Admin Quotations" }]}> 
      <div className="flex flex-col gap-4">
        <div className="flex gap-2">
          <Input
            placeholder="Search email, company, phone"
            value={search}
            onChange={(event) => {
              setSearch(event.target.value);
              setPage(1);
            }}
          />
          <Button variant="outline" onClick={() => setSearch("")}>Clear</Button>
        </div>
        {isPending ? <p>Loading quotations...</p> : null}
        {isError ? <p className="text-sm text-red-600">Failed to load quotations.</p> : null}
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Company</TableHead>
              <TableHead>Email</TableHead>
              <TableHead>Phone</TableHead>
              <TableHead>Message</TableHead>
              <TableHead>Status</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {(data?.data.data ?? []).map((quotation) => (
              <TableRow key={quotation.xid}>
                <TableCell>{quotation.company_name}</TableCell>
                <TableCell>{quotation.email}</TableCell>
                <TableCell>{quotation.phone_number ?? "-"}</TableCell>
                <TableCell className="max-w-sm truncate">{quotation.message ?? "-"}</TableCell>
                <TableCell>
                  <Select
                    value={quotation.status ?? "pending"}
                    disabled={!canUpdate || mutation.isPending}
                    onValueChange={(status) => mutation.mutate({ xid: quotation.xid, status: status as QuotationStatus })}
                  >
                    <SelectTrigger className="w-36"><SelectValue /></SelectTrigger>
                    <SelectContent>
                      {statuses.map((status) => <SelectItem key={status} value={status}>{status}</SelectItem>)}
                    </SelectContent>
                  </Select>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        <div className="flex items-center justify-between text-sm text-muted-foreground">
          <span>Total: {data?.data.total ?? 0}</span>
          <div className="flex gap-2">
            <Button variant="outline" disabled={page <= 1} onClick={() => setPage((value) => value - 1)}>Previous</Button>
            <Button variant="outline" disabled={page >= (data?.data.total_pages ?? 1)} onClick={() => setPage((value) => value + 1)}>Next</Button>
          </div>
        </div>
      </div>
    </MainLayout>
  );
}
