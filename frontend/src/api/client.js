const BASE_URL = import.meta.env.VITE_API_BASE_URL

let onUnauthorized = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  window.location.href = '/login'
}

export function setOnUnauthorized(callback) {
  onUnauthorized = callback
}

async function request(endpoint, options = {}) {
  const token = localStorage.getItem('token')
  const headers = { ...options.headers }

  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  if (!options.body && !headers['Content-Type']) {
    delete headers['Content-Type']
  } else if (options.body && !headers['Content-Type']) {
    headers['Content-Type'] = 'application/json'
  }

  const config = {
    ...options,
    headers,
  }

  let response
  try {
    response = await fetch(`${BASE_URL}${endpoint}`, config)
  } catch {
    throw new Error('Network error')
  }

  if (response.status === 401) {
    onUnauthorized()
    const data = await response.json().catch(() => ({}))
    throw new Error(data.error || 'Unauthorized')
  }

  if (!response.ok) {
    const data = await response.json().catch(() => ({}))
    throw new Error(data.error || `Request failed (${response.status})`)
  }

  return response.json()
}

export const api = {
  get(endpoint) {
    return request(endpoint, { method: 'GET' })
  },
  post(endpoint, body) {
    return request(endpoint, { method: 'POST', body: JSON.stringify(body) })
  },
  put(endpoint, body) {
    return request(endpoint, { method: 'PUT', body: JSON.stringify(body) })
  },
  delete(endpoint) {
    return request(endpoint, { method: 'DELETE' })
  },
}
