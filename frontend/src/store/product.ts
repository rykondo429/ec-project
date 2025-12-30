import { create } from 'zustand';
import { Product, ProductWithSKUs } from '@/types';
import { ProductApi } from '@/lib/api-client';

interface ProductState {
  products: Product[];
  currentProduct: Product | null;
  productWithSKUs: ProductWithSKUs | null;
  categories: string[];
  isLoading: boolean;
  error: string | null;
  searchProducts: (query: string, category?: string, sortBy?: string) => Promise<void>;
  getProductById: (id: string) => Promise<void>;
  getProductWithSKUs: (id: string) => Promise<void>;
  getCategories: () => Promise<void>;
}

const productApi = new ProductApi();

export const useProductStore = create<ProductState>((set) => ({
  products: [],
  currentProduct: null,
  productWithSKUs: null,
  categories: [],
  isLoading: false,
  error: null,

  searchProducts: async (query, category, sortBy = 'created_at') => {
    set({ isLoading: true, error: null });
    try {
      const response = await productApi.searchProducts({
        keyword: query,
        category,
        sort_by: sortBy as any,
      });
      set({ products: response.products, isLoading: false });
    } catch (error) {
      set({ error: (error as Error).message, isLoading: false });
    }
  },

  getProductById: async (id) => {
    set({ isLoading: true, error: null });
    try {
      const product = await productApi.getProduct(id);
      set({ currentProduct: product, isLoading: false });
    } catch (error) {
      set({ error: (error as Error).message, isLoading: false });
    }
  },

  getProductWithSKUs: async (id) => {
    set({ isLoading: true, error: null });
    try {
      const productWithSKUs = await productApi.getProductWithSKUs(id);
      set({ productWithSKUs, isLoading: false });
    } catch (error) {
      set({ error: (error as Error).message, isLoading: false });
    }
  },

  getCategories: async () => {
    try {
      const categories = await productApi.getCategories();
      set({ categories });
    } catch (error) {
      set({ error: (error as Error).message });
    }
  },
}));
