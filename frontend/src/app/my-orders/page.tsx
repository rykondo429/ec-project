'use client';

export const dynamic = 'force-dynamic';

import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import { Button } from '@/components/Button';
import { OrderApi } from '@/lib/api-client';
import { Order } from '@/types';

export default function OrdersPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const orderApi = new OrderApi();

  useEffect(() => {
    const fetchOrders = async () => {
      try {
        const result = await orderApi.getOrders();
        setOrders(result.orders);
      } catch (error) {
        console.error('Failed to fetch orders:', error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchOrders();
  }, []);

  if (isLoading) {
    return <div className="max-w-7xl mx-auto px-4 py-16 text-center">読み込み中...</div>;
  }

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">注文履歴</h1>

      {orders.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-600 mb-4">注文がありません</p>
          <Button asLink href="/products" variant="primary">
            ショッピングを開始
          </Button>
        </div>
      ) : (
        <div className="space-y-4">
          {orders.map((order) => (
            <div key={order.id} className="bg-white rounded-lg shadow p-6">
              <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
                <div>
                  <p className="text-sm text-gray-600">注文ID</p>
                  <p className="font-semibold">{order.id}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-600">注文日</p>
                  <p className="font-semibold">{new Date(order.created_at ?? '').toLocaleDateString('ja-JP')}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-600">ステータス</p>
                  <p className={`font-semibold ${getStatusColor(order.status ?? '')}`}>{getStatusLabel(order.status ?? '')}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-600">合計金額</p>
                  <p className="font-semibold text-primary text-lg">￥{order.total?.toLocaleString()}</p>
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm text-gray-600 mb-4">
                <p>商品数: {order.items?.length ?? 0}</p>
                <p>ポイント: +{order.points_earned ?? 0}</p>
              </div>

              <Link href={`/orders/${order.id}`}>
                <Button variant="outline" size="sm">
                  詳細を見る
                </Button>
              </Link>
            </div>
          ))}
        </div>
      )}
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

function getStatusColor(status: string): string {
  const colors: Record<string, string> = {
    pending: 'text-yellow-600',
    confirmed: 'text-blue-600',
    shipped: 'text-blue-600',
    delivered: 'text-green-600',
    cancelled: 'text-red-600',
  };
  return colors[status] || 'text-gray-600';
}
