import type {
  ApiError,
  Order,
  PaymentAccount,
  Product,
  User,
} from "./types";

const API_URL =
  process.env.NEXT_PUBLIC_API_URL?.replace(/\/$/, "") ??
  "http://localhost:8080";

export function productImageUrl(imageId: number) {
  return `${API_URL}/api/v1/public/product-images/${imageId}`;
}

export function orderReceiptUrl(orderId: string, receiptId: number, token: string) {
  return `${API_URL}/api/v1/admin/orders/${orderId}/receipts/${receiptId}/file`;
}

async function parseJson<T>(res: Response): Promise<T> {
  const data = await res.json().catch(() => ({}));
  if (!res.ok) {
    const err = data as ApiError;
    throw new Error(err?.error?.message ?? res.statusText ?? "Request failed");
  }
  return data as T;
}

function authHeaders(token?: string | null): HeadersInit {
  const h: HeadersInit = {};
  if (token) h.Authorization = `Bearer ${token}`;
  return h;
}

export const publicApi = {
  async getProducts(): Promise<Product[]> {
    const res = await fetch(`${API_URL}/api/v1/public/products`, {
      next: { revalidate: 30 },
    });
    const data = await parseJson<{ products: Product[] }>(res);
    return data.products ?? [];
  },

  async getPaymentAccounts(): Promise<PaymentAccount[]> {
    const res = await fetch(`${API_URL}/api/v1/public/payment-accounts`, {
      cache: "no-store",
    });
    const data = await parseJson<{ payment_accounts: PaymentAccount[] }>(res);
    return data.payment_accounts ?? [];
  },

  async createOrder(form: FormData) {
    const res = await fetch(`${API_URL}/api/v1/public/orders`, {
      method: "POST",
      body: form,
    });
    return parseJson<{ order_id: string; total_cents: number; status: string }>(
      res,
    );
  },
};

export const authApi = {
  async login(login: string, password: string) {
    const res = await fetch(`${API_URL}/api/v1/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ login, password }),
    });
    return parseJson<{
      access_token: string;
      expires_in: number;
      user: User;
    }>(res);
  },
};

export const adminApi = {
  async me(token: string) {
    const res = await fetch(`${API_URL}/api/v1/admin/me`, {
      headers: authHeaders(token),
      cache: "no-store",
    });
    return parseJson<User>(res);
  },

  async listOrders(token: string, params?: { page?: number; status?: string }) {
    const q = new URLSearchParams();
    if (params?.page) q.set("page", String(params.page));
    if (params?.status) q.set("status", params.status);
    const res = await fetch(`${API_URL}/api/v1/admin/orders?${q}`, {
      headers: authHeaders(token),
      cache: "no-store",
    });
    return parseJson<{
      orders: Order[];
      total: number;
      page: number;
      per_page: number;
    }>(res);
  },

  async getOrder(token: string, id: string) {
    const res = await fetch(`${API_URL}/api/v1/admin/orders/${id}`, {
      headers: authHeaders(token),
      cache: "no-store",
    });
    return parseJson<Order>(res);
  },

  async patchOrderStatus(token: string, id: string, status: string) {
    const res = await fetch(`${API_URL}/api/v1/admin/orders/${id}`, {
      method: "PATCH",
      headers: { ...authHeaders(token), "Content-Type": "application/json" },
      body: JSON.stringify({ status }),
    });
    return parseJson<Order>(res);
  },

  async listProducts(token: string) {
    const res = await fetch(`${API_URL}/api/v1/admin/products`, {
      headers: authHeaders(token),
      cache: "no-store",
    });
    const data = await parseJson<{ products: Product[] }>(res);
    return data.products ?? [];
  },

  async createProduct(
    token: string,
    body: {
      name: string;
      description: string;
      price_cents: number;
      active?: boolean;
    },
  ) {
    const res = await fetch(`${API_URL}/api/v1/admin/products`, {
      method: "POST",
      headers: { ...authHeaders(token), "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });
    return parseJson<Product>(res);
  },

  async uploadProductImage(token: string, productId: number, file: File) {
    const form = new FormData();
    form.append("file", file);
    const res = await fetch(
      `${API_URL}/api/v1/admin/products/${productId}/images`,
      { method: "POST", headers: authHeaders(token), body: form },
    );
    return parseJson<{ id: number }>(res);
  },

  async listPaymentAccounts(token: string) {
    const res = await fetch(`${API_URL}/api/v1/admin/payment-accounts`, {
      headers: authHeaders(token),
      cache: "no-store",
    });
    const data = await parseJson<{ payment_accounts: PaymentAccount[] }>(res);
    return data.payment_accounts ?? [];
  },

  async createPaymentAccount(
    token: string,
    body: Omit<PaymentAccount, "id" | "active"> & { active?: boolean },
  ) {
    const res = await fetch(`${API_URL}/api/v1/admin/payment-accounts`, {
      method: "POST",
      headers: { ...authHeaders(token), "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });
    return parseJson<PaymentAccount>(res);
  },
};
