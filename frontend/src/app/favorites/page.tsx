'use client';

import React from 'react';
import { ProductCard } from '@/components/ProductCard';
import { useFavoriteStore } from '@/store/favorite';
import { Heart } from 'lucide-react';

export default function FavoritesPage() {
  const favorites = useFavoriteStore((state) => state.favorites);

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="mb-8">
        <div className="flex items-center gap-3 mb-2">
          <Heart className="text-red-500" size={32} />
          <h1 className="text-3xl font-bold">お気に入り</h1>
        </div>
        <p className="text-gray-600">お気に入りに追加した商品一覧</p>
      </div>

      {favorites.length === 0 ? (
        <div className="text-center py-16">
          <Heart className="mx-auto mb-4 text-gray-300" size={64} />
          <p className="text-gray-500 text-lg mb-2">お気に入りの商品がありません</p>
          <p className="text-gray-400">気になる商品をお気に入りに追加してみましょう</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {favorites.map((product) => (
            <ProductCard key={product.id} product={product} />
          ))}
        </div>
      )}

      <div className="mt-8 text-center text-sm text-gray-500">
        {favorites.length > 0 && (
          <p>合計 {favorites.length} 件の商品がお気に入りに登録されています</p>
        )}
      </div>
    </div>
  );
}
