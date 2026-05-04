import type { HTMLAttributes } from 'react'
import { cn } from '../../lib/cn'

export function Card({
  children,
  className,
  ...props
}: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn(
        'rounded-card bg-surface p-6 text-ink shadow-overshoot ring-1 ring-[color:var(--color-stroke-inverse)]',
        className,
      )}
      {...props}
    >
      {children}
    </div>
  )
}

type SectionShellProps = HTMLAttributes<HTMLElement> & {
  actions?: React.ReactNode
  description: string
  eyebrow: string
  title: string
}

export function SectionShell({
  actions,
  children,
  className,
  description,
  eyebrow,
  title,
  ...props
}: SectionShellProps) {
  return (
    <section className={cn('space-y-5', className)} {...props}>
      <header className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between sm:gap-8">
        <div className="space-y-3">
          <p className="foundation-section-eyebrow">{eyebrow}</p>
          <div className="space-y-2">
            <h2 className="foundation-section-title">{title}</h2>
            <p className="foundation-body max-w-2xl">{description}</p>
          </div>
        </div>
        {actions ? <div className="flex flex-wrap gap-3">{actions}</div> : null}
      </header>
      {children}
    </section>
  )
}
