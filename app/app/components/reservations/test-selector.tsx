import { useState, useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { Loader2, Plus, Trash2 } from "lucide-react";
import type {
  UseFieldArrayAppend,
  UseFieldArrayRemove,
  FieldArrayWithId,
  FieldErrors,
  UseFormRegister,
} from "react-hook-form";

import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { FieldLabel } from "~/components/ui/field";
import { getTestCatalog } from "~/lib/api-client";
import type { TestCatalogItem } from "~/lib/types/reservation";
import type { ReservationFormData } from "~/lib/schemas/reservation";

const MAX_TEST_ITEMS = 50;

const TEST_GROUPS = [
  { id: 1, name: "Microbiology" },
  { id: 2, name: "Chemistry" },
  { id: 3, name: "Physical" },
] as const;

interface TestSelectorProps {
  fields: FieldArrayWithId<ReservationFormData, "tests", "id">[];
  append: UseFieldArrayAppend<ReservationFormData, "tests">;
  remove: UseFieldArrayRemove;
  register: UseFormRegister<ReservationFormData>;
  errors?: FieldErrors<ReservationFormData>["tests"];
}

export function TestSelector({
  fields,
  append,
  remove,
  register,
  errors,
}: TestSelectorProps) {
  const [selectedGroupId, setSelectedGroupId] = useState<number | null>(null);

  const {
    data: catalogResponse,
    isPending,
    isError,
  } = useQuery({
    queryKey: ["test-catalog", selectedGroupId],
    queryFn: () => getTestCatalog(selectedGroupId!),
    enabled: !!selectedGroupId,
    staleTime: 10 * 60 * 1000,
  });

  const catalog = catalogResponse?.data;

  const groupedCatalog = useMemo(() => {
    if (!catalog) return {};
    return catalog.reduce<Record<string, TestCatalogItem[]>>((acc, item) => {
      if (!acc[item.category]) {
        acc[item.category] = [];
      }
      acc[item.category].push(item);
      return acc;
    }, {});
  }, [catalog]);

  const catalogMap = useMemo(() => {
    if (!catalog) return new Map<number, TestCatalogItem>();
    return new Map(catalog.map((c) => [c.id, c]));
  }, [catalog]);

  const handleGroupChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const value = e.target.value;
    setSelectedGroupId(value ? Number(value) : null);
  };

  const handleAddTest = (item: TestCatalogItem) => {
    if (fields.length >= MAX_TEST_ITEMS) return;
    append({
      test_list_id: item.id,
      method: "",
      test_unit: "",
      criteria: "",
    });
  };

  const isMaxReached = fields.length >= MAX_TEST_ITEMS;

  return (
    <div className="flex flex-col gap-4">
      {/* Test Group Selection */}
      <div className="flex flex-col gap-1.5">
        <FieldLabel htmlFor="test-group-select">Test Group</FieldLabel>
        <select
          id="test-group-select"
          value={selectedGroupId ?? ""}
          onChange={handleGroupChange}
          className="h-9 w-full min-w-0 rounded-lg border border-input bg-transparent px-3 py-1.5 text-sm transition-colors outline-none focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/30"
        >
          <option value="">Select a test group...</option>
          {TEST_GROUPS.map((group) => (
            <option key={group.id} value={group.id}>
              {group.name}
            </option>
          ))}
        </select>
      </div>

      {/* Catalog Loading State */}
      {isPending && selectedGroupId && (
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <Loader2 className="size-4 animate-spin" />
          <span>Loading test catalog...</span>
        </div>
      )}

      {/* Catalog Error State */}
      {isError && (
        <div className="rounded-lg border border-destructive/30 bg-destructive/5 px-3 py-2 text-sm text-destructive">
          Test catalog is unavailable. Please try selecting the test group
          again.
        </div>
      )}

      {/* Catalog Items grouped by category */}
      {catalog && !isError && (
        <div className="flex flex-col gap-3">
          {Object.entries(groupedCatalog).map(([category, items]) => (
            <div key={category} className="flex flex-col gap-1.5">
              <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
                {category}
              </span>
              <div className="flex flex-wrap gap-1.5">
                {items.map((item) => (
                  <Button
                    key={item.id}
                    type="button"
                    variant="outline"
                    size="xs"
                    disabled={isMaxReached || isError}
                    onClick={() => handleAddTest(item)}
                    className="cursor-pointer"
                  >
                    <Plus className="size-3" data-icon="inline-start" />
                    {item.name}
                  </Button>
                ))}
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Max items warning */}
      {isMaxReached && (
        <p className="text-sm text-muted-foreground">
          Maximum of {MAX_TEST_ITEMS} test items reached.
        </p>
      )}

      {/* Test Line Items Table */}
      {fields.length > 0 && (
        <div className="flex flex-col gap-2">
          <span className="text-sm font-medium">
            Test Items ({fields.length}/{MAX_TEST_ITEMS})
          </span>
          <div className="overflow-x-auto rounded-lg border border-border">
            <table className="w-full border-collapse text-sm">
              <thead>
                <tr className="border-b border-border bg-muted/50">
                  <th className="px-3 py-2 text-left text-xs font-medium text-muted-foreground">#</th>
                  <th className="px-3 py-2 text-left text-xs font-medium text-muted-foreground">Test</th>
                  <th className="px-3 py-2 text-left text-xs font-medium text-muted-foreground">Method</th>
                  <th className="px-3 py-2 text-left text-xs font-medium text-muted-foreground">Unit</th>
                  <th className="px-3 py-2 text-left text-xs font-medium text-muted-foreground">
                    Criteria
                  </th>
                  <th className="px-3 py-2 text-left text-xs font-medium text-muted-foreground"></th>
                </tr>
              </thead>
              <tbody>
                {fields.map((field, index) => {
                  const testName = catalogMap.get(field.test_list_id)?.name;
                  const itemErrors = errors?.[index];

                  return (
                    <tr
                      key={field.id}
                      className="border-b border-border last:border-b-0 hover:bg-muted/30 transition-colors"
                    >
                      <td className="px-3 py-2 text-muted-foreground">
                        {index + 1}
                      </td>
                      <td className="px-3 py-2 whitespace-nowrap">
                        {testName ?? `Test #${field.test_list_id}`}
                      </td>
                      <td className="px-3 py-2">
                        <div className="flex flex-col gap-0.5">
                          <Input
                            placeholder="Method"
                            aria-invalid={!!itemErrors?.method}
                            {...register(`tests.${index}.method`)}
                          />
                          {itemErrors?.method && (
                            <span className="text-destructive text-xs">
                              {itemErrors.method.message}
                            </span>
                          )}
                        </div>
                      </td>
                      <td className="px-3 py-2">
                        <div className="flex flex-col gap-0.5">
                          <Input
                            placeholder="Unit"
                            aria-invalid={!!itemErrors?.test_unit}
                            {...register(`tests.${index}.test_unit`)}
                          />
                          {itemErrors?.test_unit && (
                            <span className="text-destructive text-xs">
                              {itemErrors.test_unit.message}
                            </span>
                          )}
                        </div>
                      </td>
                      <td className="px-3 py-2">
                        <div className="flex flex-col gap-0.5">
                          <Input
                            placeholder="Criteria"
                            aria-invalid={!!itemErrors?.criteria}
                            {...register(`tests.${index}.criteria`)}
                          />
                          {itemErrors?.criteria && (
                            <span className="text-destructive text-xs">
                              {itemErrors.criteria.message}
                            </span>
                          )}
                        </div>
                      </td>
                      <td className="px-3 py-2">
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon-xs"
                          onClick={() => remove(index)}
                          className="text-destructive hover:text-destructive"
                          aria-label={`Remove test item ${index + 1}`}
                        >
                          <Trash2 className="size-3.5" />
                        </Button>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* Form-level test_items error (e.g., "At least one test is required") */}
      {errors?.root?.message && (
        <p className="text-sm text-destructive">{errors.root.message}</p>
      )}
    </div>
  );
}
