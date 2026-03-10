const API_BASE = "http://localhost:3000/api"

export type Product = {
  id: number
  name: string
  price: number
  discount: number
  description: string
  image: string
}

export type ProductsResponse = {
  products: Product[]
  total: number
  limit: number
  offset: number
}

export async function getProducts(limit: number, offset: number): Promise<ProductsResponse> {
  const params = new URLSearchParams({ limit: String(limit), offset: String(offset) })
  const res = await fetch(`${API_BASE}/products/get?${params}`)
  if (!res.ok) throw new Error("Не вдалося завантажити товари")
  return res.json() as Promise<ProductsResponse>
}
