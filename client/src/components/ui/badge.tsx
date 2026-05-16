const variants: Record<string, string> = {
  pending: "bg-amber-100 text-amber-800",
  confirmed: "bg-blue-100 text-blue-800",
  ready: "bg-emerald-100 text-emerald-800",
  completed: "bg-stone-200 text-stone-700",
  cancelled: "bg-red-100 text-red-800",
  default: "bg-muted text-muted-foreground",
};

export function Badge({ status, children }: { status?: string; children: React.ReactNode }) {
  const key = status && variants[status] ? status : "default";
  return (
    <span
      className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium capitalize ${variants[key]}`}
    >
      {children}
    </span>
  );
}
