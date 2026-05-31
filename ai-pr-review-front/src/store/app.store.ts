import { create } from 'zustand';

interface AppState {
  apiHealthy: boolean | null;
  isSubmitting: boolean;
  setApiHealthy: (healthy: boolean) => void;
  setIsSubmitting: (submitting: boolean) => void;
}

export const useAppStore = create<AppState>((set) => ({
  apiHealthy: null,
  isSubmitting: false,
  setApiHealthy: (healthy) => set({ apiHealthy: healthy }),
  setIsSubmitting: (submitting) => set({ isSubmitting: submitting }),
}));
