import type { HTMLAttributes, ReactNode } from 'react'
import { cn } from '../../lib/cn'

type EmptyStateProps = HTMLAttributes<HTMLDivElement> & {
  action?: ReactNode
  eyebrow?: string
  icon?: ReactNode
  title: string
}

export function EmptyState({
  action,
  children,
  className,
  eyebrow = 'Empty state',
  icon,
  title,
  ...props
}: EmptyStateProps) {
  return (
    <div className={cn('foundation-empty-state', className)} {...props}>
      {icon ? <div className="foundation-icon">{icon}</div> : null}
      <div className="space-y-2">
        <p className="foundation-section-eyebrow">{eyebrow}</p>
        <p className="text-lg font-black tracking-tight text-ink">{title}</p>
        <p className="foundation-body max-w-md">{children}</p>
      </div>
      {action}
    </div>
  )
}
