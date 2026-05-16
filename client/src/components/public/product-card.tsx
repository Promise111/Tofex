"use client";

import { motion } from "framer-motion";
import Image from "next/image";
import { Plus } from "lucide-react";
import { productImageUrl } from "@/lib/api";
import { formatMoney } from "@/lib/format";
import type { Product } from "@/lib/types";
import { useCart } from "@/components/cart/cart-context";

export function ProductCard({ product }: { product: Product }) {
  const { add } = useCart();
  const img = product.images?.[0];

  return (
    <motion.article
      layout
      whileHover={{ y: -4 }}
      transition={{ type: "spring", stiffness: 400, damping: 28 }}
      className="group overflow-hidden rounded-2xl border border-border bg-surface shadow-sm"
    >
      <div className="relative aspect-[4/3] bg-muted overflow-hidden">
        {img ? (
          <Image
            src={productImageUrl(img.id)}
            alt={product.name}
            fill
            className="object-cover transition duration-500 group-hover:scale-105"
            sizes="(max-width: 768px) 100vw, 33vw"
            unoptimized
          />
        ) : (
          <div className="flex h-full items-center justify-center text-4xl opacity-30">🧁</div>
        )}
      </div>
      <div className="p-4">
        <h3 className="font-semibold">{product.name}</h3>
        {product.description ? (
          <p className="mt-1 line-clamp-2 text-sm text-muted-foreground">{product.description}</p>
        ) : null}
        <div className="mt-4 flex items-center justify-between">
          <span className="text-lg font-semibold text-brand">{formatMoney(product.price_cents)}</span>
          <button
            type="button"
            onClick={() => add(product)}
            className="flex items-center gap-1 rounded-xl bg-foreground px-3 py-2 text-xs font-medium text-background transition hover:opacity-90"
          >
            <Plus className="h-3.5 w-3.5" /> Add
          </button>
        </div>
      </div>
    </motion.article>
  );
}
