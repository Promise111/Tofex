"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { adminApi } from "@/lib/api";
import { getToken } from "@/lib/auth";
import { formatMoney, formatDate, orderStatusLabel } from "@/lib/format";
import type { Order } from "@/lib/types";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { FadeIn } from "@/components/motion/fade-in";

export default function AdminDashboardPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [total, setTotal] = useState(0);

  useEffect(() => {
    const token = getToken();
    if (!token) return;
    adminApi.listOrders(token).then((res) => {
      setOrders(res.orders.slice(0, 5));
      setTotal(res.total);
    });
  }, []);

  const pending = orders.filter((o) => o.status === "pending").length;

  return (
    <div className="space-y-8">
      <FadeIn>
        <h1 className="text-2xl font-semibold">Dashboard</h1>
        <p className="text-muted-foreground">Overview of recent activity</p>
      </FadeIn>

      <div className="grid gap-4 sm:grid-cols-3">
        {[
          { label: "Total orders", value: String(total) },
          { label: "Recent pending", value: String(pending) },
          { label: "Status", value: "Live" },
        ].map((s, i) => (
          <FadeIn key={s.label} delay={i * 0.05}>
            <Card>
              <p className="text-sm text-muted-foreground">{s.label}</p>
              <p className="mt-2 text-3xl font-semibold">{s.value}</p>
            </Card>
          </FadeIn>
        ))}
      </div>

      <Card>
        <div className="flex items-center justify-between">
          <h2 className="font-semibold">Recent orders</h2>
          <Link href="/admin/orders" className="text-sm text-brand hover:underline">
            View all
          </Link>
        </div>
        <ul className="mt-4 divide-y divide-border">
          {orders.length === 0 ? (
            <li className="py-8 text-center text-sm text-muted-foreground">No orders yet</li>
          ) : (
            orders.map((o) => (
              <li key={o.id} className="flex flex-wrap items-center justify-between gap-2 py-3">
                <div>
                  <p className="font-medium">{o.customer_name}</p>
                  <p className="text-xs text-muted-foreground">{formatDate(o.created_at)}</p>
                </div>
                <div className="flex items-center gap-3">
                  <Badge status={o.status}>{orderStatusLabel(o.status)}</Badge>
                  <span className="font-medium">{formatMoney(o.total_cents)}</span>
                  <Link href={`/admin/orders/${o.id}`} className="text-sm text-brand hover:underline">
                    View
                  </Link>
                </div>
              </li>
            ))
          )}
        </ul>
      </Card>
    </div>
  );
}
