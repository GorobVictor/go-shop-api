import { useEffect, useState } from "react"
import { getProducts, type Product } from "@/api/products"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Button } from "@/components/ui/button"

const PAGE_SIZE = 12

function ProductCard({ product }: { product: Product }) {
  const priceUah = (product.price / 100).toFixed(2)
  const finalPrice =
    product.discount > 0 ? product.price - product.discount : product.price
  const finalPriceUah = (finalPrice / 100).toFixed(2)
  const hasDiscount = product.discount > 0
  const discountPercent =
    product.price > 0
      ? Math.round((product.discount / product.price) * 100)
      : 0

  return (
    <Card className="group flex flex-col overflow-hidden rounded-xl border-border/80 bg-card shadow-sm transition-all duration-200 hover:-translate-y-1 hover:shadow-lg">
      <div className="relative aspect-[4/3] w-full overflow-hidden bg-muted">
        {product.image ? (
          <img
            src={product.image}
            alt={product.name}
            className="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
          />
        ) : (
          <div className="flex h-full w-full items-center justify-center text-muted-foreground">
            <span className="text-sm">Немає зображення</span>
          </div>
        )}
        {hasDiscount && discountPercent > 0 && (
          <span className="absolute right-2 top-2 rounded-full bg-primary px-2 py-0.5 text-xs font-semibold text-primary-foreground shadow">
            −{discountPercent}%
          </span>
        )}
      </div>
      <CardHeader className="space-y-1 pb-2">
        <CardTitle className="line-clamp-2 text-base font-semibold leading-tight">
          {product.name}
        </CardTitle>
        {product.description && (
          <CardDescription className="line-clamp-2 text-sm">
            {product.description}
          </CardDescription>
        )}
      </CardHeader>
      <CardContent className="flex-1 pb-2">
        <div className="flex flex-wrap items-baseline gap-2">
          {hasDiscount ? (
            <>
              <span className="text-sm text-muted-foreground line-through">
                {priceUah} грн
              </span>
              <span className="text-lg font-semibold text-primary">
                {finalPriceUah} грн
              </span>
            </>
          ) : (
            <span className="text-lg font-semibold">{priceUah} грн</span>
          )}
        </div>
      </CardContent>
      <CardFooter className="pt-0">
        <Button
          variant="outline"
          className="w-full rounded-lg"
          disabled
        >
          Деталі (скоро)
        </Button>
      </CardFooter>
    </Card>
  )
}

function SkeletonCard() {
  return (
    <Card className="overflow-hidden rounded-xl border-border/80 bg-card shadow-sm">
      <div className="aspect-[4/3] w-full animate-pulse bg-muted" />
      <CardHeader className="space-y-2">
        <div className="h-5 w-3/4 animate-pulse rounded bg-muted" />
        <div className="h-4 w-1/2 animate-pulse rounded bg-muted" />
      </CardHeader>
      <CardContent>
        <div className="h-6 w-1/3 animate-pulse rounded bg-muted" />
      </CardContent>
      <CardFooter>
        <div className="h-9 w-full animate-pulse rounded-lg bg-muted" />
      </CardFooter>
    </Card>
  )
}

export default function Catalog() {
  const [data, setData] = useState<Awaited<
    ReturnType<typeof getProducts>
  > | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [page, setPage] = useState(0)

  useEffect(() => {
    setLoading(true)
    setError(null)
    getProducts(PAGE_SIZE, page * PAGE_SIZE)
      .then(setData)
      .catch((e) => setError(e instanceof Error ? e.message : "Помилка"))
      .finally(() => setLoading(false))
  }, [page])

  if (error) {
    return (
      <div className="container mx-auto max-w-6xl px-4 py-12">
        <div className="rounded-xl border border-destructive/30 bg-destructive/5 p-6 text-center">
          <p className="font-medium text-destructive">{error}</p>
          <p className="mt-1 text-sm text-muted-foreground">
            Перевірте підключення до сервера та оновіть сторінку.
          </p>
        </div>
      </div>
    )
  }

  const totalPages = data ? Math.ceil(data.total / PAGE_SIZE) : 0
  const from = data ? page * PAGE_SIZE + 1 : 0
  const to = data
    ? Math.min((page + 1) * PAGE_SIZE, data.total)
    : 0

  return (
    <div className="container mx-auto max-w-6xl px-4 py-8">
      <section className="mb-10 text-center">
        <h1 className="text-3xl font-bold tracking-tight text-foreground sm:text-4xl">
          Каталог продуктів
        </h1>
        <p className="mt-2 text-muted-foreground">
          Відкритий каталог — переглядайте товари без реєстрації
        </p>
      </section>

      {loading ? (
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {Array.from({ length: 8 }).map((_, i) => (
            <SkeletonCard key={i} />
          ))}
        </div>
      ) : data && data.products.length > 0 ? (
        <>
          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {data.products.map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>

          {totalPages > 1 && (
            <div className="mt-10 flex flex-col items-center gap-4 sm:flex-row sm:justify-between">
              <p className="text-sm text-muted-foreground">
                Показано {from}–{to} з {data.total}
              </p>
              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  disabled={page === 0}
                  onClick={() => setPage((p) => p - 1)}
                >
                  Назад
                </Button>
                <span className="min-w-[4rem] text-center text-sm text-muted-foreground">
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
            </div>
          )}
        </>
      ) : (
        <div className="rounded-xl border border-dashed border-border bg-muted/30 py-16 text-center">
          <p className="text-muted-foreground">
            Товарів поки немає. Заходьте пізніше.
          </p>
        </div>
      )}
    </div>
  )
}
