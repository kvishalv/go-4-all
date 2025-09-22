package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetProducts(t *testing.T) {
	// Create a request to the /api/products endpoint
	req, err := http.NewRequest("GET", "/api/products", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetProducts)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the content type
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, expectedContentType)
	}

	// Parse the response
	var products []Product
	if err := json.Unmarshal(rr.Body.Bytes(), &products); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	// Check that we have 5 products
	if len(products) != 5 {
		t.Errorf("expected 5 products, got %d", len(products))
	}

	// Check that all products have required fields
	for i, product := range products {
		if product.ID == 0 {
			t.Errorf("product %d has no ID", i)
		}
		if product.Name == "" {
			t.Errorf("product %d has no name", i)
		}
		if product.Price <= 0 {
			t.Errorf("product %d has invalid price: %v", i, product.Price)
		}
	}
}

func TestGetProduct(t *testing.T) {
	// Test getting a specific product
	req, err := http.NewRequest("GET", "/api/products/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/products/{id}", GetProduct).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var product Product
	if err := json.Unmarshal(rr.Body.Bytes(), &product); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if product.ID != 1 {
		t.Errorf("expected product ID 1, got %d", product.ID)
	}
}

func TestGetProductNotFound(t *testing.T) {
	// Test getting a non-existent product
	req, err := http.NewRequest("GET", "/api/products/999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/products/{id}", GetProduct).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestCreateOrder(t *testing.T) {
	// Reset state before test
	ResetGlobalState()

	// Create a test order
	orderData := Order{
		Items: []OrderItem{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 1},
		},
	}

	jsonData, err := json.Marshal(orderData)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/orders", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateOrder)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var order Order
	if err := json.Unmarshal(rr.Body.Bytes(), &order); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	// Check that the order was created with correct data
	if order.ID == 0 {
		t.Error("order ID should not be zero")
	}
	if order.Status != "pending" {
		t.Errorf("expected status 'pending', got %s", order.Status)
	}
	if order.Total <= 0 {
		t.Errorf("order total should be positive, got %v", order.Total)
	}
	if len(order.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(order.Items))
	}
}

func TestCreateOrderInvalidJSON(t *testing.T) {
	// Test with invalid JSON
	req, err := http.NewRequest("POST", "/api/orders", bytes.NewBufferString("invalid json"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateOrder)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestGetOrders(t *testing.T) {
	// Reset state before test
	ResetGlobalState()

	// First create an order
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

	// Now test getting orders
	req, err := http.NewRequest("GET", "/api/orders", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler := http.HandlerFunc(GetOrders)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var orders []Order
	if err := json.Unmarshal(rr.Body.Bytes(), &orders); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if len(orders) == 0 {
		t.Error("expected at least one order")
	}
}

func TestProcessPayment(t *testing.T) {
	// Reset state before test
	ResetGlobalState()

	// First create an order
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

	// Now test payment processing
	paymentData := PaymentRequest{
		OrderID: order.ID,
		Amount:  order.Total,
	}

	jsonData, err := json.Marshal(paymentData)
	if err != nil {
		t.Fatal(err)
	}

	req, err = http.NewRequest("POST", "/api/payment", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	handler := http.HandlerFunc(ProcessPayment)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var payment PaymentResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &payment); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if !payment.Success {
		t.Error("payment should be successful")
	}
	if payment.OrderID != order.ID {
		t.Errorf("expected order ID %d, got %d", order.ID, payment.OrderID)
	}
}

func TestProcessPaymentOrderNotFound(t *testing.T) {
	// Reset state before test
	ResetGlobalState()

	// Test payment for non-existent order
	paymentData := PaymentRequest{
		OrderID: 999,
		Amount:  100.0,
	}

	jsonData, err := json.Marshal(paymentData)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/payment", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ProcessPayment)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestProcessPaymentInvalidJSON(t *testing.T) {
	// Test with invalid JSON
	req, err := http.NewRequest("POST", "/api/payment", bytes.NewBufferString("invalid json"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ProcessPayment)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
