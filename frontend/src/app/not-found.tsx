export const dynamic = 'force-dynamic';

export default function NotFound() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen">
      <h1 className="text-4xl font-bold mb-4">404 - ページが見つかりません</h1>
      <p className="text-gray-600 mb-8">お探しのページは存在しません。</p>
      <a href="/" className="text-blue-600 hover:underline">
        トップページに戻る
      </a>
    </div>
  );
}
