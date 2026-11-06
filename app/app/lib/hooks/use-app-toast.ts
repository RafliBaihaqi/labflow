import { toast as sonnerToast } from "sonner";

export function useAppToast() {
  return {
    success: (msg: string) => sonnerToast.success(msg),
    error: (msg: string) => sonnerToast.error(msg),
    info: (msg: string) => sonnerToast.info(msg),
    warning: (msg: string) => sonnerToast.warning(msg),
  };
}

export const toast = sonnerToast;
