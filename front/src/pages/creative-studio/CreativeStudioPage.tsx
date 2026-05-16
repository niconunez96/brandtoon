import { Sparkles, UserRoundPlus, Users } from 'lucide-react'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { useNavigate } from 'react-router-dom'
import { z } from 'zod'
import {
  useAvatarsQuery,
  useCreateAvatarMutation,
} from '../../queries/useAvatarsQuery'
import { ApiError } from '../../services/auth.api'
import { Button } from '../../shared/components/ui/button'
import { Card, SectionShell } from '../../shared/components/ui/card'
import { EmptyState } from '../../shared/components/ui/empty-state'
import { Input } from '../../shared/components/ui/field'
import { Modal } from '../../shared/components/ui/modal'
import { SidebarNav } from '../../shared/components/ui/sidebar-nav'
import { CreativeStudioSidebarFooter } from './CreativeStudioSidebarFooter'

const createAvatarSchema = z.object({
  name: z
    .string()
    .trim()
    .min(1, 'Please enter an avatar name.')
    .max(120, 'Avatar names can have up to 120 characters.'),
})

type CreateAvatarFormValues = {
  name: string
}

function getCreateAvatarErrorMessage(error: unknown) {
  if (error instanceof ApiError && error.status === 422) {
    return 'Please enter a valid avatar name before saving.'
  }

  return 'We could not create your avatar. Please try again.'
}

function AvatarComposer({
  isOpen,
  onClose,
}: {
  isOpen: boolean
  onClose: () => void
}) {
  const createAvatarMutation = useCreateAvatarMutation()
  const form = useForm<CreateAvatarFormValues>({
    defaultValues: { name: '' },
  })

  const handleSubmit = form.handleSubmit(async (values) => {
    const parsed = createAvatarSchema.safeParse(values)
    if (!parsed.success) {
      const message = parsed.error.flatten().fieldErrors.name?.[0]
      if (message) {
        form.setError('name', { message, type: 'validate' })
      }

      return
    }

    try {
      await createAvatarMutation.mutateAsync(parsed.data.name)
      form.reset()
      onClose()
    } catch {
      return
    }
  })

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <Card className="space-y-5 bg-white">
        <div className="space-y-2">
          <p className="foundation-section-eyebrow">Avatar composer</p>
          <h3 className="text-xl font-black tracking-tight text-ink">
            Name your next avatar
          </h3>
          <p className="foundation-body">
            Start simple. You can enrich each avatar later as the studio grows.
          </p>
        </div>

        <form
          className="space-y-4"
          onSubmit={(event) => void handleSubmit(event)}
        >
          <Input
            label="Avatar name"
            maxLength={120}
            message={form.formState.errors.name?.message}
            placeholder="For example: Studio Hero"
            state={form.formState.errors.name ? 'error' : 'default'}
            {...form.register('name')}
          />

          {createAvatarMutation.isError ? (
            <p className="rounded-2xl bg-error-container px-4 py-3 text-sm font-bold text-error">
              {getCreateAvatarErrorMessage(createAvatarMutation.error)}
            </p>
          ) : null}

          <div className="flex flex-wrap gap-3">
            <Button isLoading={createAvatarMutation.isPending} type="submit">
              Save avatar
            </Button>
            <Button onClick={onClose} type="button" variant="ghost">
              Cancel
            </Button>
          </div>
        </form>
      </Card>
    </Modal>
  )
}

export function CreativeStudioPage() {
  const navigate = useNavigate()
  const [isComposerOpen, setIsComposerOpen] = useState(false)
  const avatarsQuery = useAvatarsQuery()

  return (
    <div className="foundation-page min-h-screen lg:flex">
      <SidebarNav
        bottomContent={<CreativeStudioSidebarFooter />}
        items={[{ active: true, href: '/creative-studio', label: 'Avatars' }]}
        title="Brandtoon"
      />

      <div className="flex min-h-screen flex-1 flex-col">
        <main className="mx-auto flex w-full max-w-6xl flex-1 flex-col gap-6 px-4 py-6 sm:px-6 lg:px-10">
          <SectionShell
            description="This is your home for character creation. Start with one avatar, then keep expanding your roster."
            eyebrow="Avatar library"
            title="Build your cast"
          >
            <div className="space-y-6">
              {avatarsQuery.isLoading ? (
                <Card className="space-y-3 bg-white">
                  <p className="foundation-section-eyebrow">Loading avatars</p>
                  <p className="text-xl font-black tracking-tight text-ink">
                    Pulling your avatar library into the studio…
                  </p>
                </Card>
              ) : null}

              {avatarsQuery.isError ? (
                <Card className="space-y-4 bg-white">
                  <div className="space-y-2">
                    <p className="foundation-section-eyebrow">Avatar error</p>
                    <p className="text-xl font-black tracking-tight text-ink">
                      We could not load your avatars.
                    </p>
                    <p className="foundation-body">
                      Try again to reconnect with the Creative Studio service.
                    </p>
                  </div>
                  <Button
                    onClick={() => void avatarsQuery.refetch()}
                    variant="secondary"
                  >
                    Try again
                  </Button>
                </Card>
              ) : null}

              {!avatarsQuery.isLoading && !avatarsQuery.isError ? (
                avatarsQuery.data?.avatars.length ? (
                  <>
                    <div className="flex flex-wrap items-center justify-between gap-3">
                      <div>
                        <p className="foundation-section-eyebrow">
                          Existing avatars
                        </p>
                        <p className="text-lg font-black tracking-tight text-ink">
                          {avatarsQuery.data.avatars.length} avatar
                          {avatarsQuery.data.avatars.length === 1 ? '' : 's'}{' '}
                          ready to use
                        </p>
                      </div>
                      <Button
                        icon={<UserRoundPlus className="size-4" />}
                        onClick={() => setIsComposerOpen(true)}
                      >
                        Create avatar
                      </Button>
                    </div>

                    <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
                      {avatarsQuery.data.avatars.map((avatar) => (
                        <Card className="space-y-4 bg-white" key={avatar.id}>
                          <button
                            className="w-full space-y-4 text-left"
                            onClick={() =>
                              navigate(
                                `/creative-studio/avatars/${avatar.id}/avatar`,
                              )
                            }
                            type="button"
                          >
                            <div className="foundation-icon">
                              <Sparkles className="size-4" />
                            </div>
                            <div className="space-y-2">
                              <p className="foundation-section-eyebrow">
                                Avatar card
                              </p>
                              <h3 className="text-xl font-black tracking-tight text-ink">
                                {avatar.name}
                              </h3>
                              <p className="foundation-body">
                                Ready for future motion, voice, and scene
                                workflows.
                              </p>
                            </div>
                            <span className="text-sm font-extrabold text-coral cursor-pointer">
                              Customize avatar
                            </span>
                            <span className="sr-only">
                              Open {avatar.name} avatar editor
                            </span>
                          </button>
                        </Card>
                      ))}
                    </div>
                  </>
                ) : (
                  <>
                    <EmptyState
                      action={
                        <Button onClick={() => setIsComposerOpen(true)}>
                          Create avatar now
                        </Button>
                      }
                      eyebrow="No avatars yet"
                      icon={<Users className="size-4" />}
                      title="Start your avatar library"
                    >
                      Create your first avatar to begin shaping the characters
                      that represent your brand.
                    </EmptyState>
                  </>
                )
              ) : null}

              <AvatarComposer
                isOpen={isComposerOpen}
                onClose={() => setIsComposerOpen(false)}
              />
            </div>
          </SectionShell>
        </main>
      </div>
    </div>
  )
}
