interface PageHeaderProps {
  title: string;
  subtitle?: string;
  icon?: string;
}

export function PageHeader({ title, subtitle, icon = '✦' }: PageHeaderProps) {
  return (
    <div className="mb-8">
      <div className="flex items-center gap-3">
        <span className="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-sakura-400/30 to-dream-500/30 text-lg backdrop-blur-sm">
          {icon}
        </span>
        <div>
          <h1 className="text-2xl font-bold tracking-tight gradient-text">{title}</h1>
          {subtitle && (
            <p className="mt-1 text-sm text-white/50">{subtitle}</p>
          )}
        </div>
      </div>
      <div className="mt-4 h-px bg-gradient-to-r from-transparent via-white/20 to-transparent" />
    </div>
  );
}
