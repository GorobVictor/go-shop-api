import { useEffect, useState } from "react"
import { Link, Outlet } from "react-router-dom"
import { useAuth } from "@/contexts/AuthContext"
import { getProfile, type Profile } from "@/api/users"
import { Button } from "@/components/ui/button"

export default function Layout() {
  const { token, isAuthenticated, signOut } = useAuth()
  const [profile, setProfile] = useState<Profile | null>(null)

  useEffect(() => {
    if (!token) {
      setProfile(null)
      return
    }
    getProfile(token)
      .then(setProfile)
      .catch(() => setProfile(null))
  }, [token])

  return (
    <div className="flex min-h-screen flex-col bg-background">
      <header className="sticky top-0 z-10 border-b border-border/60 bg-background/90 shadow-sm backdrop-blur-sm">
        <div className="container mx-auto flex h-14 max-w-6xl items-center justify-between px-4">
          <Link
            to="/"
            className="text-lg font-bold tracking-tight text-foreground transition-colors hover:text-primary"
          >
            Shop
          </Link>
          <nav className="flex items-center gap-3">
            {isAuthenticated ? (
              <>
                <span className="hidden text-sm text-muted-foreground sm:inline">
                  {profile
                    ? `Привіт, ${profile.firstName}`
                    : "…"}
                </span>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => signOut()}
                  className="border-border"
                >
                  Вийти
                </Button>
              </>
            ) : (
              <>
                <Button variant="ghost" size="sm" asChild>
                  <Link to="/login">Увійти</Link>
                </Button>
                <Button size="sm" asChild>
                  <Link to="/register">Реєстрація</Link>
                </Button>
              </>
            )}
          </nav>
        </div>
      </header>
      <main className="flex-1">
        <Outlet />
      </main>
    </div>
  )
}
