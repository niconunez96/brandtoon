import { Outlet, useParams } from 'react-router-dom'
import { SidebarNav } from '../../../shared/components/ui/sidebar-nav'
import { CreativeStudioSidebarFooter } from '../CreativeStudioSidebarFooter'

function buildEditorItems(avatarId: string) {
  const basePath = `/creative-studio/avatars/${avatarId}`

  return [
    { end: true, href: `${basePath}/avatar`, label: 'Avatar' },
    { href: `${basePath}/voice`, label: 'Voice' },
    { href: `${basePath}/packs`, label: 'Packs' },
    { href: `${basePath}/export`, label: 'Export' },
  ]
}

export function AvatarEditorLayout() {
  const { avatarId = '' } = useParams()

  return (
    <div className="foundation-page min-h-screen lg:flex">
      <SidebarNav
        bottomContent={<CreativeStudioSidebarFooter />}
        footer={
          <div className="space-y-2">
            <p className="text-xs font-black uppercase tracking-section text-white/72">
              Avatar workflow
            </p>
            <p className="text-xl font-black tracking-tight text-white">
              Shape the character foundation before unlocking the rest.
            </p>
          </div>
        }
        items={buildEditorItems(avatarId)}
        title="Brandtoon"
      />

      <div className="flex min-h-screen flex-1 flex-col">
        <main className="mx-auto flex w-full max-w-6xl flex-1 flex-col gap-6 px-4 py-6 sm:px-6 lg:px-10">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
