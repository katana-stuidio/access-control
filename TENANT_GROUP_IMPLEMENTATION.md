# Tenant Group Implementation

This document describes the implementation of tenant groups in the access-control system, which allows for hierarchical organization of tenants (schools) under groups (departments or networks).

## Overview

The tenant group functionality enables:
- **Educational Groups**: Multiple schools under a single administrative group
- **Department Hierarchies**: Federal, state, and municipal education departments
- **Enhanced JWT Tokens**: Include both tenant and group information for better access control

## Database Schema Changes

### New Table: `tb_tenant_group`
```sql
CREATE TABLE public.tb_tenant_group (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    cnpj VARCHAR(20) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
```

### Updated Table: `tb_tenant`
```sql
-- Added group_id column
ALTER TABLE public.tb_tenant 
ADD COLUMN group_id UUID;

-- Added foreign key constraint
ALTER TABLE public.tb_tenant 
ADD CONSTRAINT fk_tenant_group 
FOREIGN KEY (group_id) 
REFERENCES public.tb_tenant_group(id) 
ON UPDATE CASCADE 
ON DELETE SET NULL;
```

## Model Changes

### New Model: `TenantGroup`
```go
type TenantGroup struct {
    ID        uuid.UUID `json:"id"`
    Name      string    `json:"name"`
    CNPJ      string    `json:"cnpj"`
    IsActive  bool      `json:"is_active"`
    CreatedAt time.Time `json:"created_at,omitempty"`
    UpdatedAt time.Time `json:"updated_at,omitempty"`
}
```

### Updated Model: `Tenant`
```go
type Tenant struct {
    ID         uuid.UUID  `json:"id"`
    GroupID    *uuid.UUID `json:"group_id,omitempty"` // optional
    CNPJ       string     `json:"cnpj"`
    Name       string     `json:"name"`
    SchemaName string     `json:"schema_name"`
    IsActive   bool       `json:"is_active"`
    CreatedAt  time.Time  `json:"created_at,omitempty"`
    UpdatedAt  time.Time  `json:"updated_at,omitempty"`
}
```

## JWT Token Enhancement

### Updated Claims Structure
```go
type Claims struct {
    Username     string `json:"username"`
    UserID       string `json:"user_id"`
    TenantID     string `json:"tenant_id"`
    TenantName   string `json:"tenant_name,omitempty"`
    GroupID      string `json:"group_id,omitempty"`
    GroupName    string `json:"group_name,omitempty"`
    Role         string `json:"role"`
    FirstAccess  bool   `json:"first_access"`
    Renew        bool   `json:"renew,omitempty"`
    TokenID      string `json:"token_id,omitempty"`
    jwt.RegisteredClaims
}
```

### Example JWT Payload
```json
{
    "username": "clovis",
    "user_id": "db77d7d0-5bec-4709-bf83-1880556ec446",
    "tenant_id": "3204fdce-560b-4f19-9bc8-875825662a4a",
    "tenant_name": "Escola Municipal Santo Antônio",
    "group_id": "7a1a5c8e-2bfb-4b5c-9b9f-1c0f2c5f30d1",
    "group_name": "Secretaria Municipal de Educação",
    "role": "Professor",
    "first_access": false,
    "token_id": "11fbab0e44e8d9489a17c4ee844cb9ca",
    "exp": 1751512962,
    "iat": 1751505762
}
```

## API Endpoints

### Tenant Group Management

#### Create Tenant Group
```http
POST /api/v1/tenant-groups/
Content-Type: application/json

{
    "name": "Secretaria Municipal de Educação",
    "cnpj": "12.345.678/0001-90"
}
```

#### Get All Tenant Groups
```http
GET /api/v1/tenant-groups/?limit=10&page=1
```

#### Get Tenant Group by ID
```http
GET /api/v1/tenant-groups/{id}
```

#### Update Tenant Group
```http
PUT /api/v1/tenant-groups/{id}
Content-Type: application/json

{
    "name": "Secretaria Municipal de Educação Atualizada",
    "cnpj": "12.345.678/0001-90",
    "is_active": true
}
```

#### Delete Tenant Group
```http
DELETE /api/v1/tenant-groups/{id}
```

### Updated Tenant Management

#### Create Tenant with Group
```http
POST /api/v1/tenants/
Content-Type: application/json

{
    "name": "Escola Municipal Santo Antônio",
    "cnpj": "98.765.432/0001-10",
    "group_id": "7a1a5c8e-2bfb-4b5c-9b9f-1c0f2c5f30d1"
}
```

## Implementation Details

### Service Layer
- **TenantGroupService**: Handles CRUD operations for tenant groups
- **Updated TenantService**: Now includes group_id in all operations
- **Updated UserService**: Enhanced authentication to fetch tenant and group data

### Handler Layer
- **TenantGroupHandler**: Manages HTTP requests for tenant group operations
- **Updated UserHandler**: Enhanced JWT generation with tenant and group information

### Database Operations
- All tenant queries now include `group_id` field
- Proper foreign key constraints ensure data integrity
- Indexes on `group_id` and `cnpj` for optimal performance

## Migration Guide

### 1. Run Database Migration
```bash
# Execute the migration script
psql -d your_database -f migrate/tenant_group_schema.sql
```

### 2. Update Application Code
The application code has been updated to include:
- New tenant group models and services
- Enhanced JWT token generation
- Updated API endpoints

### 3. Test the Implementation
```bash
# Test tenant group creation
curl -X POST http://localhost:8080/api/v1/tenant-groups/ \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Group", "cnpj": "12.345.678/0001-90"}'

# Test tenant creation with group
curl -X POST http://localhost:8080/api/v1/tenants/ \
  -H "Content-Type: application/json" \
  -d '{"name": "Test School", "cnpj": "98.765.432/0001-10", "group_id": "group-uuid-here"}'

# Test user authentication (JWT will now include group info)
curl -X POST http://localhost:8080/api/v1/user/getjwt \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "testpass"}'
```

## Benefits

1. **Hierarchical Organization**: Schools can be grouped under departments or networks
2. **Enhanced Access Control**: JWT tokens include both tenant and group context
3. **Audit Trail**: Complete visibility of user's organizational hierarchy
4. **Scalability**: Supports complex organizational structures
5. **Backward Compatibility**: Existing tenants without groups continue to work

## Security Considerations

- Group information is included in JWT tokens for authorization decisions
- Foreign key constraints prevent orphaned tenant-group relationships
- CNPJ validation ensures data integrity
- Optional group association maintains backward compatibility

## Future Enhancements

- Group-level permissions and roles
- Cross-group user management
- Group-specific configurations
- Advanced reporting and analytics by group 