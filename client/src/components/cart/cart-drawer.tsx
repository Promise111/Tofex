"use client";

import { AnimatePresence, motion } from "framer-motion";
import { Minus, Plus, ShoppingBag, X } from "lucide-react";
import Link from "next/link";
import { formatMoney } from "@/lib/format";
import { Button } from "@/components/ui/button";
import { useCart } from "./cart-context";

export function CartDrawer() {
  const { lines, open, setOpen, setQty, remove, totalCents, count } = useCart();

  return (
    <AnimatePresence>
      {open ? (
        <>
          <motion.div
            className="fixed inset-0 z-40 bg-black/40 backdrop-blur-sm"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={() => setOpen(false)}
          />
          <motion.aside
            className="fixed right-0 top-0 z-50 flex h-full w-full max-w-md flex-col bg-surface shadow-2xl"
            initial={{ x: "100%" }}
            animate={{ x: 0 }}
            exit={{ x: "100%" }}
            transition={{ type: "spring", damping: 28, stiffness: 320 }}
          >
            <motion.div className="flex items-center justify-between border-b border-border px-6 py-4">
              <div className="flex items-center gap-2">
                <ShoppingBag className="h-5 w-5 text-brand" />
                <h2 className="text-lg font-semibold">Your order</h2>
                <span className="rounded-full bg-brand/15 px-2 py-0.5 text-xs font-medium text-brand">
                  {count}
                </span>
              </div>
              <button
                type="button"
                onClick={() => setOpen(false)}
                className="rounded-lg p-2 hover:bg-muted"
                aria-label="Close cart"
              >
                <X className="h-5 w-5" />
              </button>
            </motion.div>

            <div className="flex-1 overflow-y-auto px-6 py-4">
              {lines.length === 0 ? (
                <p className="text-center text-sm text-muted-foreground py-12">
                  Your cart is empty. Browse the menu to add treats.
                </p>
              ) : (
                <ul className="space-y-4">
                  {lines.map((line) => (
                    <motion.li
                      key={line.product.id}
                      layout
                      className="flex gap-4 rounded-xl border border-border p-3"
                    >
                      <div className="flex-1">
                        <p className="font-medium">{line.product.name}</p>
                        <p className="text-sm text-muted-foreground">
                          {formatMoney(line.product.price_cents)}
                        </p>
                        <motion.div className="mt-2 flex items-center gap-2">
                          <button
                            type="button"
                            className="rounded-lg border border-border p-1 hover:bg-muted"
                            onClick={() => setQty(line.product.id, line.quantity - 1)}
                          >
                            <Minus className="h-3.5 w-3.5" />
                          </button>
                          <span className="w-6 text-center text-sm">{line.quantity}</span>
                          <button
                            type="button"
                            className="rounded-lg border border-border p-1 hover:bg-muted"
                            onClick={() => setQty(line.product.id, line.quantity + 1)}
                          >
                            <Plus className="h-3.5 w-3.5" />
                          </button>
                          <button
                            type="button"
                            className="ml-auto text-xs text-red-600 hover:underline"
                            onClick={() => remove(line.product.id)}
                          >
                            Remove
                          </button>
                        </motion.div>
                      </div>
                    </motion.li>
                  ))}
                </ul>
              )}
            </div>

            <div className="border-t border-border px-6 py-4 space-y-3">
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Subtotal</span>
                <span className="font-semibold">{formatMoney(totalCents)}</span>
              </div>
              <Link href="/checkout" onClick={() => setOpen(false)}>
                <Button className="w-full" disabled={lines.length === 0}>
                  Checkout
                </Button>
              </Link>
            </div>
          </motion.aside>
        </>
      ) : null}
    </AnimatePresence>
  );
}
