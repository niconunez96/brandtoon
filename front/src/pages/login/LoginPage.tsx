import { Navigate, useSearchParams } from 'react-router-dom'
import { useCurrentUserQuery } from '../../queries/useCurrentUserQuery'
import { buildGoogleLoginUrl } from '../../services/auth.api'
import { Badge } from '../../shared/components/ui/badge'
import { Button } from '../../shared/components/ui/button'
import { Card } from '../../shared/components/ui/card'
import { sanitizeNextPath } from '../../shared/config/api'
import { navigateToExternalUrl } from '../../shared/lib/browser'

const LOGIN_ERRORS: Record<string, string> = {
  google_auth_failed:
    'Google authentication failed. Please try again from the same browser.',
  oauth_state_error:
    'We could not prepare a secure Google login request. Please try again.',
}

export function LoginPage() {
  const [searchParams] = useSearchParams()
  const sessionQuery = useCurrentUserQuery()
  const nextPath = sanitizeNextPath(searchParams.get('next'))
  const errorCode = searchParams.get('error') ?? ''
  const errorMessage = LOGIN_ERRORS[errorCode]

  if (sessionQuery.data) {
    return <Navigate replace to={nextPath} />
  }

  return (
    <div className="foundation-page flex min-h-screen items-center justify-center px-4 py-10 sm:px-6">
      <Card className="w-full max-w-xl space-y-8 bg-white p-8 sm:p-10">
        <div className="space-y-4 text-center">
          <Badge className="bg-ink px-4 py-2 uppercase tracking-section text-white">
            Sign in to continue
          </Badge>
          <div className="space-y-3">
            <h1 className="text-4xl font-black tracking-tight text-ink sm:text-5xl">
              Log in with Google
            </h1>
            <p className="foundation-body mx-auto max-w-lg text-base sm:text-lg">
              Continue with your Google account and we will send you directly to
              the Creative Studio.
            </p>
          </div>
        </div>

        {errorMessage ? (
          <div className="rounded-3xl border border-error/12 bg-error-container px-5 py-4 text-sm font-semibold text-error">
            {errorMessage}
          </div>
        ) : null}

        <div className="space-y-4">
          <Button
            className="w-full justify-center"
            onClick={() => navigateToExternalUrl(buildGoogleLoginUrl(nextPath))}
            size="lg"
          >
            Continue with Google
          </Button>
          <p className="text-center text-sm font-medium text-ink/55">
            By continuing, you will start your authenticated Brandtoon session.
          </p>
        </div>
      </Card>
    </div>
  )
}
