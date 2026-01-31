---
title: API Specification
---

# API Specification

This document describes the REST API endpoints and their behavior.

## Authentication

All API requests require authentication using Bearer tokens:

```http
Authorization: Bearer <token>
```

## Endpoints

### GET /api/users

Returns a list of all users.

**Response:**

```json
{
  "users": [
    {
      "id": "123",
      "name": "John Doe",
      "email": "john@example.com"
    }
  ]
}
```

### POST /api/users

Creates a new user.

**Request Body:**

```json
{
  "name": "Jane Doe",
  "email": "jane@example.com"
}
```

**Response:** `201 Created`

### GET /api/users/:id

Returns a specific user by ID.

### PUT /api/users/:id

Updates a user.

### DELETE /api/users/:id

Deletes a user.

## Error Handling

All errors follow this format:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "User not found"
  }
}
```

## Rate Limiting

- 100 requests per minute per API key
- Rate limit headers included in all responses
