# Halcon Core - Backend API

Backend API for the Halcon Logistics Suite built with Go, Echo framework, and PostgreSQL.

## Features

- **Authentication & Authorization**: JWT-based auth with role-based access control (RBAC)
- **User Management**: Admin can manage users with different roles
- **Order Management**: Full CRUD operations with role-based permissions
- **Public Tracking**: No-auth endpoint for customers to track orders
- **Soft Deletes**: Orders can be soft-deleted and restored
- **Evidence Upload**: Photo evidence for delivered orders
- **Status Workflow**: Ordered → In Process → In Route → Delivered

## Roles

- **Admin**: Manage users and full system access
- **Sales**: Create and manage orders
- **Purchasing**: Read-only access to orders in process
- **Warehouse**: Update order status to In Process and In Route
- **Route**: Upload evidence and mark orders as Delivered

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 15 or higher

## Setup

1. **Install dependencies**:
```bash
go mod download
```

2. **Configure environment**:
Copy `.env.example` to `.env` and update the values:
```bash
cp .env.example .env
```

3. **Start PostgreSQL** (or use Docker):
```bash
docker run --name halcon-postgres \
  -e POSTGRES_USER=halcon_user \
  -e POSTGRES_PASSWORD=halcon_password \
  -e POSTGRES_DB=halcon_db \
  -p 5432:5432 \
  -d postgres:15
```

4. **Run the server**:
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## Default Credentials

- **Username**: `admin`
- **Password**: `admin123`

⚠️ **Change this password immediately after first login!**

## API Endpoints

### Public Endpoints

- `GET /health` - Health check
- `POST /api/auth/login` - User login
- `GET /api/track?customer_number=XXX&invoice_number=YYY` - Track order

### Protected Endpoints (Require Authentication)

#### Auth
- `GET /api/auth/me` - Get current user

#### Users (Admin only)
- `GET /api/users` - List all users
- `GET /api/users/:id` - Get user by ID
- `POST /api/users` - Create new user
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user

#### Orders
- `GET /api/orders` - List orders (with filters)
- `GET /api/orders/:id` - Get order by ID
- `POST /api/orders` - Create order (Sales only)
- `PUT /api/orders/:id` - Update order (Warehouse, Route, Sales)
- `DELETE /api/orders/:id` - Soft delete order (Admin, Sales)
- `POST /api/orders/:id/restore` - Restore deleted order (Admin, Sales)
- `POST /api/orders/:id/evidence` - Upload evidence photo (Route only)

## Project Structure

```
halcon-core/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration management
│   ├── database/
│   │   └── database.go       # Database connection and migrations
│   ├── handlers/
│   │   ├── auth.go           # Authentication handlers
│   │   ├── users.go          # User management handlers
│   │   ├── orders.go         # Order management handlers
│   │   ├── tracking.go       # Public tracking handler
│   │   └── upload.go         # File upload handler
│   ├── middleware/
│   │   ├── auth.go           # JWT authentication middleware
│   │   └── rbac.go           # Role-based access control
│   ├── models/
│   │   └── models.go         # Database models
│   └── utils/
│       └── jwt.go            # JWT utilities
├── uploads/                  # Uploaded evidence photos
├── .env                      # Environment variables
├── .env.example              # Environment template
├── go.mod                    # Go module definition
└── README.md                 # This file
```

## Development

### Run with hot reload (using air):
```bash
go install github.com/cosmtrek/air@latest
air
```

### Run tests:
```bash
go test ./...
```

## License

MIT
