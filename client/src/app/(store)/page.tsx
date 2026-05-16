import Link from "next/link";
import { ArrowRight, Sparkles } from "lucide-react";
import { FadeIn } from "@/components/motion/fade-in";
import { Button } from "@/components/ui/button";
import { publicApi } from "@/lib/api";
import { ProductCard } from "@/components/public/product-card";

export default async function HomePage() {
  let products: Awaited<ReturnType<typeof publicApi.getProducts>> = [];
  try {
    products = await publicApi.getProducts();
  } catch {
    products = [];
  }
  const featured = products.slice(0, 3);

  return (
    <div>
      <section className="relative overflow-hidden border-b border-border">
        <div className="pointer-events-none absolute -right-20 -top-20 h-72 w-72 rounded-full bg-brand/20 blur-3xl" />
        <div className="mx-auto max-w-6xl px-4 py-20 sm:px-6 sm:py-28">
          <FadeIn>
            <span className="inline-flex items-center gap-2 rounded-full border border-border bg-surface px-3 py-1 text-xs font-medium text-muted-foreground">
              <Sparkles className="h-3.5 w-3.5 text-brand" /> Now taking online orders
            </span>
            <h1 className="mt-6 max-w-2xl text-4xl font-semibold tracking-tight sm:text-5xl lg:text-6xl">
              Handcrafted treats,
              <span className="text-brand"> ready for pickup</span>
            </h1>
            <p className="mt-6 max-w-xl text-lg text-muted-foreground">
              Browse our menu, pay by bank transfer, upload your receipt, and we&apos;ll have your order ready when you arrive.
            </p>
            <div className="mt-10 flex flex-wrap gap-4">
              <Link href="/menu">
                <Button className="gap-2">
                  View menu <ArrowRight className="h-4 w-4" />
                </Button>
              </Link>
              <Link href="/checkout">
                <Button variant="secondary">Go to checkout</Button>
              </Link>
            </div>
          </FadeIn>
        </div>
      </section>

      {featured.length > 0 ? (
        <section className="mx-auto max-w-6xl px-4 py-16 sm:px-6">
          <FadeIn delay={0.1}>
            <h2 className="text-2xl font-semibold">Popular picks</h2>
            <p className="mt-2 text-muted-foreground">A taste of what we&apos;re baking today.</p>
          </FadeIn>
          <div className="mt-10 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {featured.map((p, i) => (
              <FadeIn key={p.id} delay={0.15 + i * 0.08}>
                <ProductCard product={p} />
              </FadeIn>
            ))}
          </div>
        </section>
      ) : null}

      <section className="border-t border-border bg-surface">
        <div className="mx-auto grid max-w-6xl gap-8 px-4 py-16 sm:grid-cols-3 sm:px-6">
          {[
            { step: "1", title: "Choose items", text: "Add pastries and meals from our live menu." },
            { step: "2", title: "Pay & upload receipt", text: "Transfer to our account and attach proof of payment." },
            { step: "3", title: "Pick up in store", text: "We'll confirm your order and have it ready for you." },
          ].map((item, i) => (
            <FadeIn key={item.step} delay={i * 0.1}>
              <div className="rounded-2xl border border-border bg-background p-6">
                <span className="text-3xl font-bold text-brand/40">{item.step}</span>
                <h3 className="mt-3 font-semibold">{item.title}</h3>
                <p className="mt-2 text-sm text-muted-foreground">{item.text}</p>
              </div>
            </FadeIn>
          ))}
        </div>
      </section>
    </div>
  );
}
