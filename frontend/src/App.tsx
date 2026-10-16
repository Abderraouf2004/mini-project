import { useEffect, useState } from 'react'
import { getSession, onAuthLost, type Session } from './api'
import Auth from './Auth'
import Dashboard from './Dashboard'

export default function App() {
  const [session, setSession] = useState<Session | null>(getSession)
  useEffect(() => onAuthLost(() => setSession(null)), [])
  return session ? (
    <Dashboard session={session} onLogout={() => setSession(null)} />
  ) : (
    <Auth onAuth={setSession} />
  )
}
