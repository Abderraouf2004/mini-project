import { useState, type SyntheticEvent } from 'react'
import { signin, signup, type Session } from './api'

export default function Auth({ onAuth }: { onAuth: (s: Session) => void }) {
  const [mode, setMode] = useState<'signin' | 'signup'>('signin')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirm, setConfirm] = useState('')
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)
  const isUp = mode === 'signup'

  async function submit(e: SyntheticEvent) {
    e.preventDefault()
    setError('')
    if (password.length < 8) return setError('Password must be at least 8 characters.')
    if (isUp && password !== confirm) return setError('Passwords do not match.')
    setBusy(true)
    try {
      onAuth(await (isUp ? signup : signin)(email.trim(), password))
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Something went wrong')
      setBusy(false)
    }
  }

  return (
    <div className="auth">
      <div className="auth-side">
        <form className="auth-card" onSubmit={submit} noValidate>
          <div className="logo">◆ Helpdesk</div>
          <h1>{isUp ? 'Create your account' : 'Welcome back'}</h1>
          <p className="muted">
            {isUp ? 'Start tracking support tickets in minutes.' : 'Sign in to manage your tickets.'}
          </p>

          <label>
            Email
            <input type="email" placeholder="you@company.com" value={email}
              onChange={(e) => setEmail(e.target.value)} autoComplete="email" required />
          </label>
          <label>
            Password
            <input type="password" placeholder="At least 8 characters" value={password}
              onChange={(e) => setPassword(e.target.value)}
              autoComplete={isUp ? 'new-password' : 'current-password'} required />
          </label>
          {isUp && (
            <label>
              Confirm password
              <input type="password" placeholder="Repeat password" value={confirm}
                onChange={(e) => setConfirm(e.target.value)} autoComplete="new-password" required />
            </label>
          )}

          {error && <div className="alert">{error}</div>}
          <button className="btn primary block" disabled={busy}>
            {busy ? 'Please wait…' : isUp ? 'Sign up' : 'Sign in'}
          </button>

          <p className="muted center">
            {isUp ? 'Already have an account? ' : "Don't have an account? "}
            <button type="button" className="link"
              onClick={() => { setMode(isUp ? 'signin' : 'signup'); setError('') }}>
              {isUp ? 'Sign in' : 'Sign up'}
            </button>
          </p>
        </form>
      </div>
      <div className="auth-art">
        <div className="orb" />
        <blockquote>
          Every request answered.<br />Every ticket accounted for.
        </blockquote>
      </div>
    </div>
  )
}
