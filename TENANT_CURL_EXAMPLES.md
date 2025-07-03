# Tenant API cURL Examples

This document provides comprehensive cURL examples for testing all Tenant endpoints in the access-control system, including the new group functionality.

## Base URL
```
http://localhost:8080
```

## 1. Create Tenant

### Create Tenant without Group (Backward Compatibility)
```bash
curl -X POST http://localhost:8080/api/v1/tenants/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Escola Municipal João da Silva",
    "cnpj": "11.111.111/0001-11"
  }'
```

### Create Tenant with Group
```bash
curl -X POST http://localhost:8080/api/v1/tenants/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Escola Municipal Santo Antônio",
    "cnpj": "22.222.222/0001-22",
    "group_id": "7a1a5c8e-2bfb-4b5c-9b9f-1c0f2c5f30d1"
  }'
```

### Create Multiple Tenants for Same Group
```bash
# First tenant in the group
curl -X POST http://localhost:8080/api/v1/tenants/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Escola Municipal São José",
    "cnpj": "33.333.333/0001-33",
    "group_id": "7a1a5c8e-2bfb-4b5c-9b9f-1c0f2c5f30d1"
  }'

# Second tenant in the same group
curl -X POST http://localhost:8080/api/v1/tenants/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Escola Municipal Santa Maria",
    "cnpj": "44.444.444/0001-44",
    "group_id": "7a1a5c8e-2bfb-4b5c-9b9f-1c0f2c5f30d1"
  }'
```

### Create Private School Tenant
```bash
curl -X POST http://localhost:8080/api/v1/tenants/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Colégio DioSeno - Unidade Centro",
    "cnpj": "55.555.555/0001-55",
    "group_id": "9c3c7e0g-4dhd-6d7e-1d1h-3e2h4e7h52f3"
  }'
```

