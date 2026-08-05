import { beforeEach, afterEach, describe, expect, it, vi } from 'vitest'

const BASE_URL = 'http://localhost:8080/api'

beforeEach(() => {
  vi.stubEnv('VITE_API_BASE_URL', BASE_URL)
  localStorage.clear()
  vi.restoreAllMocks()
})

let fetchMock
beforeEach(() => {
  fetchMock = vi.fn()
  vi.stubGlobal('fetch', fetchMock)
})

afterEach(() => {
  vi.unstubAllGlobals()
  vi.unstubAllEnvs()
  vi.resetModules()
})

function mockResponse(status, body) {
  return Promise.resolve({
    ok: status >= 200 && status < 300,
    status,
    json: () => Promise.resolve(body),
  })
}

describe('api.get', () => {
  it('sends GET to the correct URL', async () => {
    fetchMock.mockResolvedValueOnce(mockResponse(200, { movies: [] }))
    const { api } = await import('../api/client')

    const result = await api.get('/movies')

    expect(fetchMock).toHaveBeenCalledTimes(1)
    const [url, config] = fetchMock.mock.calls[0]
    expect(url).toBe(`${BASE_URL}/movies`)
    expect(config.method).toBe('GET')
    expect(result).toEqual({ movies: [] })
  })

  it('omits Content-Type header', async () => {
    fetchMock.mockResolvedValueOnce(mockResponse(200, {}))
    const { api } = await import('../api/client')

    await api.get('/movies')

    const [, config] = fetchMock.mock.calls[0]
    expect(config.headers['Content-Type']).toBeUndefined()
  })
})

describe('api.post', () => {
  it('sends POST with JSON body and Content-Type', async () => {
    fetchMock.mockResolvedValueOnce(mockResponse(201, { id: '1' }))
    const { api } = await import('../api/client')
    const body = { title: 'Inception' }

    const result = await api.post('/movies', body)

    const [, config] = fetchMock.mock.calls[0]
    expect(config.method).toBe('POST')
    expect(config.body).toBe(JSON.stringify(body))
    expect(config.headers['Content-Type']).toBe('application/json')
    expect(result).toEqual({ id: '1' })
  })
})

describe('api.put', () => {
  it('sends PUT with body', async () => {
    fetchMock.mockResolvedValueOnce(mockResponse(200, { id: '1' }))
    const { api } = await import('../api/client')
    const body = { title: 'Updated' }

    await api.put('/movies/1', body)

    const [, config] = fetchMock.mock.calls[0]
    expect(config.method).toBe('PUT')
    expect(config.body).toBe(JSON.stringify(body))
  })
})

describe('api.delete', () => {
  it('sends DELETE', async () => {
    fetchMock.mockResolvedValueOnce(mockResponse(200, { message: 'ok' }))
    const { api } = await import('../api/client')

    await api.delete('/movies/1')

    const [, config] = fetchMock.mock.calls[0]
    expect(config.method).toBe('DELETE')
  })
})

describe('JWT interceptor', () => {
  it('attaches Authorization header when token exists', async () => {
    localStorage.setItem('token', 'test-jwt-token')
    fetchMock.mockResolvedValueOnce(mockResponse(200, {}))
    const { api } = await import('../api/client')

    await api.get('/movies')

    const [, config] = fetchMock.mock.calls[0]
    expect(config.headers['Authorization']).toBe('Bearer test-jwt-token')
  })

  it('does not attach Authorization header when token is absent', async () => {
    fetchMock.mockResolvedValueOnce(mockResponse(200, {}))
    const { api } = await import('../api/client')

    await api.get('/movies')

    const [, config] = fetchMock.mock.calls[0]
    expect(config.headers['Authorization']).toBeUndefined()
  })
})

describe('401 handling', () => {
  it('clears localStorage on 401', async () => {
    localStorage.setItem('token', 'expired-token')
    localStorage.setItem('user', '{"name":"test"}')
    vi.stubGlobal('location', { href: '' })
    fetchMock.mockResolvedValueOnce(mockResponse(401, { error: 'Token expired' }))
    const { api } = await import('../api/client')

    await expect(api.get('/profile')).rejects.toThrow('Token expired')

    expect(localStorage.getItem('token')).toBeNull()
    expect(localStorage.getItem('user')).toBeNull()
  })

  it('is overridden by setOnUnauthorized', async () => {
    localStorage.setItem('token', 'expired-token')
    fetchMock.mockResolvedValueOnce(mockResponse(401, {}))
    const { api, setOnUnauthorized } = await import('../api/client')

    const customHandler = vi.fn()
    setOnUnauthorized(customHandler)

    await expect(api.get('/profile')).rejects.toThrow('Unauthorized')

    expect(customHandler).toHaveBeenCalledTimes(1)
    expect(localStorage.getItem('token')).toBe('expired-token')
  })
})

describe('error responses', () => {
  it('throws with body error message for non-ok responses', async () => {
    fetchMock.mockResolvedValueOnce(mockResponse(400, { error: 'Validation failed' }))
    const { api } = await import('../api/client')

    await expect(api.post('/register', {})).rejects.toThrow('Validation failed')
  })

  it('throws with status code when body has no error field', async () => {
    fetchMock.mockResolvedValueOnce(mockResponse(500, {}))
    const { api } = await import('../api/client')

    await expect(api.get('/movies')).rejects.toThrow('Request failed (500)')
  })
})

describe('network errors', () => {
  it('throws Network error on fetch rejection', async () => {
    fetchMock.mockRejectedValueOnce(new TypeError('Failed to fetch'))
    const { api } = await import('../api/client')

    await expect(api.get('/movies')).rejects.toThrow('Network error')
  })
})
