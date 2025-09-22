export interface Product {
  id: number
  name: string
  description: string
  price: number
  image: string
  category: string
}

export interface OrderItem {
  product_id: number
  quantity: number
}

export interface Order {
  id: number
  items: OrderItem[]
  total: number
  status: string
  created_at: string
}

export interface PaymentRequest {
  order_id: number
  amount: number
}

export interface PaymentResponse {
  success: boolean
  message: string
  order_id: number
}
