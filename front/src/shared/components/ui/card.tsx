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
  eyebrow: string
  title: string
  description: string
}

export function SectionShell({
  children,
  className,
  description,
  eyebrow,
  title,
  ...props
}: SectionShellProps) {
  return (
    <section className={cn('space-y-5', className)} {...props}>
      <header className="space-y-3">
        <p className="foundation-section-eyebrow">{eyebrow}</p>
        <div className="space-y-2">
          <h2 className="foundation-section-title">{title}</h2>
          <p className="foundation-body max-w-2xl">{description}</p>
        </div>
      </header>
      {children}
    </section>
  )
}
