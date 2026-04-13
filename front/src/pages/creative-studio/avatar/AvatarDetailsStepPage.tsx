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
import { Toast } from '../../../shared/components/ui/toast'

type GeneratedAvatarOption = {
  id: string
  label: string
  description: string
  gradientClassName: string
}

const avatarConfigSchema = z.object({
  artisticStyle: z.enum(['2D', '3D']),
  prompt: z
    .string()
    .max(256, 'Avatar descriptions can have up to 256 characters.'),
})

type AvatarConfigFormValues = z.infer<typeof avatarConfigSchema>

type ToastVisibility = 'hidden' | 'entering' | 'visible' | 'exiting'

const artisticStyleOptions: ArtisticStyle[] = ['2D', '3D']

const mockImageGradients = [
  'from-[#FCE7E7] via-white to-[#dfe6e9]',
  'from-[#6B4EE0]/20 via-white to-[#FCE7E7]',
  'from-[#00b179]/20 via-white to-[#dfe6e9]',
  'from-[#FF6B6B]/20 via-white to-[#FCE7E7]',
] as const

const mockImageDescriptors = {
  '2D': [
    'Editorial pose',
    'Soft sticker energy',
    'Mascot close-up',
    'Hero silhouette',
  ],
  '3D': [
    'Studio render',
    'Glow lighting pass',
    'Depth-focused angle',
    'Polished hero frame',
  ],
} as const satisfies Record<ArtisticStyle, readonly string[]>

