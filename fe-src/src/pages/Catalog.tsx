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
  const finalPrice = product.discount > 0 ? product.price - product.discount : product.price
  const finalPriceUah = (finalPrice / 100).toFixed(2)
  const hasDiscount = product.discount > 0

  return (
    <Card className="flex flex-col overflow-hidden">
      {product.image ? (
        <img
          src={product.image}
          alt={product.name}
          className="h-48 w-full object-cover"
        />
      ) : (
        <div className="flex h-48 w-full items-center justify-center bg-muted text-muted-foreground">
          Немає зображення
        </div>
      )}
      <CardHeader className="pb-2">
        <CardTitle className="line-clamp-2 text-lg">{product.name}</CardTitle>
        {product.description && (
          <CardDescription className="line-clamp-2">
            {product.description}
          </CardDescription>
        )}
      </CardHeader>
      <CardContent className="flex-1 pb-2">
        <div className="flex items-baseline gap-2">
          {hasDiscount ? (
            <>
              <span className="text-muted-foreground text-sm line-through">
                {priceUah} грн
              </span>
              <span className="font-semibold text-primary">
                {finalPriceUah} грн
              </span>
            </>
          ) : (
            <span className="font-semibold">{priceUah} грн</span>
          )}
        </div>
      </CardContent>
      <CardFooter>
        <Button variant="outline" className="w-full" disabled>
          Деталі (скоро)
        </Button>
      </CardFooter>
    </Card>
  )
}

export default function Catalog() {
  const [data, setData] = useState<Awaited<ReturnType<typeof getProducts>> | null>(null)
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
      <div className="container px-4 py-8">
        <p className="text-destructive">{error}</p>
      </div>
    )
  }

  const totalPages = data ? Math.ceil(data.total / PAGE_SIZE) : 0

  return (
    <div className="container px-4 py-6">
      <h1 className="mb-6 text-2xl font-semibold">Каталог продуктів</h1>
      {loading ? (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {Array.from({ length: 8 }).map((_, i) => (
            <Card key={i} className="overflow-hidden">
              <div className="h-48 animate-pulse bg-muted" />
              <CardHeader>
                <div className="h-5 w-3/4 animate-pulse rounded bg-muted" />
                <div className="h-4 w-1/2 animate-pulse rounded bg-muted" />
              </CardHeader>
              <CardContent>
                <div className="h-6 w-1/4 animate-pulse rounded bg-muted" />
              </CardContent>
            </Card>
          ))}
        </div>
      ) : data && data.products.length > 0 ? (
        <>
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {data.products.map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>
          {totalPages > 1 && (
            <div className="mt-8 flex justify-center gap-2">
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
        <p className="text-muted-foreground">Товарів поки немає.</p>
      )}
    </div>
  )
}
