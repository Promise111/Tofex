export interface RolePermission {
  role_id: number;
  permission: string;
}

export interface Role {
  id: number;
  name: string;
  label: string;
  description?: string;
  permissions?: RolePermission[];
}

export interface User {
  id: number;
  email: string;
  username: string;
  display_name: string;
  roles?: Role[];
}

export interface ProductImage {
  id: number;
  product_id: number;
  path: string;
  sort_order: number;
}

export interface Product {
  id: number;
  name: string;
  slug: string;
  description: string;
  price_cents: number;
  active: boolean;
  images?: ProductImage[];
}

export interface StoreBranch {
  id: number;
  name: string;
  address: string;
  city?: string;
  phone?: string;
  hours?: string;
  maps_url?: string;
  latitude?: number | null;
  longitude?: number | null;
  active: boolean;
  sort_order: number;
}

export interface PaymentAccount {
  id: number;
  bank_name: string;
  account_name: string;
  account_number: string;
  branch?: string;
  currency: string;
  display_label?: string;
  active: boolean;
}

export interface OrderItem {
  id: number;
  product_id: number;
  product_name: string;
  unit_price_cents: number;
  quantity: number;
}

export interface OrderReceipt {
  id: number;
  order_id: string;
  path: string;
  mime: string;
  size_bytes: number;
}

export interface Order {
  id: string;
  customer_name: string;
  customer_email: string;
  customer_phone?: string;
  payment_account_id: number;
  payment_snapshot: string;
  total_cents: number;
  status: string;
  customer_note?: string;
  created_at: string;
  items?: OrderItem[];
  receipts?: OrderReceipt[];
}

export interface ApiError {
  error: { code: string; message: string };
}

export interface CartLine {
  product: Product;
  quantity: number;
}
