'use client';

export const dynamic = 'force-dynamic';

import React, { useEffect } from 'react';
import { useParams } from 'next/navigation';
import { Button } from '@/components/Button';
import { OrderApi } from '@/lib/api-client';
import { Order, OrderItem } from '@/types';

export default function OrderDetailPage() {
  const params = useParams();
  const [order, setOrder] = React.useState<Order | null>(null);
  const [isLoading, setIsLoading] = React.useState(true);
  const orderApi = new OrderApi();

  useEffect(() => {
    const fetchOrder = async () => {
      try {
        const data = await orderApi.getOrder(params.id as string);
        setOrder(data);
      } catch (error) {
        console.error('Failed to fetch order:', error);
      } finally {
        setIsLoading(false);
      }
    };

    if (params.id) {
      fetchOrder();
    }
  }, [params.id]);

  if (isLoading) {
    return <div className="max-w-4xl mx-auto px-4 py-16 text-center">読み込み中...</div>;
  }

  if (!order) {
    return <div className="max-w-4xl mx-auto px-4 py-16 text-center">注文が見つかりません</div>;
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="mb-8">
        <Button asLink href="/my-orders" variant="outline">
          ← 注文履歴に戻る
        </Button>
      </div>

      <h1 className="text-3xl font-bold mb-8">注文詳細</h1>

      {/* Status */}
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-6 mb-8">
        <p className="text-sm text-gray-600 mb-2">現在のステータス</p>
        <p className="text-2xl font-bold text-blue-600">{getStatusLabel(order.status ?? '')}</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2">
          {/* Order Items */}
          <div className="bg-white rounded-lg shadow p-6 mb-8">
            <h2 className="text-xl font-bold mb-6">注文商品</h2>
            <table className="w-full">
              <thead className="bg-gray-100 border-b">
                <tr>
                  <th className="text-left px-4 py-3 font-semibold">商品</th>
                  <th className="text-center px-4 py-3 font-semibold">数量</th>
                  <th className="text-right px-4 py-3 font-semibold">価格</th>
                </tr>
              </thead>
              <tbody>
                {(order.items ?? []).map((item: OrderItem) => (
                  <tr key={item.product_id} className="border-b hover:bg-gray-50">
                    <td className="px-4 py-4">{item.product_name}</td>
                    <td className="px-4 py-4 text-center">{item.quantity}</td>
                    <td className="px-4 py-4 text-right">￥{item.subtotal?.toLocaleString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Shipping Info */}
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-xl font-bold mb-6">配送情報</h2>
            <p className="text-gray-600 whitespace-pre-wrap">{order.shipping_address}</p>
          </div>
        </div>

        {/* Order Summary */}
        <div>
          <div className="bg-white rounded-lg shadow p-6 sticky top-4">
            <h2 className="text-xl font-bold mb-6">注文概要</h2>

            <div className="space-y-4 mb-6 border-b pb-6">
              <div className="flex justify-between">
                <span>小計</span>
                <span>¥{order.subtotal?.toLocaleString()}</span>
              </div>
              <div className="flex justify-between">
                <span>配送料</span>
                <span>¥{order.shipping_cost?.toLocaleString()}</span>
              </div>
              <div className="flex justify-between">
                <span>税金</span>
                <span>¥{order.tax?.toLocaleString()}</span>
              </div>
              {(order.discount_amount ?? 0) > 0 && (
                <div className="flex justify-between text-green-600">
                  <span>割引</span>
                  <span>-￥{order.discount_amount?.toLocaleString()}</span>
                </div>
              )}
            </div>

            <div className="flex justify-between text-xl font-bold mb-6">
              <span>合計</span>
              <span className="text-primary">￥{order.total?.toLocaleString()}</span>
            </div>

            {(order.points_earned ?? 0) > 0 && (
              <div className="bg-yellow-50 border border-yellow-200 rounded p-4 mb-6">
                <p className="text-sm text-gray-600">獲得ポイント</p>
                <p className="text-lg font-bold text-yellow-600">+{order.points_earned}pt</p>
              </div>
            )}

            <div className="text-sm text-gray-600 space-y-2">
              <p>注文ID: {order.id}</p>
              <p>注文日: {new Date(order.created_at ?? '').toLocaleDateString('ja-JP')}</p>
              <p>支払い方法: {order.payment_method}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function getStatusLabel(status: string): string {
  const labels: Record<string, string> = {
    pending: '保留中',
    confirmed: '確認済み',
    shipped: '発送済み',
    delivered: '配送完了',
    cancelled: 'キャンセル',
  };
  return labels[status] || status;
}
