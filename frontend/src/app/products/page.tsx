'use client';

export const dynamic = 'force-dynamic';

import React, { useEffect, useState, Suspense } from 'react';
import { ProductCard } from '@/components/ProductCard';
import { Button } from '@/components/Button';
import { useProductStore } from '@/store/product';
import { useSearchParams } from 'next/navigation';

function ProductsPageContent() {
  const searchParams = useSearchParams();
  const { products, categories, searchProducts, getCategories, isLoading } = useProductStore();
  const [selectedCategory, setSelectedCategory] = useState(searchParams.get('category') || '');
  const [searchQuery, setSearchQuery] = useState(searchParams.get('q') || '');
  const [sortBy, setSortBy] = useState('created_at');

  useEffect(() => {
    getCategories();
    searchProducts(searchQuery, selectedCategory, sortBy);
  }, [selectedCategory, sortBy]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    searchProducts(searchQuery, selectedCategory, sortBy);
  };

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
        {/* Sidebar */}
        <div className="md:col-span-1">
          <div className="bg-white rounded-lg p-6 shadow">
            <h3 className="text-lg font-bold mb-4">フィルター</h3>

            {/* Search */}
            <form onSubmit={handleSearch} className="mb-6">
              <input
                type="text"
                placeholder="商品名で検索"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full px-3 py-2 border rounded-lg mb-3"
              />
              <Button type="submit" variant="primary" size="sm" className="w-full">
                検索
              </Button>
            </form>

            {/* Categories */}
            <div className="mb-6">
              <h4 className="font-semibold mb-3">カテゴリ</h4>
              <div className="space-y-2">
                <label className="flex items-center">
                  <input
                    type="radio"
                    name="category"
                    value=""
                    checked={selectedCategory === ''}
                    onChange={() => setSelectedCategory('')}
                    className="mr-2"
                  />
                  <span>すべて</span>
                </label>
                {categories.map((cat) => (
                  <label key={cat} className="flex items-center">
                    <input
                      type="radio"
                      name="category"
                      value={cat}
                      checked={selectedCategory === cat}
                      onChange={() => setSelectedCategory(cat)}
                      className="mr-2"
                    />
                    <span>{cat}</span>
                  </label>
                ))}
              </div>
            </div>

            {/* Sort */}
            <div>
              <h4 className="font-semibold mb-3">ソート</h4>
              <select
                value={sortBy}
                onChange={(e) => setSortBy(e.target.value)}
                className="w-full px-3 py-2 border rounded-lg"
              >
                <option value="created_at">新着順</option>
                <option value="price">価格: 低い順</option>
                <option value="rating">評価が高い順</option>
              </select>
            </div>
          </div>
        </div>

        {/* Products Grid */}
        <div className="md:col-span-3">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-2xl font-bold">商品一覧</h2>
            <span className="text-gray-600">全 {products.length} 件</span>
          </div>

          {isLoading ? (
            <div className="text-center py-12">
              <p>読み込み中...</p>
            </div>
          ) : products.length === 0 ? (
            <div className="text-center py-12">
              <p className="text-gray-600">商品が見つかりません</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {products.map((product) => (
                <ProductCard key={product.id} product={product} />
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default function ProductsPage() {
  return (
    <Suspense fallback={<div className="flex justify-center items-center min-h-screen">読み込み中...</div>}>
      <ProductsPageContent />
    </Suspense>
  );
}
