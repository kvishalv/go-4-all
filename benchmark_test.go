package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// BenchmarkGetProducts benchmarks the getProducts endpoint
func BenchmarkGetProducts(b *testing.B) {
	req, err := http.NewRequest("GET", "/api/products", nil)
	if err != nil {
		b.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetProducts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(rr, req)
		rr.Body.Reset()
	}
}

// BenchmarkCreateOrder benchmarks the createOrder endpoint
func BenchmarkCreateOrder(b *testing.B) {
	orderData := Order{
		Items: []OrderItem{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 1},
		},
	}

	jsonData, _ := json.Marshal(orderData)
	req, _ := http.NewRequest("POST", "/api/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateOrder)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset state for each iteration
		ResetGlobalState()
		handler.ServeHTTP(rr, req)
		rr.Body.Reset()
	}
}

// BenchmarkProcessPayment benchmarks the processPayment endpoint
func BenchmarkProcessPayment(b *testing.B) {
	// Setup: create an order first
	orderData := Order{
		Items: []OrderItem{
			{ProductID: 1, Quantity: 1},
		},
	}

	jsonData, _ := json.Marshal(orderData)
	req, _ := http.NewRequest("POST", "/api/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	CreateOrder(rr, req)

	var order Order
	json.Unmarshal(rr.Body.Bytes(), &order)

	// Benchmark payment processing
	paymentData := PaymentRequest{
		OrderID: order.ID,
		Amount:  order.Total,
	}

	jsonData, _ = json.Marshal(paymentData)
	req, _ = http.NewRequest("POST", "/api/payment", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	handler := http.HandlerFunc(ProcessPayment)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset state for each iteration
		ResetGlobalState()

		// Recreate order for each iteration
		orderData := Order{
			Items: []OrderItem{
				{ProductID: 1, Quantity: 1},
			},
		}
		jsonData, _ := json.Marshal(orderData)
		req, _ := http.NewRequest("POST", "/api/orders", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		CreateOrder(rr, req)

		var order Order
		json.Unmarshal(rr.Body.Bytes(), &order)

		paymentData := PaymentRequest{
			OrderID: order.ID,
			Amount:  order.Total,
		}
		jsonData, _ = json.Marshal(paymentData)
		req, _ = http.NewRequest("POST", "/api/payment", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		rr = httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		rr.Body.Reset()
	}
}

// BenchmarkOrderCalculation benchmarks order total calculation
func BenchmarkOrderCalculation(b *testing.B) {
	orderData := Order{
		Items: []OrderItem{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 1},
			{ProductID: 3, Quantity: 3},
		},
	}

	jsonData, _ := json.Marshal(orderData)
	req, _ := http.NewRequest("POST", "/api/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateOrder)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ResetGlobalState()
		handler.ServeHTTP(rr, req)
		rr.Body.Reset()
	}
}

// BenchmarkConcurrentRequests benchmarks concurrent API requests
func BenchmarkConcurrentRequests(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Simulate concurrent product requests
			req, _ := http.NewRequest("GET", "/api/products", nil)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetProducts)
			handler.ServeHTTP(rr, req)
		}
	})
}
