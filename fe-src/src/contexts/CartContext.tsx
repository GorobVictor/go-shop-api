import {
  createContext,
  useCallback,
  useContext,
  useState,
  type ReactNode,
} from "react"

export type CartItem = {
  productId: number
  quantity: number
  name: string
  price: number
  discount: number
}

type CartContextValue = {
  items: CartItem[]
  addItem: (item: Omit<CartItem, "quantity">, quantity?: number) => void
  removeItem: (productId: number) => void
  setQuantity: (productId: number, quantity: number) => void
  clearCart: () => void
  totalCount: number
  productIdsToQuantity: () => Record<number, number>
}

const CartContext = createContext<CartContextValue | null>(null)

const STORAGE_KEY = "shop_cart"

function loadCart(): CartItem[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return []
    const parsed = JSON.parse(raw) as CartItem[]
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

function saveCart(items: CartItem[]) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(items))
}

export function CartProvider({ children }: { children: ReactNode }) {
  const [items, setItems] = useState<CartItem[]>(loadCart)

  const persist = useCallback((next: CartItem[]) => {
    setItems(next)
    saveCart(next)
  }, [])

  const addItem = useCallback(
    (item: Omit<CartItem, "quantity">, quantity = 1) => {
      persist(
        (() => {
          const existing = items.find((i) => i.productId === item.productId)
          const q = (existing?.quantity ?? 0) + quantity
          if (q <= 0) {
            return items.filter((i) => i.productId !== item.productId)
          }
          const entry: CartItem = existing
            ? { ...existing, quantity: q }
            : { ...item, quantity }
          const rest = items.filter((i) => i.productId !== item.productId)
          return [...rest, entry]
        })()
      )
    },
    [items, persist]
  )

  const removeItem = useCallback(
    (productId: number) => {
      persist(items.filter((i) => i.productId !== productId))
    },
    [items, persist]
  )

  const setQuantity = useCallback(
    (productId: number, quantity: number) => {
      if (quantity <= 0) {
        removeItem(productId)
        return
      }
      const existing = items.find((i) => i.productId === productId)
      if (!existing) return
      persist(
        items.map((i) =>
          i.productId === productId ? { ...i, quantity } : i
        )
      )
    },
    [items, persist, removeItem]
  )

  const clearCart = useCallback(() => persist([]), [persist])

  const totalCount = items.reduce((acc, i) => acc + i.quantity, 0)

  const productIdsToQuantity = useCallback(
    () =>
      items.reduce<Record<number, number>>((acc, i) => {
        acc[i.productId] = i.quantity
        return acc
      }, {}),
    [items]
  )

  const value: CartContextValue = {
    items,
    addItem,
    removeItem,
    setQuantity,
    clearCart,
    totalCount,
    productIdsToQuantity,
  }

  return (
    <CartContext.Provider value={value}>{children}</CartContext.Provider>
  )
}

export function useCart() {
  const ctx = useContext(CartContext)
  if (!ctx) throw new Error("useCart must be used within CartProvider")
  return ctx
}
