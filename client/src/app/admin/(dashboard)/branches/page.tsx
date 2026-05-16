"use client";

import { useEffect, useState } from "react";
import { adminApi } from "@/lib/api";
import { getToken } from "@/lib/auth";
import type { StoreBranch } from "@/lib/types";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { FadeIn } from "@/components/motion/fade-in";

const emptyForm = {
  name: "",
  address: "",
  city: "",
  phone: "",
  hours: "",
  maps_url: "",
  latitude: "",
  longitude: "",
  sort_order: "0",
};

export default function AdminBranchesPage() {
  const [branches, setBranches] = useState<StoreBranch[]>([]);
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [form, setForm] = useState(emptyForm);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const token = getToken();

  function load() {
    if (!token) return;
    adminApi.listBranches(token).then(setBranches);
  }

  useEffect(() => {
    load();
  }, []);

  function startEdit(b: StoreBranch) {
    setEditingId(b.id);
    setShowForm(true);
    setForm({
      name: b.name,
      address: b.address,
      city: b.city ?? "",
      phone: b.phone ?? "",
      hours: b.hours ?? "",
      maps_url: b.maps_url ?? "",
      latitude: b.latitude != null ? String(b.latitude) : "",
      longitude: b.longitude != null ? String(b.longitude) : "",
      sort_order: String(b.sort_order),
    });
  }

  function resetForm() {
    setEditingId(null);
    setForm(emptyForm);
    setShowForm(false);
    setError(null);
  }

  function parseCoord(value: string): number | undefined {
    const trimmed = value.trim();
    if (!trimmed) return undefined;
    const n = Number(trimmed);
    return Number.isFinite(n) ? n : undefined;
  }

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!token) return;
    setLoading(true);
    setError(null);
    const lat = parseCoord(form.latitude);
    const lng = parseCoord(form.longitude);
    const body = {
      name: form.name,
      address: form.address,
      city: form.city,
      phone: form.phone,
      hours: form.hours,
      maps_url: form.maps_url,
      latitude: lat,
      longitude: lng,
      sort_order: Number(form.sort_order) || 0,
      active: true,
    };
    try {
      if (editingId) {
        await adminApi.patchBranch(token, editingId, body);
      } else {
        await adminApi.createBranch(token, body);
      }
      resetForm();
      load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save branch");
    } finally {
      setLoading(false);
    }
  }

  async function toggleActive(b: StoreBranch) {
    if (!token) return;
    await adminApi.patchBranch(token, b.id, { active: !b.active });
    load();
  }

  async function onDelete(id: number) {
    if (!token || !confirm("Remove this branch from the storefront?")) return;
    await adminApi.deleteBranch(token, id);
    load();
  }

  return (
    <div className="space-y-6">
      <FadeIn className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold">Store branches</h1>
          <p className="text-sm text-muted-foreground">
            Pickup locations customers see on the website
          </p>
        </div>
        <Button
          variant="secondary"
          onClick={() => {
            if (showForm) resetForm();
            else setShowForm(true);
          }}
        >
          {showForm ? "Cancel" : "Add branch"}
        </Button>
      </FadeIn>

      {showForm ? (
        <Card>
          <form onSubmit={onSubmit} className="grid gap-4 sm:grid-cols-2">
            <Input
              label="Name"
              required
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
            />
            <Input
              label="City"
              value={form.city}
              onChange={(e) => setForm({ ...form, city: e.target.value })}
            />
            <Input
              label="Address"
              required
              className="sm:col-span-2"
              value={form.address}
              onChange={(e) => setForm({ ...form, address: e.target.value })}
            />
            <Input
              label="Phone"
              value={form.phone}
              onChange={(e) => setForm({ ...form, phone: e.target.value })}
            />
            <Input
              label="Hours"
              placeholder="Mon–Sat 9am–6pm"
              value={form.hours}
              onChange={(e) => setForm({ ...form, hours: e.target.value })}
            />
            <Input
              label="Google Maps link"
              className="sm:col-span-2"
              placeholder="https://maps.google.com/..."
              value={form.maps_url}
              onChange={(e) => setForm({ ...form, maps_url: e.target.value })}
            />
            <Input
              label="Latitude"
              placeholder="6.5244"
              value={form.latitude}
              onChange={(e) => setForm({ ...form, latitude: e.target.value })}
            />
            <Input
              label="Longitude"
              placeholder="3.3792"
              value={form.longitude}
              onChange={(e) => setForm({ ...form, longitude: e.target.value })}
            />
            <Input
              label="Sort order"
              type="number"
              value={form.sort_order}
              onChange={(e) => setForm({ ...form, sort_order: e.target.value })}
            />
            {error ? (
              <p className="text-sm text-red-600 sm:col-span-2">{error}</p>
            ) : null}
            <Button type="submit" loading={loading} className="sm:col-span-2">
              {editingId ? "Update branch" : "Save branch"}
            </Button>
          </form>
        </Card>
      ) : null}

      <div className="grid gap-4 md:grid-cols-2">
        {branches.map((b) => (
          <Card key={b.id}>
            <p className="font-semibold">{b.name}</p>
            {b.city ? <p className="text-sm text-muted-foreground">{b.city}</p> : null}
            <p className="mt-2 text-sm">{b.address}</p>
            {b.phone ? <p className="text-sm">{b.phone}</p> : null}
            {b.hours ? (
              <p className="text-sm text-muted-foreground">{b.hours}</p>
            ) : null}
            {b.maps_url ? (
              <p className="mt-1 truncate text-xs text-muted-foreground">{b.maps_url}</p>
            ) : null}
            <p className="mt-2 text-xs text-muted-foreground">
              {b.active ? "Visible to customers" : "Hidden"} · sort {b.sort_order}
              {b.latitude != null && b.longitude != null
                ? ` · ${b.latitude}, ${b.longitude}`
                : ""}
            </p>
            <div className="mt-4 flex flex-wrap gap-2">
              <Button variant="secondary" onClick={() => startEdit(b)}>
                Edit
              </Button>
              <Button variant="ghost" onClick={() => toggleActive(b)}>
                {b.active ? "Hide" : "Show"}
              </Button>
              <Button variant="danger" onClick={() => onDelete(b.id)}>
                Delete
              </Button>
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
}
