import { forwardRef } from "react";

type InputProps = React.InputHTMLAttributes<HTMLInputElement> & {
  label?: string;
  error?: string;
};

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, className = "", id, ...props }, ref) => (
    <div className="space-y-1.5">
      {label ? (
        <label htmlFor={id} className="block text-sm font-medium text-muted-foreground">
          {label}
        </label>
      ) : null}
      <input
        ref={ref}
        id={id}
        className={`w-full rounded-xl border border-border bg-surface px-4 py-2.5 text-sm outline-none transition focus:border-brand focus:ring-2 focus:ring-brand/20 ${error ? "border-red-400" : ""} ${className}`}
        {...props}
      />
      {error ? <p className="text-xs text-red-600">{error}</p> : null}
    </div>
  ),
);
Input.displayName = "Input";
