import { ArrowRight, CheckCircle2, Film, Play, Sparkles } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { ActionChip } from '../../shared/components/ui/action-chip'
import { Badge } from '../../shared/components/ui/badge'
import { Button } from '../../shared/components/ui/button'
import { Card, SectionShell } from '../../shared/components/ui/card'
import { Input } from '../../shared/components/ui/field'
import { Slider } from '../../shared/components/ui/slider'
import { Topbar } from '../../shared/components/ui/topbar'

const NAV_LINKS = [
  { href: '#how-it-works', label: 'How it Works' },
  { href: '#showcase', label: 'Showcase' },
  { href: '#pricing', label: 'Pricing' },
] as const

const SHOWCASE_CARDS = [
  {
    title: 'Mascot lead',
    label: 'Hero character',
    description:
      'A confident brand mascot that anchors the landing story with a playful lead pose.',
    imageSrc: '/images/landing/characters/character-01-mascot.png',
    imageAlt: 'Brandtoon mascot character in a playful hero pose',
    className:
      'md:col-span-2 min-h-[22rem] bg-[radial-gradient(circle_at_top,_rgba(107,78,224,0.18),_transparent_44%),linear-gradient(135deg,_#ffffff_0%,_#f1f2f6_100%)] text-ink',
    imageClassName: 'max-h-[16rem] sm:max-h-[18rem] md:max-h-[19rem]',
  },
  {
    title: 'Robot energy',
    label: 'Tech-forward',
    description:
      'A transparent robot asset that keeps the cast balanced with a sharper futuristic note.',
    imageSrc: '/images/landing/characters/character-02-robot.png',
    imageAlt: 'Brandtoon robot character with a futuristic expression',
    className:
      'min-h-[22rem] bg-[linear-gradient(135deg,_#eef4ff_0%,_#dce4ff_100%)] text-ink',
    imageClassName: 'max-h-[15rem] sm:max-h-[16rem]',
  },
  {
    title: 'Executive polish',
    label: 'Presentation-ready',
    description:
      'A more formal character option for decks, proposals, and branded corporate touchpoints.',
    imageSrc: '/images/landing/characters/character-03-executive.png',
    imageAlt: 'Brandtoon executive character in a polished business outfit',
    className:
      'min-h-[22rem] bg-[linear-gradient(135deg,_#fff5f3_0%,_#fce7e7_100%)] text-ink',
    imageClassName: 'max-h-[15rem] sm:max-h-[16rem]',
  },
  {
    title: 'Youthful sidekick',
    label: 'Friendly support',
    description:
      'A softer supporting character that helps the cast stay approachable in lighter campaigns.',
    imageSrc: '/images/landing/characters/character-04-boy.png',
    imageAlt: 'Brandtoon boy character smiling with an approachable expression',
    className:
      'min-h-[20rem] bg-[linear-gradient(135deg,_#f8fbff_0%,_#e7f2ff_100%)] text-ink',
    imageClassName: 'max-h-[14rem] sm:max-h-[15rem]',
  },
  {
    title: 'Bear companion',
    label: 'Warm mascot',
    description:
      'A rounded mascot variant that adds warmth and sticker-pack charm to the landing showcase.',
    imageSrc: '/images/landing/characters/character-05-bear.png',
    imageAlt: 'Brandtoon bear character with a warm friendly pose',
    className:
      'min-h-[20rem] bg-[linear-gradient(135deg,_#fff8ef_0%,_#f6e2c9_100%)] text-ink',
    imageClassName: 'max-h-[14rem] sm:max-h-[15rem]',
  },
] as const

const HERO_CHARACTERS = [
  {
    title: 'Mascot lead',
    imageSrc: '/images/landing/characters/character-01-mascot.png',
    imageAlt: 'Brandtoon mascot character leading the hero section',
    className: 'landing-hero-character-card min-h-[18rem] sm:min-h-[24rem]',
    imageClassName: 'max-h-[15rem] sm:max-h-[21rem]',
  },
  {
    title: 'Robot spark',
    imageSrc: '/images/landing/characters/character-02-robot.png',
    imageAlt: 'Brandtoon robot character featured in the hero section',
    className: 'landing-hero-character-card min-h-[10rem] sm:min-h-[11.5rem]',
    imageClassName: 'max-h-[8.5rem] sm:max-h-[9.75rem]',
  },
  {
    title: 'Executive polish',
    imageSrc: '/images/landing/characters/character-03-executive.png',
    imageAlt: 'Brandtoon executive character featured in the hero section',
    className: 'landing-hero-character-card min-h-[10rem] sm:min-h-[11.5rem]',
    imageClassName: 'max-h-[8.5rem] sm:max-h-[9.75rem]',
  },
  {
    title: 'Youthful sidekick',
    imageSrc: '/images/landing/characters/character-04-boy.png',
    imageAlt: 'Brandtoon boy character featured in the hero section',
    className: 'landing-hero-character-card min-h-[9.5rem] sm:min-h-[10.5rem]',
    imageClassName: 'max-h-[7.5rem] sm:max-h-[8.5rem]',
  },
  {
    title: 'Bear companion',
    imageSrc: '/images/landing/characters/character-05-bear.png',
    imageAlt: 'Brandtoon bear character featured in the hero section',
    className: 'landing-hero-character-card min-h-[9.5rem] sm:min-h-[10.5rem]',
    imageClassName: 'max-h-[7.5rem] sm:max-h-[8.5rem]',
  },
] as const