### Expected Response (201 Created)
```json
{
  "id": "3204fdce-560b-4f19-9bc8-875825662a4a",
  "group_id": "7a1a5c8e-2bfb-4b5c-9b9f-1c0f2c5f30d1",
  "name": "Escola Municipal Santo Antônio",
  "cnpj": "22.222.222/0001-22",
  "schema_name": "3204fdce-560b-4f19-9bc8-875825662a4a",
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

## 2. Get All Tenants

### Get all tenants (default pagination)
```bash
curl -X GET http://localhost:8080/api/v1/tenants/
```

### Get tenants with custom pagination
```bash
curl -X GET "http://localhost:8080/api/v1/tenants/?limit=5&page=1"
```

### Get tenants with larger page size
```bash
curl -X GET "http://localhost:8080/api/v1/tenants/?limit=20&page=1"
```

### Expected Response (200 OK)
```json
{
  "total": 4,
  "current_page": 1,
  "last_page": 1,
  "data": [
    {
      "id": "11111111-1111-1111-1111-111111111111",
      "group_id": null,
      "name": "Escola Municipal João da Silva",
      "cnpj": "11.111.111/0001-11",
      "schema_name": "11111111-1111-1111-1111-111111111111",
      "is_active": true,
      "created_at": "2024-01-15T09:00:00Z",
      "updated_at": "2024-01-15T09:00:00Z"
    },
    {
      "id": "3204fdce-560b-4f19-9bc8-875825662a4a",
      "group_id": "7a1a5c8e-2bfb-4b5c-9b9f-1c0f2c5f30d1",
      "name": "Escola Municipal Santo Antônio",
      "cnpj": "22.222.222/0001-22",
      "schema_name": "3204fdce-560b-4f19-9bc8-875825662a4a",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "44444444-4444-4444-4444-444444444444",
      "group_id": "7a1a5c8e-2bfb-4b5c-9b9f-1c0f2c5f30d1",
      "name": "Escola Municipal São José",
      "cnpj": "33.333.333/0001-33",
      "schema_name": "44444444-4444-4444-4444-444444444444",
      "is_active": true,
      "created_at": "2024-01-15T11:00:00Z",
      "updated_at": "2024-01-15T11:00:00Z"
    },
    {
      "id": "55555555-5555-5555-5555-555555555555",
      "group_id": "9c3c7e0g-4dhd-6d7e-1d1h-3e2h4e7h52f3",
      "name": "Colégio DioSeno - Unidade Centro",
      "cnpj": "55.555.555/0001-55",
      "schema_name": "55555555-5555-5555-5555-555555555555",
      "is_active": true,
      "created_at": "2024-01-15T11:30:00Z",
      "updated_at": "2024-01-15T11:30:00Z"
    }
  ]
}
```

## 3. Get Tenant by ID

### Get specific tenant
```bash
curl -X GET http://localhost:8080/api/v1/tenants/3204fdce-560b-4f19-9bc8-875825662a4a
```

### Expected Response (200 OK)
```json
{
  "id": "3204fdce-560b-4f19-9bc8-875825662a4a",
  "group_id": "7a1a5c8e-2bfb-4b5c-9b9f-1c0f2c5f30d1",
  "name": "Escola Municipal Santo Antônio",
  "cnpj": "22.222.222/0001-22",
  "schema_name": "3204fdce-560b-4f19-9bc8-875825662a4a",
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Get tenant without group
```bash
curl -X GET http://localhost:8080/api/v1/tenants/11111111-1111-1111-1111-111111111111
```

### Expected Response (200 OK) - Tenant without Group
```json
{
  "id": "11111111-1111-1111-1111-111111111111",
  "group_id": null,
  "name": "Escola Municipal João da Silva",
  "cnpj": "11.111.111/0001-11",
  "schema_name": "11111111-1111-1111-1111-111111111111",
  "is_active": true,
  "created_at": "2024-01-15T09:00:00Z",
  "updated_at": "2024-01-15T09:00:00Z"
}
```

### Get non-existent tenant (404 Not Found)
```bash
curl -X GET http://localhost:8080/api/v1/tenants/00000000-0000-0000-0000-000000000000
```

### Expected Response (404 Not Found)
```json
{
  "error": "Tenant not found"
}
```

## 4. Update Tenant

### Update tenant name and CNPJ
```bash
curl -X PUT http://localhost:8080/api/v1/tenants/3204fdce-560b-4f19-9bc8-875825662a4a \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Escola Municipal Santo Antônio - Atualizada",
    "cnpj": "22.222.222/0001-22",
    "schema_name": "3204fdce-560b-4f19-9bc8-875825662a4a",
    "is_active": true
  }'
```

### Update tenant to add group (tenant was previously without group)
```bash
curl -X PUT http://localhost:8080/api/v1/tenants/11111111-1111-1111-1111-111111111111 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Escola Municipal João da Silva",
    "cnpj": "11.111.111/0001-11",
    "schema_name": "11111111-1111-1111-1111-111111111111",
    "is_active": true,
    "group_id": "7a1a5c8e-2bfb-4b5c-9b9f-1c0f2c5f30d1"
  }'
```

### Update tenant to remove group (set to null)
```bash
curl -X PUT http://localhost:8080/api/v1/tenants/3204fdce-560b-4f19-9bc8-875825662a4a \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Escola Municipal Santo Antônio",
    "cnpj": "22.222.222/0001-22",
    "schema_name": "3204fdce-560b-4f19-9bc8-875825662a4a",
    "is_active": true,
    "group_id": null
  }'
```

### Update tenant to change group
```bash
curl -X PUT http://localhost:8080/api/v1/tenants/3204fdce-560b-4f19-9bc8-875825662a4a \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Escola Municipal Santo Antônio",
    "cnpj": "22.222.222/0001-22",
    "schema_name": "3204fdce-560b-4f19-9bc8-875825662a4a",
    "is_active": true,
    "group_id": "9c3c7e0g-4dhd-6d7e-1d1h-3e2h4e7h52f3"
  }'
```

### Update tenant to inactive
```bash
curl -X PUT http://localhost:8080/api/v1/tenants/3204fdce-560b-4f19-9bc8-875825662a4a \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Escola Municipal Santo Antônio",
    "cnpj": "22.222.222/0001-22",
    "schema_name": "3204fdce-560b-4f19-9bc8-875825662a4a",
    "is_active": false
  }'
```

### Expected Response (200 OK)
```json
{
  "message": "Tenant updated successfully"
}
```

### Update with invalid CNPJ (400 Bad Request)
```bash
curl -X PUT http://localhost:8080/api/v1/tenants/3204fdce-560b-4f19-9bc8-875825662a4a \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Invalid CNPJ Test",
    "cnpj": "invalid-cnpj",
    "schema_name": "3204fdce-560b-4f19-9bc8-875825662a4a",
    "is_active": true
  }'
```

### Expected Response (400 Bad Request)
```json
{
  "error": "Invalid CNPJ format"
}
```

## 5. Delete Tenant

### Delete tenant
```bash
curl -X DELETE http://localhost:8080/api/v1/tenants/3204fdce-560b-4f19-9bc8-875825662a4a
```

### Expected Response (200 OK)
```json
{
  "message": "Tenant deleted successfully"
}
```

### Delete non-existent tenant (404 Not Found)
```bash
curl -X DELETE http://localhost:8080/api/v1/tenants/00000000-0000-0000-0000-000000000000
```

### Expected Response (404 Not Found)
```json
{
  "error": "Tenant not found"
}
```

## 6. Error Handling Examples

### Create tenant with duplicate CNPJ (409 Conflict)
```bash
curl -X POST http://localhost:8080/api/v1/tenants/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Duplicate CNPJ Test",
    "cnpj": "22.222.222/0001-22"
  }'
```

### Expected Response (409 Conflict)
```json
{
  "error": "CNPJ already exists"
}
```

### Create tenant with invalid CNPJ (400 Bad Request)
```bash
curl -X POST http://localhost:8080/api/v1/tenants/ \
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

### Create tenant with non-existent group (400 Bad Request)
```bash
curl -X POST http://localhost:8080/api/v1/tenants/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Non-existent Group Test",
    "cnpj": "66.666.666/0001-66",
    "group_id": "00000000-0000-0000-0000-000000000000"
  }'
