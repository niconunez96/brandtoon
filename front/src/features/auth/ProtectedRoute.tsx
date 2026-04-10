import { Navigate, Outlet, useLocation } from 'react-router-dom'
import { useCurrentUserQuery } from '../../queries/useCurrentUserQuery'
import { isUnauthorizedError } from '../../services/auth.api'
import { Button } from '../../shared/components/ui/button'
import { Card } from '../../shared/components/ui/card'
import { reloadBrowserWindow } from '../../shared/lib/browser'

function buildLoginUrl(pathname: string, search: string) {
  const searchParams = new URLSearchParams({ next: `${pathname}${search}` })
  return `/login?${searchParams.toString()}`
}

export function ProtectedRoute() {
  const location = useLocation()
  const sessionQuery = useCurrentUserQuery()

  if (sessionQuery.isLoading) {
    return (
      <div className="foundation-page flex items-center justify-center px-4 py-10">
        <Card className="w-full max-w-md space-y-4 bg-white text-center">
          <p className="foundation-section-eyebrow">Loading session</p>
          <p className="text-2xl font-black tracking-tight text-ink">
            Preparing your studio…
          </p>
        </Card>
      </div>
    )
  }

  if (isUnauthorizedError(sessionQuery.error)) {
    return (
      <Navigate
        replace
        to={buildLoginUrl(location.pathname, location.search)}
      />
    )
  }

  if (sessionQuery.isError) {
    return (
      <div className="foundation-page flex items-center justify-center px-4 py-10">
        <Card className="w-full max-w-md space-y-4 bg-white text-center">
          <p className="foundation-section-eyebrow">Session error</p>
          <p className="text-2xl font-black tracking-tight text-ink">
            We could not verify your session.
          </p>
          <Button onClick={reloadBrowserWindow} variant="secondary">
            Try again
          </Button>
        </Card>
      </div>
    )
  }

  return <Outlet />
}
