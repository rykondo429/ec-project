import { create } from 'zustand';
import { Product, ProductWithSKUs } from '@/types';
import { ProductApi } from '@/lib/api-client';

interface SearchParams {
  keyword?: string;
  category?: string;
  sortBy?: string;
  minPrice?: number;
  maxPrice?: number;
  minRating?: number;
}

interface ProductState {
  products: Product[];
  currentProduct: Product | null;
  productWithSKUs: ProductWithSKUs | null;
  categories: string[];
  totalProducts: number;
  isLoading: boolean;
  error: string | null;
  searchProducts: (params: SearchParams) => Promise<void>;
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
  totalProducts: 0,
  isLoading: false,
  error: null,

  searchProducts: async (params: SearchParams) => {
    set({ isLoading: true, error: null });
    try {
      const response = await productApi.searchProducts({
        keyword: params.keyword,
        category: params.category,
        sort_by: params.sortBy as 'price' | 'rating' | 'created_at' | undefined,
      });
      
      // クライアントサイドでの追加フィルタリング
      let filteredProducts = response.products || [];
      
      // 価格フィルター
      if (params.minPrice !== undefined && params.minPrice > 0) {
        filteredProducts = filteredProducts.filter(p => p.price >= params.minPrice!);
      }
      if (params.maxPrice !== undefined && params.maxPrice < 1000000) {
        filteredProducts = filteredProducts.filter(p => p.price <= params.maxPrice!);
      }
      
      // レーティングフィルター
      if (params.minRating !== undefined && params.minRating > 0) {
        filteredProducts = filteredProducts.filter(p => (p.rating || 0) >= params.minRating!);
      }
      
      set({ 
        products: filteredProducts, 
        totalProducts: filteredProducts.length,
        isLoading: false 
      });
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
