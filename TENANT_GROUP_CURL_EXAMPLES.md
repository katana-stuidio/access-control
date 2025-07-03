# TenantGroup API cURL Examples

This document provides comprehensive cURL examples for testing all TenantGroup endpoints in the access-control system.

## Base URL
```
http://localhost:8080
```

## 1. Create Tenant Group

### Create a Municipal Education Department
```bash
curl -X POST http://localhost:8080/api/v1/tenant-groups/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Secretaria Municipal de Educação de São Paulo",
    "cnpj": "12.345.678/0001-90"
  }'
```

### Create a State Education Department
```bash
curl -X POST http://localhost:8080/api/v1/tenant-groups/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Secretaria de Educação do Estado de São Paulo",
    "cnpj": "23.456.789/0001-01"
  }'
```

### Create a Private School Network
```bash
curl -X POST http://localhost:8080/api/v1/tenant-groups/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Rede de Escolas Particulares DioSeno",
    "cnpj": "34.567.890/0001-12"
  }'
```

### Expected Response (201 Created)
```json
{
  "id": "7a1a5c8e-2bfb-4b5c-9b9f-1c0f2c5f30d1",
  "name": "Secretaria Municipal de Educação de São Paulo",
  "cnpj": "12.345.678/0001-90",
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

## 2. Get All Tenant Groups

### Get all tenant groups (default pagination)
```bash
curl -X GET http://localhost:8080/api/v1/tenant-groups/
```

### Get tenant groups with custom pagination
```bash
curl -X GET "http://localhost:8080/api/v1/tenant-groups/?limit=5&page=1"
```

### Get tenant groups with larger page size
```bash
curl -X GET "http://localhost:8080/api/v1/tenant-groups/?limit=20&page=1"
```

### Expected Response (200 OK)
```json
{
  "total": 3,
  "current_page": 1,
  "last_page": 1,
  "data": [
    {
      "id": "7a1a5c8e-2bfb-4b5c-9b9f-1c0f2c5f30d1",
      "name": "Secretaria Municipal de Educação de São Paulo",
      "cnpj": "12.345.678/0001-90",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "8b2b6d9f-3cgc-5c6d-0c0g-2d1g3d6g41e2",
      "name": "Secretaria de Educação do Estado de São Paulo",
      "cnpj": "23.456.789/0001-01",
      "is_active": true,
      "created_at": "2024-01-15T11:00:00Z",
      "updated_at": "2024-01-15T11:00:00Z"
    },
    {
      "id": "9c3c7e0g-4dhd-6d7e-1d1h-3e2h4e7h52f3",
      "name": "Rede de Escolas Particulares DioSeno",
      "cnpj": "34.567.890/0001-12",
      "is_active": true,
      "created_at": "2024-01-15T11:30:00Z",
      "updated_at": "2024-01-15T11:30:00Z"
    }
  ]
}
```

## 3. Get Tenant Group by ID

### Get specific tenant group
```bash
curl -X GET http://localhost:8080/api/v1/tenant-groups/7a1a5c8e-2bfb-4b5c-9b9f-1c0f2c5f30d1
```

### Expected Response (200 OK)
```json
{
  "id": "7a1a5c8e-2bfb-4b5c-9b9f-1c0f2c5f30d1",
  "name": "Secretaria Municipal de Educação de São Paulo",
  "cnpj": "12.345.678/0001-90",
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Get non-existent tenant group (404 Not Found)
```bash
curl -X GET http://localhost:8080/api/v1/tenant-groups/00000000-0000-0000-0000-000000000000
```

### Expected Response (404 Not Found)
```json
{
  "error": "Tenant group not found"
}
```

## 4. Update Tenant Group

### Update tenant group name and status
```bash
curl -X PUT http://localhost:8080/api/v1/tenant-groups/7a1a5c8e-2bfb-4b5c-9b9f-1c0f2c5f30d1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Secretaria Municipal de Educação de São Paulo - Atualizada",
    "cnpj": "12.345.678/0001-90",
    "is_active": true
  }'
```

### Update tenant group to inactive
```bash
curl -X PUT http://localhost:8080/api/v1/tenant-groups/8b2b6d9f-3cgc-5c6d-0c0g-2d1g3d6g41e2 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Secretaria de Educação do Estado de São Paulo",
    "cnpj": "23.456.789/0001-01",
    "is_active": false
  }'
```

### Expected Response (200 OK)
```json
{
  "message": "Tenant group updated successfully"
}
```

### Update with invalid CNPJ (400 Bad Request)
```bash
curl -X PUT http://localhost:8080/api/v1/tenant-groups/7a1a5c8e-2bfb-4b5c-9b9f-1c0f2c5f30d1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Invalid CNPJ Test",
    "cnpj": "invalid-cnpj",
    "is_active": true
  }'
```

### Expected Response (400 Bad Request)
```json
{
  "error": "Invalid CNPJ format"
}
```

## 5. Delete Tenant Group

### Delete tenant group
```bash
curl -X DELETE http://localhost:8080/api/v1/tenant-groups/9c3c7e0g-4dhd-6d7e-1d1h-3e2h4e7h52f3
```

### Expected Response (200 OK)
```json
{
  "message": "Tenant group deleted successfully"
}
```

### Delete non-existent tenant group (404 Not Found)
```bash
curl -X DELETE http://localhost:8080/api/v1/tenant-groups/00000000-0000-0000-0000-000000000000
```

### Expected Response (404 Not Found)
```json
{
  "error": "Tenant group not found"
}
```

## 6. Error Handling Examples

### Create tenant group with duplicate CNPJ (409 Conflict)
```bash
curl -X POST http://localhost:8080/api/v1/tenant-groups/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Duplicate CNPJ Test",
    "cnpj": "12.345.678/0001-90"
  }'
```

### Expected Response (409 Conflict)
```json
{
  "error": "CNPJ already exists"
}
```

### Create tenant group with invalid CNPJ (400 Bad Request)
```bash
curl -X POST http://localhost:8080/api/v1/tenant-groups/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Invalid CNPJ Test",
    "cnpj": "123456789"
  }'
```

### Expected Response (400 Bad Request)
```json
{
  "error": "Invalid CNPJ format"
}
```

### Create tenant group with missing required fields (400 Bad Request)
```bash
curl -X POST http://localhost:8080/api/v1/tenant-groups/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Missing CNPJ Test"
  }'
```

### Expected Response (400 Bad Request)
```json
{
  "error": "Invalid request data"
}
```

## 7. Integration Examples

### Complete Workflow: Create Group, Create Tenant, Create User, Get JWT

#### Step 1: Create Tenant Group
```bash
curl -X POST http://localhost:8080/api/v1/tenant-groups/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Secretaria Municipal de Educação",
    "cnpj": "12.345.678/0001-90"
  }'
```

#### Step 2: Create Tenant with Group (using group ID from step 1)
```bash
curl -X POST http://localhost:8080/api/v1/tenants/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Escola Municipal Santo Antônio",
    "cnpj": "98.765.432/0001-10",
    "group_id": "7a1a5c8e-2bfb-4b5c-9b9f-1c0f2c5f30d1"
  }'
```

#### Step 3: Create User (using tenant CNPJ from step 2)
```bash
curl -X POST http://localhost:8080/api/v1/user/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Professor Clóvis Silva",
    "username": "clovis.professor",
    "password": "Senha123!",
    "cnpj": "98.765.432/0001-10",
    "email": "clovis@escola.edu.br",
    "role": "Professor"
  }'
```

#### Step 4: Get JWT Token (will include group information)
```bash
curl -X POST http://localhost:8080/api/v1/user/getjwt \
  -H "Content-Type: application/json" \
  -d '{
    "username": "clovis.professor",
    "password": "Senha123!"
  }'
```

### Expected JWT Response (200 OK)
```json
{
  "accessToken": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "tokenId": "11fbab0e44e8d9489a17c4ee844cb9ca"
}
```

### Decoded JWT Payload (example)
```json
{
  "username": "clovis.professor",
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

## 8. Testing Script

### Complete test script for all endpoints
```bash
#!/bin/bash

BASE_URL="http://localhost:8080"
API_BASE="$BASE_URL/api/v1/tenant-groups"

echo "=== Testing TenantGroup API ==="

# Test 1: Create tenant group
echo "1. Creating tenant group..."
GROUP_RESPONSE=$(curl -s -X POST "$API_BASE/" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Education Department",
    "cnpj": "12.345.678/0001-90"
  }')

echo "Response: $GROUP_RESPONSE"

# Extract group ID from response (you might need to adjust this based on your JSON parsing)
GROUP_ID=$(echo $GROUP_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

echo "Group ID: $GROUP_ID"

# Test 2: Get all tenant groups
echo "2. Getting all tenant groups..."
curl -s -X GET "$API_BASE/"

# Test 3: Get specific tenant group
echo "3. Getting specific tenant group..."
curl -s -X GET "$API_BASE/$GROUP_ID"

# Test 4: Update tenant group
echo "4. Updating tenant group..."
curl -s -X PUT "$API_BASE/$GROUP_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Test Education Department",
    "cnpj": "12.345.678/0001-90",
    "is_active": true
  }'

# Test 5: Delete tenant group
echo "5. Deleting tenant group..."
curl -s -X DELETE "$API_BASE/$GROUP_ID"

echo "=== Testing completed ==="
```

## Notes

- Replace `localhost:8080` with your actual server URL if different
- The UUIDs in the examples are placeholders; use actual UUIDs returned by your API
- All CNPJ numbers in examples are fictional; use valid CNPJ numbers for testing
- The JWT token examples show the structure; actual tokens will be much longer
- Remember to handle authentication if your API requires it 