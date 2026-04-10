import type { ButtonHTMLAttributes, ReactNode } from 'react'
import { cn } from '../../lib/cn'

type ActionChipProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  icon?: ReactNode
}

export function ActionChip({
  children,
  className,
  icon,
  type = 'button',
  ...props
}: ActionChipProps) {
  return (
    <button
      className={cn(
        'inline-flex items-center gap-2 rounded-full bg-white px-4 py-2 text-sm font-bold text-ink shadow-overshoot transition hover:-translate-y-0.5',
        className,
      )}
      type={type}
      {...props}
    >
      {icon}
      <span>{children}</span>
    </button>
  )
}
