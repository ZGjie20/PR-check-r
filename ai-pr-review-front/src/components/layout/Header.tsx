import { NavLink } from 'react-router-dom';
import { useAppStore } from '@/store/app.store';

const NAV_ITEMS = [
  { to: '/review/new', label: '新建审查', icon: '🔮' },
  { to: '/reviews', label: '历史记录', icon: '📜' },
];

export function Header() {
  const apiHealthy = useAppStore((s) => s.apiHealthy);

  const healthLabel =
    apiHealthy === null ? 'CONNECTING...' : apiHealthy ? 'ONLINE' : 'OFFLINE';
  const healthDot =
    apiHealthy === null
      ? 'bg-yellow-400 animate-pulse'
      : apiHealthy
        ? 'bg-emerald-400 animate-pulse-glow'
        : 'bg-red-400';

  return (
    <header className="sticky top-0 z-50 border-b border-white/10 bg-white/[0.06] backdrop-blur-2xl">
      <div className="mx-auto flex max-w-6xl items-center justify-between px-4 py-3">
        <div className="flex items-center gap-6">
          <NavLink to="/" className="group flex items-center gap-2.5">
            <span className="flex h-9 w-9 items-center justify-center rounded-xl bg-gradient-to-br from-sakura-400 to-dream-500 text-sm shadow-lg shadow-dream-500/30 transition-transform group-hover:scale-105">
              ✦
            </span>
            <div>
              <span className="text-base font-bold gradient-text-warm">审查之灵</span>
              <span className="block text-[10px] tracking-widest text-white/40">
                make review eazy
              </span>
            </div>
          </NavLink>

          <nav className="hidden items-center gap-1 sm:flex">
            {NAV_ITEMS.map((item) => (
              <NavLink
                key={item.to}
                to={item.to}
                className={({ isActive }) =>
                  `nav-link flex items-center gap-1.5 ${isActive ? 'nav-link-active' : ''}`
                }
              >
                <span className="text-xs">{item.icon}</span>
                {item.label}
              </NavLink>
            ))}
          </nav>
        </div>

        <div className="flex items-center gap-3">
          <div className="glass flex items-center gap-2 rounded-full px-3 py-1.5 text-xs font-medium tracking-wider text-white/60">
            <span className={`h-2 w-2 rounded-full ${healthDot}`} />
            {healthLabel}
          </div>
        </div>
      </div>

      {/* 移动端导航 */}
      <nav className="flex gap-1 border-t border-white/5 px-4 py-2 sm:hidden">
        {NAV_ITEMS.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            className={({ isActive }) =>
              `nav-link flex flex-1 items-center justify-center gap-1 text-xs ${isActive ? 'nav-link-active' : ''}`
            }
          >
            <span>{item.icon}</span>
            {item.label}
          </NavLink>
        ))}
      </nav>
    </header>
  );
}
