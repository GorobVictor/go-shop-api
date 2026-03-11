const API_BASE = "http://localhost:3000/api"

export type CreatePaymentBody = {
  products: Record<string, number>
}

export type PaymentLinkResponse = {
  link: string
}

export type ReceiptProduct = {
  productId: number
  quantity: number
  price: number
  discount: number
  name: string
}

export type Receipt = {
  id: number
  sumPrice: number
  sumDiscount: number
  createdAt: string
  stripeId: string
  stripeStatus: string
  products: ReceiptProduct[]
}

export type ReceiptsResponse = {
  receipts: Receipt[]
  total: number
  limit: number
  offset: number
}

export async function createPayment(
  token: string,
  body: CreatePaymentBody
): Promise<PaymentLinkResponse> {
  const products: Record<number, number> = {}
  for (const [id, qty] of Object.entries(body.products)) {
    const n = Number(id)
    if (Number.isInteger(n) && qty > 0) products[n] = qty
  }
  const res = await fetch(`${API_BASE}/payment/create`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ products }),
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || "Failed to create payment")
  }
  return res.json() as Promise<PaymentLinkResponse>
}

export async function getReceipts(
  token: string,
  limit: number,
  offset: number
): Promise<ReceiptsResponse> {
  const params = new URLSearchParams({
    limit: String(limit),
    offset: String(offset),
  })
  const res = await fetch(`${API_BASE}/payment/get?${params}`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!res.ok) throw new Error("Failed to load order history")
  return res.json() as Promise<ReceiptsResponse>
}
