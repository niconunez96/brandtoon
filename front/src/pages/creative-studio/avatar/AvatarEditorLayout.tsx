import { ArrowLeft, Layers3 } from 'lucide-react'
import { Outlet, useNavigate, useParams } from 'react-router-dom'
import { useCurrentUserQuery } from '../../../queries/useCurrentUserQuery'
import { Button } from '../../../shared/components/ui/button'
import { SidebarNav } from '../../../shared/components/ui/sidebar-nav'
import { Topbar } from '../../../shared/components/ui/topbar'

function buildEditorItems(avatarId: string) {
	const basePath = `/creative-studio/avatars/${avatarId}`

	return [
		{ end: true, href: `${basePath}/avatar`, label: 'Avatar' },
		{ href: `${basePath}/personality`, label: 'Personality' },
		{ href: `${basePath}/voice`, label: 'Voice' },
		{ href: `${basePath}/packs`, label: 'Packs' },
		{ href: `${basePath}/export`, label: 'Export' },
	]
}

export function AvatarEditorLayout() {
	const navigate = useNavigate()
	const { avatarId = '' } = useParams()
	const currentUserQuery = useCurrentUserQuery()

	return (
		<div className="foundation-page min-h-screen lg:flex">
			<SidebarNav
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
				<Topbar
					actions={
						<Button
							icon={<ArrowLeft className="size-4" />}
							onClick={() => navigate('/creative-studio')}
							variant="ghost"
						>
							Back to avatars
						</Button>
					}
					description="Start with the avatar fundamentals, then unlock the remaining steps later."
					eyebrow="Creative studio"
					title={currentUserQuery.data?.user.name ?? 'Avatar editor'}
				/>

				<main className="mx-auto flex w-full max-w-6xl flex-1 flex-col gap-6 px-4 py-6 sm:px-6 lg:px-10">
					<div className="flex items-center gap-3 rounded-3xl bg-white px-5 py-4 shadow-overshoot ring-1 ring-[color:var(--color-stroke-soft)]">
						<span className="foundation-icon">
							<Layers3 className="size-4" />
						</span>
						<div className="space-y-1">
							<p className="foundation-section-eyebrow">Avatar editor</p>
							<p className="text-lg font-black tracking-tight text-ink">
								Avatar ID {avatarId}
							</p>
						</div>
					</div>
					<Outlet />
				</main>
			</div>
		</div>
	)
}
