# Cozystay

> A modern, full-stack property booking platform built with Go and Next.js. Discover and book unique properties around the world with ease.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.24-blue.svg)
![Next.js](https://img.shields.io/badge/Next.js-16-black)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-18-blue)

## 🎯 Overview

Cozystay is a comprehensive property booking platform that allows users to discover, view, and book unique accommodations. The platform features a robust RESTful API backend written in Go and a modern, responsive frontend built with Next.js and TypeScript.

## ✨ Features

### Core Functionality
- **User Authentication & Authorization**
  - Secure user registration and login
  - JWT-based authentication
  - Role-based access control (Guest/Admin)
  - Protected routes and API endpoints

- **Property Management**
  - Browse and search properties
  - View detailed property information
  - Property creation and management (Admin/Host)
  - Multiple property images with captions
  - Amenities management
  - Availability checking

- **Booking System**
  - Create bookings with date ranges
  - View booking history
  - Cancel bookings
  - Automatic price calculation

- **Modern UI/UX**
  - Responsive design for all devices
  - Dark mode support
  - Beautiful, intuitive interface built with Tailwind CSS
  - shadcn/ui component library

## 🏗️ Architecture

### Tech Stack

#### Backend
- **Language**: Go 1.24
- **Framework**: Chi Router v5
- **Database**: PostgreSQL 18
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Password Hashing**: bcrypt
- **Migrations**: Goose
- **Environment**: godotenv

#### Frontend
- **Framework**: Next.js 16
- **Language**: TypeScript
- **UI Library**: React 19
- **Styling**: Tailwind CSS 4
- **Components**: shadcn/ui (Radix UI)
- **Forms**: React Hook Form + Zod
- **HTTP Client**: Axios
- **Icons**: Lucide React

## 🚀 Getting Started

### Prerequisites

- **Go** 1.24 or higher
- **Node.js** 20+ and **pnpm** (or npm/yarn)
- **PostgreSQL** 18+ (or use Docker)
- **Docker** and **Docker Compose** (optional, for containerized setup)

### Installation

#### 1. Clone the Repository

```bash
git clone https://github.com/Umesh-Tiruvalluru/BookBnb.git
cd BookBnb
```

#### 2. Backend Setup

Create a `.env` file in the root directory:

```env
# Server Configuration
PORT=4000
ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=cozystay

# DSN (optional, can be constructed from above)
DSN=host=localhost port=5432 user=your_db_user password=your_db_password dbname=cozystay sslmode=disable

# Database Connection Pool
MAX_OPEN_CONNS=25
MAX_IDLE_CONNS=5
MAX_IDLE_TIME=15m

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

Install Go dependencies:

```bash
go mod download
```

#### 3. Frontend Setup

Navigate to the frontend directory and install dependencies:

```bash
cd frontend
pnpm install
```

Create a `.env.local` file in the `frontend` directory:

```env
NEXT_PUBLIC_API_URL=http://localhost:4000/api/v1
```

#### 4. Database Setup

Run migrations using Docker:

```bash
docker-compose up migrations
```

Or manually using Goose:

```bash
goose -dir ./migrations postgres "your-dsn-here" up
```

### Running the Application

#### Development Mode

**Backend:**
```bash
# From project root
go run ./cmd/api
```

The API will be available at `http://localhost:4000`

**Frontend:**
```bash
# From frontend directory
cd frontend
pnpm dev
```

The frontend will be available at `http://localhost:3000`

#### Production Mode with Docker

Start all services:

```bash
docker-compose up -d
```

This will start:
- **Frontend** on port 3000
- **Backend API** on port 4000
- **PostgreSQL** on port 5432
- **Migrations** will run automatically

## 📡 API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login user
- `GET /api/v1/auth/me` - Get current user (Protected)

### Properties
- `GET /api/v1/properties` - Get all properties
- `GET /api/v1/properties/{id}` - Get property by ID
- `GET /api/v1/properties/{id}/availability` - Check property availability
- `POST /api/v1/properties` - Create property (Protected: Admin/Host)
- `PUT /api/v1/properties/{id}` - Update property (Protected: Admin/Host)
- `DELETE /api/v1/properties/{id}` - Delete property (Protected: Admin/Host)
- `POST /api/v1/properties/{id}/images` - Add property images (Protected)
- `DELETE /api/v1/properties/{id}/images/{imageID}` - Delete property image (Protected)

### Bookings
- `GET /api/v1/bookings` - Get user's bookings (Protected)
- `GET /api/v1/bookings/{id}` - Get booking by ID (Protected)
- `POST /api/v1/bookings` - Create booking (Protected)
- `PATCH /api/v1/bookings/{id}` - Cancel booking (Protected)

### Amenities
- `GET /api/v1/amenities` - Get all amenities
- `POST /api/v1/amenities` - Create amenity (Protected: Admin/Host)
- `POST /api/v1/amenities/{propertyID}` - Add amenities to property (Protected)

### Health Check
- `GET /api/v1/healthz` - Health check endpoint

## 🔐 Authentication

The API uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

Tokens are obtained from the `/api/v1/auth/login` endpoint.

## 🗄️ Database Schema

### Users
- User accounts with role-based access (guest/admin)
- Secure password hashing with bcrypt

### Properties
- Property listings with details (title, location, price, description)
- Linked to users (hosts)
- Support for multiple images and amenities

### Bookings
- Booking records with date ranges
- Status tracking (booked/cancelled)
- Automatic price calculation

### Property Images
- Multiple images per property
- Caption and display order support

### Amenities
- Reusable amenities catalog
- Many-to-many relationship with properties

## 🧪 Development

### Running Tests

```bash
# Backend tests
go test ./...

# Frontend tests (if configured)
cd frontend
pnpm test
```

### Code Style

The project follows Go standard formatting. Use `go fmt` before committing.

```bash
go fmt ./...
```

### Building for Production

**Backend:**
```bash
CGO_ENABLED=0 go build -o cozystay ./cmd/api
```

**Frontend:**
```bash
cd frontend
pnpm build
```

## 🐳 Docker

The project includes Docker support for easy deployment:

- **Backend Dockerfile**: Multi-stage build for optimized production image
- **Frontend Dockerfile**: Optimized Next.js production build
- **docker-compose.yml**: Orchestrates all services

### Build Images

```bash
docker-compose build
```

### View Logs

```bash
docker-compose logs -f
```

### Stop Services

```bash
docker-compose down
```

## 📝 Environment Variables

### Backend (.env)
| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | API server port | `4000` |
| `ENV` | Environment (development/production) | `development` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database user | - |
| `DB_PASSWORD` | Database password | - |
| `DB_NAME` | Database name | - |
| `DSN` | Full database connection string | - |
| `JWT_SECRET` | Secret for JWT signing | - |
| `MAX_OPEN_CONNS` | Max database connections | `25` |
| `MAX_IDLE_CONNS` | Max idle connections | `5` |
| `MAX_IDLE_TIME` | Max idle time | `15m` |

### Frontend (.env.local)
| Variable | Description | Default |
|----------|-------------|---------|
| `NEXT_PUBLIC_API_URL` | Backend API URL | `http://localhost:4000/api/v1` |

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**Umesh Tiruvalluru**

- GitHub: [@Umesh-Tiruvalluru](https://github.com/Umesh-Tiruvalluru)

## 🙏 Acknowledgments

- Built with amazing open-source libraries
- UI components from [shadcn/ui](https://ui.shadcn.com/)
- Icons from [Lucide](https://lucide.dev/)

## 📚 Additional Resources

- [Go Documentation](https://go.dev/doc/)
- [Next.js Documentation](https://nextjs.org/docs)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Chi Router Documentation](https://github.com/go-chi/chi)

---
