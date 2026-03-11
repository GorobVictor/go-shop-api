import { useEffect, useState } from "react"
import { Link } from "react-router-dom"
import { useAuth } from "@/contexts/AuthContext"
import { getReceipts, type Receipt } from "@/api/payment"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

const PAGE_SIZE = 10

const STATUS_LABEL: Record<string, string> = {
  pending: "Очікує",
  succeeded: "Оплачено",
  failed: "Помилка",
  canceled: "Скасовано",
  refunded: "Повернено",
}

function formatDate(s: string) {
  try {
    return new Date(s).toLocaleDateString("uk-UA", {
      day: "2-digit",
      month: "2-digit",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    })
  } catch {
    return s
  }
}

export default function Orders() {
  const { token, isAuthenticated } = useAuth()
  const [data, setData] = useState<Awaited<ReturnType<typeof getReceipts>> | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [page, setPage] = useState(0)

  useEffect(() => {
    if (!token) {
      setLoading(false)
      return
    }
    setLoading(true)
    getReceipts(token, PAGE_SIZE, page * PAGE_SIZE)
      .then(setData)
      .catch((e) => setError(e instanceof Error ? e.message : "Помилка"))
      .finally(() => setLoading(false))
  }, [token, page])

  if (!isAuthenticated) {
    return (
      <div className="container mx-auto max-w-2xl px-4 py-12 text-center">
        <p className="text-muted-foreground">
          Увійдіть, щоб переглядати історію покупок.
        </p>
        <Button asChild className="mt-4">
          <Link to="/login">Увійти</Link>
        </Button>
      </div>
    )
  }

  if (error) {
    return (
      <div className="container mx-auto max-w-2xl px-4 py-8">
        <div className="rounded-xl border border-destructive/30 bg-destructive/5 p-4 text-destructive">
          {error}
        </div>
      </div>
    )
  }

  const totalPages = data ? Math.ceil(data.total / PAGE_SIZE) : 0

  return (
    <div className="container mx-auto max-w-2xl px-4 py-8">
      <h1 className="mb-6 text-2xl font-semibold">Історія покупок</h1>
      {loading ? (
        <div className="space-y-4">
          {[1, 2, 3].map((i) => (
            <Card key={i} className="rounded-xl">
              <CardHeader>
                <div className="h-5 w-1/3 animate-pulse rounded bg-muted" />
                <div className="h-4 w-1/4 animate-pulse rounded bg-muted" />
              </CardHeader>
              <CardContent>
                <div className="h-16 animate-pulse rounded bg-muted" />
              </CardContent>
            </Card>
          ))}
        </div>
      ) : data && data.receipts.length > 0 ? (
        <>
          <div className="space-y-4">
            {data.receipts.map((receipt) => (
              <ReceiptCard key={receipt.id} receipt={receipt} />
            ))}
          </div>
          {totalPages > 1 && (
            <div className="mt-6 flex justify-center gap-2">
              <Button
                variant="outline"
                size="sm"
                disabled={page === 0}
                onClick={() => setPage((p) => p - 1)}
              >
                Назад
              </Button>
              <span className="flex items-center px-4 text-sm text-muted-foreground">
                {page + 1} / {totalPages}
              </span>
              <Button
                variant="outline"
                size="sm"
                disabled={page >= totalPages - 1}
                onClick={() => setPage((p) => p + 1)}
              >
                Далі
              </Button>
            </div>
          )}
        </>
      ) : (
        <Card className="rounded-xl">
          <CardContent className="flex flex-col items-center justify-center py-12">
            <p className="text-muted-foreground">Покупок поки немає</p>
            <Button asChild variant="outline" className="mt-4">
              <Link to="/">Перейти до каталогу</Link>
            </Button>
          </CardContent>
        </Card>
      )}
    </div>
  )
}

function ReceiptCard({ receipt }: { receipt: Receipt }) {
  const totalCents = receipt.sumPrice - receipt.sumDiscount
  const totalUah = (totalCents / 100).toFixed(2)
  const statusLabel = STATUS_LABEL[receipt.stripeStatus] ?? receipt.stripeStatus

  return (
    <Card className="rounded-xl">
      <CardHeader className="flex flex-row items-start justify-between space-y-0 pb-2">
        <div>
          <CardTitle className="text-base">
            Замовлення #{receipt.id}
          </CardTitle>
          <CardDescription>
            {formatDate(receipt.createdAt)} · {statusLabel}
          </CardDescription>
        </div>
        <span className="font-semibold">{totalUah} грн</span>
      </CardHeader>
      <CardContent>
        <ul className="space-y-1 text-sm text-muted-foreground">
          {receipt.products.map((p) => {
            const price = ((p.price - p.discount) / 100).toFixed(2)
            return (
              <li key={`${p.productId}-${p.name}`}>
                {p.name} × {p.quantity} — {price} грн
              </li>
            )
          })}
        </ul>
      </CardContent>
    </Card>
  )
}
