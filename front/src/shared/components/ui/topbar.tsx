import type { HTMLAttributes, ReactNode } from 'react'
import { cn } from '../../lib/cn'

type TopbarProps = HTMLAttributes<HTMLElement> & {
  actions?: ReactNode
  description?: string
  eyebrow: string
  title: string
}

export function Topbar({
  actions,
  className,
  description,
  eyebrow,
  title,
  ...props
}: TopbarProps) {
  return (
    <header
      className={cn(
        'sticky top-0 z-20 flex items-center justify-between gap-4 border-b border-[color:var(--color-stroke-soft)] bg-page/92 px-5 py-4 backdrop-blur-xl lg:px-10',
        className,
      )}
      {...props}
    >
      <div className="space-y-1">
        <p className="foundation-section-eyebrow">{eyebrow}</p>
        <p className="text-lg font-black tracking-tight text-ink">{title}</p>
        {description ? (
          <p className="foundation-body max-w-xl">{description}</p>
        ) : null}
      </div>
      {actions ? (
        <div className="flex items-center gap-3">{actions}</div>
      ) : null}
    </header>
  )
}
