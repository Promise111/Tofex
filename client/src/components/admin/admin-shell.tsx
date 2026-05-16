"use client";

import { motion } from "framer-motion";
import {
  CreditCard,
  LayoutDashboard,
  LogOut,
  MapPin,
  Package,
  ShoppingCart,
} from "lucide-react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { adminApi } from "@/lib/api";
import {
  clearSession,
  getToken,
  getStoredUser,
  hasPermission,
  isSuperAdmin,
} from "@/lib/auth";
import type { User } from "@/lib/types";

const nav = [
  { href: "/admin", label: "Dashboard", icon: LayoutDashboard, perm: null },
  { href: "/admin/orders", label: "Orders", icon: ShoppingCart, perm: "orders.read" },
  { href: "/admin/products", label: "Products", icon: Package, perm: "products.read" },
  { href: "/admin/branches", label: "Branches", icon: MapPin, perm: "branches.read" },
  {
    href: "/admin/payment-accounts",
    label: "Payment accounts",
    icon: CreditCard,
    perm: "payment_accounts.read",
  },
];

export function AdminShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [ready, setReady] = useState(false);

  useEffect(() => {
    const token = getToken();
    if (!token) {
      router.replace("/admin/login");
      return;
    }
    adminApi
      .me(token)
      .then(setUser)
      .catch(() => {
        clearSession();
        router.replace("/admin/login");
      })
      .finally(() => setReady(true));
  }, [router]);

  if (!ready) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-stone-100">
        <div className="h-8 w-8 animate-spin rounded-full border-2 border-stone-400 border-t-transparent" />
      </div>
    );
  }

  const visibleNav = nav.filter(
    (n) => !n.perm || hasPermission(user, n.perm) || isSuperAdmin(user),
  );

  return (
    <div className="flex min-h-screen bg-stone-100">
      <aside className="hidden w-64 flex-col bg-[var(--admin-sidebar)] text-[var(--admin-sidebar-fg)] lg:flex">
        <div className="border-b border-white/10 px-6 py-5">
          <Link href="/admin" className="flex items-center gap-2 font-semibold">
            <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-brand text-sm text-brand-foreground">
              T
            </span>
            Tofex Admin
          </Link>
          <p className="mt-2 truncate text-xs text-stone-400">
            {user?.display_name || user?.username}
          </p>
        </div>
        <nav className="flex-1 space-y-1 p-4">
          {visibleNav.map((item) => {
            const active = pathname === item.href || pathname.startsWith(item.href + "/");
            const Icon = item.icon;
            return (
              <Link key={item.href} href={item.href}>
                <span
                  className={`flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm transition ${
                    active ? "bg-white/10 text-white" : "text-stone-400 hover:bg-white/5 hover:text-white"
                  }`}
                >
                  <Icon className="h-4 w-4" />
                  {item.label}
                </span>
              </Link>
            );
          })}
        </nav>
        <div className="border-t border-white/10 p-4">
          <Link href="/" className="mb-2 block text-xs text-stone-500 hover:text-stone-300">
            ← View storefront
          </Link>
          <button
            type="button"
            onClick={() => {
              clearSession();
              router.push("/admin/login");
            }}
            className="flex w-full items-center gap-2 rounded-xl px-3 py-2 text-sm text-stone-400 hover:bg-white/5 hover:text-white"
          >
            <LogOut className="h-4 w-4" /> Sign out
          </button>
        </div>
      </aside>

      <div className="flex flex-1 flex-col">
        <header className="flex items-center justify-between border-b border-border bg-surface px-4 py-3 lg:px-8">
          <h1 className="text-lg font-semibold lg:hidden">Tofex Admin</h1>
          <div className="ml-auto flex gap-2 overflow-x-auto lg:hidden">
            {visibleNav.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className={`whitespace-nowrap rounded-lg px-3 py-1.5 text-xs ${
                  pathname === item.href ? "bg-brand/20 text-brand" : "bg-muted"
                }`}
              >
                {item.label}
              </Link>
            ))}
          </div>
        </header>
        <motion.main
          key={pathname}
          initial={{ opacity: 0, y: 8 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3 }}
          className="flex-1 p-4 lg:p-8"
        >
          {children}
        </motion.main>
      </div>
    </div>
  );
}
