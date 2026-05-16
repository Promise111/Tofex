"use client";

import { motion } from "framer-motion";
import { CheckCircle2 } from "lucide-react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { Suspense } from "react";
import { Button } from "@/components/ui/button";

function SuccessContent() {
  const params = useSearchParams();
  const id = params.get("id");

  return (
    <div className="mx-auto max-w-lg px-4 py-20 text-center sm:px-6">
      <motion.div
        initial={{ scale: 0.8, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ type: "spring", damping: 14 }}
      >
        <CheckCircle2 className="mx-auto h-16 w-16 text-emerald-600" />
        <h1 className="mt-6 text-3xl font-semibold">Order received</h1>
        <p className="mt-3 text-muted-foreground">
          Thank you! We&apos;ll review your payment and confirm your order for pickup.
        </p>
        {id ? (
          <p className="mt-4 rounded-xl bg-muted px-4 py-2 font-mono text-sm">Order ID: {id}</p>
        ) : null}
        <Link href="/menu" className="mt-10 inline-block">
          <Button variant="secondary">Back to menu</Button>
        </Link>
      </motion.div>
    </div>
  );
}

export default function OrderSuccessPage() {
  return (
    <Suspense>
      <SuccessContent />
    </Suspense>
  );
}
