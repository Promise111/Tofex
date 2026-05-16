"use client";

import { motion } from "framer-motion";
import { ShoppingBag } from "lucide-react";
import Link from "next/link";
import { useCart } from "@/components/cart/cart-context";

const links = [
  { href: "/", label: "Home" },
  { href: "/menu", label: "Menu" },
  { href: "/locations", label: "Locations" },
];

export function SiteHeader() {
  const { count, setOpen } = useCart();

  return (
    <motion.header
      initial={{ y: -12, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      className="sticky top-0 z-30 border-b border-border/80 bg-background/85 backdrop-blur-md"
    >
      <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-4 sm:px-6">
        <Link href="/" className="flex items-center gap-2">
          <span className="flex h-9 w-9 items-center justify-center rounded-xl bg-brand text-sm font-bold text-brand-foreground">
            T
          </span>
          <span className="text-lg font-semibold tracking-tight">Tofex</span>
        </Link>
        <nav className="hidden items-center gap-8 sm:flex">
          {links.map((l) => (
            <Link
              key={l.href}
              href={l.href}
              className="text-sm text-muted-foreground transition hover:text-foreground"
            >
              {l.label}
            </Link>
          ))}
        </nav>
        <div className="flex items-center gap-3">
          <Link
            href="/admin/login"
            className="hidden text-sm text-muted-foreground hover:text-foreground sm:block"
          >
            Staff
          </Link>
          <button
            type="button"
            onClick={() => setOpen(true)}
            className="relative flex items-center gap-2 rounded-xl bg-brand px-4 py-2 text-sm font-medium text-brand-foreground transition hover:opacity-90"
          >
            <ShoppingBag className="h-4 w-4" />
            Cart
            {count > 0 ? (
              <span className="absolute -right-1 -top-1 flex h-5 w-5 items-center justify-center rounded-full bg-foreground text-[10px] font-bold text-background">
                {count}
              </span>
            ) : null}
          </button>
        </div>
      </div>
    </motion.header>
  );
}
