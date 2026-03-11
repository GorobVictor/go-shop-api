import { useEffect, useState } from "react"
import { Link, Outlet } from "react-router-dom"
import { useAuth } from "@/contexts/AuthContext"
import { useCart } from "@/contexts/CartContext"
import { getProfile, type Profile } from "@/api/users"
import { Button } from "@/components/ui/button"

export default function Layout() {
  const { token, isAuthenticated, signOut } = useAuth()
  const { totalCount: cartCount } = useCart()
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
                <Link
                  to="/cart"
                  className="relative flex items-center gap-1.5 rounded-lg px-2 py-1.5 text-sm font-medium text-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
                >
                  Cart
                  {cartCount > 0 && (
                    <span className="flex h-5 min-w-5 items-center justify-center rounded-full bg-primary px-1.5 text-xs font-semibold text-primary-foreground">
                      {cartCount > 99 ? "99+" : cartCount}
                    </span>
                  )}
                </Link>
                <Link
                  to="/orders"
                  className="rounded-lg px-2 py-1.5 text-sm font-medium text-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
                >
                  Orders
                </Link>
                <span className="hidden text-sm text-muted-foreground sm:inline">
                  {profile ? `Hi, ${profile.firstName}` : "…"}
                </span>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => signOut()}
                  className="border-border"
                >
                  Sign out
                </Button>
              </>
            ) : (
              <>
                <Link
                  to="/cart"
                  className="relative flex items-center gap-1.5 rounded-lg px-2 py-1.5 text-sm font-medium text-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
                >
                  Cart
                  {cartCount > 0 && (
                    <span className="flex h-5 min-w-5 items-center justify-center rounded-full bg-primary px-1.5 text-xs font-semibold text-primary-foreground">
                      {cartCount > 99 ? "99+" : cartCount}
                    </span>
                  )}
                </Link>
                <Button variant="ghost" size="sm" asChild>
                  <Link to="/login">Sign in</Link>
                </Button>
                <Button size="sm" asChild>
                  <Link to="/register">Register</Link>
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
