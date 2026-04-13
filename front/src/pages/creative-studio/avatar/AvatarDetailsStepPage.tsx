import { Sparkles, WandSparkles } from 'lucide-react'
import { useEffect, useState } from 'react'
import { Controller, useForm } from 'react-hook-form'
import { useParams } from 'react-router-dom'
import { z } from 'zod'
import {
	useAvatarConfigQuery,
	useUpdateAvatarConfigMutation,
} from '../../../queries/useAvatarConfigQuery'
import type { ArtisticStyle } from '../../../services/avatar-config.api'
import { ApiError } from '../../../services/auth.api'
import { Button } from '../../../shared/components/ui/button'
import { Card, SectionShell } from '../../../shared/components/ui/card'
import { PromptField } from '../../../shared/components/ui/field'

const avatarConfigSchema = z.object({
	artisticStyle: z.enum(['2D', '3D']),
	prompt: z
		.string()
		.max(256, 'Avatar descriptions can have up to 256 characters.'),
})

type AvatarConfigFormValues = z.infer<typeof avatarConfigSchema>

const artisticStyleOptions: ArtisticStyle[] = ['2D', '3D']

function getAvatarConfigErrorMessage(error: unknown) {
	if (error instanceof ApiError && error.status === 404) {
		return 'We could not find that avatar inside your studio.'
	}

	return 'We could not load this avatar draft. Please try again.'
}

function getSaveErrorMessage(error: unknown) {
	if (error instanceof ApiError && error.status === 422) {
		return 'Please review the description and style before saving.'
	}

	if (error instanceof ApiError && error.status === 404) {
		return 'This avatar is no longer available in your studio.'
	}

	return 'We could not save your avatar draft. Please try again.'
}

export function AvatarDetailsStepPage() {
	const { avatarId = '' } = useParams()
	const [feedbackMessage, setFeedbackMessage] = useState<string | null>(null)
	const avatarConfigQuery = useAvatarConfigQuery(avatarId)
	const updateAvatarConfigMutation = useUpdateAvatarConfigMutation(avatarId)
	const form = useForm<AvatarConfigFormValues>({
		defaultValues: {
			artisticStyle: '2D',
			prompt: '',
		},
	})

	useEffect(() => {
		if (!avatarConfigQuery.data) {
			return
		}

		form.reset({
			artisticStyle:
				avatarConfigQuery.data.avatar_config?.artisticStyle ?? '2D',
			prompt: avatarConfigQuery.data.avatar_config?.prompt ?? '',
		})
	}, [avatarConfigQuery.data, form])

	const artisticStyle = form.watch('artisticStyle')

	const handleSubmit = form.handleSubmit(async (values) => {
		setFeedbackMessage(null)
		const parsed = avatarConfigSchema.safeParse(values)
		if (!parsed.success) {
			const fieldErrors = parsed.error.flatten().fieldErrors
			const promptMessage = fieldErrors.prompt?.[0]
			const styleMessage = fieldErrors.artisticStyle?.[0]

			if (promptMessage) {
				form.setError('prompt', { message: promptMessage, type: 'validate' })
			}

			if (styleMessage) {
				form.setError('artisticStyle', {
					message: styleMessage,
					type: 'validate',
				})
			}

			return
		}

		try {
			await updateAvatarConfigMutation.mutateAsync(parsed.data)
			setFeedbackMessage('Avatar draft saved.')
		} catch {
			return
		}
	})

	if (avatarConfigQuery.isLoading) {
		return (
			<Card className="space-y-3 bg-white">
				<p className="foundation-section-eyebrow">Loading avatar draft</p>
				<p className="text-xl font-black tracking-tight text-ink">
					Preparing the first editing step…
				</p>
			</Card>
		)
	}

	if (avatarConfigQuery.isError) {
		return (
			<Card className="space-y-4 bg-white">
				<div className="space-y-2">
					<p className="foundation-section-eyebrow">Avatar draft error</p>
					<p className="text-xl font-black tracking-tight text-ink">
						{getAvatarConfigErrorMessage(avatarConfigQuery.error)}
					</p>
				</div>
				<Button onClick={() => void avatarConfigQuery.refetch()} variant="secondary">
					Try again
				</Button>
			</Card>
		)
	}

	return (
		<SectionShell
			description="Define the first creative draft for this avatar. The other steps are visible so the workflow feels real, but only this foundation step is active today."
			eyebrow="Avatar step"
			title="Shape the avatar foundation"
		>
			<form className="space-y-6" onSubmit={(event) => void handleSubmit(event)}>
				<input type="hidden" {...form.register('artisticStyle')} />

				<Controller
					control={form.control}
					name="prompt"
					render={({ field }) => (
						<PromptField
							icon={<WandSparkles className="size-4" />}
							maxLength={256}
							name={field.name}
							onBlur={field.onBlur}
							onChange={field.onChange}
							placeholder="Describe the personality, silhouette, and visual energy you want this avatar to carry."
							title="Avatar description"
							value={field.value}
						/>
					)}
				/>

				<div className="space-y-3 rounded-card bg-white p-6 shadow-overshoot ring-1 ring-[color:var(--color-stroke-soft)]">
					<div className="space-y-2">
						<p className="foundation-section-eyebrow">Artistic style</p>
						<p className="text-lg font-black tracking-tight text-ink">
							Choose the visual direction
						</p>
					</div>
					<div className="flex flex-wrap gap-3">
						{artisticStyleOptions.map((option) => {
							const isActive = artisticStyle === option
							return (
								<button
									aria-pressed={isActive}
									className={`inline-flex min-h-11 items-center gap-2 rounded-full border px-5 py-2.5 text-sm font-extrabold tracking-[-0.02em] transition ${
										isActive
											? 'border-transparent bg-coral text-white shadow-sticker'
											: 'border-[color:var(--color-stroke-soft)] bg-surface text-ink hover:bg-white'
									}`}
									key={option}
									onClick={() => {
										form.clearErrors('artisticStyle')
										form.setValue('artisticStyle', option, {
											shouldDirty: true,
											shouldTouch: true,
										})
									}}
									type="button"
								>
									<Sparkles className="size-4" />
									<span>{option}</span>
								</button>
							)
						})}
					</div>
					{form.formState.errors.artisticStyle?.message ? (
						<p className="text-xs font-bold text-error">
							{form.formState.errors.artisticStyle.message}
						</p>
					) : null}
				</div>

				{form.formState.errors.prompt?.message ? (
					<p className="rounded-2xl bg-error-container px-4 py-3 text-sm font-bold text-error">
						{form.formState.errors.prompt.message}
					</p>
				) : null}

				{updateAvatarConfigMutation.isError ? (
					<p className="rounded-2xl bg-error-container px-4 py-3 text-sm font-bold text-error">
						{getSaveErrorMessage(updateAvatarConfigMutation.error)}
					</p>
				) : null}

				{feedbackMessage ? (
					<p className="rounded-2xl bg-secondary/10 px-4 py-3 text-sm font-bold text-secondary">
						{feedbackMessage}
					</p>
				) : null}

				<div className="flex flex-wrap gap-3">
					<Button isLoading={updateAvatarConfigMutation.isPending} type="submit">
						Save avatar draft
					</Button>
					<Button
						onClick={() =>
							form.reset({
								artisticStyle:
									avatarConfigQuery.data?.avatar_config?.artisticStyle ?? '2D',
								prompt: avatarConfigQuery.data?.avatar_config?.prompt ?? '',
							})
						}
						type="button"
						variant="ghost"
					>
						Reset draft
					</Button>
				</div>
			</form>
		</SectionShell>
	)
}
