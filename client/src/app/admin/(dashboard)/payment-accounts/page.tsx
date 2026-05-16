"use client";

import { useEffect, useState } from "react";
import { adminApi } from "@/lib/api";
import { getToken } from "@/lib/auth";
import type { PaymentAccount } from "@/lib/types";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { FadeIn } from "@/components/motion/fade-in";

export default function AdminPaymentAccountsPage() {
  const [accounts, setAccounts] = useState<PaymentAccount[]>([]);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({
    bank_name: "",
    account_name: "",
    account_number: "",
    branch: "",
    display_label: "",
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const token = getToken();

  function load() {
    if (!token) return;
    adminApi.listPaymentAccounts(token).then(setAccounts);
  }

  useEffect(() => {
    load();
  }, []);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!token) return;
    setLoading(true);
    setError(null);
    try {
      await adminApi.createPaymentAccount(token, {
        bank_name: form.bank_name,
        account_name: form.account_name,
        account_number: form.account_number,
        branch: form.branch,
        display_label: form.display_label,
        currency: "NGN",
        active: true,
      });
      setForm({ bank_name: "", account_name: "", account_number: "", branch: "", display_label: "" });
      setShowForm(false);
      load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create account");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="space-y-6">
      <FadeIn className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold">Payment accounts</h1>
          <p className="text-sm text-muted-foreground">Bank details customers use when paying</p>
        </div>
        <Button variant="secondary" onClick={() => setShowForm(!showForm)}>
          {showForm ? "Cancel" : "Add account"}
        </Button>
      </FadeIn>

      {showForm ? (
        <Card>
          <form onSubmit={onSubmit} className="grid gap-4 sm:grid-cols-2">
            <Input label="Bank name" required value={form.bank_name} onChange={(e) => setForm({ ...form, bank_name: e.target.value })} />
            <Input label="Account name" required value={form.account_name} onChange={(e) => setForm({ ...form, account_name: e.target.value })} />
            <Input label="Account number" required value={form.account_number} onChange={(e) => setForm({ ...form, account_number: e.target.value })} />
            <Input label="Branch" value={form.branch} onChange={(e) => setForm({ ...form, branch: e.target.value })} />
            <Input label="Display label" className="sm:col-span-2" value={form.display_label} onChange={(e) => setForm({ ...form, display_label: e.target.value })} />
            {error ? <p className="text-sm text-red-600 sm:col-span-2">{error}</p> : null}
            <Button type="submit" loading={loading} className="sm:col-span-2">
              Save account
            </Button>
          </form>
        </Card>
      ) : null}

      <div className="grid gap-4 md:grid-cols-2">
        {accounts.map((a) => (
          <Card key={a.id}>
            <p className="font-semibold">{a.display_label || a.bank_name}</p>
            <p className="mt-2 text-sm">{a.account_name}</p>
            <p className="font-mono text-sm">{a.account_number}</p>
            {a.branch ? <p className="text-sm text-muted-foreground">{a.branch}</p> : null}
            <p className="mt-2 text-xs text-muted-foreground">{a.active ? "Active" : "Inactive"}</p>
          </Card>
        ))}
      </div>
    </div>
  );
}
