export function SiteFooter() {
  return (
    <footer className="mt-auto border-t border-border bg-surface">
      <div className="mx-auto flex max-w-6xl flex-col gap-2 px-4 py-10 text-sm text-muted-foreground sm:flex-row sm:items-center sm:justify-between sm:px-6">
        <p>© {new Date().getFullYear()} Tofex. Fresh from our kitchen to you.</p>
        <p>
          <a href="/locations" className="hover:text-foreground">
            Find a branch
          </a>
          {" · "}
          Order online · Pay by transfer · Pick up in store
        </p>
      </div>
    </footer>
  );
}
