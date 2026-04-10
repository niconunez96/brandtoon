import { Clapperboard } from 'lucide-react'
import type { HTMLAttributes, ReactNode } from 'react'
import { cn } from '../../lib/cn'
import { Badge } from './badge'
import { EmptyState } from './empty-state'

type Column<T extends Record<string, ReactNode>> = {
  align?: 'left' | 'right'
  header: string
  key: keyof T
}

type RowData = Record<string, ReactNode> & { id: string }

type DataTableProps<T extends RowData> = HTMLAttributes<HTMLDivElement> & {
  columns: ReadonlyArray<Column<T>>
  emptyDescription?: string
  emptyTitle?: string
  rows: ReadonlyArray<T>
}

export function DataTable<T extends RowData>({
  className,
  columns,
  emptyDescription = 'There is no data yet for this surface.',
  emptyTitle = 'No rows available',
  rows,
  ...props
}: DataTableProps<T>) {
  return (
    <div
      className={cn(
        'overflow-hidden rounded-card bg-white shadow-overshoot foundation-outline-soft',
        className,
      )}
      {...props}
    >
      <div className="border-b border-[color:var(--color-stroke-soft)] bg-surface px-5 py-4">
        <p className="foundation-section-eyebrow">Data list surface</p>
      </div>
      {rows.length === 0 ? (
        <div className="p-5">
          <EmptyState title={emptyTitle}>{emptyDescription}</EmptyState>
        </div>
      ) : (
        <table className="min-w-full border-separate border-spacing-0 text-left">
          <thead>
            <tr>
              {columns.map((column) => (
                <th
                  className={cn(
                    'border-b border-[color:var(--color-stroke-soft)] px-5 py-4 text-[11px] font-extrabold uppercase tracking-section text-ink/45',
                    column.align === 'right' && 'text-right',
                  )}
                  key={String(column.key)}
                  scope="col"
                >
                  {column.header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {rows.map((row) => (
              <tr
                className="bg-white transition hover:bg-surface/60"
                key={row.id}
              >
                {columns.map((column) => (
                  <td
                    className={cn(
                      'border-b border-[color:var(--color-stroke-soft)] px-5 py-4 text-sm font-medium text-ink',
                      column.align === 'right' && 'text-right',
                    )}
                    key={String(column.key)}
                  >
                    {row[column.key]}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  )
}

export function AssetRow({
  label,
  meta,
  status,
}: {
  label: string
  meta: string
  status: 'Ready' | 'Processing' | 'Needs review'
}) {
  const tone =
    status === 'Ready'
      ? 'success'
      : status === 'Needs review'
        ? 'warning'
        : 'info'

  return (
    <div className="space-y-1">
      <div className="flex items-center gap-3">
        <span className="foundation-icon size-10 rounded-2xl">
          <Clapperboard aria-hidden="true" className="size-4" strokeWidth={2} />
        </span>
        <div>
          <p className="font-bold text-ink">{label}</p>
          <p className="text-xs font-medium text-ink/52">{meta}</p>
        </div>
      </div>
      <Badge tone={tone}>{status}</Badge>
    </div>
  )
}

export type { Column }
