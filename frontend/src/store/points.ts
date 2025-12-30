import { create } from 'zustand';
import { PointBalance, PointTransaction } from '@/types';
import { PointApi } from '@/lib/api-client';

interface PointState {
  balance: PointBalance | null;
  transactions: PointTransaction[];
  isLoading: boolean;
  error: string | null;
  fetchBalance: () => Promise<void>;
  fetchTransactions: (page?: number) => Promise<void>;
}

const pointApi = new PointApi();

export const usePointStore = create<PointState>((set) => ({
  balance: null,
  transactions: [],
  isLoading: false,
  error: null,

  fetchBalance: async () => {
    set({ isLoading: true, error: null });
    try {
      const balance = await pointApi.getBalance();
      set({ balance, isLoading: false });
    } catch (error) {
      set({ error: (error as Error).message, isLoading: false });
    }
  },

  fetchTransactions: async (page = 1) => {
    set({ isLoading: true, error: null });
    try {
      const result = await pointApi.getTransactions(page);
      set({ transactions: result.transactions, isLoading: false });
    } catch (error) {
      set({ error: (error as Error).message, isLoading: false });
    }
  },
}));
