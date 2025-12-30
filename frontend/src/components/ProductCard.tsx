import React from 'react';
import Link from 'next/link';
import { Star, ShoppingCart, Heart } from 'lucide-react';
import { Product } from '@/types';
import { useCartStore } from '@/store/cart';
import { useFavoriteStore } from '@/store/favorite';

interface ProductCardProps {
  product: Product;
}

export const ProductCard: React.FC<ProductCardProps> = ({ product }) => {
  const addItem = useCartStore((state) => state.addItem);
  const { toggleFavorite, isFavorite } = useFavoriteStore();
  const [isAdding, setIsAdding] = React.useState(false);
  const isProductFavorite = isFavorite(product.id);

  const handleAddToCart = async () => {
    setIsAdding(true);
    try {
      await addItem(product.id, product.name, product.price, 1);
    } catch (error) {
      console.error('Failed to add item to cart:', error);
    } finally {
      setIsAdding(false);
    }
  };

  const displayPrice = product.is_sale && product.sale_price ? product.sale_price : product.price;
  const originalPrice = product.price;

  return (
    <div className="bg-white rounded-lg shadow hover:shadow-lg transition-shadow overflow-hidden">
      <div className="aspect-video bg-gray-200 relative overflow-hidden">
        <div className="w-full h-full bg-gradient-to-br from-gray-300 to-gray-400 flex items-center justify-center">
          <span className="text-gray-500">No Image</span>
        </div>
        {product.is_sale && product.sale_price && (
          <div className="absolute top-2 right-2 bg-red-500 text-white px-3 py-1 rounded text-sm font-bold">
            Sale
          </div>
        )}
        <button
          onClick={() => toggleFavorite(product)}
          className="absolute top-2 left-2 p-2 bg-white/90 hover:bg-white rounded-full shadow-md transition-all"
          title={isProductFavorite ? 'お気に入りから削除' : 'お気に入りに追加'}
        >
          <Heart
            size={20}
            className={isProductFavorite ? 'fill-red-500 text-red-500' : 'text-gray-600'}
          />
        </button>
      </div>

      <div className="p-4">
        <h3 className="font-semibold text-lg mb-2 line-clamp-2">{product.name}</h3>

        <div className="flex items-center gap-2 mb-3">
          <div className="flex items-center">
            {[...Array(5)].map((_, i) => (
              <Star
                key={i}
                size={14}
                className={i < Math.floor(product.rating ?? 0) ? 'fill-yellow-400 text-yellow-400' : 'text-gray-300'}
              />
            ))}
          </div>
          <span className="text-sm text-gray-600">({product.review_count})</span>
        </div>

        <div className="mb-4">
          <div className="flex items-baseline gap-2">
            <span className="text-2xl font-bold text-primary">¥{displayPrice.toLocaleString()}</span>
            {product.is_sale && product.sale_price && (
              <span className="text-sm line-through text-gray-500">¥{originalPrice.toLocaleString()}</span>
            )}
          </div>
        </div>

        <div className="flex gap-2">
          <Link
            href={`/products/${product.id}`}
            className="flex-1 text-center py-2 border border-primary text-primary rounded hover:bg-blue-50"
          >
            詳細
          </Link>
          <button
            onClick={handleAddToCart}
            disabled={isAdding || product.stock === 0}
            className="flex-1 py-2 bg-primary text-white rounded hover:bg-blue-600 disabled:bg-gray-300 flex items-center justify-center gap-2"
          >
            <ShoppingCart size={16} />
            {product.stock === 0 ? '在庫なし' : 'カートに追加'}
          </button>
        </div>
      </div>
    </div>
  );
};
