import { useCallback, useEffect, useMemo, useState, type SyntheticEvent } from 'react'
import {
  createTicket, deleteTicket, listAllTickets, signout, updateTicket,
  type Session, type Status, type Ticket,
} from './api'

const LABEL: Record<Status, string> = { OPEN: 'Open', IN_PROGRESS: 'In progress', CLOSED: 'Closed' }
const STATUSES: Status[] = ['OPEN', 'IN_PROGRESS', 'CLOSED']
const PAGE = 8
type Modal = { ticket?: Ticket } | null

const fmt = (s: string) => new Date(s).toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' })
const msg = (e: unknown) => (e instanceof Error ? e.message : 'Something went wrong')

export default function Dashboard({ session, onLogout }: { session: Session; onLogout: () => void }) {
  const [tickets, setTickets] = useState<Ticket[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [query, setQuery] = useState('')
  const [filter, setFilter] = useState<'ALL' | Status>('ALL')
  const [page, setPage] = useState(1)
  const [modal, setModal] = useState<Modal>(null)

  const load = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      setTickets(await listAllTickets())
    } catch (e) {
      setError(msg(e))
    } finally {
      setLoading(false)
    }
  }, [])
  useEffect(() => { void load() }, [load])

  const counts = useMemo(() => {
    const c = { OPEN: 0, IN_PROGRESS: 0, CLOSED: 0 }
    tickets.forEach((t) => c[t.status]++)
    return c
  }, [tickets])

  const visible = useMemo(() => {
    const q = query.trim().toLowerCase()
    return tickets.filter(
      (t) => (filter === 'ALL' || t.status === filter) &&
        (!q || t.title.toLowerCase().includes(q) || t.description.toLowerCase().includes(q)),
    )
  }, [tickets, query, filter])

  const pages = Math.max(1, Math.ceil(visible.length / PAGE))
  const current = Math.min(page, pages)
  const rows = visible.slice((current - 1) * PAGE, current * PAGE)
  const rate = tickets.length ? Math.round((counts.CLOSED / tickets.length) * 100) : 0

  async function changeStatus(t: Ticket, status: Status) {
    try {
      const u = await updateTicket(t.id, { status })
      setTickets((all) => all.map((x) => (x.id === u.id ? u : x)))
    } catch (e) { setError(msg(e)) }
  }
  async function remove(t: Ticket) {
    if (!window.confirm(`Delete "${t.title}"?`)) return
    try {
      await deleteTicket(t.id)
      setTickets((all) => all.filter((x) => x.id !== t.id))
    } catch (e) { setError(msg(e)) }
  }
  async function logout() {
    await signout()
    onLogout()
  }

  return (
    <div className="shell">
      <aside className="sidebar">
        <div className="logo">◆ Helpdesk</div>
        <nav>
          <a className="nav active">▦ Tickets</a>
        </nav>
        <div className="user">
          <div className="avatar">{session.email[0]?.toUpperCase()}</div>
          <div className="user-mail" title={session.email}>{session.email}</div>
          <button className="btn ghost sm" onClick={logout}>Sign out</button>
        </div>
      </aside>

      <main className="main">
        <header className="top">
          <div>
            <h1>Tickets</h1>
            <p className="muted">Track and resolve your support requests.</p>
          </div>
          <button className="btn primary" onClick={() => setModal({})}>+ New ticket</button>
        </header>

        {error && <div className="alert" onClick={() => setError('')}>{error} <span className="muted">(dismiss)</span></div>}

        <section className="stats">
          <Stat label="Total" value={tickets.length} />
          <Stat label="Open" value={counts.OPEN} tone="open" />
          <Stat label="In progress" value={counts.IN_PROGRESS} tone="prog" />
          <Stat label="Closed" value={counts.CLOSED} tone="done" />
          <div className="card stat">
            <span className="muted">Resolution rate</span>
            <strong>{rate}%</strong>
            <div className="bar"><div style={{ width: `${rate}%` }} /></div>
          </div>
        </section>

        <section className="card">
          <div className="toolbar">
            <input className="search" placeholder="Search tickets…" value={query}
              onChange={(e) => { setQuery(e.target.value); setPage(1) }} />
            <div className="tabs">
              {(['ALL', ...STATUSES] as const).map((f) => (
                <button key={f} className={filter === f ? 'tab on' : 'tab'}
                  onClick={() => { setFilter(f); setPage(1) }}>
                  {f === 'ALL' ? 'All' : LABEL[f]}
                </button>
              ))}
            </div>
          </div>

          {loading ? (
            <div className="empty">Loading tickets…</div>
          ) : rows.length === 0 ? (
            <div className="empty">
              {tickets.length === 0 ? 'No tickets yet. Create your first one.' : 'No tickets match your filters.'}
            </div>
          ) : (
            <div className="scroll">
              <table>
                <thead>
                  <tr><th>Ticket</th><th>Status</th><th>Created</th><th>Updated</th><th /></tr>
                </thead>
                <tbody>
                  {rows.map((t) => (
                    <tr key={t.id}>
                      <td>
                        <div className="t-title">{t.title}</div>
                        <div className="t-desc">#{t.id.slice(0, 8)} · {t.description}</div>
                      </td>
                      <td>
                        <select className={`pill ${t.status}`} value={t.status}
                          onChange={(e) => void changeStatus(t, e.target.value as Status)}>
                          {STATUSES.map((s) => <option key={s} value={s}>{LABEL[s]}</option>)}
                        </select>
                      </td>
                      <td className="muted">{fmt(t.created_at)}</td>
                      <td className="muted">{fmt(t.updated_at)}</td>
                      <td className="actions">
                        <button className="btn ghost sm" onClick={() => setModal({ ticket: t })}>Edit</button>
                        <button className="btn danger sm" onClick={() => void remove(t)}>Delete</button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}

          {pages > 1 && (
            <div className="pager">
              <span className="muted">{visible.length} tickets</span>
              <div>
                <button className="btn ghost sm" disabled={current === 1} onClick={() => setPage(current - 1)}>← Prev</button>
                <span className="muted"> {current} / {pages} </span>
                <button className="btn ghost sm" disabled={current === pages} onClick={() => setPage(current + 1)}>Next →</button>
              </div>
            </div>
          )}
        </section>
      </main>

      {modal && (
        <TicketModal
          ticket={modal.ticket}
          onClose={() => setModal(null)}
          onSaved={(t) => {
            setTickets((all) => (all.some((x) => x.id === t.id) ? all.map((x) => (x.id === t.id ? t : x)) : [t, ...all]))
            setModal(null)
          }}
        />
      )}
    </div>
  )
}

function Stat({ label, value, tone }: { label: string; value: number; tone?: string }) {
  return (
    <div className={`card stat ${tone ?? ''}`}>
      <span className="muted">{label}</span>
      <strong>{value}</strong>
    </div>
  )
}

function TicketModal({ ticket, onClose, onSaved }: {
  ticket?: Ticket
  onClose: () => void
  onSaved: (t: Ticket) => void
}) {
  const [title, setTitle] = useState(ticket?.title ?? '')
  const [description, setDescription] = useState(ticket?.description ?? '')
  const [status, setStatus] = useState<Status>(ticket?.status ?? 'OPEN')
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)

  async function submit(e: SyntheticEvent) {
    e.preventDefault()
    if (!title.trim() || !description.trim()) return setError('Title and description are required.')
    setBusy(true)
    setError('')
    try {
      onSaved(
        ticket
          ? await updateTicket(ticket.id, { title, description, status })
          : await createTicket({ title, description }),
      )
    } catch (err) {
      setError(msg(err))
      setBusy(false)
    }
  }

  return (
    <div className="overlay" onMouseDown={onClose}>
      <form className="modal" onSubmit={submit} onMouseDown={(e) => e.stopPropagation()}>
        <h2>{ticket ? 'Edit ticket' : 'New ticket'}</h2>
        <label>
          Title
          <input value={title} maxLength={200} onChange={(e) => setTitle(e.target.value)} autoFocus />
        </label>
        <label>
          Description <span className="muted">{description.length}/500</span>
          <textarea rows={5} value={description} maxLength={500} onChange={(e) => setDescription(e.target.value)} />
        </label>
        {ticket && (
          <label>
            Status
            <select value={status} onChange={(e) => setStatus(e.target.value as Status)}>
              {STATUSES.map((s) => <option key={s} value={s}>{LABEL[s]}</option>)}
            </select>
          </label>
        )}
        {error && <div className="alert">{error}</div>}
        <div className="modal-actions">
          <button type="button" className="btn ghost" onClick={onClose}>Cancel</button>
          <button className="btn primary" disabled={busy}>{busy ? 'Saving…' : ticket ? 'Save changes' : 'Create ticket'}</button>
        </div>
      </form>
    </div>
  )
}
