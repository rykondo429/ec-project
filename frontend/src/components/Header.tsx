'use client';

import React from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { ShoppingCart, Heart, Search, User, LogOut } from 'lucide-react';
import { useCartStore } from '@/store/cart';
import { useFavoriteStore } from '@/store/favorite';
import { useAuth } from '@/contexts/AuthContext';

export const Header: React.FC = () => {
  const router = useRouter();
  const cart = useCartStore((state) => state.cart);
  const favorites = useFavoriteStore((state) => state.favorites);
  const { user, loading, signOut } = useAuth();

  const handleSearchClick = () => {
    router.push('/products');
  };

  return (
    <header className="bg-white shadow">
      <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
        <Link href="/" className="text-2xl font-bold text-primary">
          EC Sample
        </Link>

        <nav className="hidden md:flex gap-6 items-center">
          <Link href="/products" className="text-gray-600 hover:text-primary">
            商品一覧
          </Link>
          <Link href="/my-orders" className="text-gray-600 hover:text-primary">
            注文履歴
          </Link>
          <Link href="/my-points" className="text-gray-600 hover:text-primary">
            ポイント
          </Link>
        </nav>

        <div className="flex items-center gap-4">
          <button 
            onClick={handleSearchClick}
            className="p-2 hover:bg-gray-100 rounded-lg"
            title="商品を検索"
          >
            <Search size={20} />
          </button>
          <Link href="/favorites" className="relative p-2 hover:bg-gray-100 rounded-lg" title="お気に入り">
            <Heart size={20} />
            {favorites.length > 0 && (
              <span className="absolute top-0 right-0 bg-red-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
                {favorites.length}
              </span>
            )}
          </Link>
          <Link href="/cart" className="relative p-2 hover:bg-gray-100 rounded-lg">
            <ShoppingCart size={20} />
            {cart && (cart.item_count ?? 0) > 0 && (
              <span className="absolute top-0 right-0 bg-red-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
                {cart.item_count}
              </span>
            )}
          </Link>

          {/* 認証関連 */}
          {!loading && (
            <>
              {user ? (
                <div className="flex items-center gap-2">
                  <Link
                    href="/my-page"
                    className="flex items-center gap-2 px-3 py-2 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors"
                  >
                    <User size={16} />
                    <span className="text-sm font-medium">{user.email}</span>
                  </Link>
                  <button
                    onClick={signOut}
                    className="p-2 hover:bg-gray-100 rounded-lg text-gray-600 hover:text-red-600"
                    title="ログアウト"
                  >
                    <LogOut size={20} />
                  </button>
                </div>
              ) : (
                <div className="flex items-center gap-2">
                  <Link
                    href="/login"
                    className="px-4 py-2 text-sm font-medium text-gray-700 hover:text-primary"
                  >
                    ログイン
                  </Link>
                  <Link
                    href="/signup"
                    className="px-4 py-2 text-sm font-medium text-white bg-primary hover:bg-primary-dark rounded-lg"
                  >
                    新規登録
                  </Link>
                </div>
              )}
            </>
          )}
        </div>
      </div>
    </header>
  );
};
