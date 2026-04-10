import {
  type CSSProperties,
  type ChangeEvent,
  type InputHTMLAttributes,
  useState,
} from 'react'

type SliderMarker = {
  label: string
  value: string
}

type SliderProps = Omit<InputHTMLAttributes<HTMLInputElement>, 'type'> & {
  label: string
  markers?: SliderMarker[]
  valueLabel: string
}

export function Slider({
  defaultValue,
  label,
  markers,
  max = 100,
  min = 0,
  onChange,
  value,
  valueLabel,
  ...props
}: SliderProps) {
  const [internalValue, setInternalValue] = useState(
    Number(defaultValue ?? min),
  )
  const currentValue = Number(value ?? internalValue)
  const range = Number(max) - Number(min)
  const fill = `${range === 0 ? 0 : ((currentValue - Number(min)) / range) * 100}%`

  const handleChange = (event: ChangeEvent<HTMLInputElement>) => {
    if (value === undefined) {
      setInternalValue(Number(event.target.value))
    }

    onChange?.(event)
  }

  return (
    <label className="flex flex-col gap-4">
      <span className="flex items-center justify-between gap-3">
        <span className="text-[11px] font-extrabold uppercase tracking-section text-ink/48">
          {label}
        </span>
        <span className="rounded-full bg-coral px-4 py-1.5 text-[11px] font-black tracking-[0.18em] text-white shadow-sticker">
          {valueLabel}
        </span>
      </span>
      <input
        aria-valuemax={Number(max)}
        aria-valuemin={Number(min)}
        aria-valuenow={currentValue}
        className="brandtoon-slider"
        defaultValue={defaultValue}
        max={max}
        min={min}
        onChange={handleChange}
        style={{ '--slider-fill': fill } as CSSProperties}
        type="range"
        value={value}
        {...props}
      />
      {markers?.length ? (
        <span className="grid grid-cols-3 text-xs font-bold text-ink/48">
          {markers.map((marker) => (
            <span
              className="flex flex-col gap-1 first:items-start last:items-end"
              key={marker.label}
            >
              <span>{marker.label}</span>
              <span className="text-[11px] font-medium text-ink/34">
                {marker.value}
              </span>
            </span>
          ))}
        </span>
      ) : null}
    </label>
  )
}
