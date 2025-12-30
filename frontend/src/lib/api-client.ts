import axios, { AxiosInstance } from 'axios';
import * as Types from '@/types';
import { getIdToken } from './cognito';

// API クライアント設定
const createApiClient = (baseURL: string): AxiosInstance => {
  const client = axios.create({
    baseURL,
    headers: {
      'Content-Type': 'application/json',
    },
  });

  // Cognitoトークンを自動的に付与
  client.interceptors.request.use(async (config) => {
    try {
      const token = await getIdToken();
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    } catch (error) {
      // トークン取得失敗時はanonymousとして処理
      console.warn('Failed to get auth token:', error);
    }
    return config;
  });

  return client;
};

// 製品API
export class ProductApi {
  private client: AxiosInstance;

  constructor(baseURL: string = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8001') {
    this.client = createApiClient(baseURL);
  }

  async searchProducts(params: Types.SearchProductsRequest): Promise<Types.SearchProductsResponse> {
    const response = await this.client.get<Types.SearchProductsResponse>('/api/v1/products/search', {
      params,
    });
    return response.data;
  }

  async getProduct(id: string): Promise<Types.Product> {
    const response = await this.client.get<Types.Product>(`/api/v1/products/${id}`);
    return response.data;
  }

  async getProductWithSKUs(id: string): Promise<Types.ProductWithSKUs> {
    const response = await this.client.get<Types.ProductWithSKUs>(`/api/v1/products/${id}/with-skus`);
    return response.data;
  }

  async getSKUsByProduct(productId: string): Promise<Types.SKU[]> {
    const response = await this.client.get<Types.SKU[]>(`/api/v1/products/${productId}/skus`);
    return response.data;
  }

  async getCategories(): Promise<string[]> {
    const response = await this.client.get<{ categories: string[] }>('/api/v1/products/categories');
    return response.data.categories;
  }
}

// カートAPI
export class CartApi {
  private client: AxiosInstance;

  constructor(baseURL: string = process.env.NEXT_PUBLIC_CART_API_URL || 'http://localhost:8002') {
    this.client = createApiClient(baseURL);
  }

  async getCart(): Promise<Types.Cart> {
    const response = await this.client.get<Types.Cart>('/api/v1/carts');
    return response.data;
  }

  async addItem(item: Types.AddItemRequest): Promise<Types.Cart> {
    const response = await this.client.post<Types.Cart>('/api/v1/carts/items', item);
    return response.data;
  }

  async updateItem(productId: string, quantity: number): Promise<Types.Cart> {
    const response = await this.client.put<Types.Cart>(`/api/v1/carts/items/${productId}`, {
      quantity,
    });
    return response.data;
  }

  async removeItem(productId: string): Promise<Types.Cart> {
    const response = await this.client.delete<Types.Cart>(`/api/v1/carts/items/${productId}`);
    return response.data;
  }

  async clearCart(): Promise<void> {
    await this.client.delete('/api/v1/carts');
  }
}

// 注文API
export class OrderApi {
  private client: AxiosInstance;

  constructor(baseURL: string = process.env.NEXT_PUBLIC_ORDER_API_URL || 'http://localhost:8003') {
    this.client = createApiClient(baseURL);
    this.setupInterceptors();
  }

  private setupInterceptors() {
    this.client.interceptors.request.use((config) => {
      const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });
  }

  async createOrder(orderData: {
    items: Array<{ product_id: string; quantity: number }>;
    shipping_address: string;
    payment_method: string;
    coupon_code?: string;
  }): Promise<Types.Order> {
    const response = await this.client.post<Types.Order>('/api/v1/orders', orderData);
    return response.data;
  }

  async getOrder(id: string): Promise<Types.Order> {
    const response = await this.client.get<Types.Order>(`/api/v1/orders/${id}`);
    return response.data;
  }

  async getOrders(page: number = 1, pageSize: number = 10): Promise<{ orders: Types.Order[]; total: number }> {
    const response = await this.client.get<{ orders: Types.Order[]; total: number }>('/api/v1/orders', {
      params: { page, page_size: pageSize },
    });
    return response.data;
  }
}

// ポイントAPI
export class PointApi {
  private client: AxiosInstance;

  constructor(baseURL: string = process.env.NEXT_PUBLIC_POINT_API_URL || 'http://localhost:8004') {
    this.client = createApiClient(baseURL);
    this.setupInterceptors();
  }

  private setupInterceptors() {
    this.client.interceptors.request.use((config) => {
      const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });
  }

  async getBalance(): Promise<Types.PointBalance> {
    const response = await this.client.get<Types.PointBalance>('/api/v1/points/balance');
    return response.data;
  }

  async getTransactions(page: number = 1, pageSize: number = 20): Promise<{
    transactions: Types.PointTransaction[];
    total: number;
  }> {
    const response = await this.client.get<{ transactions: Types.PointTransaction[]; total: number }>(
      '/api/v1/points/transactions',
      {
        params: { page, page_size: pageSize },
      }
    );
    return response.data;
  }
}

// プロモーションAPI
export class PromotionApi {
  private client: AxiosInstance;

  constructor(baseURL: string = process.env.NEXT_PUBLIC_PROMOTION_API_URL || 'http://localhost:8005') {
    this.client = createApiClient(baseURL);
  }

  async validateCoupon(code: string, orderTotal: number): Promise<{ valid: boolean; discount_amount: number }> {
    const response = await this.client.post<{ valid: boolean; discount_amount: number }>(
      '/api/v1/promotions/validate-coupon',
      {
        code,
        order_total: orderTotal,
      }
    );
    return response.data;
  }

  async getActiveCoupons(): Promise<Types.Coupon[]> {
    const response = await this.client.get<Types.Coupon[]>('/api/v1/promotions/coupons/active');
    return response.data;
  }

  async getActiveSales(): Promise<Types.Sale[]> {
    const response = await this.client.get<Types.Sale[]>('/api/v1/promotions/sales/active');
    return response.data;
  }
}

// API クライアント生成関数
export const createApiClients = () => {
  return {
    product: new ProductApi(),
    cart: new CartApi(),
    order: new OrderApi(),
    point: new PointApi(),
    promotion: new PromotionApi(),
  };
};

export type ApiClients = ReturnType<typeof createApiClients>;
