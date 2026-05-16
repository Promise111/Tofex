export function Card({
  children,
  className = "",
}: {
  children: React.ReactNode;
  className?: string;
}) {
  return (
    <div
      className={`rounded-2xl border border-border bg-surface p-6 shadow-sm shadow-black/5 ${className}`}
    >
      {children}
    </div>
  );
}
