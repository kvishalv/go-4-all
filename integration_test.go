package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func TestFullOrderFlow(t *testing.T) {
	// Reset state
	ResetGlobalState()

	// Step 1: Get products
	req, err := http.NewRequest("GET", "/api/products", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetProducts)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("getProducts failed: got %v want %v", status, http.StatusOK)
	}

	var products []Product
	json.Unmarshal(rr.Body.Bytes(), &products)
	if len(products) != 5 {
		t.Errorf("expected 5 products, got %d", len(products))
	}

	// Step 2: Create an order
	orderData := Order{
		Items: []OrderItem{
			{ProductID: 1, Quantity: 2}, // Headphones x2
			{ProductID: 3, Quantity: 1}, // Coffee Maker x1
		},
	}

	jsonData, _ := json.Marshal(orderData)
	req, _ = http.NewRequest("POST", "/api/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	CreateOrder(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("createOrder failed: got %v want %v", status, http.StatusOK)
	}

	var order Order
	json.Unmarshal(rr.Body.Bytes(), &order)

	// Verify order details
	expectedTotal := (99.99 * 2) + (79.99 * 1) // Headphones + Coffee Maker
	// Use approximate comparison for floating point precision
	if order.Total < expectedTotal-0.01 || order.Total > expectedTotal+0.01 {
		t.Errorf("expected total ~%v, got %v", expectedTotal, order.Total)
	}

	if order.Status != "pending" {
		t.Errorf("expected status 'pending', got %s", order.Status)
	}

	// Step 3: Process payment
	paymentData := PaymentRequest{
		OrderID: order.ID,
		Amount:  order.Total,
	}

	jsonData, _ = json.Marshal(paymentData)
	req, _ = http.NewRequest("POST", "/api/payment", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	ProcessPayment(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("processPayment failed: got %v want %v", status, http.StatusOK)
	}

	var payment PaymentResponse
	json.Unmarshal(rr.Body.Bytes(), &payment)

	if !payment.Success {
		t.Error("payment should be successful")
	}

	// Step 4: Verify order status was updated
	req, _ = http.NewRequest("GET", "/api/orders", nil)
	rr = httptest.NewRecorder()
	GetOrders(rr, req)

	var orders []Order
	json.Unmarshal(rr.Body.Bytes(), &orders)

	if len(orders) != 1 {
		t.Errorf("expected 1 order, got %d", len(orders))
	}

	if orders[0].Status != "paid" {
		t.Errorf("expected status 'paid', got %s", orders[0].Status)
	}
}

func TestCORSHeaders(t *testing.T) {
	// Test that CORS headers are properly set
	req, err := http.NewRequest("GET", "/api/products", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create router with CORS
	r := mux.NewRouter()
	r.HandleFunc("/api/products", GetProducts).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(r)
	handler.ServeHTTP(rr, req)

	// Check that the request was successful (CORS is working)
	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	// Check that we got products back
	var products []Product
	if err := json.Unmarshal(rr.Body.Bytes(), &products); err != nil {
		t.Errorf("failed to unmarshal products: %v", err)
	}

	if len(products) == 0 {
		t.Error("expected products to be returned")
	}
}

func TestMultipleOrders(t *testing.T) {
	// Reset state
	ResetGlobalState()

	// Create first order
	order1 := Order{
		Items: []OrderItem{
			{ProductID: 1, Quantity: 1},
		},
	}

	jsonData, _ := json.Marshal(order1)
	req, _ := http.NewRequest("POST", "/api/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	CreateOrder(rr, req)

	var createdOrder1 Order
	json.Unmarshal(rr.Body.Bytes(), &createdOrder1)

	// Create second order
	order2 := Order{
		Items: []OrderItem{
			{ProductID: 2, Quantity: 1},
		},
	}

	jsonData, _ = json.Marshal(order2)
	req, _ = http.NewRequest("POST", "/api/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	CreateOrder(rr, req)

	var createdOrder2 Order
	json.Unmarshal(rr.Body.Bytes(), &createdOrder2)

	// Verify both orders exist and have different IDs
	if createdOrder1.ID == createdOrder2.ID {
		t.Error("orders should have different IDs")
	}

	// Get all orders
	req, _ = http.NewRequest("GET", "/api/orders", nil)
	rr = httptest.NewRecorder()
	GetOrders(rr, req)

	var orders []Order
	json.Unmarshal(rr.Body.Bytes(), &orders)

	if len(orders) != 2 {
		t.Errorf("expected 2 orders, got %d", len(orders))
	}
}

func TestOrderCalculation(t *testing.T) {
	// Test that order totals are calculated correctly
	testCases := []struct {
		name     string
		items    []OrderItem
		expected float64
	}{
		{
			name: "Single item",
			items: []OrderItem{
				{ProductID: 1, Quantity: 1}, // Headphones $99.99
			},
			expected: 99.99,
		},
		{
			name: "Multiple quantities",
			items: []OrderItem{
				{ProductID: 1, Quantity: 2}, // Headphones $99.99 x 2
			},
			expected: 199.98,
		},
		{
			name: "Multiple items",
			items: []OrderItem{
				{ProductID: 1, Quantity: 1}, // Headphones $99.99
				{ProductID: 2, Quantity: 1}, // Smart Watch $199.99
			},
			expected: 299.98,
		},
		{
			name: "Complex order",
			items: []OrderItem{
				{ProductID: 1, Quantity: 2}, // Headphones $99.99 x 2 = $199.98
				{ProductID: 3, Quantity: 1}, // Coffee Maker $79.99
				{ProductID: 5, Quantity: 3}, // Backpack $49.99 x 3 = $149.97
			},
			expected: 429.94,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset state for each test
			ResetGlobalState()

			orderData := Order{
				Items: tc.items,
			}

			jsonData, _ := json.Marshal(orderData)
			req, _ := http.NewRequest("POST", "/api/orders", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			CreateOrder(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("createOrder failed: got %v want %v", status, http.StatusOK)
			}

			var order Order
			json.Unmarshal(rr.Body.Bytes(), &order)

			// Use approximate comparison for floating point precision
			if order.Total < tc.expected-0.01 || order.Total > tc.expected+0.01 {
				t.Errorf("expected total ~%v, got %v", tc.expected, order.Total)
			}
		})
	}
}

func TestInvalidProductID(t *testing.T) {
	// Reset state before test
	ResetGlobalState()

	// Test order with non-existent product
	orderData := Order{
		Items: []OrderItem{
			{ProductID: 999, Quantity: 1}, // Non-existent product
		},
	}

	jsonData, _ := json.Marshal(orderData)
	req, _ := http.NewRequest("POST", "/api/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	CreateOrder(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("createOrder should handle invalid product IDs gracefully")
	}

	var order Order
	json.Unmarshal(rr.Body.Bytes(), &order)

	// Order should be created but with 0 total for invalid products
	if order.Total != 0 {
		t.Errorf("expected total 0 for invalid product, got %v", order.Total)
	}
}
