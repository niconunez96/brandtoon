import { create } from 'zustand'

type GreetingState = {
  name: string
  setName: (name: string) => void
}

export const useGreetingStore = create<GreetingState>((set) => ({
  name: 'world',
  setName: (name) => set({ name }),
}))
