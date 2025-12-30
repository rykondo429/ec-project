'use client';

export const dynamic = 'force-dynamic';

import React, { useEffect } from 'react';
import Link from 'next/link';
import { Button } from '@/components/Button';
import { useCartStore } from '@/store/cart';
import { CartItem } from '@/types';
import { Trash2, ArrowRight } from 'lucide-react';

export default function CartPage() {
  const { cart, fetchCart, updateItem, removeItem, isLoading } = useCartStore();

  useEffect(() => {
    fetchCart();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  if (isLoading) {
    return <div className="max-w-7xl mx-auto px-4 py-16 text-center">読み込み中...</div>;
  }

  if (!cart || !cart.items || cart.items.length === 0) {
    return (
      <div className="max-w-7xl mx-auto px-4 py-16 text-center">
        <h1 className="text-2xl font-bold mb-4">ショッピングカートが空です</h1>
        <p className="text-gray-600 mb-8">商品を追加してください</p>
        <Button asLink href="/products" variant="primary" size="lg">
          ショッピングを続ける <ArrowRight size={20} />
        </Button>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">ショッピングカート</h1>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Cart Items */}
        <div className="lg:col-span-2">
          <div className="bg-white rounded-lg shadow overflow-hidden">
            <table className="w-full">
              <thead className="bg-gray-100 border-b">
                <tr>
                  <th className="text-left px-4 py-3 font-semibold">商品</th>
                  <th className="text-center px-4 py-3 font-semibold">数量</th>
                  <th className="text-right px-4 py-3 font-semibold">価格</th>
                  <th className="text-right px-4 py-3 font-semibold"></th>
                </tr>
              </thead>
              <tbody>
                {cart.items.map((item: CartItem) => (
                  <tr key={item.product_id} className="border-b hover:bg-gray-50">
                    <td className="px-4 py-4">
                      <Link href={`/products/${item.product_id}`} className="text-primary hover:underline">
                        {item.product_name}
                      </Link>
                    </td>
                    <td className="px-4 py-4 text-center">
                      <div className="flex items-center justify-center border rounded w-24 mx-auto">
                        <button
                          onClick={() => updateItem(item.product_id!, Math.max(1, (item.quantity ?? 1) - 1))}
                          className="px-3 py-1 hover:bg-gray-100"
                        >
                          −
                        </button>
                        <span className="px-3 py-1">{item.quantity}</span>
                        <button
                          onClick={() => updateItem(item.product_id!, (item.quantity ?? 0) + 1)}
                          className="px-3 py-1 hover:bg-gray-100"
                        >
                          ＋
                        </button>
                      </div>
                    </td>
                    <td className="px-4 py-4 text-right">
                      ¥{item.subtotal?.toLocaleString()}
                    </td>
                    <td className="px-4 py-4 text-right">
                      <button
                        onClick={() => removeItem(item.product_id!)}
                        className="text-red-600 hover:text-red-800"
                      >
                        <Trash2 size={18} />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <Link href="/products" className="text-primary hover:underline mt-4 inline-block">
            ← 続けてショッピング
          </Link>
        </div>

        {/* Cart Summary */}
        <div>
          <div className="bg-white rounded-lg shadow p-6 sticky top-4">
            <h2 className="text-xl font-bold mb-6">注文確認</h2>

            <div className="space-y-4 mb-6 border-b pb-6">
              <div className="flex justify-between">
                <span>小計</span>
                <span className="font-semibold">¥{cart.total?.toLocaleString()}</span>
              </div>
              <div className="flex justify-between text-gray-600">
                <span>配送料</span>
                <span>¥1,000</span>
              </div>
              <div className="flex justify-between text-gray-600">
                <span>税金</span>
                <span>¥{Math.round((cart.total ?? 0) * 0.1).toLocaleString()}</span>
              </div>
            </div>

            <div className="flex justify-between text-xl font-bold mb-6">
              <span>合計</span>
              <span className="text-primary">¥{((cart.total ?? 0) + 1000 + Math.round((cart.total ?? 0) * 0.1)).toLocaleString()}</span>
            </div>

            <Button asLink href="/checkout" variant="primary" size="lg" className="w-full mb-3">
              会計に進む
            </Button>
            <Button asLink href="/products" variant="outline" size="lg" className="w-full">
              続けてショッピング
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
