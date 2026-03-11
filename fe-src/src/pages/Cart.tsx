import { useState } from "react"
import { Link } from "react-router-dom"
import { useAuth } from "@/contexts/AuthContext"
import { useCart } from "@/contexts/CartContext"
import { createPayment } from "@/api/payment"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"

export default function Cart() {
  const { token, isAuthenticated } = useAuth()
  const { items, setQuantity, removeItem, productIdsToQuantity, clearCart } = useCart()
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  if (!isAuthenticated) {
    return (
      <div className="container mx-auto max-w-2xl px-4 py-12 text-center">
        <p className="text-muted-foreground">
          Увійдіть, щоб переглядати кошик та оформляти замовлення.
        </p>
        <Button asChild className="mt-4">
          <Link to="/login">Увійти</Link>
        </Button>
      </div>
    )
  }

  const totalCents = items.reduce(
    (acc, i) => acc + (i.price - i.discount) * i.quantity,
    0
  )
  const totalUah = (totalCents / 100).toFixed(2)

  async function handleCheckout() {
    if (items.length === 0 || !token) return
    setError(null)
    setLoading(true)
    try {
      const { link } = await createPayment(token, {
        products: Object.fromEntries(
          Object.entries(productIdsToQuantity()).map(([k, v]) => [String(k), v])
        ),
      })
      clearCart()
      window.location.href = link
    } catch (e) {
      setError(e instanceof Error ? e.message : "Помилка оформлення")
    } finally {
      setLoading(false)
    }
  }

  if (items.length === 0 && !loading) {
    return (
      <div className="container mx-auto max-w-2xl px-4 py-12">
        <h1 className="mb-4 text-2xl font-semibold">Кошик</h1>
        <Card className="rounded-xl">
          <CardContent className="flex flex-col items-center justify-center py-12">
            <p className="text-muted-foreground">Кошик порожній</p>
            <Button asChild variant="outline" className="mt-4">
              <Link to="/">Перейти до каталогу</Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="container mx-auto max-w-2xl px-4 py-8">
      <h1 className="mb-6 text-2xl font-semibold">Кошик</h1>
      {error && (
        <div className="mb-4 rounded-lg bg-destructive/10 px-3 py-2 text-sm text-destructive">
          {error}
        </div>
      )}
      <div className="space-y-4">
        {items.map((item) => {
          const finalPrice = ((item.price - item.discount) / 100).toFixed(2)
          const lineTotal = ((item.price - item.discount) * item.quantity / 100).toFixed(2)
          return (
            <Card key={item.productId} className="rounded-xl">
              <CardHeader className="flex flex-row items-start justify-between space-y-0 pb-2">
                <CardTitle className="text-base">{item.name}</CardTitle>
                <Button
                  variant="ghost"
                  size="sm"
                  className="text-muted-foreground -mr-2"
                  onClick={() => removeItem(item.productId)}
                >
                  Видалити
                </Button>
              </CardHeader>
              <CardContent className="space-y-2">
                <div className="flex items-center gap-4">
                  <div className="flex items-center gap-2">
                    <label className="text-sm text-muted-foreground">Кількість</label>
                    <Input
                      type="number"
                      min={1}
                      value={item.quantity}
                      onChange={(e) => {
                        const v = parseInt(e.target.value, 10)
                        setQuantity(item.productId, Number.isNaN(v) ? 1 : v)
                      }}
                      className="h-9 w-20 rounded-lg"
                    />
                  </div>
                  <span className="text-muted-foreground">
                    {finalPrice} грн × {item.quantity} = {lineTotal} грн
                  </span>
                </div>
              </CardContent>
            </Card>
          )
        })}
      </div>
      <Card className="mt-6 rounded-xl">
        <CardContent className="flex flex-col gap-4 pt-6">
          <div className="flex justify-between text-lg font-semibold">
            <span>Разом</span>
            <span>{totalUah} грн</span>
          </div>
          <Button
            className="w-full rounded-lg"
            size="lg"
            disabled={loading}
            onClick={handleCheckout}
          >
            {loading ? "Перенаправлення на оплату…" : "Оформити замовлення"}
          </Button>
          <p className="text-center text-xs text-muted-foreground">
            Оплата через Stripe. Після натискання ви будете перенаправлені на безпечну сторінку оплати.
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
