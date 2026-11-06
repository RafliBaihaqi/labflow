import { useUIStore } from "~/lib/store/ui-store";

import { EditTestModal } from "../reservations/edit-test-modal";
import { StatusUpdateModal } from "../reservations/status-update-modal";

export function GlobalModalContainer() {
  const activeModal = useUIStore((s) => s.activeModal);
  const data = useUIStore((s) => s.modalData);
  const closeModal = useUIStore((s) => s.closeModal);

  if (!activeModal || !data) return null;

  return (
    <>
      {activeModal === "edit-test" && (
        <EditTestModal
          open={true}
          onOpenChange={(open) => {
            if (!open) closeModal();
          }}
          test={(data as any).test}
          reservationXid={(data as any).reservationXid}
        />
      )}
      {activeModal === "status-update" && (
        <StatusUpdateModal
          open={true}
          onOpenChange={(open) => {
            if (!open) closeModal();
          }}
          xid={(data as any).xid}
          currentStatus={(data as any).currentStatus ?? null}
        />
      )}
    </>
  );
}
