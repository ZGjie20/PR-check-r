import type { InputHTMLAttributes } from 'react';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  error?: string;
}

export function Input({ error, className = '', ...props }: InputProps) {
  return (
    <div className="w-full">
      <input
        className={`input-glass ${error ? 'input-glass-error' : ''} ${className}`}
        {...props}
      />
      {error && (
        <p className="mt-2 flex items-center gap-1 text-sm text-red-300">
          <span>⚠</span> {error}
        </p>
      )}
    </div>
  );
}
