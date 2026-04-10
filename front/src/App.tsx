import { Route, Routes } from 'react-router-dom'
import { ProtectedRoute } from './features/auth/ProtectedRoute'
import { CreativeStudioPage } from './pages/creative-studio/CreativeStudioPage'
import { LandingPage } from './pages/landing/LandingPage'
import { LoginPage } from './pages/login/LoginPage'

export function App() {
  return (
    <Routes>
      <Route element={<LandingPage />} path="/" />
      <Route element={<LoginPage />} path="/login" />
      <Route element={<ProtectedRoute />}>
        <Route element={<CreativeStudioPage />} path="/creative-studio" />
      </Route>
    </Routes>
  )
}
