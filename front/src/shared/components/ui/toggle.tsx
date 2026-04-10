import { useId, useState } from 'react'
import { cn } from '../../lib/cn'

type ToggleProps = {
  checked?: boolean
  defaultChecked?: boolean
  disabled?: boolean
  description: string
  label: string
  onCheckedChange?: (checked: boolean) => void
}

export function Toggle({
  checked,
  defaultChecked = false,
  description,
  disabled = false,
  label,
  onCheckedChange,
}: ToggleProps) {
  const [internalChecked, setInternalChecked] = useState(defaultChecked)
  const descriptionId = useId()
  const labelId = useId()
  const isChecked = checked ?? internalChecked

  const handleClick = () => {
    if (disabled) {
      return
    }

    const nextValue = !isChecked
    if (checked === undefined) {
      setInternalChecked(nextValue)
    }
    onCheckedChange?.(nextValue)
  }

  return (
    <button
      aria-checked={isChecked}
      aria-describedby={descriptionId}
      aria-labelledby={labelId}
      className={cn(
        'flex w-full items-center justify-between gap-4 rounded-panel bg-white px-4 py-4 text-left shadow-overshoot transition ring-1 ring-[color:var(--color-stroke-soft)]',
        !disabled &&
          'hover:-translate-y-0.5 focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-coral/12',
        disabled && 'cursor-not-allowed opacity-60',
      )}
      disabled={disabled}
      onClick={handleClick}
      role="switch"
      type="button"
    >
      <div className="space-y-1.5">
        <p
          className="text-sm font-extrabold tracking-[-0.02em] text-ink"
          id={labelId}
        >
          {label}
        </p>
        <p
          className="text-sm font-medium leading-6 text-ink/48"
          id={descriptionId}
        >
          {description}
        </p>
      </div>
      <div
        aria-hidden="true"
        className={cn(
          'flex h-8 w-14 items-center rounded-full p-1 transition',
          isChecked ? 'justify-end bg-coral' : 'justify-start bg-surface-high',
        )}
      >
        <span className="size-6 rounded-full bg-white shadow-md ring-1 ring-black/5" />
      </div>
    </button>
  )
}