const PROCESS_STEPS = [
  {
    step: '1. Describe',
    title: 'Describe your vision',
    body: 'Say what the brand should feel like and anchor the mascot energy in a few words.',
  },
  {
    step: '2. Tune',
    title: 'Personalize the vibe',
    body: 'Use the existing sliders and controls to push the character toward playful or polished.',
  },
  {
    step: '3. Animate',
    title: 'Export motion assets',
    body: 'Deliver PNG, GIF, and MP4 outputs that are ready for campaigns, decks, and social.',
  },
] as const

const PRICING_FEATURES = [
  'Commercial license',
  'High-res 4K exports',
  'SVG vector options',
  '3D source files',
] as const

export function LandingPage() {
  const navigate = useNavigate()
  const goToLogin = () => navigate('/login?next=%2Fcreative-studio')

  return (
    <div className="foundation-page landing-shell">
      <header className="sticky top-0 z-30 px-4 pt-4 sm:px-6 lg:px-8">
        <Topbar
          actions={
            <>
              <nav
                aria-label="Landing navigation"
                className="hidden items-center gap-6 md:flex"
              >
                {NAV_LINKS.map((link) => (
                  <a
                    className="text-sm font-semibold text-ink/64 transition hover:text-coral"
                    href={link.href}
                    key={link.href}
                  >
                    {link.label}
                  </a>
                ))}
              </nav>
              <Button
                className="min-h-10 px-5 py-2 text-[11px] uppercase tracking-section"
                onClick={goToLogin}
              >
                Start creating
              </Button>
            </>
          }
          className="mx-auto max-w-6xl rounded-full border border-white/70 bg-white/80 px-5 py-3 shadow-overshoot backdrop-blur-xl lg:px-6"
          eyebrow="Brandtoon"
          title="BRANDTOON"
        />
      </header>

      <main className="mx-auto flex max-w-6xl flex-col gap-24 px-4 pb-16 pt-10 sm:px-6 lg:px-8 lg:pt-14">
        <section className="grid gap-10 lg:grid-cols-[1fr_1.05fr] lg:items-center">
          <div className="space-y-7">
            <Badge className="bg-ink px-4 py-2 uppercase tracking-section">
              Beta access now open
            </Badge>

            <div className="space-y-5">
              <h1 className="max-w-xl text-5xl font-black leading-[0.95] tracking-[-0.06em] text-ink sm:text-[3.75rem]">
                Your brand, <span className="text-coral">animated.</span>
              </h1>
              <p className="max-w-xl text-lg font-medium leading-8 text-ink/68">
                BRANDTOON turns your logo, mascot, and visual identity into
                animated sticker packs, expressive avatars, and motion-ready
                brand assets from one playful workflow.
              </p>
            </div>

            <div className="flex flex-wrap gap-4">
              <Button onClick={goToLogin} size="lg">
                Create your avatar
              </Button>
              <Button
                icon={<Play aria-hidden="true" className="size-4" />}
                size="lg"
                variant="secondary"
              >
                Watch showcase
              </Button>
            </div>

            <div className="flex flex-wrap gap-3">
              <ActionChip
                icon={<Sparkles aria-hidden="true" className="size-4" />}
              >
                2D + 3D outputs
              </ActionChip>
              <ActionChip icon={<Film aria-hidden="true" className="size-4" />}>
                Campaign-ready exports
              </ActionChip>
            </div>
          </div>

          <div className="relative">
            <Card className="landing-hero-card relative overflow-hidden bg-white p-4 sm:p-5">
              <div className="landing-hero-art relative overflow-hidden rounded-[2rem] p-4 sm:p-5">
                <div className="landing-orb absolute left-6 top-6 h-28 w-28" />
                <div className="landing-orb absolute bottom-10 right-8 h-36 w-36 bg-[radial-gradient(circle,_rgba(252,231,231,0.92)_0%,_rgba(252,231,231,0)_72%)]" />
                <div className="relative grid gap-3 sm:grid-cols-[1.15fr_0.85fr] sm:gap-4">
                  <article
                    className={`${HERO_CHARACTERS[0].className} sm:row-span-2`}
                  >
                    <Badge className="bg-white/92 text-ink">
                      {HERO_CHARACTERS[0].title}
                    </Badge>
                    <img
                      alt={HERO_CHARACTERS[0].imageAlt}
                      className={`landing-hero-character-image ${HERO_CHARACTERS[0].imageClassName}`}
                      loading="eager"
                      src={HERO_CHARACTERS[0].imageSrc}
                    />
                  </article>

                  <div className="grid gap-3 sm:grid-cols-1 sm:gap-4">
                    {HERO_CHARACTERS.slice(1, 3).map((character) => (
                      <article
                        className={character.className}
                        key={character.title}
                      >
                        <Badge className="bg-white/92 text-ink">
                          {character.title}
                        </Badge>
                        <img
                          alt={character.imageAlt}
                          className={`landing-hero-character-image ${character.imageClassName}`}
                          loading="eager"
                          src={character.imageSrc}
                        />
                      </article>
                    ))}
                  </div>

                  <div className="grid grid-cols-2 gap-3 sm:col-span-2 sm:gap-4">
                    {HERO_CHARACTERS.slice(3).map((character) => (
                      <article
                        className={character.className}
                        key={character.title}
                      >
                        <Badge className="bg-white/92 text-ink">
                          {character.title}
                        </Badge>
                        <img
                          alt={character.imageAlt}
                          className={`landing-hero-character-image ${character.imageClassName}`}
                          loading="lazy"
                          src={character.imageSrc}
                        />
                      </article>
                    ))}
                  </div>
                </div>
              </div>
            </Card>

            <img
              alt="Brandtoon decorative sticker swoosh"
              className="pointer-events-none absolute -bottom-14 left-0 hidden w-48 rotate-[-7deg] select-none md:block"
              loading="lazy"
              src="/images/landing/stickers/sticker-swoosh.png"
            />
          </div>
        </section>

        <SectionShell
          className="scroll-mt-28"
          description="From clean mascot systems to glossy 3D characters, BRANDTOON shows how one brand can expand into a full animated asset library without changing the shared system."
          eyebrow="The BRANDTOON playground"
          id="showcase"
          title="Diverse styles, infinite identity."
        >
          <div className="space-y-6">
            <div className="flex flex-wrap gap-3">
              <Badge variant="outline">2D Cartoon</Badge>
              <Badge variant="outline">3D Avatar</Badge>
              <Badge variant="outline">Motion Pack</Badge>
            </div>

            <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
              {SHOWCASE_CARDS.map((card) => (
                <Card className={card.className} key={card.title}>
                  <div className="flex h-full flex-col gap-6">
                    <div className="flex items-start justify-between gap-3">
                      <Badge className="w-fit bg-white/88 text-ink">
                        {card.label}
                      </Badge>
                    </div>
                    <div className="flex flex-1 items-center justify-center overflow-hidden rounded-[1.75rem] bg-white/46 p-4 sm:p-5">
                      <img
                        alt={card.imageAlt}
                        className={`w-full object-contain ${card.imageClassName}`}
                        loading="lazy"
                        src={card.imageSrc}
                      />
                    </div>
                    <div className="space-y-2">
                      <p className="text-2xl font-black tracking-tight">
                        {card.title}
                      </p>
                      <p className="text-sm font-medium opacity-80">
                        {card.description}
                      </p>
                    </div>
                  </div>
                </Card>
              ))}
            </div>
          </div>
        </SectionShell>

        <SectionShell
          className="scroll-mt-28"
          description="The BRANDTOON flow stays intentionally simple: capture your brand intent, tune the personality, then export motion-ready assets."
          eyebrow="The process"
          id="how-it-works"
          title="A guided creative journey."
        >
          <div className="grid gap-6 lg:grid-cols-3">
            <Card className="space-y-5 bg-white">
              <div className="space-y-3">
                <p className="foundation-section-eyebrow">
                  {PROCESS_STEPS[0].step}
                </p>
                <p className="text-xl font-black tracking-tight text-ink">
                  {PROCESS_STEPS[0].title}
                </p>
                <p className="foundation-body">{PROCESS_STEPS[0].body}</p>
              </div>
              <Input
                defaultValue="A playful burger mascot with high-energy motion accents"
                hint="Keep it short and specific."
                label="Describe your vision"
              />
            </Card>

            <Card className="space-y-5 bg-white lg:translate-y-10">
              <div className="space-y-3">
                <p className="foundation-section-eyebrow">
                  {PROCESS_STEPS[1].step}
                </p>
                <p className="text-xl font-black tracking-tight text-ink">
                  {PROCESS_STEPS[1].title}
                </p>
                <p className="foundation-body">{PROCESS_STEPS[1].body}</p>
              </div>
              <Slider
                defaultValue={30}
                label="Funny ↔ Serious"
                markers={[
                  { label: 'Funny', value: 'Lighter' },
                  { label: 'Balanced', value: 'Editorial' },
                  { label: 'Serious', value: 'Sharper' },
                ]}
                valueLabel="TUNED"
              />
            </Card>

            <Card className="space-y-5 bg-white">
              <div className="space-y-3">
                <p className="foundation-section-eyebrow">
                  {PROCESS_STEPS[2].step}
                </p>
                <p className="text-xl font-black tracking-tight text-ink">
                  {PROCESS_STEPS[2].title}
                </p>
                <p className="foundation-body">{PROCESS_STEPS[2].body}</p>
              </div>
              <div className="space-y-4 rounded-panel bg-surface p-4">
                <div className="flex aspect-video items-center justify-center rounded-[1.5rem] bg-[linear-gradient(135deg,_rgba(45,52,54,0.12)_0%,_rgba(45,52,54,0.02)_100%)]">
                  <span className="inline-flex size-14 items-center justify-center rounded-full bg-white/75 shadow-overshoot">
                    <Play aria-hidden="true" className="size-5 text-ink" />
                  </span>
                </div>
                <div className="flex gap-3">
                  <ActionChip className="flex-1 justify-center">MP4</ActionChip>
                  <ActionChip className="flex-1 justify-center">GIF</ActionChip>
                </div>
              </div>
            </Card>
          </div>
        </SectionShell>

        <section className="scroll-mt-28" id="pricing">
          <Card className="landing-pricing-card space-y-6 text-white">
            <div className="space-y-4 text-center">
              <p className="landing-pricing-kicker">On-demand pricing</p>
              <h2 className="text-4xl font-black tracking-tight">
                Pay per pack, ship faster.
              </h2>
              <p className="landing-pricing-copy mx-auto max-w-sm">
                Get a ready-to-use mascot asset pack for campaigns, social, and
                product marketing, without a subscription.
              </p>
            </div>

            <div className="rounded-panel bg-white/97 px-6 py-7 text-center text-ink shadow-overshoot">
              <div className="flex items-end justify-center gap-2">
                <span className="text-5xl font-black tracking-tight">$29</span>
                <span className="landing-pricing-meta">per pack</span>
              </div>
              <p className="landing-pricing-highlight">
                Includes 50+ custom variations
              </p>
            </div>

            <ul className="grid gap-3 sm:grid-cols-2">
              {PRICING_FEATURES.map((feature) => (
                <li className="landing-pricing-feature" key={feature}>
                  <CheckCircle2
                    aria-hidden="true"
                    className="size-4 text-coral"
                  />
                  {feature}
                </li>
              ))}
            </ul>

            <Button
              className="w-full justify-center bg-coral text-white"
              onClick={goToLogin}
              size="lg"
            >
              Purchase credits
            </Button>
          </Card>
        </section>

        <section>
          <Card className="landing-cta-card overflow-hidden bg-ink px-6 py-12 text-center text-white sm:px-10 sm:py-16">
            <div className="relative z-10 mx-auto max-w-2xl space-y-6">
              <Badge className="bg-white/10 px-4 py-2 text-white">
                Start today
              </Badge>
              <h2 className="text-4xl font-black tracking-tight sm:text-5xl">
                Bring your brand characters to life.
              </h2>
              <p className="mx-auto max-w-xl text-base font-medium leading-7 text-white/78">
                Create animated mascots and sticker-ready exports your team can
                publish right away.
              </p>
              <div className="flex justify-center">
                <Button
                  icon={<ArrowRight aria-hidden="true" className="size-4" />}
                  onClick={goToLogin}
                  size="lg"
                >
                  Launch generator
                </Button>
              </div>
            </div>
          </Card>
        </section>
      </main>

      <footer className="mx-auto mt-8 max-w-6xl rounded-t-[2.5rem] bg-surface px-4 py-10 sm:px-6 lg:px-8">
        <div className="flex flex-col gap-6 md:flex-row md:items-end md:justify-between">
          <div className="space-y-3">
            <p className="text-lg font-black tracking-tight text-ink">
              BRANDTOON
            </p>
            <p className="max-w-md text-[11px] font-extrabold uppercase tracking-section text-ink/45">
              © 2026 BRANDTOON. Built for animated brand storytelling.
            </p>
          </div>

          <nav
            aria-label="Footer navigation"
            className="flex flex-wrap gap-5 text-[11px] font-extrabold uppercase tracking-section text-ink/45"
          >
            <a href="/privacy">Privacy Policy</a>
            <a href="/terms">Terms of Service</a>
            <a href="/contact">Contact</a>
            <a href="/twitter">Twitter</a>
            <a href="/discord">Discord</a>
          </nav>
        </div>
      </footer>
    </div>
  )
}
