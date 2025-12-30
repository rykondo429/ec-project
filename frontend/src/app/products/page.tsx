'use client';

export const dynamic = 'force-dynamic';

import React, { useEffect, useState, Suspense, useCallback } from 'react';
import { ProductCard } from '@/components/ProductCard';
import { SearchFilter } from '@/components/SearchFilter';
import { SearchResultsHeader } from '@/components/SearchResultsHeader';
import { useProductStore } from '@/store/product';
import { useSearchParams } from 'next/navigation';

function ProductsPageContent() {
  const searchParams = useSearchParams();
  const { products, categories, totalProducts, searchProducts, getCategories, isLoading } = useProductStore();
  
  const [selectedCategory, setSelectedCategory] = useState(searchParams.get('category') || '');
  const [searchQuery, setSearchQuery] = useState(searchParams.get('q') || '');
  const [sortBy, setSortBy] = useState('created_at');
  const [minPrice, setMinPrice] = useState(0);
  const [maxPrice, setMaxPrice] = useState(1000000);
  const [minRating, setMinRating] = useState(0);
  const [searchDebounce, setSearchDebounce] = useState<NodeJS.Timeout | null>(null);

  useEffect(() => {
    getCategories();
  }, []);

  // 検索実行
  const executeSearch = useCallback(() => {
    searchProducts({
      keyword: searchQuery,
      category: selectedCategory,
      sortBy,
      minPrice,
      maxPrice,
      minRating,
    });
  }, [searchQuery, selectedCategory, sortBy, minPrice, maxPrice, minRating]);

  // 初回ロードと依存関係変更時の検索
  useEffect(() => {
    executeSearch();
  }, [selectedCategory, sortBy, minPrice, maxPrice, minRating]);

  // デバウンス付き検索
  const handleSearchChange = (value: string) => {
    setSearchQuery(value);
    
    if (searchDebounce) {
      clearTimeout(searchDebounce);
    }
    
    const timeout = setTimeout(() => {
      searchProducts({
        keyword: value,
        category: selectedCategory,
        sortBy,
        minPrice,
        maxPrice,
        minRating,
      });
    }, 500);
    
    setSearchDebounce(timeout);
  };

  const handleClearFilters = () => {
    setSearchQuery('');
    setSelectedCategory('');
    setSortBy('created_at');
    setMinPrice(0);
    setMaxPrice(1000000);
    setMinRating(0);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      <div className="max-w-7xl mx-auto px-4 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          {/* Sidebar */}
          <div className="lg:col-span-1">
            <SearchFilter
              searchQuery={searchQuery}
              onSearchChange={handleSearchChange}
              selectedCategory={selectedCategory}
              onCategoryChange={setSelectedCategory}
              categories={categories}
              sortBy={sortBy}
              onSortChange={setSortBy}
              minPrice={minPrice}
              maxPrice={maxPrice}
              onPriceChange={(min, max) => {
                setMinPrice(min);
                setMaxPrice(max);
              }}
              minRating={minRating}
              onRatingChange={setMinRating}
              onClearFilters={handleClearFilters}
              totalProducts={totalProducts}
              isLoading={isLoading}
            />
          </div>

          {/* Products Grid */}
          <div className="lg:col-span-3">
            <SearchResultsHeader
              totalProducts={totalProducts}
              isLoading={isLoading}
              searchQuery={searchQuery}
              selectedCategory={selectedCategory}
            />

            {isLoading ? (
              <div className="text-center py-20">
                <div className="inline-flex flex-col items-center">
                  <svg
                    className="animate-spin h-12 w-12 text-blue-600 mb-4"
                    fill="none"
                    viewBox="0 0 24 24"
                  >
                    <circle
                      className="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      strokeWidth="4"
                    ></circle>
                    <path
                      className="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    ></path>
                  </svg>
                  <p className="text-gray-600 font-medium">商品を検索中...</p>
                </div>
              </div>
            ) : products.length === 0 ? (
              <div className="bg-white rounded-lg shadow-md p-12 text-center">
                <svg
                  className="mx-auto h-24 w-24 text-gray-400 mb-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
                <h3 className="text-xl font-semibold text-gray-700 mb-2">
                  商品が見つかりませんでした
                </h3>
                <p className="text-gray-500 mb-4">
                  検索条件を変更してもう一度お試しください
                </p>
                <button
                  onClick={handleClearFilters}
                  className="inline-flex items-center px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
                >
                  <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                  </svg>
                  フィルターをリセット
                </button>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
                {products.map((product) => (
                  <ProductCard key={product.id} product={product} />
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default function ProductsPage() {
  return (
    <Suspense fallback={
      <div className="flex justify-center items-center min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
        <div className="text-center">
          <svg
            className="animate-spin h-16 w-16 text-blue-600 mx-auto mb-4"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              className="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              strokeWidth="4"
            ></circle>
            <path
              className="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
          </svg>
          <p className="text-gray-600 font-medium text-lg">読み込み中...</p>
        </div>
      </div>
    }>
      <ProductsPageContent />
    </Suspense>
  );
}
