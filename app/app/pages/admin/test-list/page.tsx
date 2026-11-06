import { useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { MainLayout } from "~/components/main/main-layout";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "~/components/ui/table";
import { createTestList, deactivateTestList, listTestLists, updateTestList } from "~/lib/api-client";
import { queryKeys } from "~/lib/query-keys";
import { useAuthStore } from "~/lib/store/auth-store";
import type { TestListInput, TestListItem } from "~/lib/types/test-list";

const emptyForm: TestListInput = { category_id: 1, value: "", label: "" };

function numberOrUndefined(value: string) {
  return value === "" ? undefined : Number(value);
}

export default function AdminTestListPage() {
  const hasPermission = useAuthStore((state) => state.hasPermission);
  const queryClient = useQueryClient();
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const [editing, setEditing] = useState<TestListItem | null>(null);
  const [form, setForm] = useState<TestListInput>(emptyForm);
  const filters = useMemo(() => ({ page, limit: 10, search }), [page, search]);

  const { data, isPending, isError } = useQuery({
    queryKey: queryKeys.testLists(filters),
    queryFn: () => listTestLists(filters),
    enabled: hasPermission("read_test"),
  });

  const invalidate = () => queryClient.invalidateQueries({ queryKey: ["test-lists"] });
  const createMutation = useMutation({ mutationFn: createTestList, onSuccess: invalidate });
  const updateMutation = useMutation({ mutationFn: ({ xid, input }: { xid: string; input: TestListInput }) => updateTestList(xid, input), onSuccess: invalidate });
  const deleteMutation = useMutation({ mutationFn: deactivateTestList, onSuccess: invalidate });

  const resetForm = () => {
    setEditing(null);
    setForm(emptyForm);
  };

  const submit = () => {
    if (editing) {
      updateMutation.mutate({ xid: editing.xid, input: form });
    } else {
      createMutation.mutate(form);
    }
    resetForm();
  };

  if (!hasPermission("read_test")) {
    return <MainLayout breadcrumbs={[{ label: "Test List" }]}>Forbidden</MainLayout>;
  }

  return (
    <MainLayout breadcrumbs={[{ label: "Test List" }]}> 
      <div className="grid gap-6">
        <div className="grid gap-3 rounded-lg border p-4">
          <h2 className="font-heading text-lg font-semibold">{editing ? "Create New Version" : "Create Test"}</h2>
          <div className="grid gap-3 md:grid-cols-3">
            <div className="space-y-1">
              <Label>Category ID</Label>
              <Input type="number" value={form.category_id} onChange={(event) => setForm({ ...form, category_id: Number(event.target.value) })} />
            </div>
            <div className="space-y-1">
              <Label>Value</Label>
              <Input value={form.value} onChange={(event) => setForm({ ...form, value: event.target.value })} />
            </div>
            <div className="space-y-1">
              <Label>Label</Label>
              <Input value={form.label} onChange={(event) => setForm({ ...form, label: event.target.value })} />
            </div>
            <div className="space-y-1">
              <Label>Liquid Amount</Label>
              <Input type="number" value={form.sample_amount_liquid ?? ""} onChange={(event) => setForm({ ...form, sample_amount_liquid: numberOrUndefined(event.target.value) })} />
            </div>
            <div className="space-y-1">
              <Label>Solid Amount</Label>
              <Input type="number" value={form.sample_amount_solid ?? ""} onChange={(event) => setForm({ ...form, sample_amount_solid: numberOrUndefined(event.target.value) })} />
            </div>
            <div className="space-y-1">
              <Label>Powder Amount</Label>
              <Input type="number" value={form.sample_amount_powder ?? ""} onChange={(event) => setForm({ ...form, sample_amount_powder: numberOrUndefined(event.target.value) })} />
            </div>
            <div className="space-y-1">
              <Label>Leadtime</Label>
              <Input type="number" value={form.leadtime ?? ""} onChange={(event) => setForm({ ...form, leadtime: numberOrUndefined(event.target.value) })} />
            </div>
          </div>
          <div className="flex gap-2">
            <Button disabled={!hasPermission(editing ? "update_test" : "create_test")} onClick={submit}>{editing ? "Save New Version" : "Create"}</Button>
            {editing && <Button variant="outline" onClick={resetForm}>Cancel</Button>}
          </div>
        </div>

        <div className="flex gap-2">
          <Input placeholder="Search value or label" value={search} onChange={(event) => { setSearch(event.target.value); setPage(1); }} />
          <Button variant="outline" onClick={() => setSearch("")}>Clear</Button>
        </div>
        {isPending ? <p>Loading test list...</p> : null}
        {isError ? <p className="text-sm text-red-600">Failed to load test list.</p> : null}
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Label</TableHead>
              <TableHead>Category</TableHead>
              <TableHead>Version</TableHead>
              <TableHead>Current</TableHead>
              <TableHead>Leadtime</TableHead>
              <TableHead>Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {(data?.data.data ?? []).map((item) => (
              <TableRow key={item.xid}>
                <TableCell>{item.label}</TableCell>
                <TableCell>{item.category_name}</TableCell>
                <TableCell>{item.version}</TableCell>
                <TableCell>{item.is_current ? "Yes" : "No"}</TableCell>
                <TableCell>{item.leadtime ?? "-"}</TableCell>
                <TableCell className="flex gap-2">
                  <Button size="sm" variant="outline" disabled={!hasPermission("update_test")} onClick={() => { setEditing(item); setForm({ category_id: item.category_id, value: item.value, label: item.label, sample_amount_liquid: item.sample_amount_liquid, sample_amount_solid: item.sample_amount_solid, sample_amount_powder: item.sample_amount_powder, leadtime: item.leadtime }); }}>Edit</Button>
                  <Button size="sm" variant="destructive" disabled={!hasPermission("delete_test") || !item.is_current} onClick={() => deleteMutation.mutate(item.xid)}>Deactivate</Button>
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
