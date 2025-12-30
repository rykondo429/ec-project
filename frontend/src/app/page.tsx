'use client';

export const dynamic = 'force-dynamic';

import React, { useEffect } from 'react';
import Link from 'next/link';
import { Button } from '@/components/Button';
import { ArrowRight } from 'lucide-react';
import { useProductStore } from '@/store/product';

export default function Home() {
  const { products, searchProducts, getCategories } = useProductStore();

  useEffect(() => {
    searchProducts({ keyword: '', sortBy: 'created_at' });
    getCategories();
  }, []);

  return (
    <div>
      {/* Hero Section */}
      <section className="bg-gradient-to-r from-primary to-blue-600 text-white py-16">
        <div className="max-w-7xl mx-auto px-4 text-center">
          <h1 className="text-4xl md:text-5xl font-bold mb-4">Welcome to EC Sample</h1>
          <p className="text-lg mb-8 opacity-90">最高品質の商品を手頃な価格でお届け</p>
          <Button asLink href="/products" variant="secondary" size="lg">
            ショッピングを開始 <ArrowRight size={20} />
          </Button>
        </div>
      </section>

      {/* Featured Products */}
      <section className="max-w-7xl mx-auto px-4 py-16">
        <h2 className="text-3xl font-bold mb-8">新着商品</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {products.slice(0, 4).map((product) => (
            <Link key={product.id} href={`/products/${product.id}`}>
              <div className="bg-white rounded-lg shadow hover:shadow-lg transition-shadow p-4 cursor-pointer">
                <div className="aspect-video bg-gray-200 rounded mb-4 flex items-center justify-center">
                  <span className="text-gray-500">No Image</span>
                </div>
                <h3 className="font-semibold mb-2 line-clamp-2">{product.name}</h3>
                <p className="text-primary font-bold text-lg">¥{product.price.toLocaleString()}</p>
              </div>
            </Link>
          ))}
        </div>
      </section>

      {/* Categories Section */}
      <section className="bg-white py-16 border-t">
        <div className="max-w-7xl mx-auto px-4">
          <h2 className="text-3xl font-bold mb-8">カテゴリから探す</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {['Electronics', 'Fashion', 'Home'].map((category) => (
              <Link key={category} href={`/products?category=${category}`}>
                <div className="bg-gray-100 hover:bg-gray-200 rounded-lg p-8 text-center cursor-pointer transition-colors">
                  <h3 className="text-xl font-bold">{category}</h3>
                </div>
              </Link>
            ))}
          </div>
        </div>
      </section>
    </div>
  );
}
