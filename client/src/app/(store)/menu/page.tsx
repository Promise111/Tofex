import { FadeIn, Stagger, StaggerItem } from "@/components/motion/fade-in";
import { ProductCard } from "@/components/public/product-card";
import { publicApi } from "@/lib/api";

export const metadata = { title: "Menu · Tofex" };

export default async function MenuPage() {
  let products: Awaited<ReturnType<typeof publicApi.getProducts>> = [];
  let error: string | null = null;
  try {
    products = await publicApi.getProducts();
  } catch (e) {
    error = e instanceof Error ? e.message : "Could not load menu";
  }

  return (
    <div className="mx-auto max-w-6xl px-4 py-12 sm:px-6">
      <FadeIn>
        <h1 className="text-3xl font-semibold tracking-tight">Our menu</h1>
        <p className="mt-2 text-muted-foreground">
          Everything is made fresh. Prices shown in NGN.
        </p>
      </FadeIn>

      {error ? (
        <p className="mt-8 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-800">
          {error}. Is the API running at {process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080"}?
        </p>
      ) : null}

      {products.length === 0 && !error ? (
        <p className="mt-12 text-center text-muted-foreground">No items available yet. Check back soon.</p>
      ) : (
        <Stagger className="mt-10 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {products.map((p) => (
            <StaggerItem key={p.id}>
              <ProductCard product={p} />
            </StaggerItem>
          ))}
        </Stagger>
      )}
    </div>
  );
}
