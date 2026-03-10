import { Link, Outlet } from "react-router-dom"
import { useAuth } from "@/contexts/AuthContext"
import { Button } from "@/components/ui/button"

export default function Layout() {
  const { isAuthenticated, signOut } = useAuth()

  return (
    <div className="flex min-h-screen flex-col">
      <header className="sticky top-0 z-10 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container flex h-14 items-center justify-between px-4">
          <Link to="/" className="font-semibold">
            Shop
          </Link>
          <nav className="flex items-center gap-2">
            {isAuthenticated ? (
              <>
                <span className="text-muted-foreground text-sm">Ви увійшли</span>
                <Button variant="ghost" size="sm" onClick={() => signOut()}>
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
