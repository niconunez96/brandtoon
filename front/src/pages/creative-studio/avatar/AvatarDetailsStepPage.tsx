import { Sparkles, WandSparkles } from 'lucide-react'
import { useEffect } from 'react'
import { Controller, useForm } from 'react-hook-form'
import { useParams } from 'react-router-dom'
import { z } from 'zod'
import {
  useAvatarConfigQuery,
  useGenerateAvatarOptionsMutation,
  useUpdateAvatarConfigMutation,
} from '../../../queries/useAvatarConfigQuery'
import { ApiError } from '../../../services/auth.api'
import type {
  ArtisticStyle,
  Personality,
} from '../../../services/avatar-config.api'
import { Button } from '../../../shared/components/ui/button'
import { Card, SectionShell } from '../../../shared/components/ui/card'
import { PromptField } from '../../../shared/components/ui/field'

const avatarConfigSchema = z.object({
  artisticStyle: z.enum(['2D', '3D']),
  personality: z.enum(['Friendly', 'Bold', 'Playful']),
  prompt: z
    .string()
    .max(256, 'Avatar descriptions can have up to 256 characters.'),
})

type AvatarConfigFormValues = z.infer<typeof avatarConfigSchema>

const artisticStyleOptions: ArtisticStyle[] = ['2D', '3D']
const personalityOptions: Personality[] = ['Friendly', 'Bold', 'Playful']

function getAvatarConfigErrorMessage(error: unknown) {
  if (error instanceof ApiError && error.status === 404) {
    return 'We could not find that avatar inside your studio.'
  }

  return 'We could not load this avatar draft. Please try again.'
}

function getSaveErrorMessage(error: unknown) {
  if (error instanceof ApiError && error.status === 422) {
    return 'Please review the description, style, and personality before saving.'
  }

  if (error instanceof ApiError && error.status === 404) {
    return 'This avatar is no longer available in your studio.'
  }

  return 'We could not save your avatar draft. Please try again.'
}

