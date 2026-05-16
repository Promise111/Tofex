"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Image from "next/image";
import { adminApi, orderReceiptUrl } from "@/lib/api";
import { getToken } from "@/lib/auth";
import { formatMoney, formatDate, orderStatusLabel } from "@/lib/format";
import type { Order } from "@/lib/types";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { FadeIn } from "@/components/motion/fade-in";

const nextStatuses = ["pending", "confirmed", "ready", "completed", "cancelled"];

export default function AdminOrderDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [order, setOrder] = useState<Order | null>(null);
  const [loading, setLoading] = useState(false);
  const token = getToken();

  function load() {
    if (!token || !id) return;
    adminApi.getOrder(token, id).then(setOrder);
  }

  useEffect(() => {
    load();
  }, [id]);

  async function updateStatus(status: string) {
    if (!token || !id) return;
    setLoading(true);
    try {
      const updated = await adminApi.patchOrderStatus(token, id, status);
      setOrder(updated);
    } finally {
      setLoading(false);
    }
  }

  if (!order) {
    return <p className="text-muted-foreground">Loading order…</p>;
  }

  let paymentSnap: Record<string, string> = {};
  try {
    paymentSnap = JSON.parse(order.payment_snapshot);
  } catch {
    /* ignore */
  }

  return (
    <div className="mx-auto max-w-3xl space-y-6">
      <FadeIn>
        <div className="flex flex-wrap items-start justify-between gap-4">
          <div>
            <h1 className="text-2xl font-semibold">{order.customer_name}</h1>
            <p className="text-muted-foreground">{formatDate(order.created_at)}</p>
          </div>
          <Badge status={order.status}>{orderStatusLabel(order.status)}</Badge>
        </div>
      </FadeIn>

      <Card>
        <h2 className="font-semibold">Contact</h2>
        <p className="mt-2 text-sm">{order.customer_email}</p>
        {order.customer_phone ? <p className="text-sm">{order.customer_phone}</p> : null}
        {order.customer_note ? (
          <p className="mt-2 text-sm text-muted-foreground">Note: {order.customer_note}</p>
        ) : null}
      </Card>

      <Card>
        <h2 className="font-semibold">Items</h2>
        <ul className="mt-3 space-y-2">
          {order.items?.map((it) => (
            <li key={it.id} className="flex justify-between text-sm">
              <span>
                {it.product_name} × {it.quantity}
              </span>
              <span>{formatMoney(it.unit_price_cents * it.quantity)}</span>
            </li>
          ))}
        </ul>
        <p className="mt-4 border-t border-border pt-3 text-right font-semibold">
          Total {formatMoney(order.total_cents)}
        </p>
      </Card>

      <Card>
        <h2 className="font-semibold">Payment destination</h2>
        <p className="mt-2 text-sm">{paymentSnap.bank_name}</p>
        <p className="text-sm">{paymentSnap.account_name}</p>
        <p className="font-mono text-sm">{paymentSnap.account_number}</p>
      </Card>

      {order.receipts?.length ? (
        <Card>
          <h2 className="font-semibold">Receipt</h2>
          {order.receipts.map((r) =>
            token ? (
              <div key={r.id} className="relative mt-4 aspect-video max-h-80 overflow-hidden rounded-xl bg-muted">
                <Image
                  src={orderReceiptUrl(order.id, r.id, token)}
                  alt="Payment receipt"
                  fill
                  className="object-contain"
                  unoptimized
                />
              </div>
            ) : null,
          )}
        </Card>
      ) : null}

      <Card>
        <h2 className="mb-3 font-semibold">Update status</h2>
        <div className="flex flex-wrap gap-2">
          {nextStatuses.map((s) => (
            <Button
              key={s}
              variant={order.status === s ? "primary" : "secondary"}
              disabled={loading || order.status === s}
              onClick={() => updateStatus(s)}
            >
              {orderStatusLabel(s)}
            </Button>
          ))}
        </div>
      </Card>
    </div>
  );
}
