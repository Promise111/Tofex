"use client";

import { motion } from "framer-motion";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { authApi } from "@/lib/api";
import { setSession } from "@/lib/auth";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card } from "@/components/ui/card";

export default function AdminLoginPage() {
  const router = useRouter();
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError(null);
    try {
      const res = await authApi.login(login, password);
      setSession(res.access_token, res.user);
      router.push("/admin");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <motion.div
      className="flex min-h-screen items-center justify-center bg-stone-100 px-4"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
    >
      <Card className="w-full max-w-md">
        <Link href="/" className="text-sm text-muted-foreground hover:text-foreground">
          ← Back to store
        </Link>
        <h1 className="mt-4 text-2xl font-semibold">Staff sign in</h1>
        <p className="mt-1 text-sm text-muted-foreground">Tofex admin dashboard</p>
        <form onSubmit={onSubmit} className="mt-8 space-y-4">
          <Input
            label="Email or username"
            required
            value={login}
            onChange={(e) => setLogin(e.target.value)}
          />
          <Input
            label="Password"
            type="password"
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          {error ? <p className="text-sm text-red-600">{error}</p> : null}
          <Button type="submit" className="w-full" loading={loading}>
            Sign in
          </Button>
        </form>
      </Card>
    </motion.div>
  );
}
