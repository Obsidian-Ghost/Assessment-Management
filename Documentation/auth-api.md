# Authentication API

This document describes the Authentication API endpoints for the Assessment Management System.

## Base URL

All endpoints are relative to the base URL `/api`.

## Authentication Endpoints

### Login

Authenticates a user and returns a JWT token.

**Endpoint:** `POST /auth/login`

**Authentication Required:** No

**Request Body:**

```json
{
  "email": "admin@example.com",
  "password": "admin123"
}
```

**Response:**

Status Code: 200 OK

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZTQxNzU2ZmUtNTI0Mi00Njk4LTgxMDItNzM1NjIxYjY5OWQ4Iiwib3JnYW5pemF0aW9uX2lkIjoiOTI2NDFhN2QtOTY2ZS00ZTI5LThkNTItMWQ4ZWE2Y2FjNTMwIiwiZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsImZpcnN0X25hbWUiOiJTeXN0ZW0iLCJsYXN0X25hbWUiOiJBZG1pbmlzdHJhdG9yIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzQzMzMyMTU2LCJuYmYiOjE3NDMyNDU3NTYsImlhdCI6MTc0MzI0NTc1Nn0.0XCv8NMimhX9Y27gBYLauJ0KglE8mkWpegfDNWAynCo",
  "refresh_token": "pjVT1AV_vlPwURRMqIE-6LSWGEAr8WCxhiTLQvEXdBs",
  "expires_in": 86400,
  "user": {
    "id": "e41756fe-5242-4698-8102-735621b699d8",
    "organization_id": "92641a7d-966e-4e29-8d52-1d8ea6cac530",
    "email": "admin@example.com",
    "first_name": "System",
    "last_name": "Administrator",
    "role": "admin",
    "created_at": "2025-03-28T11:02:23.175179Z",
    "updated_at": "2025-03-28T11:12:19.589734Z"
  }
}
```

**Error Responses:**

Status Code: 400 Bad Request - Invalid request body
Status Code: 401 Unauthorized - Invalid credentials

```json
{
  "message": "Invalid credentials"
}
```

### Get Current User

Returns information about the authenticated user.

**Endpoint:** `GET /auth/me`

**Authentication Required:** Yes

**Response:**

Status Code: 200 OK

```json
{
  "id": "e41756fe-5242-4698-8102-735621b699d8",
  "organization_id": "92641a7d-966e-4e29-8d52-1d8ea6cac530",
  "email": "admin@example.com",
  "first_name": "System",
  "last_name": "Administrator",
  "role": "admin",
  "created_at": "2025-03-28T11:02:23.175179Z",
  "updated_at": "2025-03-28T11:12:19.589734Z"
}
```

**Error Responses:**

Status Code: 401 Unauthorized - Missing or invalid token

```json
{
  "message": "Unauthorized"
}
```

### Change Password

Allows a user to change their password.

**Endpoint:** `POST /auth/change-password`

**Authentication Required:** Yes

**Request Body:**

```json
{
  "current_password": "admin123",
  "new_password": "newSecurePassword123",
  "confirm_password": "newSecurePassword123"
}
```

**Response:**

Status Code: 200 OK

```json
{
  "message": "Password updated successfully"
}
```

**Error Responses:**

Status Code: 400 Bad Request - Invalid request body or passwords don't match

```json
{
  "message": "New password and confirm password do not match"
}
```

Status Code: 401 Unauthorized - Current password is incorrect

```json
{
  "message": "Current password is incorrect"
}
```

## Using the Authentication Token

After successful login, the JWT token should be included in the Authorization header for all authenticated requests:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZTQxNzU2ZmUtNTI0Mi00Njk4LTgxMDItNzM1NjIxYjY5OWQ4Iiwib3JnYW5pemF0aW9uX2lkIjoiOTI2NDFhN2QtOTY2ZS00ZTI5LThkNTItMWQ4ZWE2Y2FjNTMwIiwiZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsImZpcnN0X25hbWUiOiJTeXN0ZW0iLCJsYXN0X25hbWUiOiJBZG1pbmlzdHJhdG9yIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzQzMzMyMTU2LCJuYmYiOjE3NDMyNDU3NTYsImlhdCI6MTc0MzI0NTc1Nn0.0XCv8NMimhX9Y27gBYLauJ0KglE8mkWpegfDNWAynCo
```

## Token Structure

The JWT token contains the following claims:

| Claim           | Description                           |
|-----------------|---------------------------------------|
| user_id         | Unique identifier for the user        |
| organization_id | Organization the user belongs to      |
| email           | User's email address                  |
| first_name      | User's first name                     |
| last_name       | User's last name                      |
| role            | User's role (admin, teacher, student) |
| exp             | Token expiration timestamp            |
| nbf             | Not before timestamp                  |
| iat             | Issued at timestamp                   |

## Token Expiration

Access tokens expire 24 hours after issuance. Refresh tokens expire 7 days after issuance. When an access token expires, it can be refreshed using the refresh token without requiring the user to log in again.

## Refresh Token Endpoints

### Refresh Access Token

Generates a new access token using a refresh token.

**Endpoint:** `POST /auth/token/refresh`

**Authentication Required:** No

**Request Body:**

```json
{
  "refresh_token": "pjVT1AV_vlPwURRMqIE-6LSWGEAr8WCxhiTLQvEXdBs"
}
```

**Response:**

Status Code: 200 OK

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZTQxNzU2ZmUtNTI0Mi00Njk4LTgxMDItNzM1NjIxYjY5OWQ4Iiwib3JnYW5pemF0aW9uX2lkIjoiOTI2NDFhN2QtOTY2ZS00ZTI5LThkNTItMWQ4ZWE2Y2FjNTMwIiwiZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsImZpcnN0X25hbWUiOiJTeXN0ZW0iLCJsYXN0X25hbWUiOiJBZG1pbmlzdHJhdG9yIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzQzMzMyMTU2LCJuYmYiOjE3NDMyNDU3NTYsImlhdCI6MTc0MzI0NTc1Nn0.0XCv8NMimhX9Y27gBYLauJ0KglE8mkWpegfDNWAynCo",
  "refresh_token": "pjVT1AV_vlPwURRMqIE-6LSWGEAr8WCxhiTLQvEXdBs",
  "expires_in": 86400
}
```

**Error Responses:**

Status Code: 400 Bad Request - Invalid request body
Status Code: 401 Unauthorized - Invalid, expired, or revoked refresh token

```json
{
  "message": "Invalid refresh token"
}
```

### Revoke Refresh Token

Revokes a specific refresh token so it can no longer be used.

**Endpoint:** `POST /auth/token/revoke`

**Authentication Required:** Yes

**Request Body:**

```json
{
  "refresh_token": "pjVT1AV_vlPwURRMqIE-6LSWGEAr8WCxhiTLQvEXdBs"
}
```

**Response:**

Status Code: 200 OK

```json
{
  "message": "Token revoked successfully"
}
```

**Error Responses:**

Status Code: 400 Bad Request - Invalid request body
Status Code: 401 Unauthorized - User not authenticated

### Revoke All Refresh Tokens

Revokes all refresh tokens for the authenticated user.

**Endpoint:** `POST /auth/token/revoke-all`

**Authentication Required:** Yes

**Response:**

Status Code: 200 OK

```json
{
  "message": "All tokens revoked successfully"
}
```

**Error Responses:**

Status Code: 401 Unauthorized - User not authenticated