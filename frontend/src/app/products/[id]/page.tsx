'use client';

export const dynamic = 'force-dynamic';

import React, { useEffect } from 'react';
import { useParams } from 'next/navigation';
import { Star, ShoppingCart, Check } from 'lucide-react';
import { Button } from '@/components/Button';
import { useProductStore } from '@/store/product';
import { useCartStore } from '@/store/cart';
import type { SKU } from '@/types';

export default function ProductDetailPage() {
  const params = useParams();
  const { productWithSKUs, getProductWithSKUs, isLoading } = useProductStore();
  const { addItem } = useCartStore();
  const [quantity, setQuantity] = React.useState(1);
  const [isAdding, setIsAdding] = React.useState(false);
  const [selectedSKU, setSelectedSKU] = React.useState<SKU | null>(null);

  useEffect(() => {
    if (params.id) {
      getProductWithSKUs(params.id as string);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [params.id]);

  // SKUが読み込まれたら最初のSKUを選択
  useEffect(() => {
    if (productWithSKUs?.skus && productWithSKUs.skus.length > 0 && !selectedSKU) {
      setSelectedSKU(productWithSKUs.skus[0]);
    }
  }, [productWithSKUs, selectedSKU]);

  const handleAddToCart = async () => {
    if (!productWithSKUs) return;
    
    setIsAdding(true);
    try {
      if (selectedSKU) {
        // SKU選択時
        await addItem(
          productWithSKUs.id,
          `${productWithSKUs.name} - ${selectedSKU.name}`,
          selectedSKU.price,
          quantity,
          selectedSKU.id,
          selectedSKU.sku_code,
          selectedSKU.attributes
        );
      } else {
        // SKUなし（従来の商品のみ）
        await addItem(
          productWithSKUs.id,
          productWithSKUs.name,
          productWithSKUs.price,
          quantity
        );
      }
      alert('カートに追加しました');
    } catch {
      alert('カートへの追加に失敗しました');
    } finally {
      setIsAdding(false);
    }
  };

  // 選択されたSKUまたは商品から価格と在庫を取得
  const displayPrice = selectedSKU ? selectedSKU.price : 
    (productWithSKUs?.is_sale && productWithSKUs?.sale_price ? productWithSKUs.sale_price : productWithSKUs?.price ?? 0);
  const displayStock = selectedSKU ? selectedSKU.stock : 
    (productWithSKUs?.skus?.reduce((sum, sku) => sum + (sku.stock || 0), 0) ?? 0);
  const hasSKUs = productWithSKUs?.skus && productWithSKUs.skus.length > 0;

  // 属性でグルーピング（例: サイズごと、色ごと）
  const getAttributeOptions = (attrKey: string): string[] => {
    if (!productWithSKUs?.skus) return [];
    const values = new Set<string>();
    productWithSKUs.skus.forEach((sku: SKU) => {
      if (sku.attributes && sku.attributes[attrKey]) {
        values.add(sku.attributes[attrKey]);
      }
    });
    return Array.from(values);
  };

  // 選択された属性に基づいてSKUをフィルタリング
  const [selectedAttributes, setSelectedAttributes] = React.useState<Record<string, string>>({});

  // 利用可能な属性キーを取得
  const attributeKeys = React.useMemo(() => {
    if (!productWithSKUs?.skus || productWithSKUs.skus.length === 0) return [];
    const keys = new Set<string>();
    productWithSKUs.skus.forEach((sku: SKU) => {
      if (sku.attributes) {
        Object.keys(sku.attributes).forEach(key => keys.add(key));
      }
    });
    return Array.from(keys);
  }, [productWithSKUs]);

  // 属性選択時にSKUを更新
  useEffect(() => {
    if (!productWithSKUs?.skus) return;
    
    // すべての属性が選択されている場合、対応するSKUを見つける
    const matchingSKU = productWithSKUs.skus.find((sku: SKU) => {
      if (!sku.attributes) return false;
      return attributeKeys.every(key => {
        return !selectedAttributes[key] || sku.attributes![key] === selectedAttributes[key];
      });
    });

    if (matchingSKU) {
      setSelectedSKU(matchingSKU);
    }
  }, [selectedAttributes, productWithSKUs, attributeKeys]);

  const handleAttributeSelect = (attrKey: string, value: string) => {
    setSelectedAttributes(prev => ({
      ...prev,
      [attrKey]: value
    }));
  };

  if (isLoading) {
    return <div className="max-w-7xl mx-auto px-4 py-16 text-center">読み込み中...</div>;
  }

  if (!productWithSKUs) {
    return <div className="max-w-7xl mx-auto px-4 py-16 text-center">商品が見つかりません</div>;
  }

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        {/* Product Image */}
        <div>
          <div className="aspect-square bg-gray-200 rounded-lg flex items-center justify-center relative">
            <span className="text-gray-500">No Image</span>
            {productWithSKUs.is_sale && (
              <div className="absolute top-4 right-4 bg-red-500 text-white px-4 py-2 rounded font-bold">
                Sale
              </div>
            )}
          </div>
        </div>

        {/* Product Info */}
        <div>
          <div className="mb-4">
            <span className="text-sm text-gray-500 bg-gray-100 px-3 py-1 rounded">
              {productWithSKUs.category}
            </span>
          </div>

          <h1 className="text-3xl font-bold mb-4">{productWithSKUs.name}</h1>

          {/* Rating */}
          <div className="flex items-center gap-2 mb-6">
            <div className="flex">
              {[...Array(5)].map((_, i) => (
                <Star
                  key={i}
                  size={18}
                  className={i < Math.floor(productWithSKUs.rating ?? 0) ? 'fill-yellow-400 text-yellow-400' : 'text-gray-300'}
                />
              ))}
            </div>
            <span className="text-gray-600">({productWithSKUs.review_count} レビュー)</span>
          </div>

          {/* SKU Selection */}
          {hasSKUs && (
            <div className="mb-6 space-y-4">
              {attributeKeys.map(attrKey => {
                const options = getAttributeOptions(attrKey);
                const displayName = attrKey === 'size' ? 'サイズ' : 
                                   attrKey === 'color' ? 'カラー' : attrKey;
                
                return (
                  <div key={attrKey}>
                    <label className="block text-sm font-semibold mb-2 capitalize">
                      {displayName}
                    </label>
                    <div className="flex gap-2 flex-wrap">
                      {options.map(option => {
                        const isSelected = selectedAttributes[attrKey] === option;
                        const matchingSKU = productWithSKUs.skus?.find((sku: SKU) => 
                          sku.attributes && sku.attributes[attrKey] === option
                        );
                        const isAvailable = matchingSKU && matchingSKU.stock > 0;
                        
                        return (
                          <button
                            key={option}
                            onClick={() => handleAttributeSelect(attrKey, option)}
                            disabled={!isAvailable}
                            className={`
                              px-4 py-2 border-2 rounded-lg font-medium transition-all
                              ${isSelected 
                                ? 'border-primary bg-primary text-white' 
                                : 'border-gray-300 hover:border-primary'}
                              ${!isAvailable && 'opacity-40 cursor-not-allowed line-through'}
                              relative
                            `}
                          >
                            {option}
                            {isSelected && (
                              <Check size={16} className="absolute top-1 right-1" />
                            )}
                          </button>
                        );
                      })}
                    </div>
                  </div>
                );
              })}
              
              {selectedSKU && (
                <div className="text-sm text-gray-600 bg-gray-50 p-3 rounded">
                  <p className="font-semibold">選択中: {selectedSKU.name}</p>
                  <p>SKUコード: {selectedSKU.sku_code}</p>
                </div>
              )}
            </div>
          )}

          {/* Price */}
          <div className="mb-6">
            <div className="flex items-baseline gap-3 mb-2">
              <span className="text-4xl font-bold text-primary">¥{displayPrice.toLocaleString()}</span>
              {productWithSKUs.is_sale && productWithSKUs.sale_price && !selectedSKU && (
                <span className="text-xl line-through text-gray-500">￥{productWithSKUs.price?.toLocaleString()}</span>
              )}
            </div>
            {productWithSKUs.is_sale && productWithSKUs.sale_price && !selectedSKU && (
              <p className="text-red-600 font-semibold">
                {Math.round((((productWithSKUs.price ?? 0) - (productWithSKUs.sale_price ?? 0)) / (productWithSKUs.price ?? 1)) * 100)}% OFF
              </p>
            )}
          </div>

          {/* Stock */}
          <div className="mb-6">
            <p className={`font-semibold ${displayStock > 0 ? 'text-green-600' : 'text-red-600'}`}>
              {displayStock > 0 ? `在庫あり (${displayStock}個)` : '在庫なし'}
            </p>
          </div>

          {/* Description */}
          <div className="mb-8">
            <h3 className="font-semibold mb-2">商品説明</h3>
            <p className="text-gray-600 whitespace-pre-wrap">{productWithSKUs.description}</p>
          </div>

          {/* Add to Cart */}
          {displayStock > 0 && (
            <div className="flex gap-4 mb-8">
              <div className="flex items-center border rounded-lg">
                <button
                  onClick={() => setQuantity(Math.max(1, quantity - 1))}
                  className="px-4 py-2 hover:bg-gray-100"
                >
                  −
                </button>
                <input
                  type="number"
                  value={quantity}
                  onChange={(e) => setQuantity(Math.max(1, parseInt(e.target.value) || 1))}
                  min="1"
                  max={displayStock}
                  className="w-12 text-center border-none"
                />
                <button
                  onClick={() => setQuantity(Math.min(displayStock, quantity + 1))}
                  className="px-4 py-2 hover:bg-gray-100"
                >
                  ＋
                </button>
              </div>
              <Button
                onClick={handleAddToCart}
                disabled={isAdding || (hasSKUs && !selectedSKU)}
                variant="primary"
                size="lg"
                className="flex-1"
              >
                <ShoppingCart size={20} />
                カートに追加
              </Button>
            </div>
          )}

          {/* Related Info */}
          <div className="border-t pt-6">
            <p className="text-sm text-gray-600">
              商品ID: <span className="font-mono">{productWithSKUs.id}</span>
            </p>
            {hasSKUs && (
              <p className="text-sm text-gray-600 mt-1">
                {productWithSKUs.skus?.length}種類のバリエーション
              </p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
