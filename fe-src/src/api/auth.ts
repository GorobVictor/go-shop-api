const API_BASE = "http://localhost:3000/api"

export type SignInBody = {
  email: string
  password: string
}

export type SignUpBody = {
  firstName: string
  lastName: string
  email: string
  password: string
}

export type TokenResponse = {
  token: string
}

export async function signIn(body: SignInBody): Promise<TokenResponse> {
  const res = await fetch(`${API_BASE}/users/signin`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || "Sign in failed")
  }
  return res.json() as Promise<TokenResponse>
}

export async function signUp(body: SignUpBody): Promise<TokenResponse> {
  const res = await fetch(`${API_BASE}/users/signup`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || "Registration failed")
  }
  return res.json() as Promise<TokenResponse>
}
