// API Schema Types Usage Examples
// このファイルは使用例です - 実際のコードではありません

/**
 * 例1: Product Service の型を使用
 */
import type { components as ProductComponents } from '../../api-schema/types/product-service';

type Product = ProductComponents['schemas']['Product'];
type SearchResponse = ProductComponents['schemas']['SearchResponse'];

const product: Product = {
  id: 'prod-123',
  name: 'スニーカー',
  description: 'おしゃれなスニーカー',
  price: 5980,
  category: 'shoes',
  imageUrl: '/images/sneaker.jpg',
  stock: 50,
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString(),
};

/**
 * 例2: Cart Service の型を使用
 */
import type { components as CartComponents } from '../../api-schema/types/cart-service';

type Cart = CartComponents['schemas']['Cart'];
type CartItem = CartComponents['schemas']['CartItem'];

const cart: Cart = {
  userId: 'user-123',
  items: [
    {
      productId: 'prod-123',
      productName: 'スニーカー',
      quantity: 2,
      price: 5980,
      imageUrl: '/images/sneaker.jpg',
    },
  ],
  total: 11960,
  updatedAt: new Date().toISOString(),
};

/**
 * 例3: Order Service の型を使用
 */
import type { components as OrderComponents } from '../../api-schema/types/order-service';

type Order = OrderComponents['schemas']['Order'];
type CreateOrderRequest = OrderComponents['schemas']['CreateOrderRequest'];

const orderRequest: CreateOrderRequest = {
  items: [
    {
      productId: 'prod-123',
      quantity: 2,
      price: 5980,
    },
  ],
  couponCode: 'SAVE10',
  usePoints: 500,
  shippingAddress: {
    postalCode: '100-0001',
    prefecture: '東京都',
    city: '千代田区',
    address1: '丸の内1-1-1',
    address2: 'ビル101',
  },
  paymentMethod: 'credit_card',
};

/**
 * 例4: パス型を使用
 */
import type { paths as ProductPaths } from '../../api-schema/types/product-service';

// GET /products/{id} のレスポンス型
type GetProductResponse =
  ProductPaths['/products/{id}']['get']['responses']['200']['content']['application/json'];

// GET /products/search のクエリパラメータ型
type SearchQueryParams = ProductPaths['/products/search']['get']['parameters']['query'];

const searchParams: SearchQueryParams = {
  keyword: 'スニーカー',
  category: 'shoes',
  price_min: 1000,
  price_max: 10000,
  page: 1,
  page_size: 20,
  sort_by: 'price',
  sort_order: 'asc',
};

/**
 * 例5: API クライアント関数で型を使用
 */
import type { components } from '../../api-schema/types/product-service';

async function getProduct(id: string): Promise<components['schemas']['Product']> {
  const response = await fetch(`http://localhost:8001/api/v1/products/${id}`);
  if (!response.ok) {
    throw new Error('Failed to fetch product');
  }
  return response.json();
}

async function searchProducts(
  params: ProductPaths['/products/search']['get']['parameters']['query']
): Promise<components['schemas']['SearchResponse']> {
  const query = new URLSearchParams(params as any).toString();
  const response = await fetch(`http://localhost:8001/api/v1/products/search?${query}`);
  if (!response.ok) {
    throw new Error('Failed to search products');
  }
  return response.json();
}

/**
 * 例6: すべてのサービスの型を統合インポート
 */
import type {
  ProductService,
  CartService,
  OrderService,
  PointService,
  PromotionService,
} from '../../api-schema/types';

type AllProducts = ProductService.components['schemas']['Product'];
type AllCarts = CartService.components['schemas']['Cart'];
type AllOrders = OrderService.components['schemas']['Order'];
type PointBalance = PointService.components['schemas']['PointBalance'];
type Coupon = PromotionService.components['schemas']['Coupon'];

/**
 * 例7: React コンポーネントで使用
 */
import { useState, useEffect } from 'react';
import type { components } from '../../api-schema/types/product-service';

type Product = components['schemas']['Product'];

function ProductDetail({ productId }: { productId: string }) {
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`http://localhost:8001/api/v1/products/${productId}`)
      .then((res) => res.json())
      .then((data: Product) => {
        setProduct(data);
        setLoading(false);
      });
  }, [productId]);

  if (loading) return <div>Loading...</div>;
  if (!product) return <div>Product not found</div>;

  return (
    <div>
      <h1>{product.name}</h1>
      <p>{product.description}</p>
      <p>¥{product.price.toLocaleString()}</p>
    </div>
  );
}

/**
 * 例8: Zustand ストアで使用
 */
import { create } from 'zustand';
import type { components } from '../../api-schema/types/cart-service';

type Cart = components['schemas']['Cart'];

interface CartState {
  cart: Cart | null;
  loading: boolean;
  fetchCart: () => Promise<void>;
  addItem: (productId: string, quantity: number) => Promise<void>;
}

const useCartStore = create<CartState>((set) => ({
  cart: null,
  loading: false,
  fetchCart: async () => {
    set({ loading: true });
    const response = await fetch('http://localhost:8002/api/v1/carts', {
      headers: { Authorization: `Bearer ${localStorage.getItem('token')}` },
    });
    const cart: Cart = await response.json();
    set({ cart, loading: false });
  },
  addItem: async (productId: string, quantity: number) => {
    await fetch('http://localhost:8002/api/v1/carts/items', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${localStorage.getItem('token')}`,
      },
      body: JSON.stringify({ productId, quantity }),
    });
  },
}));
