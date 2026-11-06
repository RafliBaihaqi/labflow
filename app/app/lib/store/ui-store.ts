import { create } from "zustand";
import type { TestItemResponse } from "~/lib/types/reservation";

export type ModalPayloads = {
  "edit-test": {
    test: TestItemResponse;
    reservationXid: string;
  };
  "status-update": {
    xid: string;
    currentStatus: string | null;
  };
  "certificate-preview": {
    xid: string;
  };
  "feedback-form": {
    certificateId: string;
    reservationId: string;
  };
  export: {
    filters: Record<string, string>;
  };
};

export type ModalType = keyof ModalPayloads | null;

type ModalData<T extends ModalType> = T extends keyof ModalPayloads ? ModalPayloads[T] : null;

interface UIState {
  activeModal: ModalType;
  modalData: ModalData<ModalType> | null;
  openModal: <T extends ModalType>(
    type: T,
    data: T extends keyof ModalPayloads ? ModalPayloads[T] : null
  ) => void;
  closeModal: () => void;
}

export const useUIStore = create<UIState>((set) => ({
  activeModal: null,
  modalData: null,
  openModal: <T extends ModalType>(
    activeModal: T,
    modalData: (T extends keyof ModalPayloads ? ModalPayloads[T] : null) | null = null
  ) => set({ activeModal, modalData: modalData as ModalData<ModalType> }),
  closeModal: () => set({ activeModal: null, modalData: null }),
}));
