import { X } from 'lucide-react'
import type { HTMLAttributes } from 'react'
import { cn } from '../../lib/cn'

type BadgeTone = 'default' | 'success' | 'muted' | 'error' | 'info' | 'warning'
type BadgeVariant = 'solid' | 'outline' | 'dark' | 'removable'

const toneClassNames: Record<BadgeTone, string> = {
  default: 'bg-white text-ink ring-1 ring-[color:var(--color-stroke-soft)]',
  success: 'bg-success-container text-secondary',
  muted: 'bg-surface-high text-ink/70',
  error: 'bg-error-container text-error',
  info: 'bg-info-container text-info',
  warning: 'bg-warning-container text-warning',
}

type BadgeProps = HTMLAttributes<HTMLSpanElement> & {
  onRemove?: () => void
  tone?: BadgeTone
  variant?: BadgeVariant
}

export function Badge({
  children,
  className,
  onRemove,
  tone = 'default',
  variant = 'solid',
  ...props
}: BadgeProps) {
  const variantClassNames: Record<BadgeVariant, string> = {
    solid: toneClassNames[tone],
    outline:
      'bg-white text-ink/72 ring-1 ring-[color:var(--color-stroke-strong)]',
    dark: 'bg-ink text-white',
    removable:
      'bg-[color:var(--color-coral-soft)] text-[color:var(--color-coral-ink)]',
  }

  return (
    <span
      className={cn(
        'inline-flex items-center gap-2 rounded-full px-3.5 py-1.5 text-[11px] font-bold tracking-[-0.01em]',
        variantClassNames[variant],
        className,
      )}
      {...props}
    >
      {children}
      {variant === 'removable' ? (
        <button
          aria-label="Remove tag"
          className="inline-flex size-4 items-center justify-center rounded-full bg-white/70 text-[10px] font-black text-[color:var(--color-coral-ink)]"
          onClick={onRemove}
          type="button"
        >
          <X aria-hidden="true" className="size-3" strokeWidth={2.8} />
        </button>
      ) : null}
    </span>
  )
}
