import type { ButtonHTMLAttributes, ReactNode } from 'react'
import { cn } from '../../lib/cn'

type ButtonVariant = 'primary' | 'secondary' | 'ghost'
type ButtonSize = 'md' | 'lg'

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  isLoading?: boolean
  variant?: ButtonVariant
  icon?: ReactNode
  size?: ButtonSize
}

const baseClassName =
  'inline-flex min-h-11 items-center justify-center gap-2 rounded-full border text-sm font-extrabold tracking-[-0.02em] transition duration-200 disabled:pointer-events-none disabled:opacity-60 focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-coral/20'

const sizeClassNames: Record<ButtonSize, string> = {
  md: 'px-5 py-2.5',
  lg: 'px-6 py-3 text-base',
}

const variantClassNames: Record<ButtonVariant, string> = {
  primary:
    'border-transparent bg-coral text-white shadow-sticker hover:-translate-y-0.5 hover:bg-coral-strong',
  secondary:
    'border-[color:var(--color-coral-outline)] bg-[color:var(--color-coral-soft)] text-[color:var(--color-coral-ink)] shadow-overshoot hover:-translate-y-0.5 hover:bg-white',
  ghost:
    'border-[color:var(--color-coral-outline)] bg-white text-[color:var(--color-coral-ink)] hover:bg-[color:var(--color-coral-soft)]',
}

function LoadingDots() {
  return (
    <span aria-hidden="true" className="inline-flex items-center gap-1">
      {[0, 1, 2].map((dot) => (
        <span
          className="size-1.5 animate-pulse rounded-full bg-current opacity-80"
          key={dot}
        />
      ))}
    </span>
  )
}

export function Button({
  children,
  className,
  icon,
  isLoading = false,
  size = 'md',
  type = 'button',
  variant = 'primary',
  ...props
}: ButtonProps) {
  return (
    <button
      className={cn(
        baseClassName,
        sizeClassNames[size],
        variantClassNames[variant],
        className,
      )}
      disabled={isLoading || props.disabled}
      type={type}
      {...props}
    >
      {isLoading ? <LoadingDots /> : icon}
      <span>{children}</span>
    </button>
  )
}
