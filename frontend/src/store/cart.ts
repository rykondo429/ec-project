import { create } from 'zustand';
import { Cart } from '@/types';
import { CartApi } from '@/lib/api-client';

interface CartState {
  cart: Cart | null;
  isLoading: boolean;
  error: string | null;
  fetchCart: () => Promise<void>;
  addItem: (
    productId: string,
    productName: string,
    price: number,
    quantity: number,
    skuId?: string,
    skuCode?: string,
    skuAttributes?: Record<string, string>
  ) => Promise<void>;
  updateItem: (productId: string, quantity: number) => Promise<void>;
  removeItem: (productId: string) => Promise<void>;
  clearCart: () => Promise<void>;
}

const cartApi = new CartApi();

export const useCartStore = create<CartState>((set) => ({
  cart: null,
  isLoading: false,
  error: null,

  fetchCart: async () => {
    set({ isLoading: true, error: null });
    try {
      const cart = await cartApi.getCart();
      set({ cart, isLoading: false });
    } catch (error) {
      set({ error: (error as Error).message, isLoading: false });
    }
  },

  addItem: async (productId, productName, price, quantity, skuId?, skuCode?, skuAttributes?) => {
    set({ isLoading: true, error: null });
    try {
      const cart = await cartApi.addItem({
        product_id: productId,
        product_name: productName,
        price,
        quantity,
        ...(skuId && { sku_id: skuId }),
        ...(skuCode && { sku_code: skuCode }),
        ...(skuAttributes && { sku_attributes: skuAttributes }),
      });
      set({ cart, isLoading: false });
    } catch (error) {
      set({ error: (error as Error).message, isLoading: false });
    }
  },

  updateItem: async (productId, quantity) => {
    set({ isLoading: true, error: null });
    try {
      const cart = await cartApi.updateItem(productId, quantity);
      set({ cart, isLoading: false });
    } catch (error) {
      set({ error: (error as Error).message, isLoading: false });
    }
  },

  removeItem: async (productId) => {
    set({ isLoading: true, error: null });
    try {
      const cart = await cartApi.removeItem(productId);
      set({ cart, isLoading: false });
    } catch (error) {
      set({ error: (error as Error).message, isLoading: false });
    }
  },

  clearCart: async () => {
    set({ isLoading: true, error: null });
    try {
      await cartApi.clearCart();
      set({ cart: null, isLoading: false });
    } catch (error) {
      set({ error: (error as Error).message, isLoading: false });
    }
  },
}));
