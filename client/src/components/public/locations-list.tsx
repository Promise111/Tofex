"use client";

import { useMemo, useState } from "react";
import { ExternalLink, MapPin, Navigation } from "lucide-react";
import type { StoreBranch } from "@/lib/types";
import { branchMapsUrl, sortBranchesByDistance } from "@/lib/geo";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

type Props = {
  branches: StoreBranch[];
};

export function LocationsList({ branches }: Props) {
  const [sorted, setSorted] = useState(branches);
  const [locating, setLocating] = useState(false);
  const [locationError, setLocationError] = useState<string | null>(null);
  const [userCoords, setUserCoords] = useState<{ lat: number; lng: number } | null>(
    null,
  );

  const hasGeocoded = useMemo(
    () => branches.some((b) => b.latitude != null && b.longitude != null),
    [branches],
  );

  function findNearest() {
    if (!hasGeocoded) {
      setLocationError(
        "Branches need latitude and longitude for distance sorting. You can still open directions via the map links.",
      );
      return;
    }
    if (!navigator.geolocation) {
      setLocationError("Your browser does not support location access.");
      return;
    }
    setLocating(true);
    setLocationError(null);
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        const lat = pos.coords.latitude;
        const lng = pos.coords.longitude;
        setUserCoords({ lat, lng });
        setSorted(sortBranchesByDistance(branches, lat, lng));
        setLocating(false);
      },
      () => {
        setLocationError(
          "Could not get your location. Allow location access or pick a branch from the list.",
        );
        setLocating(false);
      },
      { enableHighAccuracy: true, timeout: 15000 },
    );
  }

  if (branches.length === 0) {
    return (
      <Card className="text-center text-muted-foreground">
        <p>No locations listed yet. Check back soon.</p>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center gap-3">
        <Button
          type="button"
          variant="secondary"
          className="gap-2"
          onClick={findNearest}
          loading={locating}
          disabled={!hasGeocoded}
        >
          <Navigation className="h-4 w-4" />
          Find nearest to me
        </Button>
        {!hasGeocoded ? (
          <p className="text-sm text-muted-foreground">
            Map links work for all branches; add coordinates in admin to sort by distance.
          </p>
        ) : null}
      </div>

      {locationError ? (
        <p className="text-sm text-amber-700 dark:text-amber-400">{locationError}</p>
      ) : null}

      <div className="grid gap-4 md:grid-cols-2">
        {sorted.map((branch, index) => {
          const mapsUrl = branchMapsUrl(branch);
          const showNearest =
            userCoords != null &&
            index === 0 &&
            branch.latitude != null &&
            branch.longitude != null;

          return (
            <Card key={branch.id} className={showNearest ? "ring-2 ring-brand/40" : ""}>
              <div className="flex items-start justify-between gap-3">
                <div>
                  <div className="flex items-center gap-2">
                    <MapPin className="h-4 w-4 shrink-0 text-brand" />
                    <h2 className="font-semibold">{branch.name}</h2>
                  </div>
                  {showNearest ? (
                    <span className="mt-1 inline-block rounded-full bg-brand/15 px-2 py-0.5 text-xs font-medium text-brand">
                      Nearest to you
                    </span>
                  ) : null}
                </div>
              </div>
              <p className="mt-3 text-sm text-muted-foreground">{branch.address}</p>
              {branch.city ? (
                <p className="text-sm text-muted-foreground">{branch.city}</p>
              ) : null}
              {branch.phone ? (
                <p className="mt-2 text-sm">
                  <a href={`tel:${branch.phone}`} className="hover:text-brand">
                    {branch.phone}
                  </a>
                </p>
              ) : null}
              {branch.hours ? (
                <p className="mt-1 text-sm text-muted-foreground">Hours: {branch.hours}</p>
              ) : null}
              <a
                href={mapsUrl}
                target="_blank"
                rel="noopener noreferrer"
                className="mt-4 inline-flex items-center gap-2 text-sm font-medium text-brand hover:underline"
              >
                Open in Google Maps
                <ExternalLink className="h-3.5 w-3.5" />
              </a>
            </Card>
          );
        })}
      </div>
    </div>
  );
}