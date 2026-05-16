export function formatMoney(cents: number, currency = "NGN") {
  return new Intl.NumberFormat("en-NG", {
    style: "currency",
    currency,
    minimumFractionDigits: 0,
  }).format(cents / 100);
}

export function formatDate(iso: string) {
  return new Intl.DateTimeFormat("en-NG", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(new Date(iso));
}

export function orderStatusLabel(status: string) {
  const map: Record<string, string> = {
    pending: "Pending",
    confirmed: "Confirmed",
    ready: "Ready for pickup",
    completed: "Completed",
    cancelled: "Cancelled",
  };
  return map[status] ?? status;
}
