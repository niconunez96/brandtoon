import {
  type InputHTMLAttributes,
  type PropsWithChildren,
  type ReactNode,
  type TextareaHTMLAttributes,
  useId,
  useState,
} from 'react'
import { cn } from '../../lib/cn'

type FieldShellProps = {
  describedBy?: string
  hint?: string
  id: string
  label: string
  message?: string
  state?: 'default' | 'error' | 'success'
}

type FieldProps = {
  hint?: string
  label: string
  message?: string
  state?: 'default' | 'error' | 'success'
}

const fieldClassName =
  'w-full rounded-2xl border bg-white/80 px-4 py-3 text-sm font-medium text-ink outline-none transition placeholder:text-ink/30 focus:ring-4'

const fieldStateClassNames = {
  default:
    'border-[color:var(--color-stroke-soft)] focus:border-coral/50 focus:ring-coral/12',
  error:
    'border-error/30 bg-error-container/45 focus:border-error/50 focus:ring-error/10',
  success:
    'border-secondary/25 bg-white focus:border-secondary/40 focus:ring-secondary/10',
}

const messageStateClassNames = {
  default: 'text-ink/45',
  error: 'text-error',
  success: 'text-secondary',
}

function FieldShell({
  children,
  describedBy,
  hint,
  id,
  label,
  message,
  state = 'default',
}: PropsWithChildren<FieldShellProps>) {
  return (
    <div className="flex flex-col gap-2.5">
      <label
        className="text-[11px] font-extrabold uppercase tracking-section text-ink/48"
        htmlFor={id}
      >
        {label}
      </label>
      {children}
      {hint ? (
        <span className="text-xs font-medium text-ink/45" id={describedBy}>
          {hint}
        </span>
      ) : null}
      {message ? (
        <span
          className={cn('text-xs font-bold', messageStateClassNames[state])}
        >
          {message}
        </span>
      ) : null}
    </div>
  )
}

type InputProps = InputHTMLAttributes<HTMLInputElement> & FieldProps

export function Input({
  className,
  hint,
  label,
  message,
  state = 'default',
  ...props
}: InputProps) {
  const id = useId()
  const describedBy = hint ? `${id}-hint` : undefined

  return (
    <FieldShell
      describedBy={describedBy}
      hint={hint}
      id={id}
      label={label}
      message={message}
      state={state}
    >
      <input
        aria-describedby={describedBy}
        aria-invalid={state === 'error' || undefined}
        className={cn(fieldClassName, fieldStateClassNames[state], className)}
        id={id}
        {...props}
      />
    </FieldShell>
  )
}

type TextareaProps = TextareaHTMLAttributes<HTMLTextAreaElement> & FieldProps

export function Textarea({
  className,
  hint,
  label,
  message,
  state = 'default',
  ...props
}: TextareaProps) {
  const id = useId()
  const describedBy = hint ? `${id}-hint` : undefined

  return (
    <FieldShell
      describedBy={describedBy}
      hint={hint}
      id={id}
      label={label}
      message={message}
      state={state}
    >
      <textarea
        aria-describedby={describedBy}
        aria-invalid={state === 'error' || undefined}
        className={cn(
          fieldClassName,
          fieldStateClassNames[state],
          'min-h-32 resize-none',
          className,
        )}
        id={id}
        {...props}
      />
    </FieldShell>
  )
}

type PromptFieldProps = Omit<
  TextareaHTMLAttributes<HTMLTextAreaElement>,
  'children'
> & {
  icon?: ReactNode
  title: string
}

export function PromptField({
  className,
  defaultValue,
  icon,
  maxLength = 500,
  onChange,
  placeholder,
  title,
  value,
  ...props
}: PromptFieldProps) {
  const id = useId()
  const [internalValue, setInternalValue] = useState(String(defaultValue ?? ''))
  const currentValue = String(value ?? internalValue)

  return (
    <div className="space-y-3 rounded-panel bg-white p-4 shadow-overshoot ring-1 ring-[color:var(--color-stroke-soft)] sm:p-5">
      <div className="flex items-center gap-3">
        <span className="inline-flex size-10 items-center justify-center rounded-2xl bg-surface text-ink shadow-overshoot ring-1 ring-[color:var(--color-stroke-soft)]">
          {icon}
        </span>
        <label
          className="text-sm font-extrabold tracking-[-0.02em] text-ink"
          htmlFor={id}
        >
          {title}
        </label>
      </div>
      <div className="space-y-2">
        <textarea
          className={cn(
            'min-h-44 w-full resize-none rounded-[1.75rem] border border-[color:var(--color-stroke-strong)] bg-white px-5 py-4 text-sm font-medium text-ink outline-none transition placeholder:text-ink/34 focus:border-coral/45 focus:ring-4 focus:ring-coral/10',
            className,
          )}
          id={id}
          maxLength={maxLength}
          onChange={(event) => {
            if (value === undefined) {
              setInternalValue(event.target.value)
            }

            onChange?.(event)
          }}
          placeholder={placeholder}
          value={currentValue}
          {...props}
        />
        <div className="flex justify-end text-xs font-medium text-ink/42">
          {currentValue.length} / {maxLength}
        </div>
      </div>
    </div>
  )
}
