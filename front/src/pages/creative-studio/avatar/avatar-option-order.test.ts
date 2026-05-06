// @vitest-environment node

import { describe, expect, it } from 'vitest'
import { orderAvatarOptionsBySelection } from './avatar-option-order'

describe('orderAvatarOptionsBySelection', () => {
  it('preserves backend order even when one option is selected', () => {
    const result = orderAvatarOptionsBySelection([
      { id: 'option-1', href: '/1.png', selected: false },
      { id: 'option-2', href: '/2.png', selected: true },
      { id: 'option-3', href: '/3.png', selected: false },
    ])

    expect(result.map((option) => option.id)).toEqual([
      'option-1',
      'option-2',
      'option-3',
    ])
  })

  it('preserves the original order when no option is selected', () => {
    const result = orderAvatarOptionsBySelection([
      { id: 'option-1', href: '/1.png', selected: false },
      { id: 'option-2', href: '/2.png', selected: false },
    ])

    expect(result.map((option) => option.id)).toEqual(['option-1', 'option-2'])
  })
})
