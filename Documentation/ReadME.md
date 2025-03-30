# Assessment Management System Documentation

Welcome to the Assessment Management System documentation. This document serves as an entry point to all aspects of the system's architecture, API endpoints, database schema, and more.

## Overview

The Assessment Management System is a comprehensive platform designed to manage organizations, users, courses, and assessments. The system implements Role-Based Access Control (RBAC) with three roles:

1. **Admin** - Can manage all aspects of the system including organizations, users, courses, and assessments
2. **Teacher** - Can manage courses they are assigned to, create assessments, and grade student submissions
3. **Student** - Can enroll in courses, view and submit assessments, and view their grades

## System Architecture

The system follows a clean architecture pattern with the following layers:

1. **Handlers** - HTTP request handlers that parse requests and call appropriate services
2. **Services** - Business logic layer that implements core functionality
3. **Repositories** - Data access layer that interacts with the database
4. **Models** - Data structures used throughout the application

## Key Features

- **Multi-organization Support** - System supports multiple organizations with isolated data
- **Role-Based Access Control** - Granular permissions based on user roles
- **Course Management** - Create, update, and manage courses
- **Assessment System** - Create assessments, submit answers, and grade submissions
- **User Management** - Manage users with different roles
- **Email Uniqueness** - Ensures global email uniqueness across all organizations

## Documentation Sections

### API Documentation

- [Authentication API](auth-api.md) - Authentication endpoints for login and session management
- [Admin API](admin-api.md) - Endpoints for administrative functions
- [Teacher API](teacher-api.md) - Endpoints for teacher operations
- [Student API](student-api.md) - Endpoints for student operations

### Database and Implementation

- [Database Schema](database-schema.md) - Database tables, relationships, and constraints
- [Deployment Guide](deployment.md) - Instructions for deploying the system

## Getting Started

To get started with the system:

1. Review the [Deployment Guide](deployment.md) for setup instructions
2. Use the default admin credentials (email: admin@example.com, password: admin123) to log in
3. Create organizations, users, courses, and assessments as needed

## Authentication Flow

The system uses JWT (JSON Web Tokens) for authentication:

1. Users authenticate via the `/api/auth/login` endpoint
2. A JWT token is issued containing user information and permissions
3. This token is included in the Authorization header for all subsequent requests
4. The middleware validates the token and extracts user context for authorization

## Security Considerations

- All passwords are securely hashed using bcrypt
- JWT tokens have a configurable expiration time
- Users can only access resources within their organization
- Role-based permissions restrict access to sensitive operations
- Email addresses are globally unique across all organizations

## Environment Variables

The system uses the following environment variables:

- `PORT` - Server port (default: 5000)
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - Secret key for JWT signing
- `JWT_EXPIRATION` - JWT token expiration time in hours (default: 24)

## Data Isolation

One of the key architectural decisions is organization-based data isolation:

1. Each user belongs to exactly one organization
2. Users can only see and manipulate data within their organization
3. The middleware enforces organization-based access control
4. Database queries include organization filtering to ensure data isolation

## API Design Principles

1. **Consistency** - All endpoints follow the same patterns for requests and responses
2. **Pagination** - List endpoints support pagination for large data sets
3. **Filtering** - Many endpoints support filtering to narrow down results
4. **Error Handling** - Clear error responses with appropriate HTTP status codes
5. **Validation** - Input validation at all levels (handler, service, repository)

## Contributing

For information on contributing to this project, please contact the system administrator.

## License

This project is proprietary and confidential. Unauthorized copying, distribution, or use is strictly prohibited.