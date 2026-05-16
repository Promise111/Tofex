import { FadeIn } from "@/components/motion/fade-in";
import { LocationsList } from "@/components/public/locations-list";
import { publicApi } from "@/lib/api";

export default async function LocationsPage() {
  let branches: Awaited<ReturnType<typeof publicApi.getBranches>> = [];
  try {
    branches = await publicApi.getBranches();
  } catch {
    branches = [];
  }

  return (
    <div className="mx-auto max-w-6xl px-4 py-12 sm:px-6">
      <FadeIn>
        <h1 className="text-3xl font-semibold tracking-tight">Our locations</h1>
        <p className="mt-3 max-w-2xl text-muted-foreground">
          Find a Tofex branch near you for pickup. Use the map link on each location for directions.
        </p>
      </FadeIn>
      <FadeIn delay={0.1} className="mt-10">
        <LocationsList branches={branches} />
      </FadeIn>
    </div>
  );
}
