interface EmptyStateProps {
  title: string;
  description?: string;
  icon?: string;
}

export function EmptyState({ title, description, icon = '📭' }: EmptyStateProps) {
  return (
    <div className="glass-card flex flex-col items-center justify-center py-20 text-center">
      <span className="mb-4 text-5xl opacity-60">{icon}</span>
      <p className="text-lg font-medium gradient-text">{title}</p>
      {description && (
        <p className="mt-2 max-w-sm text-sm text-white/40">{description}</p>
      )}
    </div>
  );
}
