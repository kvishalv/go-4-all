# Testing Guide for eCommerce Backend

This document provides comprehensive information about testing the Go backend for the eCommerce application.

## Test Structure

The test suite is organized into several files in the root directory:

- **`unit_test.go`** - Unit tests for individual functions
- **`integration_test.go`** - Integration tests for complete workflows
- **`benchmark_test.go`** - Performance benchmarks

## Running Tests

### Basic Test Commands

```bash
# Run all tests
go test

# Run tests with verbose output
go test -v

# Run tests with coverage
go test -cover

# Run tests with race detection
go test -race

# Run specific test
go test -run TestGetProducts

# Run benchmarks
go test -bench=.
```

### Using the Makefile

```bash
# Install dependencies
make deps

# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage report
make test-coverage

# Run tests with race detection
make test-race

# Run all tests (unit + integration + coverage)
make test-all

# Generate coverage report
make coverage
```

### Using the Test Script

```bash
# Make script executable
chmod +x run_tests.sh

# Run the test script
./run_tests.sh
```

## Test Coverage

Current test coverage: **81.2%**

### Coverage Report

To generate a detailed coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

Open `coverage.html` in your browser to see detailed coverage information.

## Test Categories

### 1. Unit Tests (`unit_test.go`)

Tests individual functions in isolation:

- **`TestGetProducts`** - Tests product listing endpoint
- **`TestGetProduct`** - Tests single product retrieval
- **`TestGetProductNotFound`** - Tests error handling for non-existent products
- **`TestCreateOrder`** - Tests order creation
- **`TestCreateOrderInvalidJSON`** - Tests error handling for invalid JSON
- **`TestGetOrders`** - Tests order listing
- **`TestProcessPayment`** - Tests payment processing
- **`TestProcessPaymentOrderNotFound`** - Tests error handling for non-existent orders
- **`TestProcessPaymentInvalidJSON`** - Tests error handling for invalid payment JSON

### 2. Integration Tests (`integration_test.go`)

Tests complete workflows and system interactions:

- **`TestFullOrderFlow`** - Complete order workflow from product selection to payment
- **`TestCORSHeaders`** - Tests CORS configuration
- **`TestMultipleOrders`** - Tests handling of multiple orders
- **`TestOrderCalculation`** - Tests order total calculations with various scenarios
- **`TestInvalidProductID`** - Tests graceful handling of invalid product IDs

### 3. Benchmark Tests (`benchmark_test.go`)

Performance testing for critical operations:

- **`BenchmarkGetProducts`** - Product listing performance
- **`BenchmarkCreateOrder`** - Order creation performance
- **`BenchmarkProcessPayment`** - Payment processing performance
- **`BenchmarkOrderCalculation`** - Order calculation performance
- **`BenchmarkConcurrentRequests`** - Concurrent request handling

## Test Data

### Sample Products

The test suite uses 5 sample products:

1. **Wireless Headphones** - $99.99 (Electronics)
2. **Smart Watch** - $199.99 (Electronics)
3. **Coffee Maker** - $79.99 (Kitchen)
4. **Running Shoes** - $129.99 (Sports)
5. **Laptop Backpack** - $49.99 (Accessories)

### Test Scenarios

#### Order Calculation Tests

- Single item orders
- Multiple quantity orders
- Multiple item orders
- Complex orders with various combinations

#### Error Handling Tests

- Invalid JSON input
- Non-existent product IDs
- Non-existent order IDs
- Invalid payment requests

#### Integration Workflows

- Complete order flow (products → cart → order → payment)
- Multiple concurrent orders
- CORS header validation

## Performance Benchmarks

Current benchmark results (Apple M4):

```
BenchmarkGetProducts-10           	  951819	      1245 ns/op
BenchmarkCreateOrder-10           	 4628532	       261.9 ns/op
BenchmarkProcessPayment-10        	       1	1000707125 ns/op
BenchmarkOrderCalculation-10      	 4414767	       277.9 ns/op
BenchmarkConcurrentRequests-10    	 1636618	       732.1 ns/op
```

## Test Configuration

### Environment Variables

- `TEST_MODE=true` - Enables test mode
- `VERBOSE=true` - Enables verbose test output

### Test Utilities

- **`resetGlobalState()`** - Resets global variables for clean test runs
- **`GetTestConfig()`** - Retrieves test configuration
- **`SetupTestEnvironment()`** - Sets up test environment
- **`CleanupTestEnvironment()`** - Cleans up after tests

## Best Practices

### Writing Tests

1. **Use descriptive test names** that explain what is being tested
2. **Test both success and failure cases**
3. **Use table-driven tests** for multiple scenarios
4. **Reset state** between tests to avoid interference
5. **Use approximate comparisons** for floating-point values

### Test Organization

1. **Unit tests** should test individual functions
2. **Integration tests** should test complete workflows
3. **Benchmark tests** should measure performance
4. **Use helper functions** to reduce code duplication

### Error Handling

1. **Test error conditions** thoroughly
2. **Verify error messages** are appropriate
3. **Test edge cases** and boundary conditions
4. **Use proper HTTP status codes**

## Continuous Integration

The test suite is designed to run in CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
- name: Run Tests
  run: |
    go test -v -cover ./...
    go test -bench=. ./...
```

## Troubleshooting

### Common Issues

1. **Floating-point precision** - Use approximate comparisons for monetary values
2. **Test isolation** - Always reset global state between tests
3. **CORS testing** - Test CORS configuration separately from business logic
4. **Performance testing** - Payment processing includes artificial delay

### Debug Tips

1. Use `-v` flag for verbose output
2. Run individual tests with `-run` flag
3. Use `-race` flag to detect race conditions
4. Generate coverage reports for detailed analysis

## Future Improvements

1. **Increase test coverage** to 80%+
2. **Add property-based testing** for edge cases
3. **Add load testing** for high-traffic scenarios
4. **Add database integration tests** when migrating from in-memory storage
5. **Add API contract testing** for frontend integration
