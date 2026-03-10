import { Link } from "react-router-dom"
import { useAuth } from "@/contexts/AuthContext"
import { Button } from "@/components/ui/button"

export default function Home() {
  const { isAuthenticated, signOut } = useAuth()

  return (
    <div className="flex min-h-screen flex-col items-center justify-center gap-6 p-4">
      <h1 className="text-2xl font-semibold">Shop</h1>
      {isAuthenticated ? (
        <div className="flex gap-4">
          <Button variant="outline" asChild>
            <Link to="/login">Перейти до входу</Link>
          </Button>
          <Button variant="destructive" onClick={() => signOut()}>
            Вийти
          </Button>
        </div>
      ) : (
        <div className="flex gap-4">
          <Button asChild>
            <Link to="/login">Увійти</Link>
          </Button>
          <Button variant="outline" asChild>
            <Link to="/register">Реєстрація</Link>
          </Button>
        </div>
      )}
    </div>
  )
}
