import { useState, useEffect } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { AxiosError } from "axios";
import { Save, Loader2, Edit } from "lucide-react";
import { toast } from "sonner";

import { Button } from "~/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "~/components/ui/dialog";
import { Label } from "~/components/ui/label";
import { Input } from "~/components/ui/input";
import { Textarea } from "~/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "~/components/ui/select";
import { updateTest } from "~/lib/api-client";
import type { TestItemResponse } from "~/lib/types/reservation";

interface EditTestModalProps {
  test: TestItemResponse;
  reservationXid: string;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
}

export function EditTestModal({
  test,
  reservationXid,
  open: controlledOpen,
  onOpenChange,
}: EditTestModalProps) {
  const [uncontrolledOpen, setUncontrolledOpen] = useState(false);
  const open = controlledOpen !== undefined ? controlledOpen : uncontrolledOpen;
  const setOpen = onOpenChange !== undefined ? onOpenChange : setUncontrolledOpen;

  const [testResult, setTestResult] = useState(test?.test_result || "");
  const [testSummary, setTestSummary] = useState(test?.test_summary || "");
  const [notes, setNotes] = useState(test?.notes || "");
  const queryClient = useQueryClient();

  useEffect(() => {
    if (open && test) {
      setTestResult(test.test_result || "");
      setTestSummary(test.test_summary || "");
      setNotes(test.notes || "");
    }
  }, [open, test]);

  const mutation = useMutation({
    mutationFn: async () => {
      if (!test?.xid) throw new Error("Test is required");
      return updateTest(test.xid, {
        test_result: testResult || null,
        test_summary: testSummary || null,
        notes: notes || null,
      });
    },
    onSuccess: () => {
      toast.success("Test updated successfully");
      queryClient.invalidateQueries({ queryKey: ["reservation", reservationXid] });
      setOpen(false);
    },
    onError: (error) => {
      const message =
        error instanceof AxiosError
          ? error.response?.data?.message || "Failed to update test"
          : "An unexpected error occurred";
      toast.error(message);
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    mutation.mutate();
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      {controlledOpen === undefined && (
        <DialogTrigger asChild>
          <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
            <span className="sr-only">Edit</span>
            <Edit className="h-4 w-4" />
          </Button>
        </DialogTrigger>
      )}
      <DialogContent className="sm:max-w-[425px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Edit Test Item</DialogTitle>
            <DialogDescription>
              Update the test result and summary for {test?.test_list_name}.
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="test_result" className="text-right">
                Result
              </Label>
              <Input
                id="test_result"
                value={testResult}
                onChange={(e) => setTestResult(e.target.value)}
                className="col-span-3"
                placeholder="e.g., 5.4"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="test_summary" className="text-right">
                Summary
              </Label>
              <div className="col-span-3">
                <Select value={testSummary} onValueChange={setTestSummary}>
                  <SelectTrigger id="test_summary">
                    <SelectValue placeholder="Select summary" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="Comply">Comply</SelectItem>
                    <SelectItem value="Not Comply">Not Comply</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
            <div className="grid grid-cols-4 items-start gap-4">
              <Label htmlFor="notes" className="text-right pt-2">
                Notes
              </Label>
              <Textarea
                id="notes"
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                className="col-span-3 resize-none"
                placeholder="Optional notes"
                rows={3}
              />
            </div>
          </div>
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => setOpen(false)}
              disabled={mutation.isPending}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={mutation.isPending}>
              {mutation.isPending ? (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              ) : (
                <Save className="mr-2 h-4 w-4" />
              )}
              Save changes
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
