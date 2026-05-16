"use client";

import { useEffect, useState } from "react";
import { adminApi } from "@/lib/api";
import { getToken } from "@/lib/auth";
import { formatMoney } from "@/lib/format";
import type { Product } from "@/lib/types";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { FadeIn } from "@/components/motion/fade-in";

export default function AdminProductsPage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ name: "", description: "", price: "" });
  const [imageFile, setImageFile] = useState<File | null>(null);
  const [targetId, setTargetId] = useState<number | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const token = getToken();

  function load() {
    if (!token) return;
    adminApi.listProducts(token).then(setProducts);
  }

  useEffect(() => {
    load();
  }, []);

  async function createProduct(e: React.FormEvent) {
    e.preventDefault();
    if (!token) return;
    const cents = Math.round(parseFloat(form.price) * 100);
    if (!form.name || !cents) {
      setError("Name and price are required");
      return;
    }
    setLoading(true);
    setError(null);
    try {
      const p = await adminApi.createProduct(token, {
        name: form.name,
        description: form.description,
        price_cents: cents,
        active: true,
      });
      if (imageFile) {
        await adminApi.uploadProductImage(token, p.id, imageFile);
      }
      setForm({ name: "", description: "", price: "" });
      setImageFile(null);
      setShowForm(false);
      load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create product");
    } finally {
      setLoading(false);
    }
  }

  async function uploadFor(id: number) {
    if (!token || !imageFile) return;
    setLoading(true);
    try {
      await adminApi.uploadProductImage(token, id, imageFile);
      setImageFile(null);
      setTargetId(null);
      load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Upload failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="space-y-6">
      <FadeIn className="flex flex-wrap items-center justify-between gap-4">
        <h1 className="text-2xl font-semibold">Products</h1>
        <Button onClick={() => setShowForm(!showForm)} variant="secondary">
          {showForm ? "Cancel" : "Add product"}
        </Button>
      </FadeIn>

      {showForm ? (
        <Card>
          <form onSubmit={createProduct} className="space-y-4">
            <Input label="Name" required value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
            <Input label="Price (NGN)" required type="number" step="0.01" value={form.price} onChange={(e) => setForm({ ...form, price: e.target.value })} />
            <div className="space-y-1.5">
              <label className="text-sm font-medium text-muted-foreground">Description</label>
              <textarea
                className="w-full rounded-xl border border-border bg-surface px-4 py-2.5 text-sm"
                rows={3}
                value={form.description}
                onChange={(e) => setForm({ ...form, description: e.target.value })}
              />
            </div>
            <div>
              <label className="text-sm font-medium text-muted-foreground">Image (optional)</label>
              <input type="file" accept="image/*" className="mt-1 block text-sm" onChange={(e) => setImageFile(e.target.files?.[0] ?? null)} />
            </div>
            {error ? <p className="text-sm text-red-600">{error}</p> : null}
            <Button type="submit" loading={loading}>Save product</Button>
          </form>
        </Card>
      ) : null}

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {products.map((p) => (
          <Card key={p.id}>
            <p className="font-semibold">{p.name}</p>
            <p className="text-sm text-muted-foreground line-clamp-2">{p.description}</p>
            <p className="mt-2 font-medium text-brand">{formatMoney(p.price_cents)}</p>
            <p className="mt-1 text-xs text-muted-foreground">
              {p.active ? "Active" : "Inactive"} · {p.images?.length ?? 0} image(s)
            </p>
            <div className="mt-3 flex flex-col gap-2">
              <input
                type="file"
                accept="image/*"
                className="text-xs"
                onChange={(e) => {
                  setTargetId(p.id);
                  setImageFile(e.target.files?.[0] ?? null);
                }}
              />
              {targetId === p.id && imageFile ? (
                <Button variant="secondary" loading={loading} onClick={() => uploadFor(p.id)}>
                  Upload image
                </Button>
              ) : null}
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
}
