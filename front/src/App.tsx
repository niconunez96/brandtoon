import { Route, Routes } from 'react-router-dom'
import { ProtectedRoute } from './features/auth/ProtectedRoute'
import { CreativeStudioPage } from './pages/creative-studio/CreativeStudioPage'
import { AvatarDetailsStepPage } from './pages/creative-studio/avatar/AvatarDetailsStepPage'
import { AvatarEditorLayout } from './pages/creative-studio/avatar/AvatarEditorLayout'
import { AvatarPlaceholderStepPage } from './pages/creative-studio/avatar/AvatarPlaceholderStepPage'
import { LandingPage } from './pages/landing/LandingPage'
import { LoginPage } from './pages/login/LoginPage'

export function App() {
  return (
    <Routes>
      <Route element={<LandingPage />} path="/" />
      <Route element={<LoginPage />} path="/login" />
      <Route element={<ProtectedRoute />}>
        <Route element={<CreativeStudioPage />} path="/creative-studio" />
        <Route
          element={<AvatarEditorLayout />}
          path="/creative-studio/avatars/:avatarId"
        >
          <Route element={<AvatarDetailsStepPage />} path="avatar" />
          <Route
            element={<AvatarPlaceholderStepPage stepLabel="Personality" />}
            path="personality"
          />
          <Route
            element={<AvatarPlaceholderStepPage stepLabel="Voice" />}
            path="voice"
          />
          <Route
            element={<AvatarPlaceholderStepPage stepLabel="Packs" />}
            path="packs"
          />
          <Route
            element={<AvatarPlaceholderStepPage stepLabel="Export" />}
            path="export"
          />
        </Route>
      </Route>
    </Routes>
  )
}
