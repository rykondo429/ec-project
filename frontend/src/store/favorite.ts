import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { Product } from '@/types';

interface FavoriteState {
  favorites: Product[];
  addFavorite: (product: Product) => void;
  removeFavorite: (productId: string) => void;
  isFavorite: (productId: string) => boolean;
  toggleFavorite: (product: Product) => void;
}

export const useFavoriteStore = create<FavoriteState>()(
  persist(
    (set, get) => ({
      favorites: [],

      addFavorite: (product) => {
        set((state) => {
          if (state.favorites.some((p) => p.id === product.id)) {
            return state;
          }
          return { favorites: [...state.favorites, product] };
        });
      },

      removeFavorite: (productId) => {
        set((state) => ({
          favorites: state.favorites.filter((p) => p.id !== productId),
        }));
      },

      isFavorite: (productId) => {
        return get().favorites.some((p) => p.id === productId);
      },

      toggleFavorite: (product) => {
        const state = get();
        if (state.isFavorite(product.id)) {
          state.removeFavorite(product.id);
        } else {
          state.addFavorite(product);
        }
      },
    }),
    {
      name: 'favorite-storage',
    }
  )
);
