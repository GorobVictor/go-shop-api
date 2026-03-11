import { Link, useLocation } from "react-router-dom"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

type Variant = "success" | "failed" | "canceled" | "something-went-wrong"

const CONFIG: Record<
  Variant,
  { title: string; description: string; icon: string }
> = {
  success: {
    title: "Оплата успішна",
    description: "Дякуємо за покупку. Деталі замовлення можна переглянути в історії покупок.",
    icon: "✓",
  },
  failed: {
    title: "Оплата не виконана",
    description: "Щось пішло не так під час оплати. Спробуйте ще раз або оберіть інший спосіб оплати.",
    icon: "✕",
  },
  canceled: {
    title: "Оплату скасовано",
    description: "Ви скасували оплату. Кошик збережено — можете продовжити покупки.",
    icon: "○",
  },
  "something-went-wrong": {
    title: "Щось пішло не так",
    description: "Виникла помилка. Спробуйте пізніше або зверніться до підтримки.",
    icon: "?",
  },
}

export default function PaymentResult() {
  const path = useLocation().pathname
  const variant: Variant = path.includes("success")
    ? "success"
    : path.includes("failed")
      ? "failed"
      : path.includes("canceled")
        ? "canceled"
        : "something-went-wrong"

  const config = CONFIG[variant]
  const isSuccess = variant === "success"

  return (
    <div className="container mx-auto flex max-w-md flex-col items-center justify-center px-4 py-16">
      <Card className="w-full rounded-2xl text-center">
        <CardHeader>
          <div
            className={`mx-auto flex h-14 w-14 items-center justify-center rounded-full text-2xl font-bold ${
              isSuccess ? "bg-primary/15 text-primary" : "bg-muted text-muted-foreground"
            }`}
          >
            {config.icon}
          </div>
          <CardTitle className="mt-4">{config.title}</CardTitle>
          <CardDescription>{config.description}</CardDescription>
        </CardHeader>
        <CardContent className="flex flex-col gap-3">
          <Button asChild className="rounded-lg">
            <Link to="/">{isSuccess ? "На головну" : "Повернутися в каталог"}</Link>
          </Button>
          {isSuccess && (
            <Button asChild variant="outline" className="rounded-lg">
              <Link to="/orders">Історія покупок</Link>
            </Button>
          )}
          {variant === "canceled" && (
            <Button asChild variant="outline" className="rounded-lg">
              <Link to="/cart">Перейти до кошика</Link>
            </Button>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
