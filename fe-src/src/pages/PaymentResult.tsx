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
    title: "Payment successful",
    description: "Thank you for your purchase. You can view order details in your order history.",
    icon: "✓",
  },
  failed: {
    title: "Payment failed",
    description: "Something went wrong during payment. Please try again or use a different payment method.",
    icon: "✕",
  },
  canceled: {
    title: "Payment canceled",
    description: "You canceled the payment. Your cart has been saved — you can continue shopping.",
    icon: "○",
  },
  "something-went-wrong": {
    title: "Something went wrong",
    description: "An error occurred. Please try again later or contact support.",
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
            <Link to="/">{isSuccess ? "Home" : "Back to catalog"}</Link>
          </Button>
          {isSuccess && (
            <Button asChild variant="outline" className="rounded-lg">
              <Link to="/orders">Order history</Link>
            </Button>
          )}
          {variant === "canceled" && (
            <Button asChild variant="outline" className="rounded-lg">
              <Link to="/cart">Go to cart</Link>
            </Button>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
