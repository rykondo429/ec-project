'use client';

export const dynamic = 'force-dynamic';

import React, { useState } from 'react';
import { Button } from '@/components/Button';
import { useCartStore } from '@/store/cart';
import { OrderApi } from '@/lib/api-client';
import { useRouter } from 'next/navigation';
import { CartItem } from '@/types';

export default function CheckoutPage() {
  const router = useRouter();
  const { cart, clearCart } = useCartStore();
  const [formData, setFormData] = useState({
    shippingAddress: '',
    paymentMethod: 'credit_card',
    couponCode: '',
  });
  const [isProcessing, setIsProcessing] = useState(false);
  const [error, setError] = useState('');

  const orderApi = new OrderApi();

  if (!cart || !cart.items || cart.items.length === 0) {
    return (
      <div className="max-w-7xl mx-auto px-4 py-16 text-center">
        <p>カートが空です。商品を追加してください。</p>
        <Button asLink href="/products" variant="primary" className="mt-4">
          ショッピングに戻る
        </Button>
      </div>
    );
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsProcessing(true);
    setError('');

    try {
      const order = await orderApi.createOrder({
        items: (cart.items ?? []).map((item: CartItem) => ({
          product_id: item.product_id!,
          quantity: item.quantity!,
        })),
        shipping_address: formData.shippingAddress,
        payment_method: formData.paymentMethod,
        coupon_code: formData.couponCode,
      });

      await clearCart();
      router.push(`/orders/${order.id}`);
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setIsProcessing(false);
    }
  };

  const totalAmount = (cart.total ?? 0) + 1000 + Math.round((cart.total ?? 0) * 0.1);

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">注文確認</h1>

      {error && <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">{error}</div>}

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <form onSubmit={handleSubmit} className="lg:col-span-2 bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-bold mb-6">配送先</h2>

          <div className="mb-6">
            <label className="block font-semibold mb-2">配送先住所 *</label>
            <textarea
              value={formData.shippingAddress}
              onChange={(e) => setFormData({ ...formData, shippingAddress: e.target.value })}
              required
              rows={4}
              className="w-full px-4 py-2 border rounded-lg"
              placeholder="都道府県市区町村番地"
            />
          </div>

          <h2 className="text-xl font-bold mb-6 mt-8">支払い方法</h2>

          <div className="space-y-4 mb-6">
            {[
              { value: 'credit_card', label: 'クレジットカード' },
              { value: 'bank_transfer', label: '銀行振込' },
              { value: 'convenience_store', label: 'コンビニ決済' },
            ].map((method) => (
              <label key={method.value} className="flex items-center border rounded p-4 cursor-pointer">
                <input
                  type="radio"
                  name="paymentMethod"
                  value={method.value}
                  checked={formData.paymentMethod === method.value}
                  onChange={(e) => setFormData({ ...formData, paymentMethod: e.target.value })}
                  className="mr-3"
                />
                <span>{method.label}</span>
              </label>
            ))}
          </div>

          <h2 className="text-xl font-bold mb-6 mt-8">クーポンコード</h2>

          <div className="mb-6">
            <input
              type="text"
              value={formData.couponCode}
              onChange={(e) => setFormData({ ...formData, couponCode: e.target.value })}
              placeholder="クーポンコードを入力（オプション）"
              className="w-full px-4 py-2 border rounded-lg"
            />
          </div>

          <Button type="submit" disabled={isProcessing} variant="primary" size="lg" className="w-full">
            {isProcessing ? '処理中...' : '注文を確定する'}
          </Button>
        </form>

        {/* Order Summary */}
        <div>
          <div className="bg-white rounded-lg shadow p-6 sticky top-4">
            <h2 className="text-xl font-bold mb-6">注文内容</h2>

            <div className="space-y-4 mb-6 border-b pb-6">
              {(cart.items ?? []).map((item: CartItem) => (
                <div key={item.product_id} className="flex justify-between text-sm">
                  <span>{item.product_name}</span>
                  <span>¥{item.subtotal?.toLocaleString()}</span>
                </div>
              ))}
            </div>

            <div className="space-y-4 mb-6 border-b pb-6">
              <div className="flex justify-between">
                <span>小計</span>
                <span>￥{cart.total?.toLocaleString()}</span>
              </div>
              <div className="flex justify-between">
                <span>配送料</span>
                <span>¥1,000</span>
              </div>
              <div className="flex justify-between">
                <span>税金</span>
                <span>￥{Math.round((cart.total ?? 0) * 0.1).toLocaleString()}</span>
              </div>
            </div>

            <div className="flex justify-between text-xl font-bold">
              <span>合計</span>
              <span className="text-primary">¥{totalAmount.toLocaleString()}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
