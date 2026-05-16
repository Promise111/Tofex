import type { StoreBranch } from "./types";

export function distanceKm(
  lat1: number,
  lon1: number,
  lat2: number,
  lon2: number,
): number {
  const R = 6371;
  const toRad = (d: number) => (d * Math.PI) / 180;
  const dLat = toRad(lat2 - lat1);
  const dLon = toRad(lon2 - lon1);
  const a =
    Math.sin(dLat / 2) ** 2 +
    Math.cos(toRad(lat1)) * Math.cos(toRad(lat2)) * Math.sin(dLon / 2) ** 2;
  return R * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
}

export function branchMapsUrl(branch: StoreBranch): string {
  if (branch.maps_url?.trim()) {
    return branch.maps_url.trim();
  }
  if (branch.latitude != null && branch.longitude != null) {
    return `https://www.google.com/maps?q=${branch.latitude},${branch.longitude}`;
  }
  const query = [branch.address, branch.city].filter(Boolean).join(", ");
  return `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(query)}`;
}

export function sortBranchesByDistance(
  branches: StoreBranch[],
  lat: number,
  lng: number,
): StoreBranch[] {
  const withCoords = branches.filter(
    (b) => b.latitude != null && b.longitude != null,
  );
  const withoutCoords = branches.filter(
    (b) => b.latitude == null || b.longitude == null,
  );
  withCoords.sort((a, b) => {
    const da = distanceKm(lat, lng, a.latitude!, a.longitude!);
    const db = distanceKm(lat, lng, b.latitude!, b.longitude!);
    return da - db;
  });
  return [...withCoords, ...withoutCoords];
}