```

### Expected Response (400 Bad Request)
```json
{
  "error": "Group not found"
}
```

### Create tenant with missing required fields (400 Bad Request)
```bash
curl -X POST http://localhost:8080/api/v1/tenants/ \
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

#### Step 1: Create Tenant Group (if not exists)
```bash
curl -X POST http://localhost:8080/api/v1/tenant-groups/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Secretaria Municipal de Educação",
    "cnpj": "12.345.678/0001-90"
  }'
```

#### Step 2: Create Tenant with Group
```bash
curl -X POST http://localhost:8080/api/v1/tenants/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Escola Municipal Santo Antônio",
    "cnpj": "98.765.432/0001-10",
    "group_id": "7a1a5c8e-2bfb-4b5c-9b9f-1c0f2c5f30d1"
  }'
```

#### Step 3: Create User (using tenant CNPJ)
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

## 8. Group Management Examples

### Get all tenants in a specific group
```bash
# First, get all tenants and filter by group_id in your application
curl -X GET "http://localhost:8080/api/v1/tenants/?limit=100&page=1"
```

### Move tenant between groups
```bash
# Move tenant from one group to another
curl -X PUT http://localhost:8080/api/v1/tenants/3204fdce-560b-4f19-9bc8-875825662a4a \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Escola Municipal Santo Antônio",
    "cnpj": "22.222.222/0001-22",
    "schema_name": "3204fdce-560b-4f19-9bc8-875825662a4a",
    "is_active": true,
    "group_id": "9c3c7e0g-4dhd-6d7e-1d1h-3e2h4e7h52f3"
  }'
```

### Remove tenant from group (set group_id to null)
```bash
curl -X PUT http://localhost:8080/api/v1/tenants/3204fdce-560b-4f19-9bc8-875825662a4a \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Escola Municipal Santo Antônio",
    "cnpj": "22.222.222/0001-22",
    "schema_name": "3204fdce-560b-4f19-9bc8-875825662a4a",
    "is_active": true,
    "group_id": null
  }'
```

## 9. Testing Script

### Complete test script for all tenant endpoints
```bash
#!/bin/bash

BASE_URL="http://localhost:8080"
API_BASE="$BASE_URL/api/v1/tenants"

echo "=== Testing Tenant API ==="

# Test 1: Create tenant without group
echo "1. Creating tenant without group..."
TENANT_RESPONSE=$(curl -s -X POST "$API_BASE/" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test School Without Group",
    "cnpj": "11.111.111/0001-11"
  }')

echo "Response: $TENANT_RESPONSE"

# Extract tenant ID from response
TENANT_ID=$(echo $TENANT_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

echo "Tenant ID: $TENANT_ID"

# Test 2: Get all tenants
echo "2. Getting all tenants..."
curl -s -X GET "$API_BASE/"

# Test 3: Get specific tenant
echo "3. Getting specific tenant..."
curl -s -X GET "$API_BASE/$TENANT_ID"

# Test 4: Update tenant
echo "4. Updating tenant..."
curl -s -X PUT "$API_BASE/$TENANT_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Test School",
    "cnpj": "11.111.111/0001-11",
    "schema_name": "'$TENANT_ID'",
    "is_active": true
  }'

# Test 5: Delete tenant
echo "5. Deleting tenant..."
curl -s -X DELETE "$API_BASE/$TENANT_ID"

echo "=== Testing completed ==="
```

## 10. Advanced Scenarios

### Bulk Tenant Creation for a Group
```bash
# Create multiple tenants for the same group
for i in {1..5}; do
  curl -X POST http://localhost:8080/api/v1/tenants/ \
    -H "Content-Type: application/json" \
    -d '{
      "name": "Escola Municipal Teste '$i'",
      "cnpj": "'$i$i'.'$i$i$i'.'$i$i$i'/'0001-'$i$i",
      "group_id": "7a1a5c8e-2bfb-4b5c-9b9f-1c0f2c5f30d1"
    }'
done
```

### Tenant Migration Between Groups
```bash
# Move all tenants from one group to another
# First, get all tenants in the source group
TENANTS=$(curl -s -X GET "http://localhost:8080/api/v1/tenants/?limit=100&page=1" | jq -r '.data[] | select(.group_id == "7a1a5c8e-2bfb-4b5c-9b9f-1c0f2c5f30d1") | .id')

# Then update each tenant to the new group
for tenant_id in $TENANTS; do
  curl -X PUT http://localhost:8080/api/v1/tenants/$tenant_id \
    -H "Content-Type: application/json" \
    -d '{
      "group_id": "9c3c7e0g-4dhd-6d7e-1d1h-3e2h4e7h52f3"
    }'
done
```

## Notes

- Replace `localhost:8080` with your actual server URL if different
- The UUIDs in the examples are placeholders; use actual UUIDs returned by your API
- All CNPJ numbers in examples are fictional; use valid CNPJ numbers for testing
- The `group_id` field is optional - tenants can exist without being part of a group
- When updating a tenant, you can change its group association or remove it entirely
- The `schema_name` field is automatically generated from the tenant ID
- Remember to handle authentication if your API requires it
- The JWT token examples show the structure; actual tokens will be much longer 