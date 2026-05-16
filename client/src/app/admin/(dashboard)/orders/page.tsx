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

const statuses = ["", "pending", "confirmed", "ready", "completed", "cancelled"];

export default function AdminOrdersPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [status, setStatus] = useState("");
  const [loading, setLoading] = useState(true);

  function load() {
    const token = getToken();
    if (!token) return;
    setLoading(true);
    adminApi
      .listOrders(token, status ? { status } : undefined)
      .then((r) => setOrders(r.orders))
      .finally(() => setLoading(false));
  }

  useEffect(() => {
    load();
  }, [status]);

  return (
    <div className="space-y-6">
      <FadeIn>
        <h1 className="text-2xl font-semibold">Orders</h1>
      </FadeIn>
      <select
        className="rounded-xl border border-border bg-surface px-4 py-2 text-sm"
        value={status}
        onChange={(e) => setStatus(e.target.value)}
      >
        {statuses.map((s) => (
          <option key={s || "all"} value={s}>
            {s ? orderStatusLabel(s) : "All statuses"}
          </option>
        ))}
      </select>
      <Card className="overflow-hidden p-0">
        {loading ? (
          <p className="p-8 text-center text-sm text-muted-foreground">Loading…</p>
        ) : orders.length === 0 ? (
          <p className="p-8 text-center text-sm text-muted-foreground">No orders found</p>
        ) : (
          <table className="w-full text-sm">
            <thead className="border-b border-border bg-muted/50 text-left text-muted-foreground">
              <tr>
                <th className="px-4 py-3 font-medium">Customer</th>
                <th className="px-4 py-3 font-medium">Date</th>
                <th className="px-4 py-3 font-medium">Total</th>
                <th className="px-4 py-3 font-medium">Status</th>
                <th className="px-4 py-3" />
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {orders.map((o) => (
                <tr key={o.id} className="hover:bg-muted/30">
                  <td className="px-4 py-3">
                    <p className="font-medium">{o.customer_name}</p>
                    <p className="text-xs text-muted-foreground">{o.customer_email}</p>
                  </td>
                  <td className="px-4 py-3 text-muted-foreground">{formatDate(o.created_at)}</td>
                  <td className="px-4 py-3 font-medium">{formatMoney(o.total_cents)}</td>
                  <td className="px-4 py-3">
                    <Badge status={o.status}>{orderStatusLabel(o.status)}</Badge>
                  </td>
                  <td className="px-4 py-3 text-right">
                    <Link href={`/admin/orders/${o.id}`} className="text-brand hover:underline">
                      Details
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </Card>
    </div>
  );
}
