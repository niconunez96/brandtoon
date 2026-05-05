import type { HTMLAttributes, ReactNode } from 'react'
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
} from 'react'
import { createPortal } from 'react-dom'
import { cn } from '../../lib/cn'

type ModalContextValue = {
  isOpen: boolean
  close: () => void
}

const ModalContext = createContext<ModalContextValue | null>(null)

function useModalContext() {
  const context = useContext(ModalContext)
  if (!context) {
    throw new Error('Modal components must be used within a Modal provider')
  }
  return context
}

type ModalProps = {
  children: ReactNode
  isOpen: boolean
  onClose: () => void
}

export function Modal({ children, isOpen, onClose }: ModalProps) {
  const [isMounted, setIsMounted] = useState(false)

  useEffect(() => {
    setIsMounted(true)
  }, [])

  const handleClose = useCallback(() => {
    onClose()
  }, [onClose])

  useEffect(() => {
    if (!isOpen) return

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        handleClose()
      }
    }

    document.addEventListener('keydown', handleKeyDown)
    return () => document.removeEventListener('keydown', handleKeyDown)
  }, [isOpen, handleClose])

  useEffect(() => {
    if (!isOpen) return

    // Lock body scroll
    const original = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    return () => {
      document.body.style.overflow = original
    }
  }, [isOpen])

  if (!isMounted || !isOpen) return null

  return createPortal(
    <ModalContext.Provider value={{ isOpen, close: handleClose }}>
      <div className="fixed inset-0 z-50 flex items-center justify-center">
        <Backdrop onClick={handleClose} />
        <Panel>{children}</Panel>
      </div>
    </ModalContext.Provider>,
    document.body,
  )
}

type BackdropProps = HTMLAttributes<HTMLDivElement>

function Backdrop({ className, ...props }: BackdropProps) {
  return (
    <div
      className={cn(
        'animate-fade-in bg-ink/48 fixed inset-0 backdrop-blur-sm',
        className,
      )}
      {...props}
    />
  )
}

type PanelProps = HTMLAttributes<HTMLDivElement>

function Panel({ className, ...props }: PanelProps) {
  return (
    <div
      className={cn(
        'animate-scale-in relative z-10 w-full max-w-md p-6',
        className,
      )}
      {...props}
    />
  )
}

type ModalCloseProps = HTMLAttributes<HTMLButtonElement>

export function ModalClose({ className, ...props }: ModalCloseProps) {
  const { close } = useModalContext()

  return (
    <button
      className={cn(
        'absolute right-4 top-4 flex size-8 items-center justify-center rounded-xl text-ink-soft transition-colors hover:bg-surface-high focus:outline-none focus:ring-2 focus:ring-coral focus:ring-offset-2',
        className,
      )}
      onClick={close}
      type="button"
      aria-label="Close modal"
      {...props}
    >
      <svg
        aria-hidden="true"
        className="size-5"
        fill="none"
        stroke="currentColor"
        strokeWidth={2}
        viewBox="0 0 24 24"
      >
        <title>Close</title>
        <path d="M6 18L18 6M6 6l12 12" />
      </svg>
    </button>
  )
}
