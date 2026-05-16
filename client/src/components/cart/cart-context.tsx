"use client";

import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useState,
} from "react";
import type { CartLine, Product } from "@/lib/types";

type CartContextValue = {
  lines: CartLine[];
  count: number;
  totalCents: number;
  add: (product: Product, qty?: number) => void;
  setQty: (productId: number, qty: number) => void;
  remove: (productId: number) => void;
  clear: () => void;
  open: boolean;
  setOpen: (v: boolean) => void;
};

const CartContext = createContext<CartContextValue | null>(null);

export function CartProvider({ children }: { children: React.ReactNode }) {
  const [lines, setLines] = useState<CartLine[]>([]);
  const [open, setOpen] = useState(false);

  const add = useCallback((product: Product, qty = 1) => {
    setLines((prev) => {
      const i = prev.findIndex((l) => l.product.id === product.id);
      if (i >= 0) {
        const next = [...prev];
        next[i] = { ...next[i], quantity: next[i].quantity + qty };
        return next;
      }
      return [...prev, { product, quantity: qty }];
    });
    setOpen(true);
  }, []);

  const setQty = useCallback((productId: number, qty: number) => {
    if (qty <= 0) {
      setLines((p) => p.filter((l) => l.product.id !== productId));
      return;
    }
    setLines((p) =>
      p.map((l) => (l.product.id === productId ? { ...l, quantity: qty } : l)),
    );
  }, []);

  const remove = useCallback((productId: number) => {
    setLines((p) => p.filter((l) => l.product.id !== productId));
  }, []);

  const clear = useCallback(() => setLines([]), []);

  const count = useMemo(
    () => lines.reduce((s, l) => s + l.quantity, 0),
    [lines],
  );
  const totalCents = useMemo(
    () => lines.reduce((s, l) => s + l.product.price_cents * l.quantity, 0),
    [lines],
  );

  const value = useMemo(
    () => ({
      lines,
      count,
      totalCents,
      add,
      setQty,
      remove,
      clear,
      open,
      setOpen,
    }),
    [lines, count, totalCents, add, setQty, remove, clear, open],
  );

  return <CartContext.Provider value={value}>{children}</CartContext.Provider>;
}

export function useCart() {
  const ctx = useContext(CartContext);
  if (!ctx) throw new Error("useCart must be used within CartProvider");
  return ctx;
}
