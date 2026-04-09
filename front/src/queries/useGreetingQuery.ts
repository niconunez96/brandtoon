import { useQuery } from '@tanstack/react-query'
import { fetchGreeting } from '../services/greeting.api'

export const useGreetingQuery = (name: string) =>
  useQuery({
    queryKey: ['greeting', name],
    queryFn: () => fetchGreeting(name),
    enabled: Boolean(name),
  })
