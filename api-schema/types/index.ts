/**
 * EC-Sample API Types
 * 
 * OpenAPIから自動生成されたTypeScript型定義
 * このファイルを編集しないでください - 自動生成されます
 */

// Product Service Types (Port 8001)
export type * as ProductService from './product-service';

// Cart Service Types (Port 8002)
export type * as CartService from './cart-service';

// Order Service Types (Port 8003)
export type * as OrderService from './order-service';

// Point Service Types (Port 8004)
export type * as PointService from './point-service';

// Promotion Service Types (Port 8005)
export type * as PromotionService from './promotion-service';

/**
 * 使用例:
 * 
 * import type { ProductService } from '@ec-sample/api-schema/types';
 * 
 * type Product = ProductService.components['schemas']['Product'];
 * type SearchResponse = ProductService.components['schemas']['SearchResponse'];
 * 
 * // パスパラメータとレスポンス型
 * type GetProductResponse = ProductService.paths['/products/{id}']['get']['responses']['200']['content']['application/json'];
 */
