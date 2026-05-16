"use client";

import { motion } from "framer-motion";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { Upload } from "lucide-react";
import { useCart } from "@/components/cart/cart-context";
import { FadeIn } from "@/components/motion/fade-in";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card } from "@/components/ui/card";
import { publicApi } from "@/lib/api";
import { formatMoney } from "@/lib/format";
import type { PaymentAccount } from "@/lib/types";

export default function CheckoutPage() {
  const router = useRouter();
  const { lines, totalCents, clear } = useCart();
  const [accounts, setAccounts] = useState<PaymentAccount[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [form, setForm] = useState({
    customer_name: "",
    customer_email: "",
    customer_phone: "",
    customer_note: "",
    payment_account_id: "",
  });
  const [receipt, setReceipt] = useState<File | null>(null);

  useEffect(() => {
    publicApi.getPaymentAccounts().then(setAccounts).catch(() => setError("Could not load payment accounts"));
  }, []);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (lines.length === 0) {
      setError("Your cart is empty");
      return;
    }
    if (!receipt) {
      setError("Please upload your payment receipt");
      return;
    }
    if (!form.payment_account_id) {
      setError("Please select a payment account");
      return;
    }
    setLoading(true);
    setError(null);
    try {
      const fd = new FormData();
      fd.append("customer_name", form.customer_name);
      fd.append("customer_email", form.customer_email);
      if (form.customer_phone) fd.append("customer_phone", form.customer_phone);
      if (form.customer_note) fd.append("customer_note", form.customer_note);
      fd.append("payment_account_id", form.payment_account_id);
      fd.append(
        "items",
        JSON.stringify(lines.map((l) => ({ product_id: l.product.id, quantity: l.quantity }))),
      );
      fd.append("receipt", receipt);
      const res = await publicApi.createOrder(fd);
      clear();
      router.push(`/order/success?id=${res.order_id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Order failed");
    } finally {
      setLoading(false);
    }
  }

  const selected = accounts.find((a) => String(a.id) === form.payment_account_id);

  return (
    <div className="mx-auto max-w-2xl px-4 py-12 sm:px-6">
      <FadeIn>
        <h1 className="text-3xl font-semibold">Checkout</h1>
        <p className="mt-2 text-muted-foreground">Pay by transfer, then upload your receipt to complete the order.</p>
      </FadeIn>

      <form onSubmit={onSubmit} className="mt-10 space-y-8">
        <Card>
          <h2 className="font-semibold">Your details</h2>
          <div className="mt-4 space-y-4">
            <Input
              label="Full name"
              required
              value={form.customer_name}
              onChange={(e) => setForm({ ...form, customer_name: e.target.value })}
            />
            <Input
              label="Email"
              type="email"
              required
              value={form.customer_email}
              onChange={(e) => setForm({ ...form, customer_email: e.target.value })}
            />
            <Input
              label="Phone (optional)"
              value={form.customer_phone}
              onChange={(e) => setForm({ ...form, customer_phone: e.target.value })}
            />
            <div className="space-y-1.5">
              <label className="block text-sm font-medium text-muted-foreground">Note (optional)</label>
              <textarea
                className="w-full rounded-xl border border-border bg-surface px-4 py-2.5 text-sm outline-none focus:border-brand focus:ring-2 focus:ring-brand/20"
                rows={2}
                value={form.customer_note}
                onChange={(e) => setForm({ ...form, customer_note: e.target.value })}
              />
            </div>
          </div>
        </Card>

        <Card>
          <h2 className="font-semibold">Payment account</h2>
          <p className="mt-1 text-sm text-muted-foreground">Transfer the total to one of these accounts.</p>
          <select
            required
            className="mt-4 w-full rounded-xl border border-border bg-surface px-4 py-2.5 text-sm"
            value={form.payment_account_id}
            onChange={(e) => setForm({ ...form, payment_account_id: e.target.value })}
          >
            <option value="">Select account…</option>
            {accounts.map((a) => (
              <option key={a.id} value={a.id}>
                {a.display_label || a.bank_name} — {a.account_number}
              </option>
            ))}
          </select>
          {selected ? (
            <motion.div
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: "auto" }}
              className="mt-4 rounded-xl bg-muted p-4 text-sm"
            >
              <p><strong>{selected.bank_name}</strong></p>
              <p>{selected.account_name}</p>
              <p className="font-mono">{selected.account_number}</p>
              {selected.branch ? <p className="text-muted-foreground">{selected.branch}</p> : null}
            </motion.div>
          ) : null}
        </Card>

        <Card>
          <h2 className="font-semibold">Payment receipt</h2>
          <p className="mt-1 text-sm text-muted-foreground">Upload a screenshot or photo of your transfer receipt.</p>
          <label className="mt-4 flex cursor-pointer flex-col items-center justify-center rounded-xl border-2 border-dashed border-border bg-muted/50 px-6 py-10 transition hover:border-brand">
            <Upload className="h-8 w-8 text-muted-foreground" />
            <span className="mt-2 text-sm font-medium">
              {receipt ? receipt.name : "Click to upload image"}
            </span>
            <input
              type="file"
              accept="image/*"
              className="hidden"
              onChange={(e) => setReceipt(e.target.files?.[0] ?? null)}
            />
          </label>
        </Card>

        <Card>
          <div className="flex justify-between">
            <span className="text-muted-foreground">{lines.length} item(s)</span>
            <span className="text-xl font-semibold">{formatMoney(totalCents)}</span>
          </div>
        </Card>

        {error ? (
          <p className="rounded-xl bg-red-50 px-4 py-3 text-sm text-red-800">{error}</p>
        ) : null}

        <Button type="submit" className="w-full" loading={loading} disabled={lines.length === 0}>
          Place order
        </Button>
      </form>
    </div>
  );
}
