'use client';

export const dynamic = 'force-dynamic';

import React, { useEffect } from 'react';
import { usePointStore } from '@/store/points';
import { Gift } from 'lucide-react';

export default function MyPointsPage() {
  const { balance, transactions, fetchBalance, fetchTransactions, isLoading } = usePointStore();

  useEffect(() => {
    fetchBalance();
    fetchTransactions();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  if (isLoading) {
    return <div className="max-w-7xl mx-auto px-4 py-16 text-center">読み込み中...</div>;
  }

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">ポイント残高</h1>

      {/* Points Summary */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div className="bg-gradient-to-br from-blue-500 to-blue-600 text-white rounded-lg shadow-lg p-8">
          <p className="text-sm opacity-90 mb-2">利用可能ポイント</p>
          <p className="text-4xl font-bold">{balance?.available_points?.toLocaleString() ?? 0}pt</p>
        </div>

        <div className="bg-gradient-to-br from-green-500 to-green-600 text-white rounded-lg shadow-lg p-8">
          <p className="text-sm opacity-90 mb-2">総ポイント</p>
          <p className="text-4xl font-bold">{balance?.total_points?.toLocaleString() ?? 0}pt</p>
        </div>

        <div className="bg-gradient-to-br from-orange-500 to-orange-600 text-white rounded-lg shadow-lg p-8">
          <p className="text-sm opacity-90 mb-2">有効期限切れ</p>
          <p className="text-4xl font-bold">{balance?.expired_points?.toLocaleString() ?? 0}pt</p>
        </div>
      </div>

      {/* Transactions */}
      <div className="bg-white rounded-lg shadow">
        <div className="p-6 border-b">
          <h2 className="text-2xl font-bold">ポイント履歴</h2>
        </div>

        {transactions.length === 0 ? (
          <div className="p-12 text-center text-gray-600">
            <Gift size={48} className="mx-auto mb-4 opacity-50" />
            <p>ポイントの履歴がありません</p>
          </div>
        ) : (
          <table className="w-full">
            <thead className="bg-gray-100 border-b">
              <tr>
                <th className="text-left px-6 py-3 font-semibold">日付</th>
                <th className="text-left px-6 py-3 font-semibold">種類</th>
                <th className="text-left px-6 py-3 font-semibold">説明</th>
                <th className="text-right px-6 py-3 font-semibold">ポイント</th>
                <th className="text-right px-6 py-3 font-semibold">残高</th>
              </tr>
            </thead>
            <tbody>
              {transactions.map((tx) => (
                <tr key={tx.id} className="border-b hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm">
                    {new Date(tx.created_at ?? '').toLocaleDateString('ja-JP')}
                  </td>
                  <td className="px-6 py-4">
                    <span
                      className={`px-3 py-1 rounded text-sm font-semibold ${getTypeStyles(tx.type ?? '')}`}
                    >
                      {getTypeLabel(tx.type ?? '')}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm">{tx.reason}</td>
                  <td className={`px-6 py-4 text-right font-semibold ${(tx.points ?? 0) > 0 ? 'text-green-600' : 'text-red-600'}`}>
                    {(tx.points ?? 0) > 0 ? '+' : ''}{tx.points}pt
                  </td>
                  <td className="px-6 py-4 text-right">{tx.balance?.toLocaleString() ?? 0}pt</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {/* Points Usage Tips */}
      <div className="mt-8 bg-blue-50 border border-blue-200 rounded-lg p-6">
        <h3 className="font-bold mb-4 flex items-center gap-2">
          <Gift size={20} className="text-blue-600" />
          ポイント利用について
        </h3>
        <ul className="text-sm text-gray-700 space-y-2">
          <li>• 1ポイント = 1円で注文時に使用できます</li>
          <li>• 商品購入時に自動的にポイントが付与されます</li>
          <li>• ポイント有効期限: 取得から1年</li>
          <li>• 返品時はポイントも返却されます</li>
        </ul>
      </div>
    </div>
  );
}

function getTypeLabel(type: string): string {
  const labels: Record<string, string> = {
    earn: 'ポイント獲得',
    use: 'ポイント利用',
    expire: '有効期限切れ',
    refund: '返金',
  };
  return labels[type] || type;
}

function getTypeStyles(type: string): string {
  const styles: Record<string, string> = {
    earn: 'bg-green-100 text-green-800',
    use: 'bg-blue-100 text-blue-800',
    expire: 'bg-gray-100 text-gray-800',
    refund: 'bg-orange-100 text-orange-800',
  };
  return styles[type] || 'bg-gray-100 text-gray-800';
}
