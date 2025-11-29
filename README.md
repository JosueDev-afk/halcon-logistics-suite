# Halcon Logistics Suite

A full-stack logistics management system built with Go (backend) and Vue 3 (frontend).

## 🚀 Features

- **Public Order Tracking**: No-auth tracking page for customers
- **Role-Based Access Control (RBAC)**: 5 distinct roles with specific permissions
- **Order Management**: Complete order lifecycle from creation to delivery
- **Evidence Upload**: Photo evidence for delivered orders
- **Soft Delete**: Recycle bin for deleted orders
- **User Management**: Admin panel for managing users

## 👥 User Roles

1. **Admin**: Manage users and full system access
2. **Sales**: Create and manage orders
3. **Purchasing**: Read-only access to orders in process
4. **Warehouse**: Update order status to In Process and In Route
5. **Route**: Upload evidence and mark orders as Delivered

## 📦 Order Lifecycle

```
Ordered → In Process → In Route → Delivered
```

- **Sales** creates order → Status: `Ordered`
- **Warehouse** prepares → Status: `In Process`
- **Warehouse** dispatches → Status: `In Route`
- **Route** delivers + uploads photo → Status: `Delivered`

## 🛠️ Tech Stack

### Backend (halcon-core)
- Go 1.21+
- Echo framework
- GORM (PostgreSQL)
- JWT authentication
- Docker

### Frontend (halcon-client)
- Vue 3 + TypeScript
- Vite
- Pinia (state management)
- Vue Router
- TailwindCSS
- Axios

### Database
- PostgreSQL 15

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+ (for local development)
- Node.js 20+ & pnpm (for local development)

### Using Docker Compose (Recommended)

1. **Clone the repository**
```bash
cd halcon-logistics-suite
```

2. **Start all services**
```bash
docker-compose up -d
```

This will start:
- PostgreSQL on port 5432
- Backend API on port 8080
- Frontend on port 5173

3. **Access the application**
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- Public Tracking: http://localhost:5173/track

### Default Credentials

```
Username: admin
Password: admin123
```

⚠️ **Change this password immediately after first login!**

## 📁 Project Structure

```
halcon-logistics-suite/
├── halcon-core/              # Go backend
│   ├── cmd/server/           # Main application
│   ├── internal/
│   │   ├── config/           # Configuration
│   │   ├── database/         # DB connection & migrations
│   │   ├── handlers/         # HTTP handlers
│   │   ├── middleware/       # Auth & RBAC middleware
│   │   ├── models/           # Database models
│   │   └── utils/            # Utilities (JWT, etc.)
│   └── uploads/              # Evidence photos
│
├── halcon-client/            # Vue 3 frontend
│   ├── src/
│   │   ├── api/              # API client
│   │   ├── stores/           # Pinia stores
│   │   ├── router/           # Vue Router
│   │   ├── views/            # Page components
│   │   └── assets/           # Static assets
│   └── public/
│
└── docker-compose.yml        # Docker orchestration
```

## 🔧 Local Development

### Backend

```bash
cd halcon-core

# Install dependencies
go mod download

# Start PostgreSQL (or use Docker)
docker run --name halcon-postgres \
  -e POSTGRES_USER=halcon_user \
  -e POSTGRES_PASSWORD=halcon_password \
  -e POSTGRES_DB=halcon_db \
  -p 5432:5432 \
  -d postgres:15

# Run the server
go run cmd/server/main.go
```

### Frontend

```bash
cd halcon-client

# Install dependencies
pnpm install

# Start dev server
pnpm dev
```

## 📚 API Documentation

### Public Endpoints
- `GET /health` - Health check
- `POST /api/auth/login` - User login
- `GET /api/track` - Track order (no auth)

### Protected Endpoints

#### Authentication
- `GET /api/auth/me` - Get current user

#### Users (Admin only)
- `GET /api/users` - List users
- `POST /api/users` - Create user
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user

#### Orders
- `GET /api/orders` - List orders (with filters)
- `GET /api/orders/:id` - Get order details
- `POST /api/orders` - Create order (Sales)
- `PUT /api/orders/:id` - Update order
- `DELETE /api/orders/:id` - Soft delete (Admin, Sales)
- `POST /api/orders/:id/restore` - Restore deleted order
- `POST /api/orders/:id/evidence` - Upload evidence (Route)

## 🔐 Environment Variables

### Backend (.env)
```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=halcon_user
DB_PASSWORD=halcon_password
DB_NAME=halcon_db
JWT_SECRET=your-secret-key
CORS_ALLOWED_ORIGINS=http://localhost:5173
```

### Frontend (.env)
```env
VITE_API_URL=http://localhost:8080
```

## 📝 License

MIT

## 👨‍💻 Author

Built for Tecmilenio - Web Applications Design Course
