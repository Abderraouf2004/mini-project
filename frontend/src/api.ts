export type Status = 'OPEN' | 'IN_PROGRESS' | 'CLOSED'

export interface Ticket {
  id: string
  title: string
  description: string
  status: Status
  user_id: string
  created_at: string
  updated_at: string
}

export interface Session {
  accessToken: string
  refreshToken: string
  email: string
}

export class ApiError extends Error {
  status: number
  constructor(status: number, message: string) {
    super(message)
    this.status = status
  }
}

const KEY = 'helpdesk_session'
let lostCb: (() => void) | null = null

export function getSession(): Session | null {
  try {
    const raw = localStorage.getItem(KEY)
    return raw ? (JSON.parse(raw) as Session) : null
  } catch {
    return null
  }
}
function saveSession(s: Session | null) {
  if (s) localStorage.setItem(KEY, JSON.stringify(s))
  else localStorage.removeItem(KEY)
}
export function onAuthLost(cb: () => void) {
  lostCb = cb
  return () => {
    lostCb = null
  }
}


async function toError(res: Response): Promise<ApiError> {
  const text = await res.text()
  let msg = text.trim() || res.statusText
  try {
    const j = JSON.parse(text)
    msg = j.message || j.error || msg
    if (j.details) msg += `: ${j.details}`
  } catch {
    /* plain text */
  }
  if (res.status === 429) {
    const wait = res.headers.get('Retry-After')
    msg = `Too many requests. Try again${wait ? ` in ${wait}s` : ' shortly'}.`
  }
  return new ApiError(res.status, msg)
}

async function send(path: string, method: string, body?: unknown, token?: string) {
  let res: Response
  try {
    res = await fetch(`/api${path}`, {
      method,
      headers: {
        ...(body !== undefined ? { 'Content-Type': 'application/json' } : {}),
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
      body: body !== undefined ? JSON.stringify(body) : undefined,
    })
  } catch {
    throw new ApiError(0, 'Cannot reach the server. Is the backend running on :8080?')
  }
  if (!res.ok) throw await toError(res)
  return res
}

async function refresh(s: Session): Promise<Session | null> {
  try {
    const res = await send('/auth/refresh', 'POST', undefined, s.refreshToken)
    const json = await res.json()
    const next = { ...s, accessToken: json.data.accessToken as string }
    saveSession(next)
    return next
  } catch {
    return null
  }
}

async function authed<T>(path: string, method = 'GET', body?: unknown): Promise<T> {
  let s = getSession()
  if (!s) throw new ApiError(401, 'Not signed in')
  try {
    return (await (await send(path, method, body, s.accessToken)).json()) as T
  } catch (e) {
    if (!(e instanceof ApiError) || e.status !== 401) throw e
    s = await refresh(s)
    if (!s) {
      saveSession(null)
      lostCb?.()
      throw new ApiError(401, 'Session expired. Please sign in again.')
    }
    return (await (await send(path, method, body, s.accessToken)).json()) as T
  }
}

async function authenticate(path: string, email: string, password: string) {
  const res = await send(path, 'POST', { email, password })
  const { data } = await res.json()
  const s: Session = { accessToken: data.accessToken, refreshToken: data.refreshToken, email }
  saveSession(s)
  return s
}

export const signin = (email: string, password: string) => authenticate('/auth/signin', email, password)
export const signup = (email: string, password: string) => authenticate('/auth/signup', email, password)

export async function signout() {
  const s = getSession()
  saveSession(null)
  if (!s) return
  await Promise.allSettled([
    send('/auth/signout', 'POST', undefined, s.accessToken),
    send('/auth/signout', 'POST', undefined, s.refreshToken),
  ])
}

export async function listAllTickets(): Promise<Ticket[]> {
  const all: Ticket[] = []
  for (let page = 1; ; page++) {
    const r = await authed<{ data: Ticket[]; pagination: { total_pages: number } }>(
      `/tickets?page=${page}&page_size=100`,
    )
    all.push(...r.data)
    if (page >= r.pagination.total_pages) return all
  }
}

type Fields = { title: string; description: string; status?: Status }
export const createTicket = (f: Fields) =>
  authed<{ data: Ticket }>('/tickets', 'POST', { title: f.title, description: f.description }).then((r) => r.data)
export const updateTicket = (id: string, f: Partial<Fields>) =>
  authed<{ data: Ticket }>(`/tickets/${id}`, 'PUT', f).then((r) => r.data)
export const deleteTicket = (id: string) => authed<unknown>(`/tickets/${id}`, 'DELETE')
