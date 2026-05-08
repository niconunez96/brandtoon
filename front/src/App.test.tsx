import { render, screen } from '@testing-library/react'
import { MemoryRouter, Outlet } from 'react-router-dom'
import { describe, expect, it, vi } from 'vitest'
import { App } from './App'

vi.mock('./features/auth/ProtectedRoute', () => ({
  ProtectedRoute: () => <Outlet />,
}))

vi.mock('./pages/creative-studio/CreativeStudioPage', () => ({
  CreativeStudioPage: () => <div>Creative Studio page</div>,
}))

vi.mock('./pages/creative-studio/avatar/AvatarDetailsStepPage', () => ({
  AvatarDetailsStepPage: () => <div>Avatar details step</div>,
}))

vi.mock('./pages/creative-studio/avatar/AvatarEditorLayout', () => ({
  AvatarEditorLayout: () => (
    <div>
      <span>Avatar editor layout</span>
      <Outlet />
    </div>
  ),
}))

vi.mock('./pages/creative-studio/avatar/AvatarPlaceholderStepPage', () => ({
  AvatarPlaceholderStepPage: ({ stepLabel }: { stepLabel: string }) => (
    <div>{stepLabel} placeholder</div>
  ),
}))

vi.mock('./pages/landing/LandingPage', () => ({
  LandingPage: () => <div>Landing page</div>,
}))

vi.mock('./pages/login/LoginPage', () => ({
  LoginPage: () => <div>Login page</div>,
}))

function renderApp(initialEntry: string) {
  return render(
    <MemoryRouter
      future={{ v7_relativeSplatPath: true, v7_startTransition: true }}
      initialEntries={[initialEntry]}
    >
      <App />
    </MemoryRouter>,
  )
}

describe('App routes', () => {
  it('renders the nested avatar editor route inside the editor layout', () => {
    renderApp('/creative-studio/avatars/avatar-1/avatar')

    expect(screen.getByText('Avatar editor layout')).toBeInTheDocument()
    expect(screen.getByText('Avatar details step')).toBeInTheDocument()
  })
})
