import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom"
import { AuthProvider } from "@/contexts/AuthContext"
import { CartProvider } from "@/contexts/CartContext"
import Layout from "@/components/Layout"
import Catalog from "@/pages/Catalog"
import Cart from "@/pages/Cart"
import Orders from "@/pages/Orders"
import Login from "@/pages/Login"
import Register from "@/pages/Register"
import PaymentResult from "@/pages/PaymentResult"

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <CartProvider>
          <Routes>
            <Route element={<Layout />}>
              <Route path="/" element={<Catalog />} />
              <Route path="/cart" element={<Cart />} />
              <Route path="/orders" element={<Orders />} />
              <Route path="/login" element={<Login />} />
              <Route path="/register" element={<Register />} />
              <Route path="/payment/success" element={<PaymentResult />} />
              <Route path="/payment/failed" element={<PaymentResult />} />
              <Route path="/payment/canceled" element={<PaymentResult />} />
              <Route path="/payment/something-went-wrong" element={<PaymentResult />} />
            </Route>
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </CartProvider>
      </AuthProvider>
    </BrowserRouter>
  )
}

export default App
