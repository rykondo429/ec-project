'use client';

import React, { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';
import { User, Package, CreditCard, Heart, Settings, LogOut } from 'lucide-react';

// 動的レンダリングを強制
export const dynamic = 'force-dynamic';

export default function MyPage() {
  const router = useRouter();
  const { user, loading, signOut } = useAuth();

  useEffect(() => {
    if (!loading && !user) {
      router.push('/login?redirect=/my-page');
    }
  }, [user, loading, router]);

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto"></div>
          <p className="mt-4 text-gray-600">読み込み中...</p>
        </div>
      </div>
    );
  }

  if (!user) {
    return null;
  }

  const menuItems = [
    {
      icon: Package,
      title: '注文履歴',
      description: 'これまでのご注文を確認できます',
      href: '/my-orders',
    },
    {
      icon: CreditCard,
      title: 'ポイント',
      description: '保有ポイントと履歴を確認',
      href: '/my-points',
    },
    {
      icon: Heart,
      title: 'お気に入り',
      description: 'お気に入り商品の管理',
      href: '/favorites',
    },
    {
      icon: Settings,
      title: '設定',
      description: 'アカウント設定の変更',
      href: '/settings',
    },
  ];

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      {/* ユーザー情報カード */}
      <div className="bg-white rounded-lg shadow-md p-6 mb-8">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <div className="w-16 h-16 bg-primary rounded-full flex items-center justify-center">
              <User size={32} className="text-white" />
            </div>
            <div>
              <h1 className="text-2xl font-bold text-gray-900">
                {user.email}
              </h1>
              <p className="text-sm text-gray-500">会員ID: {user.sub}</p>
            </div>
          </div>
          <button
            onClick={signOut}
            className="flex items-center gap-2 px-4 py-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
          >
            <LogOut size={20} />
            <span>ログアウト</span>
          </button>
        </div>
      </div>

      {/* メニューグリッド */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {menuItems.map((item) => {
          const Icon = item.icon;
          return (
            <Link
              key={item.href}
              href={item.href}
              className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow"
            >
              <div className="flex flex-col items-center text-center">
                <div className="w-12 h-12 bg-primary bg-opacity-10 rounded-full flex items-center justify-center mb-4">
                  <Icon size={24} className="text-primary" />
                </div>
                <h3 className="text-lg font-semibold text-gray-900 mb-2">
                  {item.title}
                </h3>
                <p className="text-sm text-gray-600">
                  {item.description}
                </p>
              </div>
            </Link>
          );
        })}
      </div>

      {/* 最近の活動 */}
      <div className="mt-8 bg-white rounded-lg shadow-md p-6">
        <h2 className="text-xl font-bold text-gray-900 mb-4">最近の活動</h2>
        <div className="space-y-4">
          <div className="flex items-center justify-between py-3 border-b">
            <div>
              <p className="font-medium text-gray-900">新規登録完了</p>
              <p className="text-sm text-gray-500">アカウントが作成されました</p>
            </div>
            <span className="text-sm text-gray-400">今日</span>
          </div>
        </div>
      </div>
    </div>
  );
}
