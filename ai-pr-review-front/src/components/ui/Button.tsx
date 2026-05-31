import type { ButtonHTMLAttributes } from 'react';

type ButtonVariant = 'primary' | 'secondary' | 'ghost';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant;
  loading?: boolean;
}

const variantClasses: Record<ButtonVariant, string> = {
  primary:
    'btn-glow bg-gradient-to-r from-sakura-500 to-dream-500 text-white shadow-lg shadow-dream-500/25 hover:shadow-dream-500/40 hover:brightness-110 disabled:from-sakura-500/50 disabled:to-dream-500/50',
  secondary:
    'glass text-white/80 hover:bg-white/15 hover:text-white disabled:opacity-40',
  ghost: 'text-white/60 hover:bg-white/10 hover:text-white disabled:opacity-40',
};

export function Button({
  variant = 'primary',
  loading = false,
  disabled,
  className = '',
  children,
  ...props
}: ButtonProps) {
  return (
    <button
      className={`inline-flex items-center justify-center rounded-xl px-5 py-2.5 text-sm font-medium transition-all duration-300 disabled:cursor-not-allowed ${variantClasses[variant]} ${className}`}
      disabled={disabled || loading}
      {...props}
    >
      {loading ? (
        <span className="flex items-center gap-2">
          <span className="h-4 w-4 animate-spin rounded-full border-2 border-white/30 border-t-white" />
          处理中...
        </span>
      ) : (
        children
      )}
    </button>
  );
}
