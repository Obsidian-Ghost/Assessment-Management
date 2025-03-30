# Authentication and Authorization

This document describes the authentication and authorization mechanisms used in the Assessment Management System.

## Authentication

The system uses JSON Web Tokens (JWT) for authentication:

### JWT Structure

Each JWT contains the following claims:

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

### Authentication Flow

1. User sends credentials (email and password) to `/api/auth/login`
2. System validates credentials and generates both an access token and a refresh token
3. Client stores both tokens (typically in localStorage, with appropriate security measures)
4. Client includes the access token in the Authorization header of subsequent requests
5. Server validates the token and extracts user information for each request
6. When the access token expires, client can use the refresh token to obtain a new access token without re-authentication
7. If needed, clients can explicitly revoke refresh tokens for security purposes

### Example Authentication Header

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZTQxNzU2ZmUtNTI0Mi00Njk4LTgxMDItNzM1NjIxYjY5OWQ4Iiwib3JnYW5pemF0aW9uX2lkIjoiOTI2NDFhN2QtOTY2ZS00ZTI5LThkNTItMWQ4ZWE2Y2FjNTMwIiwiZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsImZpcnN0X25hbWUiOiJTeXN0ZW0iLCJsYXN0X25hbWUiOiJBZG1pbmlzdHJhdG9yIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzQzMzMyMTU2LCJuYmYiOjE3NDMyNDU3NTYsImlhdCI6MTc0MzI0NTc1Nn0.0XCv8NMimhX9Y27gBYLauJ0KglE8mkWpegfDNWAynCo
```

## Authorization (RBAC)

The system implements Role-Based Access Control (RBAC) with three primary roles:

### Admin Role

Administrators have organization-level access:

- Manage organizations (create, update, delete)
- Manage all users (create, update, delete)
- Manage all courses (create, update, delete)
- View all assessments and submissions (read-only)
- View statistics for the organization, teachers, and students

### Teacher Role

Teachers have course-level access for their assigned courses:

- View assigned courses
- Create, update, and delete assessments for assigned courses
- Grade student submissions
- View student enrollments in their courses
- Cannot modify organizational structure or users

### Student Role

Students have limited access focused on their educational journey:

- Enroll in available courses
- View enrolled courses and their details
- Submit assessments for enrolled courses
- View grades and feedback for their submissions
- Cannot modify course structure or user data

## Implementation Details

### Middleware

The system uses middleware to enforce authentication and role-based access:

1. **AuthMiddleware**: Validates the JWT token and attaches user context
2. **AdminOnly**: Ensures the user has the admin role
3. **TeacherOnly**: Ensures the user has the teacher role
4. **StudentOnly**: Ensures the user has the student role

### Organization-Level Isolation

All data access is restricted to the user's organization:

- Admins can only manage users and courses within their organization
- Teachers can only access courses and assessments within their organization
- Students can only access courses and assessments within their organization

### Authorization Checks

Beyond role-based middleware, the system implements additional authorization checks:

- Teachers can only manage assessments for courses they are assigned to
- Students can only view and submit assessments for courses they are enrolled in
- Course enrollment status is enforced for student enrollment

## Password Management

The system implements secure password practices:

- All passwords are hashed using bcrypt with appropriate cost factors
- Passwords are never stored in plain text
- Password change operations require user authentication
- Initial user creation assigns a default password that users should change

## Token Security

The authentication system uses a two-token approach for enhanced security:

### Access Tokens
- Short-lived JWT tokens used for API authentication (24-hour validity)
- Signed using HMAC SHA-256 with a secret key
- Include expiration, issued at, and not before timestamps
- Contain the minimal required user information
- Include the user's organization ID to enforce organization-level isolation

### Refresh Tokens
- Long-lived tokens used to obtain new access tokens (7-day validity)
- Securely stored in the database with user association
- Can be revoked individually or all at once for security purposes
- Implementing token rotation for enhanced security (new refresh token issued when used)
- Protected against token reuse attacks

### Security Benefits
- Limited access token lifespan reduces the risk if tokens are compromised
- Ability to revoke authentication without requiring password change
- Improved user experience by reducing the need for frequent logins
- Fine-grained control over active sessions