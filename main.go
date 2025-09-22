package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Product represents a product in our store
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Image       string  `json:"image"`
	Category    string  `json:"category"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// Order represents a customer order
type Order struct {
	ID        int         `json:"id"`
	Items     []OrderItem `json:"items"`
	Total     float64     `json:"total"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
}

// PaymentRequest represents a payment request
type PaymentRequest struct {
	OrderID int     `json:"order_id"`
	Amount  float64 `json:"amount"`
}

// PaymentResponse represents a payment response
type PaymentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	OrderID int    `json:"order_id"`
}

// In-memory storage (in production, use a database)
var products []Product
var orders []Order
var nextOrderID = 1

func init() {
	// Initialize with 5 sample products
	products = []Product{
		{
			ID:          1,
			Name:        "Wireless Headphones",
			Description: "High-quality wireless headphones with noise cancellation",
			Price:       99.99,
			Image:       "https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=300&h=200&fit=crop",
			Category:    "Electronics",
		},
		{
			ID:          2,
			Name:        "Smart Watch",
			Description: "Fitness tracking smartwatch with heart rate monitor",
			Price:       199.99,
			Image:       "https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=300&h=200&fit=crop",
			Category:    "Electronics",
		},
		{
			ID:          3,
			Name:        "Coffee Maker",
			Description: "Automatic drip coffee maker with programmable timer",
			Price:       79.99,
			Image:       "https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=300&h=200&fit=crop",
			Category:    "Kitchen",
		},
		{
			ID:          4,
			Name:        "Running Shoes",
			Description: "Comfortable running shoes with breathable mesh",
			Price:       129.99,
			Image:       "https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=300&h=200&fit=crop",
			Category:    "Sports",
		},
		{
			ID:          5,
			Name:        "Laptop Backpack",
			Description: "Durable laptop backpack with multiple compartments",
			Price:       49.99,
			Image:       "https://images.unsplash.com/photo-1553062407-98eeb64c6a62?w=300&h=200&fit=crop",
			Category:    "Accessories",
		},
	}
}

// Get all products
func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// Get a single product by ID
func GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	for _, product := range products {
		if product.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(product)
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

// Create a new order
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Calculate total
	total := 0.0
	for _, item := range order.Items {
		for _, product := range products {
			if product.ID == item.ProductID {
				total += product.Price * float64(item.Quantity)
				break
			}
		}
	}

	// Create order
	order.ID = nextOrderID
	nextOrderID++
	order.Total = total
	order.Status = "pending"
	order.CreatedAt = time.Now()

	orders = append(orders, order)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// Get all orders
func GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// Process payment
func ProcessPayment(w http.ResponseWriter, r *http.Request) {
	var paymentReq PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&paymentReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find the order
	var order *Order
	for i := range orders {
		if orders[i].ID == paymentReq.OrderID {
			order = &orders[i]
			break
		}
	}

	if order == nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Simulate payment processing
	// In a real application, integrate with a payment processor like Stripe
	time.Sleep(1 * time.Second) // Simulate processing time

	// Update order status
	order.Status = "paid"

	response := PaymentResponse{
		Success: true,
		Message: "Payment processed successfully",
		OrderID: order.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/api/products", GetProducts).Methods("GET")
	r.HandleFunc("/api/products/{id}", GetProduct).Methods("GET")
	r.HandleFunc("/api/orders", CreateOrder).Methods("POST")
	r.HandleFunc("/api/orders", GetOrders).Methods("GET")
	r.HandleFunc("/api/payment", ProcessPayment).Methods("POST")

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(r)

	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

// ResetGlobalState resets the global state for testing
func ResetGlobalState() {
	products = []Product{
		{
			ID:          1,
			Name:        "Wireless Headphones",
			Description: "High-quality wireless headphones with noise cancellation",
			Price:       99.99,
			Image:       "https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=300&h=200&fit=crop",
			Category:    "Electronics",
		},
		{
			ID:          2,
			Name:        "Smart Watch",
			Description: "Fitness tracking smartwatch with heart rate monitor",
			Price:       199.99,
			Image:       "https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=300&h=200&fit=crop",
			Category:    "Electronics",
		},
		{
			ID:          3,
			Name:        "Coffee Maker",
			Description: "Automatic drip coffee maker with programmable timer",
			Price:       79.99,
			Image:       "https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=300&h=200&fit=crop",
			Category:    "Kitchen",
		},
		{
			ID:          4,
			Name:        "Running Shoes",
			Description: "Comfortable running shoes with breathable mesh",
			Price:       129.99,
			Image:       "https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=300&h=200&fit=crop",
			Category:    "Sports",
		},
		{
			ID:          5,
			Name:        "Laptop Backpack",
			Description: "Durable laptop backpack with multiple compartments",
			Price:       49.99,
			Image:       "https://images.unsplash.com/photo-1553062407-98eeb64c6a62?w=300&h=200&fit=crop",
			Category:    "Accessories",
		},
	}
	orders = []Order{}
	nextOrderID = 1
}
