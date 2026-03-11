import {
  createContext,
  useCallback,
  useContext,
  useState,
  type ReactNode,
} from "react"
import { signIn as apiSignIn, signUp as apiSignUp } from "@/api/auth"

const TOKEN_KEY = "auth_token"

type AuthContextValue = {
  token: string | null
  isAuthenticated: boolean
  signIn: (email: string, password: string) => Promise<void>
  signUp: (
    firstName: string,
    lastName: string,
    email: string,
    password: string
  ) => Promise<void>
  signOut: () => void
  error: string | null
  clearError: () => void
}

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(() =>
    localStorage.getItem(TOKEN_KEY)
  )
  const [error, setError] = useState<string | null>(null)

  const signIn = useCallback(async (email: string, password: string) => {
    setError(null)
    try {
      const { token: t } = await apiSignIn({ email, password })
      localStorage.setItem(TOKEN_KEY, t)
      setToken(t)
    } catch (e) {
      const msg = e instanceof Error ? e.message : "Sign in failed"
      setError(msg)
      throw e
    }
  }, [])

  const signUp = useCallback(
    async (
      firstName: string,
      lastName: string,
      email: string,
      password: string
    ) => {
      setError(null)
      try {
        const { token: t } = await apiSignUp({
          firstName,
          lastName,
          email,
          password,
        })
        localStorage.setItem(TOKEN_KEY, t)
        setToken(t)
      } catch (e) {
        const msg = e instanceof Error ? e.message : "Registration failed"
        setError(msg)
        throw e
      }
    },
    []
  )

  const signOut = useCallback(() => {
    localStorage.removeItem(TOKEN_KEY)
    setToken(null)
  }, [])

  const clearError = useCallback(() => setError(null), [])

  const value: AuthContextValue = {
    token,
    isAuthenticated: !!token,
    signIn,
    signUp,
    signOut,
    error,
    clearError,
  }

  return (
    <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error("useAuth must be used within AuthProvider")
  return ctx
}
