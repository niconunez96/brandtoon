import { useGreetingQuery } from './queries/useGreetingQuery'
import { useGreetingStore } from './store/greeting.store'

export function App() {
  const name = useGreetingStore((s) => s.name)
  const setName = useGreetingStore((s) => s.setName)
  const { data, isLoading, error, refetch, isFetching } = useGreetingQuery(name)

  return (
    <main className="layout">
      <h1>Brandtoon FE + BE Hello World</h1>

      <div className="controls">
        <label htmlFor="name">Name</label>
        <input
          id="name"
          value={name}
          onChange={(event) => setName(event.target.value)}
          placeholder="world"
        />
        <button type="button" onClick={() => refetch()} disabled={isFetching}>
          {isFetching ? 'Loading...' : 'Refresh greeting'}
        </button>
      </div>

      {isLoading && <p>Loading greeting...</p>}
      {error && <p className="error">Could not fetch greeting.</p>}
      {data && <p className="message">{data.message}</p>}
    </main>
  )
}
