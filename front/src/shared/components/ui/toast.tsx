import { X } from 'lucide-react'
import type { HTMLAttributes, ReactNode } from 'react'
import { cn } from '../../lib/cn'

type Tone = 'default' | 'success' | 'error' | 'warning' | 'info'

type FeedbackProps = HTMLAttributes<HTMLDivElement> & {
  action?: ReactNode
  dismissLabel?: string
  icon?: ReactNode
  onDismiss?: () => void
  title: string
  tone?: Tone
  role?: 'status' | 'alert'
}

const toneClassNames: Record<Tone, string> = {
  default: 'bg-white text-ink',
  success: 'bg-success-container text-secondary',
  error: 'bg-error-container text-error',
  warning: 'bg-warning-container text-warning',
  info: 'bg-info-container text-info',
}

function FeedbackCallout({
  action,
  children,
  className,
  dismissLabel = 'Dismiss notification',
  icon,
  onDismiss,
  role = 'status',
  title,
  tone = 'default',
  ...props
}: FeedbackProps) {
  const isToast = role === 'status'

  return (
    <div
      aria-live={role === 'alert' ? 'assertive' : 'polite'}
      className={cn(
        isToast
          ? 'flex items-start gap-4 rounded-panel bg-[color:var(--color-toast-surface)] px-4 py-4 text-white shadow-overshoot'
          : 'flex items-start justify-between gap-4 rounded-well px-4 py-4 shadow-overshoot foundation-outline-soft',
        !isToast && toneClassNames[tone],
        className,
      )}
      role={role}
      {...props}
    >
      <div className="flex min-w-0 flex-1 items-start gap-3">
        {icon ? (
          <div
            className={cn(
              'inline-flex size-10 items-center justify-center rounded-2xl',
              isToast
                ? 'bg-[color:var(--color-success-deep)] text-[color:var(--color-green-accent)]'
                : 'foundation-icon',
            )}
          >
            {icon}
          </div>
        ) : null}
        <div className="min-w-0 space-y-1">
          <p className="text-sm font-extrabold tracking-[-0.02em]">{title}</p>
          <p
            className={cn(
              'text-sm font-medium',
              isToast ? 'text-white/58' : 'opacity-80',
            )}
          >
            {children}
          </p>
        </div>
      </div>
      <div className="flex items-start gap-2">
        {action}
        {isToast ? (
          <button
            aria-label={dismissLabel}
            className="inline-flex size-9 items-center justify-center rounded-full bg-[color:var(--color-toast-close)] text-base font-black text-white/70 transition hover:text-white focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-white/15"
            onClick={onDismiss}
            type="button"
          >
            <X aria-hidden="true" className="size-4" strokeWidth={2.6} />
          </button>
        ) : null}
      </div>
    </div>
  )
}

export type ToastProps = Omit<FeedbackProps, 'role'>

export function Toast(props: ToastProps) {
  return <FeedbackCallout {...props} role="status" />
}

export function Alert(props: ToastProps) {
  return <FeedbackCallout {...props} role="alert" />
}