export function AvatarDetailsStepPage() {
  const { avatarId = '' } = useParams()
  const avatarConfigQuery = useAvatarConfigQuery(avatarId)
  const updateAvatarConfigMutation = useUpdateAvatarConfigMutation(avatarId)
  const generateAvatarOptionsMutation =
    useGenerateAvatarOptionsMutation(avatarId)
  const form = useForm<AvatarConfigFormValues>({
    defaultValues: {
      artisticStyle: '2D',
      personality: 'Friendly',
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
      personality:
        avatarConfigQuery.data.avatar_config?.personality ?? 'Friendly',
      prompt: avatarConfigQuery.data.avatar_config?.prompt ?? '',
    })
  }, [avatarConfigQuery.data, form])

  const artisticStyle = form.watch('artisticStyle')
  const personality = form.watch('personality')

  const avatarOptions =
    avatarConfigQuery.data?.avatar_config?.avatarOptions ?? []
  const selectedGeneratedOption =
    avatarOptions.find((option) => option.selected) ?? avatarOptions[0] ?? null

  async function persistDraft(values: AvatarConfigFormValues) {
    const parsed = avatarConfigSchema.safeParse(values)
    if (!parsed.success) {
      const fieldErrors = parsed.error.flatten().fieldErrors
      const promptMessage = fieldErrors.prompt?.[0]
      const styleMessage = fieldErrors.artisticStyle?.[0]
      const personalityMessage = fieldErrors.personality?.[0]

      if (promptMessage) {
        form.setError('prompt', { message: promptMessage, type: 'validate' })
      }

      if (styleMessage) {
        form.setError('artisticStyle', {
          message: styleMessage,
          type: 'validate',
        })
      }

      if (personalityMessage) {
        form.setError('personality', {
          message: personalityMessage,
          type: 'validate',
        })
      }

      return false
    }

    try {
      await updateAvatarConfigMutation.mutateAsync(parsed.data)
      return true
    } catch {
      return false
    }
  }

  const handleSaveSubmit = form.handleSubmit(async (values) => {
    await persistDraft(values)
  })

  const handleGenerate = form.handleSubmit(async (values) => {
    const saved = await persistDraft(values)
    if (!saved) {
      return
    }

    await generateAvatarOptionsMutation.mutateAsync()
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
        <Button
          onClick={() => void avatarConfigQuery.refetch()}
          variant="secondary"
        >
          Try again
        </Button>
      </Card>
    )
  }

  return (
    <SectionShell
      actions={
        <>
          <Button
            isLoading={generateAvatarOptionsMutation.isPending}
            onClick={() => void handleGenerate()}
          >
            Generate
          </Button>
          <Button
            isLoading={
              updateAvatarConfigMutation.isPending &&
              !generateAvatarOptionsMutation.isPending
            }
            onClick={() => void handleSaveSubmit()}
            variant="secondary"
          >
            Save as draft
          </Button>
        </>
      }
      description="Define the first creative draft for this avatar. The other steps are visible so the workflow feels real, but only this foundation step is active today."
      eyebrow="Avatar step"
      title="Shape the avatar foundation"
    >
      <form
        className="space-y-6"
        onSubmit={(event) => void handleSaveSubmit(event)}
      >
        <input type="hidden" {...form.register('artisticStyle')} />
        <input type="hidden" {...form.register('personality')} />

        <div className="grid gap-6 xl:grid-cols-[minmax(0,1.45fr)_minmax(360px,1fr)]">
          <div className="space-y-5">
            <Card className="space-y-4 bg-white p-4 md:p-6">
              <div className="relative flex aspect-square items-end overflow-hidden rounded-[2.25rem] bg-gradient-to-br from-[#FCE7E7] via-white to-[#dfe6e9] p-5 shadow-overshoot">
                <span className="rounded-full bg-coral px-4 py-2 text-[11px] font-black uppercase tracking-[0.2em] text-white shadow-sticker">
                  Active prototype
                </span>
              </div>

              <div className="space-y-3">
                <div>
                  <p className="foundation-section-eyebrow">
                    Generation Results
                  </p>
                  <p className="text-sm text-ink/70">
                    Choose a generated direction
                  </p>
                </div>

                {avatarOptions.length === 0 ? (
                  <Card className="bg-surface p-4 text-sm text-ink/70">
                    Save your draft, then generate avatar options from the
                    backend.
                  </Card>
                ) : (
                  <fieldset>
                    <legend className="sr-only">
                      Generated avatar image options
                    </legend>
                    <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
                      {avatarOptions.map((option, index) => {
                        const isSelected = option.selected

                        return (
                          <label
                            className="block cursor-pointer"
                            key={option.href}
                          >
                            <input
                              checked={isSelected}
                              className="sr-only"
                              name="generated-avatar-option"
                              readOnly
                              type="radio"
                              value={option.href}
                            />
                            <div
                              className={`space-y-2 rounded-3xl border p-2 transition focus-within:ring-2 focus-within:ring-coral focus-within:ring-offset-2 focus-within:ring-offset-white ${
                                isSelected
                                  ? 'border-coral bg-coral/5 shadow-sticker'
                                  : 'border-[color:var(--color-stroke-soft)] bg-surface hover:bg-white'
                              }`}
                            >
                              <div className="aspect-square rounded-2xl bg-gradient-to-br from-[#FCE7E7] via-white to-[#dfe6e9]" />
                              <p className="px-1 text-[11px] font-extrabold uppercase tracking-[0.18em] text-ink/65">
                                Option {index + 1}
                              </p>
                            </div>
                          </label>
                        )
                      })}
                    </div>
                  </fieldset>
                )}

                {selectedGeneratedOption ? (
                  <p className="text-sm text-ink/70">
                    Latest generated asset: {selectedGeneratedOption.href}
                  </p>
                ) : null}
              </div>
            </Card>
          </div>

          <div className="space-y-4">
            <Card className="space-y-5 bg-white p-6">
              <div className="space-y-2">
                <p className="foundation-section-eyebrow">Create avatar</p>
                <p className="text-xl font-black tracking-tight text-ink">
                  Shape your brand's face
                </p>
                </div>

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

              <div className="space-y-3">
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
                        className={`cursor-pointer inline-flex min-h-11 items-center gap-2 rounded-full border px-5 py-2.5 text-sm font-extrabold tracking-[-0.02em] transition ${
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

              <div className="space-y-3">
                <div className="space-y-2">
                  <p className="foundation-section-eyebrow">Personality</p>
                  <p className="text-lg font-black tracking-tight text-ink">
                    Choose the brand energy
                  </p>
                </div>
                <div className="flex flex-wrap gap-3">
                  {personalityOptions.map((option) => {
                    const isActive = personality === option

                    return (
                      <button
                        aria-pressed={isActive}
                        className={`cursor-pointer inline-flex min-h-11 items-center rounded-full border px-5 py-2.5 text-sm font-extrabold tracking-[-0.02em] transition ${
                          isActive
                            ? 'border-transparent bg-coral text-white shadow-sticker'
                            : 'border-[color:var(--color-stroke-soft)] bg-surface text-ink hover:bg-white'
                        }`}
                        key={option}
                        onClick={() => {
                          form.clearErrors('personality')
                          form.setValue('personality', option, {
                            shouldDirty: true,
                            shouldTouch: true,
                          })
                        }}
                        type="button"
                      >
                        <span>{option}</span>
                      </button>
                    )
                  })}
                </div>
                {form.formState.errors.personality?.message ? (
                  <p className="text-xs font-bold text-error">
                    {form.formState.errors.personality.message}
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

              {generateAvatarOptionsMutation.isError ? (
                <p className="rounded-2xl bg-error-container px-4 py-3 text-sm font-bold text-error">
                  We could not start avatar generation. Please try again.
                </p>
              ) : null}

              <div className="flex flex-wrap gap-3">
                <Button
                  onClick={() =>
                    form.reset({
                      artisticStyle:
                        avatarConfigQuery.data?.avatar_config?.artisticStyle ??
                        '2D',
                      personality:
                        avatarConfigQuery.data?.avatar_config?.personality ??
                        'Friendly',
                      prompt:
                        avatarConfigQuery.data?.avatar_config?.prompt ?? '',
                    })
                  }
                  type="button"
                  variant="ghost"
                >
                  Reset draft
                </Button>
              </div>
            </Card>
          </div>
        </div>
      </form>
    </SectionShell>
  )
}
