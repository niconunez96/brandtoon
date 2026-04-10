import type { HTMLAttributes } from 'react'
import { cn } from '../../lib/cn'
import { Card } from './card'

type MetricProgressCardProps = HTMLAttributes<HTMLDivElement> & {
  helper: string
  label: string
  progress: number
  trend?: string
  value: string
}

export function MetricProgressCard({
  className,
  helper,
  label,
  progress,
  trend,
  value,
  ...props
}: MetricProgressCardProps) {
  const clampedProgress = Math.max(0, Math.min(100, progress))

  return (
    <Card className={cn('bg-white', className)} {...props}>
      <div className="flex items-start justify-between gap-3">
        <div className="space-y-2">
          <p className="foundation-section-eyebrow">Metric card</p>
          <p className="text-base font-black tracking-tight text-ink">
            {label}
          </p>
        </div>
        {trend ? (
          <p className="text-sm font-bold text-secondary">{trend}</p>
        ) : null}
      </div>
      <div className="mt-6 flex items-end justify-between gap-3">
        <p className="text-4xl font-black tracking-[-0.05em] text-ink">
          {value}
        </p>
        <p className="max-w-44 text-right text-sm font-medium text-ink/52">
          {helper}
        </p>
      </div>
      <div
        aria-label={label}
        aria-valuemax={100}
        aria-valuemin={0}
        aria-valuenow={clampedProgress}
        className="foundation-progress-track mt-5"
        role="progressbar"
        tabIndex={0}
      >
        <div
          className="foundation-progress-bar"
          style={{ width: `${clampedProgress}%` }}
        />
      </div>
    </Card>
  )
}
