import { useState, useEffect } from 'react'
import Head from 'next/head'
import ProductCard from '../components/ProductCard'
import Cart from '../components/Cart'
import Checkout from '../components/Checkout'
import { Product, OrderItem } from '../types'

const API_BASE = 'http://localhost:8080/api'

export default function Home() {
  const [products, setProducts] = useState<Product[]>([])
  const [cart, setCart] = useState<OrderItem[]>([])
  const [showCart, setShowCart] = useState(false)
  const [showCheckout, setShowCheckout] = useState(false)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchProducts()
  }, [])

  const fetchProducts = async () => {
    try {
      const response = await fetch(`${API_BASE}/products`)
      const data = await response.json()
      setProducts(data)
    } catch (error) {
      console.error('Error fetching products:', error)
    } finally {
      setLoading(false)
    }
  }

  const addToCart = (productId: number) => {
    const existingItem = cart.find(item => item.product_id === productId)
    if (existingItem) {
      setCart(cart.map(item => 
        item.product_id === productId 
          ? { ...item, quantity: item.quantity + 1 }
          : item
      ))
    } else {
      setCart([...cart, { product_id: productId, quantity: 1 }])
    }
  }

  const removeFromCart = (productId: number) => {
    setCart(cart.filter(item => item.product_id !== productId))
  }

  const updateQuantity = (productId: number, quantity: number) => {
    if (quantity <= 0) {
      removeFromCart(productId)
    } else {
      setCart(cart.map(item => 
        item.product_id === productId 
          ? { ...item, quantity }
          : item
      ))
    }
  }

  const getCartTotal = () => {
    return cart.reduce((total, item) => {
      const product = products.find(p => p.id === item.product_id)
      return total + (product ? product.price * item.quantity : 0)
    }, 0)
  }

  const getCartItemCount = () => {
    return cart.reduce((total, item) => total + item.quantity, 0)
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="text-xl">Loading...</div>
      </div>
    )
  }

  return (
    <>
      <Head>
        <title>eCommerce Store</title>
        <meta name="description" content="A simple eCommerce store" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main className="min-h-screen bg-gray-50">
        {/* Header */}
        <header className="bg-white shadow-sm border-b">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex justify-between items-center py-4">
              <h1 className="text-2xl font-bold text-gray-900">eCommerce Store</h1>
              <button
                onClick={() => setShowCart(true)}
                className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 flex items-center gap-2"
              >
                <span>Cart ({getCartItemCount()})</span>
                <span>${getCartTotal().toFixed(2)}</span>
              </button>
            </div>
          </div>
        </header>

        {/* Products Grid */}
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {products.map((product) => (
              <ProductCard
                key={product.id}
                product={product}
                onAddToCart={addToCart}
              />
            ))}
          </div>
        </div>

        {/* Cart Modal */}
        {showCart && (
          <Cart
            cart={cart}
            products={products}
            onClose={() => setShowCart(false)}
            onUpdateQuantity={updateQuantity}
            onRemoveItem={removeFromCart}
            onCheckout={() => {
              setShowCart(false)
              setShowCheckout(true)
            }}
            total={getCartTotal()}
          />
        )}

        {/* Checkout Modal */}
        {showCheckout && (
          <Checkout
            cart={cart}
            products={products}
            onClose={() => setShowCheckout(false)}
            onOrderSuccess={() => {
              setCart([])
              setShowCheckout(false)
            }}
            total={getCartTotal()}
          />
        )}
      </main>
    </>
  )
}
