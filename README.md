# eCommerce Store

A simple eCommerce application with a Go backend and Next.js frontend.

## Features

- Product catalog with 5 sample products
- Shopping cart functionality
- Order placement
- Payment processing (simulated)
- Responsive design

## Tech Stack

- **Backend**: Go with Gorilla Mux router
- **Frontend**: Next.js with TypeScript and Tailwind CSS
- **API**: RESTful API with CORS support

## Setup Instructions

### Prerequisites

- Go 1.21 or later
- Node.js 18 or later
- npm or yarn

### Backend Setup

1. Install Go dependencies:
```bash
go mod tidy
```

2. Run the Go server:
```bash
go run main.go
```

The backend will be available at `http://localhost:8080`

### Frontend Setup

1. Install Node.js dependencies:
```bash
npm install
```

2. Run the Next.js development server:
```bash
npm run dev
```

The frontend will be available at `http://localhost:3000`

## API Endpoints

- `GET /api/products` - Get all products
- `GET /api/products/{id}` - Get a specific product
- `POST /api/orders` - Create a new order
- `GET /api/orders` - Get all orders
- `POST /api/payment` - Process payment

## Usage

1. Start both the backend and frontend servers
2. Open `http://localhost:3000` in your browser
3. Browse products and add them to your cart
4. Click the cart button to view your items
5. Proceed to checkout and fill in the form
6. Complete the order (payment is simulated)

## Testing

The project includes comprehensive tests with 81.2% coverage:

```bash
# Run all tests
go test

# Run with coverage
go test -cover

# Run benchmarks
go test -bench=.

# Using Makefile
make test-all
```

See `TESTING.md` for detailed testing documentation.

## Sample Products

The store comes with 5 sample products:
- Wireless Headphones ($99.99)
- Smart Watch ($199.99)
- Coffee Maker ($79.99)
- Running Shoes ($129.99)
- Laptop Backpack ($49.99)