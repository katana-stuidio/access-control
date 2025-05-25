# Access Control System

Welcome to the Access Control System documentation. This system provides a robust solution for managing user access and authentication.

## Overview

The Access Control System is a Go-based service that provides:
- User authentication and authorization
- JWT token management
- Tenant-based multi-tenancy
- Role-based access control

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Docker
- PostgreSQL

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/access-control.git
```

2. Build the project:
```bash
go build
```

3. Run with Docker:
```bash
docker-compose up
```

## API Documentation

### Authentication Endpoints
- POST `/api/v1/user/getjwt` - Get JWT token
- POST `/api/v1/user/validatejwt` - Validate JWT token
- POST `/api/v1/user/refreshjwt` - Refresh JWT token

### User Management
- GET `/api/v1/user/` - Get all users
- POST `/api/v1/user/` - Create user
- GET `/api/v1/user/{id}` - Get user by ID
- PATCH `/api/v1/user/{id}` - Update user
- DELETE `/api/v1/user/{id}` - Delete user

### Tenant Management
- GET `/api/v1/tenant/` - Get all tenants
- POST `/api/v1/tenant/` - Create tenant
- GET `/api/v1/tenant/{id}` - Get tenant by ID
- PATCH `/api/v1/tenant/{id}` - Update tenant
- DELETE `/api/v1/tenant/{id}` - Delete tenant

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 