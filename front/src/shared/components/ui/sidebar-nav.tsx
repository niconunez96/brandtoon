import { Dot } from 'lucide-react'
import type { HTMLAttributes, ReactNode } from 'react'
import { cn } from '../../lib/cn'
import { Card } from './card'

type SidebarNavItem = {
  active?: boolean
  badge?: ReactNode
  href: string
  label: string
  shortLabel?: string
}

type SidebarNavProps = HTMLAttributes<HTMLElement> & {
  footer?: ReactNode
  items: SidebarNavItem[]
  title: string
}

export function SidebarNav({
  className,
  footer,
  items,
  title,
  ...props
}: SidebarNavProps) {
  return (
    <aside
      className={cn(
        'hidden w-72 flex-col border-r border-[color:var(--color-stroke-soft)] bg-white/70 px-6 py-8 backdrop-blur-xl lg:flex',
        className,
      )}
      {...props}
    >
      <div className="space-y-2">
        <p className="foundation-section-eyebrow">Brandtoon shell</p>
        <h1 className="text-2xl font-black tracking-tight">{title}</h1>
      </div>
      <nav aria-label="Primary sidebar" className="mt-10 space-y-2">
        {items.map((item) => (
          <a
            className={cn(
              'flex items-center justify-between rounded-2xl px-4 py-3 text-sm font-bold transition',
              item.active
                ? 'bg-surface text-ink shadow-overshoot'
                : 'text-ink/55 hover:bg-white',
            )}
            href={item.href}
            key={item.label}
          >
            <span>{item.label}</span>
            {item.badge ?? (
              <Dot aria-hidden="true" className="size-4" strokeWidth={2.6} />
            )}
          </a>
        ))}
      </nav>
      {footer ? (
        <Card className="mt-auto rounded-hero bg-gradient-to-br from-coral via-coral-strong to-pink-sample text-white shadow-sticker">
          {footer}
        </Card>
      ) : null}
    </aside>
  )
}

export type { SidebarNavItem }
