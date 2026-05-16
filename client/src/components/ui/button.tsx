"use client";

import { motion } from "framer-motion";

type Variant = "primary" | "secondary" | "ghost" | "danger";

const styles: Record<Variant, string> = {
  primary:
    "bg-brand text-brand-foreground shadow-md shadow-brand/25 hover:shadow-lg hover:shadow-brand/30",
  secondary:
    "bg-surface border border-border text-foreground hover:bg-muted",
  ghost: "text-foreground hover:bg-muted",
  danger: "bg-red-600 text-white hover:bg-red-700",
};

type ButtonProps = {
  variant?: Variant;
  loading?: boolean;
  className?: string;
  children: React.ReactNode;
  type?: "button" | "submit" | "reset";
  disabled?: boolean;
  onClick?: () => void;
};

export function Button({
  variant = "primary",
  loading,
  className = "",
  children,
  disabled,
  type = "button",
  onClick,
}: ButtonProps) {
  return (
    <motion.div
      className="inline-block"
      whileHover={{ scale: disabled || loading ? 1 : 1.02 }}
      whileTap={{ scale: disabled || loading ? 1 : 0.98 }}
    >
      <button
        type={type}
        onClick={onClick}
        className={`inline-flex w-full items-center justify-center gap-2 rounded-xl px-5 py-2.5 text-sm font-medium transition-colors disabled:cursor-not-allowed disabled:opacity-50 ${styles[variant]} ${className}`}
        disabled={disabled || loading}
      >
        {loading ? (
          <span className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
        ) : null}
        {children}
      </button>
    </motion.div>
  );
}
