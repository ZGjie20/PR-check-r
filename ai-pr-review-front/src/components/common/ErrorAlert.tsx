interface ErrorAlertProps {
  message: string;
  onDismiss?: () => void;
}

export function ErrorAlert({ message, onDismiss }: ErrorAlertProps) {
  return (
    <div className="flex items-start justify-between gap-3 rounded-xl border border-red-400/30 bg-red-500/10 px-4 py-3 text-sm text-red-200 backdrop-blur-sm">
      <span className="flex items-start gap-2">
        <span className="mt-0.5 shrink-0">⚠️</span>
        {message}
      </span>
      {onDismiss && (
        <button
          type="button"
          onClick={onDismiss}
          className="shrink-0 text-red-300 transition-colors hover:text-red-100"
          aria-label="关闭"
        >
          ×
        </button>
      )}
    </div>
  );
}
