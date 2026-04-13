import { Card, SectionShell } from '../../../shared/components/ui/card'

type AvatarPlaceholderStepPageProps = {
	stepLabel: string
}

export function AvatarPlaceholderStepPage({
	stepLabel,
}: AvatarPlaceholderStepPageProps) {
	return (
		<SectionShell
			description="This step is intentionally blank for now so the full workflow is visible while implementation continues in future changes."
			eyebrow="Coming soon"
			title={`${stepLabel} step`}
		>
			<Card className="space-y-3 bg-white">
				<p className="foundation-section-eyebrow">Placeholder</p>
				<p className="text-xl font-black tracking-tight text-ink">
					{stepLabel} will render here soon.
				</p>
			</Card>
		</SectionShell>
	)
}
