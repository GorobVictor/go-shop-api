const API_BASE = "http://localhost:3000/api"

export type Profile = {
  ID: number
  firstName: string
  lastName: string
  email: string
  userRole: string
  createdAt: string
}

export async function getProfile(token: string): Promise<Profile> {
  const res = await fetch(`${API_BASE}/users/me`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!res.ok) throw new Error("Failed to load profile")
  return res.json() as Promise<Profile>
}
