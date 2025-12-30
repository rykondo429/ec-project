'use client';

import React from 'react';

interface SearchResultsHeaderProps {
  totalProducts: number;
  isLoading: boolean;
  searchQuery?: string;
  selectedCategory?: string;
}

export function SearchResultsHeader({
  totalProducts,
  isLoading,
  searchQuery,
  selectedCategory,
}: SearchResultsHeaderProps) {
  return (
    <div className="bg-white rounded-lg shadow-md p-6 mb-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-gray-800 mb-2">商品一覧</h2>
          <div className="flex items-center gap-2 text-sm text-gray-600">
            {searchQuery && (
              <span className="inline-flex items-center">
                <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
                「{searchQuery}」の検索結果
              </span>
            )}
            {selectedCategory && (
              <span className="inline-flex items-center ml-2">
                <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                </svg>
                カテゴリ: {selectedCategory}
              </span>
            )}
          </div>
        </div>
        <div className="text-right">
          <div className="inline-flex items-center bg-blue-50 px-4 py-2 rounded-lg">
            <svg className="w-5 h-5 text-blue-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
            </svg>
            <span className="text-2xl font-bold text-blue-600">
              {isLoading ? '...' : totalProducts.toLocaleString()}
            </span>
            <span className="ml-1 text-sm text-gray-600">件</span>
          </div>
        </div>
      </div>

      {/* Stats Bar */}
      {!isLoading && totalProducts > 0 && (
        <div className="mt-4 pt-4 border-t border-gray-200">
          <div className="grid grid-cols-3 gap-4 text-center">
            <div className="bg-gradient-to-br from-purple-50 to-purple-100 rounded-lg p-3">
              <div className="text-xs text-purple-600 font-medium mb-1">見つかりました</div>
              <div className="text-lg font-bold text-purple-700">{totalProducts}件</div>
            </div>
            <div className="bg-gradient-to-br from-green-50 to-green-100 rounded-lg p-3">
              <div className="text-xs text-green-600 font-medium mb-1">平均価格</div>
              <div className="text-lg font-bold text-green-700">---</div>
            </div>
            <div className="bg-gradient-to-br from-amber-50 to-amber-100 rounded-lg p-3">
              <div className="text-xs text-amber-600 font-medium mb-1">平均評価</div>
              <div className="text-lg font-bold text-amber-700">---</div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
