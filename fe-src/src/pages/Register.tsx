import { useState } from "react"
import { Link, useNavigate } from "react-router-dom"
import { useAuth } from "@/contexts/AuthContext"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

export default function Register() {
  const navigate = useNavigate()
  const { signUp, error, clearError } = useAuth()
  const [firstName, setFirstName] = useState("")
  const [lastName, setLastName] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    clearError()
    setLoading(true)
    try {
      await signUp(firstName, lastName, email, password)
      navigate("/", { replace: true })
    } catch {
      // error already set in context
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="relative flex min-h-[calc(100vh-3.5rem)] items-center justify-center bg-gradient-to-b from-accent/30 to-background px-4 py-12">
      <Card className="w-full max-w-sm rounded-2xl border-border/80 shadow-xl">
        <CardHeader className="space-y-1 text-center">
          <CardTitle className="text-2xl">Реєстрація</CardTitle>
          <CardDescription>
            Створіть акаунт, заповнивши форму нижче
          </CardDescription>
        </CardHeader>
        <form onSubmit={handleSubmit}>
          <CardContent className="space-y-4">
            {error && (
              <div className="rounded-lg bg-destructive/10 px-3 py-2.5 text-sm text-destructive">
                {error}
              </div>
            )}
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="firstName">Ім'я</Label>
                <Input
                  id="firstName"
                  placeholder="Іван"
                  value={firstName}
                  onChange={(e) => setFirstName(e.target.value)}
                  required
                  autoComplete="given-name"
                  className="rounded-lg"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="lastName">Прізвище</Label>
                <Input
                  id="lastName"
                  placeholder="Петренко"
                  value={lastName}
                  onChange={(e) => setLastName(e.target.value)}
                  required
                  autoComplete="family-name"
                  className="rounded-lg"
                />
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="you@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                autoComplete="email"
                className="rounded-lg"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">Пароль</Label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                autoComplete="new-password"
                className="rounded-lg"
              />
            </div>
          </CardContent>
          <CardFooter className="flex flex-col gap-4">
            <Button
              type="submit"
              className="w-full rounded-lg"
              disabled={loading}
            >
              {loading ? "Реєстрація…" : "Зареєструватися"}
            </Button>
            <p className="text-center text-sm text-muted-foreground">
              Вже є акаунт?{" "}
              <Link
                to="/login"
                className="font-medium text-primary underline-offset-4 hover:underline"
              >
                Увійти
              </Link>
            </p>
          </CardFooter>
        </form>
      </Card>
    </div>
  )
}
