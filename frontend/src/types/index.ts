/**
 * API型定義
 * api-schema/types から自動生成された型をre-exportします
 */

// OpenAPI自動生成型のインポート
import type { components as ProductServiceAPI } from '../../../api-schema/types/product-service';
import type { components as CartServiceAPI } from '../../../api-schema/types/cart-service';
import type { components as OrderServiceAPI } from '../../../api-schema/types/order-service';
import type { components as PointServiceAPI } from '../../../api-schema/types/point-service';
import type { components as PromotionServiceAPI } from '../../../api-schema/types/promotion-service';

// 使いやすいようにエイリアスをエクスポート

// Product Service
export type Product = ProductServiceAPI['schemas']['Product'];
export type SKU = ProductServiceAPI['schemas']['SKU'];
export type ProductWithSKUs = ProductServiceAPI['schemas']['ProductWithSKUs'];
export type SearchProductsResponse = ProductServiceAPI['schemas']['SearchResult'];
export type SearchProductsRequest = {
  keyword?: string;
  category?: string;
  price_min?: number;
  price_max?: number;
  sort_by?: 'price' | 'rating' | 'created_at';
  sort_order?: 'asc' | 'desc';
  page?: number;
  page_size?: number;
};

// Cart Service
export type Cart = CartServiceAPI['schemas']['Cart'];
export type CartItem = CartServiceAPI['schemas']['CartItem'];
export type AddItemRequest = CartServiceAPI['schemas']['AddItemRequest'];
export type UpdateItemRequest = CartServiceAPI['schemas']['UpdateItemRequest'];

// Order Service
export type Order = OrderServiceAPI['schemas']['Order'];
export type OrderItem = OrderServiceAPI['schemas']['OrderItem'];
export type CreateOrderRequest = OrderServiceAPI['schemas']['CreateOrderRequest'];

// Point Service
export type PointBalance = PointServiceAPI['schemas']['PointBalance'];
export type PointTransaction = PointServiceAPI['schemas']['PointTransaction'];
export type EarnPointsRequest = PointServiceAPI['schemas']['EarnPointsRequest'];
export type UsePointsRequest = PointServiceAPI['schemas']['UsePointsRequest'];

// Promotion Service
export type Coupon = PromotionServiceAPI['schemas']['Coupon'];
export type Sale = PromotionServiceAPI['schemas']['Sale'];
export type ValidateCouponRequest = PromotionServiceAPI['schemas']['ValidateCouponRequest'];
export type ValidateCouponResponse = PromotionServiceAPI['schemas']['CouponValidationResult'];

// 元のAPIコンポーネント型もエクスポート（必要に応じて）
export type { ProductServiceAPI, CartServiceAPI, OrderServiceAPI, PointServiceAPI, PromotionServiceAPI };
