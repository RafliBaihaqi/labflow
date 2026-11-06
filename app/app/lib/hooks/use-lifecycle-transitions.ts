import { useQuery } from "@tanstack/react-query";
import { getLifecycleTransitions } from "~/lib/api-client";
import { queryKeys } from "~/lib/query-keys";

export function useLifecycleTransitions() {
  return useQuery({
    queryKey: queryKeys.lifecycleTransitions,
    queryFn: getLifecycleTransitions,
    staleTime: 10 * 60 * 1000,
  });
}