function createGeneratedAvatarOptions({
  seed,
  style,
  prompt,
}: {
  seed: number
  style: ArtisticStyle
  prompt: string
}): GeneratedAvatarOption[] {
  const trimmedPrompt = prompt.trim()
  const promptDescriptor =
    trimmedPrompt.length > 0 ? trimmedPrompt : 'Creative direction'

  return Array.from({ length: 4 }, (_, index) => {
    const descriptor =
      mockImageDescriptors[style][
        (seed + index) % mockImageDescriptors[style].length
      ]
    const promptSnippet = promptDescriptor.slice(0, 48)

    return {
      description: `${style} preview ${index + 1}`,
      gradientClassName:
        mockImageGradients[(seed + index) % mockImageGradients.length],
      id: `${style.toLowerCase()}-${seed}-${index}`,
      label: `${descriptor} · ${promptSnippet}`,
    }
  })
}

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
  const [toastVisibility, setToastVisibility] =
    useState<ToastVisibility>('hidden')
  const [generationSeed, setGenerationSeed] = useState(0)
  const [generatedOptions, setGeneratedOptions] = useState<
    GeneratedAvatarOption[]
  >(() => createGeneratedAvatarOptions({ prompt: '', seed: 0, style: '2D' }))
  const [selectedGeneratedOptionId, setSelectedGeneratedOptionId] = useState(
    () =>
      createGeneratedAvatarOptions({ prompt: '', seed: 0, style: '2D' })[0]
        ?.id ?? '',
  )
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

  useEffect(() => {
    if (!feedbackMessage || toastVisibility === 'hidden') {
      return
    }

    if (toastVisibility === 'entering') {
      const enterTimer = window.setTimeout(() => {
        setToastVisibility('visible')
      }, 16)

      return () => window.clearTimeout(enterTimer)
    }

    if (toastVisibility === 'visible') {
      const visibleTimer = window.setTimeout(() => {
        setToastVisibility('exiting')
      }, 2400)

      return () => window.clearTimeout(visibleTimer)
    }

    const exitTimer = window.setTimeout(() => {
      setToastVisibility('hidden')
      setFeedbackMessage(null)
    }, 250)

    return () => window.clearTimeout(exitTimer)
  }, [feedbackMessage, toastVisibility])

  const artisticStyle = form.watch('artisticStyle')
  const prompt = form.watch('prompt')
  const selectedGeneratedOption =
    generatedOptions.find(
      (option) => option.id === selectedGeneratedOptionId,
    ) ?? generatedOptions[0]

  function showSaveSuccessToast(message: string) {
    setFeedbackMessage(message)
    setToastVisibility('entering')
  }

  function dismissToast() {
    setToastVisibility((currentVisibility) =>
      currentVisibility === 'hidden' ? currentVisibility : 'exiting',
    )
  }

  function handleRegenerate() {
    const nextSeed = generationSeed + 1
    const nextOptions = createGeneratedAvatarOptions({
      prompt,
      seed: nextSeed,
      style: artisticStyle,
    })

    setGenerationSeed(nextSeed)
    setGeneratedOptions(nextOptions)
    setSelectedGeneratedOptionId(nextOptions[0]?.id ?? '')
  }

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
      showSaveSuccessToast('Avatar draft saved.')
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
      description="Define the first creative draft for this avatar. The other steps are visible so the workflow feels real, but only this foundation step is active today."
      eyebrow="Avatar step"
      title="Shape the avatar foundation"
    >
      {feedbackMessage && toastVisibility !== 'hidden' ? (
        <div className="pointer-events-none fixed right-4 top-4 z-50 sm:right-6 sm:top-6">
          <div
            className={`pointer-events-auto transition-all duration-250 ease-out ${
              toastVisibility === 'visible'
                ? 'translate-x-0 opacity-100'
                : '-translate-x-6 opacity-0'
            }`}
          >
            <Toast onDismiss={dismissToast} title={feedbackMessage}>
              Your latest avatar draft was synced.
            </Toast>
          </div>
        </div>
      ) : null}

      <form
        className="space-y-6"
        onSubmit={(event) => void handleSubmit(event)}
      >
        <input type="hidden" {...form.register('artisticStyle')} />

        <div className="grid gap-6 xl:grid-cols-[minmax(0,1.45fr)_minmax(360px,1fr)]">
          <div className="space-y-5">
            <Card className="space-y-4 bg-white p-4 md:p-6">
              <div
                className={`relative flex aspect-square items-end overflow-hidden rounded-[2.25rem] bg-gradient-to-br p-5 shadow-overshoot ${selectedGeneratedOption?.gradientClassName ?? mockImageGradients[0]}`}
              >
                <span className="rounded-full bg-coral px-4 py-2 text-[11px] font-black uppercase tracking-[0.2em] text-white shadow-sticker">
                  Active prototype
                </span>
              </div>

              <div className="space-y-3">
                <div className="flex items-center justify-between gap-3">
                  <div>
                    <p className="foundation-section-eyebrow">
                      Generation Results
                    </p>
                    <p className="text-sm text-ink/70">
                      Choose a generated direction
                    </p>
                  </div>
                  <Button
                    onClick={handleRegenerate}
                    type="button"
                    variant="ghost"
                  >
                    Regenerate
                  </Button>
                </div>

                <fieldset>
                  <legend className="sr-only">
                    Generated avatar image options
                  </legend>
                  <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
                    {generatedOptions.map((option, index) => {
                      const isSelected = selectedGeneratedOptionId === option.id

                      return (
                        <label className="block cursor-pointer" key={option.id}>
                          <input
                            checked={isSelected}
                            className="sr-only"
                            name="generated-avatar-option"
                            onChange={() =>
                              setSelectedGeneratedOptionId(option.id)
                            }
                            type="radio"
                            value={option.id}
                          />
                          <div
                            className={`space-y-2 rounded-3xl border p-2 transition focus-within:ring-2 focus-within:ring-coral focus-within:ring-offset-2 focus-within:ring-offset-white ${
                              isSelected
                                ? 'border-coral bg-coral/5 shadow-sticker'
                                : 'border-[color:var(--color-stroke-soft)] bg-surface hover:bg-white'
                            }`}
                          >
                            <div
                              className={`aspect-square rounded-2xl bg-gradient-to-br ${option.gradientClassName}`}
                            />
                            <p className="px-1 text-[11px] font-extrabold uppercase tracking-[0.18em] text-ink/65">
                              Option {index + 1}
                            </p>
                          </div>
                        </label>
                      )
                    })}
                  </div>
                </fieldset>
              </div>
            </Card>
          </div>

          <div className="space-y-4">
            <Card className="space-y-5 bg-white p-6">
              <div className="space-y-2">
                <p className="foundation-section-eyebrow">Create avatar</p>
                <p className="text-3xl font-black tracking-tight text-ink">
                  Shape your brand's face
                </p>
                <p className="text-sm text-ink/70">
                  Craft your first draft, then iterate with generated options on
                  the left.
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

              <div className="flex flex-wrap gap-3">
                <Button
                  isLoading={updateAvatarConfigMutation.isPending}
                  type="submit"
                >
                  Save avatar draft
                </Button>
                <Button
                  onClick={() =>
                    form.reset({
                      artisticStyle:
                        avatarConfigQuery.data?.avatar_config?.artisticStyle ??
                        '2D',
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
