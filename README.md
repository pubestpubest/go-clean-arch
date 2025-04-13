# Order Management System

A Go-based order management system built with clean architecture principles. This system provides a robust API for managing users, shops, products, and orders.

## Features

### User Management

- User registration
- User login
- User logout
- Get user details by ID
- Update user information

### Shop Management

- Create new products
- Get product details by ID
- Update product information
- Delete products
- List all products
- Get paginated product list
- Manage order statuses
- Get products by order ID
- Get order details by ID

### Order Management

- Create new orders
- Get order details by ID
- Get orders by user ID
- Get orders by shop ID

## Technical Stack

- Go (Golang)
- Echo Framework
- PostgreSQL
- GORM
- Docker
- Google Cloud Build

## Project Structure

```
.
├── configs/           # Configuration files
├── domain/           # Domain interfaces
├── entity/           # Database entities
├── features/         # Feature modules
│   ├── order/       # Order management
│   ├── product/     # Product management
│   ├── shop/        # Shop management
│   └── user/        # User management
├── middleware/       # HTTP middleware
├── seeders/         # Database seeders
└── utils/           # Utility functions
```

## Setup and Installation

1. Clone the repository
2. Copy `configs.example` to `configs` and configure your environment variables
3. Run the application using Docker:
   ```bash
   docker-compose up
   ```
4. Or run locally:
   ```bash
   go run main.go
   ```

## API Endpoints

### User Endpoints

- `POST /users/register` - Register a new user
- `POST /users/login` - User login
- `POST /users/logout` - User logout
- `GET /users/:id` - Get user by ID
- `PUT /users/:id` - Update user information

### Shop Endpoints

- `POST /shops/products` - Create a new product
- `GET /shops/products/:id` - Get product by ID
- `PUT /shops/products/:id` - Update product
- `DELETE /shops/products/:id` - Delete product
- `GET /shops/products` - Get all products
- `GET /shops/products/list` - Get paginated product list
- `PUT /shops/orders/:id/status` - Update order status
- `GET /shops/orders/:id/products` - Get products by order ID
- `GET /shops/orders/:id` - Get order by ID

### Order Endpoints

- `POST /orders` - Create a new order
- `GET /orders/:id` - Get order by ID
- `GET /users/:id/orders` - Get orders by user ID
- `GET /shops/:id/orders` - Get orders by shop ID

## Development

The project follows clean architecture principles with clear separation of concerns:

- Domain layer for business logic
- Repository layer for data access
- Usecase layer for application logic
- Delivery layer for API endpoints

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request
