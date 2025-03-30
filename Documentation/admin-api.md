# Admin API

This document describes the Admin API endpoints for the Assessment Management System.

## Base URL

All admin endpoints are relative to the base URL `/api/admin`.

## Authentication

All admin endpoints require authentication with a valid JWT token and admin role.

## Organization Management

### Create Organization

Creates a new organization.

**Endpoint:** `POST /organizations`

**Request Body:**

```json
{
  "name": "Example University",
  "slogan": "Knowledge for All"
}
```

**Response:**

Status Code: 201 Created

```json
{
  "id": "92641a7d-966e-4e29-8d52-1d8ea6cac530",
  "name": "Example University",
  "slogan": "Knowledge for All",
  "created_at": "2025-03-29T12:30:45.123456Z",
  "updated_at": "2025-03-29T12:30:45.123456Z"
}
```

### Get All Organizations

Retrieves all organizations.

**Endpoint:** `GET /organizations`

**Query Parameters:**

- `page` (optional): Page number for pagination (default: 1)
- `limit` (optional): Number of items per page (default: 10)

**Response:**

Status Code: 200 OK

```json
{
  "organizations": [
    {
      "id": "92641a7d-966e-4e29-8d52-1d8ea6cac530",
      "name": "Example University",
      "slogan": "Knowledge for All",
      "created_at": "2025-03-29T12:30:45.123456Z",
      "updated_at": "2025-03-29T12:30:45.123456Z"
    },
    {
      "id": "a69f9822-36c6-4282-8b34-27323d20315f",
      "name": "Tech Institute",
      "slogan": "Innovation Through Education",
      "created_at": "2025-03-29T12:35:45.123456Z",
      "updated_at": "2025-03-29T12:35:45.123456Z"
    }
  ],
  "pagination": {
    "total": 2,
    "page": 1,
    "limit": 10,
    "pages": 1
  }
}
```

### Get Organization by ID

Retrieves a specific organization by ID.

**Endpoint:** `GET /organizations/:id`

**URL Parameters:**

- `id`: Organization ID

**Response:**

Status Code: 200 OK

```json
{
  "id": "92641a7d-966e-4e29-8d52-1d8ea6cac530",
  "name": "Example University",
  "slogan": "Knowledge for All",
  "created_at": "2025-03-29T12:30:45.123456Z",
  "updated_at": "2025-03-29T12:30:45.123456Z"
}
```

### Update Organization

Updates an existing organization.

**Endpoint:** `PUT /organizations/:id`

**URL Parameters:**

- `id`: Organization ID

**Request Body:**

```json
{
  "name": "Example University Updated",
  "slogan": "Excellence in Education"
}
```

**Response:**

Status Code: 200 OK

```json
{
  "id": "92641a7d-966e-4e29-8d52-1d8ea6cac530",
  "name": "Example University Updated",
  "slogan": "Excellence in Education",
  "created_at": "2025-03-29T12:30:45.123456Z",
  "updated_at": "2025-03-29T12:40:45.123456Z"
}
```

### Delete Organization

Deletes an organization and all associated data.

**Endpoint:** `DELETE /organizations/:id`

**URL Parameters:**

- `id`: Organization ID

**Response:**

Status Code: 204 No Content

### Get Organization Statistics

Retrieves statistics for a specific organization.

**Endpoint:** `GET /organizations/:id/stats`

**URL Parameters:**

- `id`: Organization ID

**Response:**

Status Code: 200 OK

```json
{
  "organization_id": "92641a7d-966e-4e29-8d52-1d8ea6cac530",
  "organization_name": "Example University",
  "user_count": {
    "total": 45,
    "admin": 3,
    "teacher": 12,
    "student": 30
  },
  "course_count": 8,
  "assessment_count": 24,
  "submission_count": 150,
  "average_grade": 85.7
}
```

## User Management

### Create User

Creates a new user within the organization.

**Endpoint:** `POST /users`

**Request Body:**

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "password": "password123",
  "role": "teacher",
  "organization_id": "92641a7d-966e-4e29-8d52-1d8ea6cac530"
}
```

**Response:**

Status Code: 201 Created

```json
{
  "id": "7f8d4e1c-9b0a-4e2d-8c7f-6b5a3d2e1c0b",
  "organization_id": "92641a7d-966e-4e29-8d52-1d8ea6cac530",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "role": "teacher",
  "created_at": "2025-03-29T12:50:45.123456Z",
  "updated_at": "2025-03-29T12:50:45.123456Z"
}
```

### Get All Users

Retrieves all users within the organization.

**Endpoint:** `GET /users`

**Query Parameters:**

- `page` (optional): Page number for pagination (default: 1)
- `limit` (optional): Number of items per page (default: 10)
- `role` (optional): Filter by role (admin, teacher, student)
- `search` (optional): Search by name or email

**Response:**

Status Code: 200 OK

```json
{
  "users": [
    {
      "id": "e41756fe-5242-4698-8102-735621b699d8",
      "organization_id": "92641a7d-966e-4e29-8d52-1d8ea6cac530",
      "first_name": "System",
      "last_name": "Administrator",
      "email": "admin@example.com",
      "role": "admin",
      "created_at": "2025-03-28T11:02:23.175179Z",
      "updated_at": "2025-03-28T11:12:19.589734Z"
    },
    {
      "id": "7f8d4e1c-9b0a-4e2d-8c7f-6b5a3d2e1c0b",
      "organization_id": "92641a7d-966e-4e29-8d52-1d8ea6cac530",
      "first_name": "John",
      "last_name": "Doe",
      "email": "john.doe@example.com",
      "role": "teacher",
      "created_at": "2025-03-29T12:50:45.123456Z",
      "updated_at": "2025-03-29T12:50:45.123456Z"
    }
  ],
  "pagination": {
    "total": 45,
    "page": 1,
    "limit": 10,
    "pages": 5
  }
}
```

### Get User by ID

Retrieves a specific user by ID.

**Endpoint:** `GET /users/:id`

**URL Parameters:**

- `id`: User ID

**Response:**

Status Code: 200 OK

```json
{
  "id": "7f8d4e1c-9b0a-4e2d-8c7f-6b5a3d2e1c0b",
  "organization_id": "92641a7d-966e-4e29-8d52-1d8ea6cac530",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "role": "teacher",
  "created_at": "2025-03-29T12:50:45.123456Z",
  "updated_at": "2025-03-29T12:50:45.123456Z"
}
```

### Update User

Updates an existing user.

**Endpoint:** `PUT /users/:id`

**URL Parameters:**

- `id`: User ID

**Request Body:**

```json
{
  "first_name": "John",
  "last_name": "Smith",
  "email": "john.smith@example.com",
  "role": "teacher"
}
```

**Response:**

Status Code: 200 OK

```json
{
  "id": "7f8d4e1c-9b0a-4e2d-8c7f-6b5a3d2e1c0b",
  "organization_id": "92641a7d-966e-4e29-8d52-1d8ea6cac530",
  "first_name": "John",
  "last_name": "Smith",
  "email": "john.smith@example.com",
  "role": "teacher",
  "created_at": "2025-03-29T12:50:45.123456Z",
  "updated_at": "2025-03-29T13:05:45.123456Z"
}
```

### Delete User

Deletes a user.

**Endpoint:** `DELETE /users/:id`

**URL Parameters:**

- `id`: User ID

**Response:**

Status Code: 204 No Content

### Bulk Upload Users

Uploads multiple users from a CSV file.

**Endpoint:** `POST /users/bulk`

**Request Body:**

Multipart form with a `file` field containing a CSV file with the following format:
```
first_name,last_name,email,role
Jane,Doe,jane.doe@example.com,teacher
John,Smith,john.smith@example.com,student
```

**Response:**

Status Code: 200 OK

```json
{
  "message": "Successfully imported 2 users",
  "users": [
    {
      "id": "9e8d7c6b-5a4b-3c2d-1e0f-9a8b7c6d5e4f",
      "email": "jane.doe@example.com",
      "status": "created"
    },
    {
      "id": "1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p",
      "email": "john.smith@example.com",
      "status": "created"
    }
  ],
  "errors": []
}
```

### Get Teacher Statistics

Retrieves statistics for a specific teacher.

**Endpoint:** `GET /users/teachers/:id/stats`

**URL Parameters:**

- `id`: Teacher ID

**Response:**

Status Code: 200 OK

```json
{
  "teacher": {
    "id": "7f8d4e1c-9b0a-4e2d-8c7f-6b5a3d2e1c0b",
    "first_name": "John",
    "last_name": "Smith",
    "email": "john.smith@example.com"
  },
  "courses": {
    "assigned_count": 3,
    "active_count": 2
  },
  "assessments": {
    "created_count": 15,
    "graded_count": 120,
    "pending_count": 5
  },
  "students": {
    "total_count": 75
  }
}
```

### Get Student Statistics

Retrieves statistics for a specific student.

**Endpoint:** `GET /users/students/:id/stats`

**URL Parameters:**

- `id`: Student ID

**Response:**

Status Code: 200 OK

```json
{
  "student": {
    "id": "1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p",
    "first_name": "John",
    "last_name": "Smith",
    "email": "john.smith@example.com"
  },
  "courses": {
    "enrolled_count": 4,
    "active_count": 3
  },
  "assessments": {
    "assigned_count": 20,
    "completed_count": 15,
    "pending_count": 5,
    "overdue_count": 0
  },
  "grades": {
    "average_score": 87.5,
    "highest_score": 98,
    "lowest_score": 75
  }
}
```

## Course Management

### Create Course

Creates a new course within the organization.

**Endpoint:** `POST /courses`

**Request Body:**

```json
{
  "name": "Introduction to Computer Science",
  "description": "A beginner's guide to computer science principles",
  "enrollment_open": true
}
```

**Response:**

Status Code: 201 Created

```json
{
  "id": "3c4d5e6f-7g8h-9i0j-1k2l-3m4n5o6p7q8r",
  "organization_id": "92641a7d-966e-4e29-8d52-1d8ea6cac530",
  "name": "Introduction to Computer Science",
  "description": "A beginner's guide to computer science principles",
  "enrollment_open": true,
  "created_at": "2025-03-29T13:15:45.123456Z",
  "updated_at": "2025-03-29T13:15:45.123456Z"
}
```

### Get All Courses

Retrieves all courses within the organization.

**Endpoint:** `GET /courses`

**Query Parameters:**

- `page` (optional): Page number for pagination (default: 1)
- `limit` (optional): Number of items per page (default: 10)
- `search` (optional): Search by name or description

**Response:**

Status Code: 200 OK

```json
{
  "courses": [
    {
      "course": {
        "id": "3c4d5e6f-7g8h-9i0j-1k2l-3m4n5o6p7q8r",
        "organization_id": "92641a7d-966e-4e29-8d52-1d8ea6cac530",
        "name": "Introduction to Computer Science",
        "description": "A beginner's guide to computer science principles",
        "enrollment_open": true,
        "created_at": "2025-03-29T13:15:45.123456Z",
        "updated_at": "2025-03-29T13:15:45.123456Z"
      },
      "student_count": 15,
      "teacher_count": 1
    }
  ],
  "pagination": {
    "total": 8,
    "page": 1,
    "limit": 10,
    "pages": 1
  }
}
```

### Get Course by ID

Retrieves a specific course by ID.

**Endpoint:** `GET /courses/:id`

**URL Parameters:**

- `id`: Course ID

**Response:**

Status Code: 200 OK

```json
{
  "course": {
    "id": "3c4d5e6f-7g8h-9i0j-1k2l-3m4n5o6p7q8r",
    "organization_id": "92641a7d-966e-4e29-8d52-1d8ea6cac530",
    "name": "Introduction to Computer Science",
    "description": "A beginner's guide to computer science principles",
    "enrollment_open": true,
    "created_at": "2025-03-29T13:15:45.123456Z",
    "updated_at": "2025-03-29T13:15:45.123456Z"
  },
  "organization_id": "92641a7d-966e-4e29-8d52-1d8ea6cac530",
  "organization_name": "Example University",
  "teachers": [
    {
      "id": "7f8d4e1c-9b0a-4e2d-8c7f-6b5a3d2e1c0b",
      "first_name": "John",
      "last_name": "Smith",
      "email": "john.smith@example.com",
      "role": "teacher"
    }
  ],
  "student_count": 15
}
```

### Update Course

Updates an existing course.

**Endpoint:** `PUT /courses/:id`

**URL Parameters:**

- `id`: Course ID

**Request Body:**

```json
{
  "name": "Advanced Computer Science",
  "description": "An in-depth exploration of computer science principles",
  "enrollment_open": false
}
```

**Response:**

Status Code: 200 OK

```json
{
  "id": "3c4d5e6f-7g8h-9i0j-1k2l-3m4n5o6p7q8r",
  "organization_id": "92641a7d-966e-4e29-8d52-1d8ea6cac530",
  "name": "Advanced Computer Science",
  "description": "An in-depth exploration of computer science principles",
  "enrollment_open": false,
  "created_at": "2025-03-29T13:15:45.123456Z",
  "updated_at": "2025-03-29T13:25:45.123456Z"
}
```

### Delete Course

Deletes a course.

**Endpoint:** `DELETE /courses/:id`

**URL Parameters:**

- `id`: Course ID

**Response:**

Status Code: 204 No Content

**Error Responses:**

Status Code: 400 Bad Request - Course has enrolled students

```json
{
  "message": "Cannot delete course with enrolled students. Archive the course instead."
}
```

### Assign Teacher to Course

Assigns a teacher to a course.

**Endpoint:** `POST /courses/:id/teachers`

**URL Parameters:**

- `id`: Course ID

**Request Body:**

```json
{
  "teacher_id": "7f8d4e1c-9b0a-4e2d-8c7f-6b5a3d2e1c0b"
}
```

**Response:**

Status Code: 200 OK

```json
{
  "message": "Teacher assigned successfully",
  "course_id": "3c4d5e6f-7g8h-9i0j-1k2l-3m4n5o6p7q8r",
  "teacher_id": "7f8d4e1c-9b0a-4e2d-8c7f-6b5a3d2e1c0b"
}
```

### Remove Teacher from Course

Removes a teacher from a course.

**Endpoint:** `DELETE /courses/:id/teachers/:teacherId`

**URL Parameters:**

- `id`: Course ID
- `teacherId`: Teacher ID

**Response:**

Status Code: 204 No Content

### Get Course Teachers

Retrieves all teachers assigned to a course.

**Endpoint:** `GET /courses/:id/teachers`

**URL Parameters:**

- `id`: Course ID

**Response:**

Status Code: 200 OK

```json
{
  "teachers": [
    {
      "id": "7f8d4e1c-9b0a-4e2d-8c7f-6b5a3d2e1c0b",
      "first_name": "John",
      "last_name": "Smith",
      "email": "john.smith@example.com",
      "role": "teacher"
    }
  ]
}
```

### Toggle Course Enrollment

Opens or closes enrollment for a course.

**Endpoint:** `PUT /courses/:id/enrollment`

**URL Parameters:**

- `id`: Course ID

**Request Body:**

```json
{
  "enrollment_open": true
}
```

**Response:**

Status Code: 200 OK

```json
{
  "id": "3c4d5e6f-7g8h-9i0j-1k2l-3m4n5o6p7q8r",
  "organization_id": "92641a7d-966e-4e29-8d52-1d8ea6cac530",
  "name": "Advanced Computer Science",
  "description": "An in-depth exploration of computer science principles",
  "enrollment_open": true,
  "created_at": "2025-03-29T13:15:45.123456Z",
  "updated_at": "2025-03-29T13:35:45.123456Z"
}
```

### Manage Student Enrollment

Enrolls a student in a course.

**Endpoint:** `POST /courses/:id/students`

**URL Parameters:**

- `id`: Course ID

**Request Body:**

```json
{
  "student_id": "1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p"
}
```

**Response:**

Status Code: 200 OK

```json
{
  "message": "Student enrolled successfully",
  "course_id": "3c4d5e6f-7g8h-9i0j-1k2l-3m4n5o6p7q8r",
  "student_id": "1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p"
}
```

### Bulk Enroll Students

Enrolls multiple students in a course.

**Endpoint:** `POST /courses/:id/students/bulk`

**URL Parameters:**

- `id`: Course ID

**Request Body:**

```json
{
  "student_ids": [
    "1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p",
    "2b3c4d5e-6f7g-8h9i-0j1k-2l3m4n5o6p7q"
  ]
}
```

**Response:**

Status Code: 200 OK

```json
{
  "message": "Successfully enrolled 2 students",
  "course_id": "3c4d5e6f-7g8h-9i0j-1k2l-3m4n5o6p7q8r",
  "enrolled_count": 2,
  "errors": []
}
```

### Get Course Students

Retrieves all students enrolled in a course.

**Endpoint:** `GET /courses/:id/students`

**URL Parameters:**

- `id`: Course ID

**Query Parameters:**

- `page` (optional): Page number for pagination (default: 1)
- `limit` (optional): Number of items per page (default: 10)
- `search` (optional): Search by name or email

**Response:**

Status Code: 200 OK

```json
{
  "students": [
    {
      "id": "1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p",
      "first_name": "John",
      "last_name": "Smith",
      "email": "john.smith@example.com",
      "role": "student",
      "enrolled_at": "2025-03-29T13:35:45.123456Z"
    }
  ],
  "pagination": {
    "total": 15,
    "page": 1,
    "limit": 10,
    "pages": 2
  }
}
```

## Assessment Management (Read-Only)

### Get All Assessments

Retrieves all assessments across all courses.

**Endpoint:** `GET /assessments`

**Query Parameters:**

- `page` (optional): Page number for pagination (default: 1)
- `limit` (optional): Number of items per page (default: 10)
- `course_id` (optional): Filter by course ID
- `type` (optional): Filter by assessment type

**Response:**

Status Code: 200 OK

```json
{
  "assessments": [
    {
      "assessment": {
        "id": "4d5e6f7g-8h9i-0j1k-2l3m-4n5o6p7q8r9s",
        "course_id": "3c4d5e6f-7g8h-9i0j-1k2l-3m4n5o6p7q8r",
        "teacher_id": "7f8d4e1c-9b0a-4e2d-8c7f-6b5a3d2e1c0b",
        "title": "Midterm Exam",
        "description": "Comprehensive examination of course material",
        "type": "exam",
        "max_score": 100,
        "due_date": "2025-04-15T23:59:59Z",
        "created_at": "2025-03-29T14:00:00Z",
        "updated_at": "2025-03-29T14:00:00Z"
      },
      "submission_count": 10,
      "graded_count": 8,
      "ungraded_count": 2
    }
  ],
  "pagination": {
    "total": 24,
    "page": 1,
    "limit": 10,
    "pages": 3
  }
}
```

### Get Assessment by ID

Retrieves a specific assessment by ID.

**Endpoint:** `GET /assessments/:id`

**URL Parameters:**

- `id`: Assessment ID

**Response:**

Status Code: 200 OK

```json
{
  "assessment": {
    "id": "4d5e6f7g-8h9i-0j1k-2l3m-4n5o6p7q8r9s",
    "course_id": "3c4d5e6f-7g8h-9i0j-1k2l-3m4n5o6p7q8r",
    "teacher_id": "7f8d4e1c-9b0a-4e2d-8c7f-6b5a3d2e1c0b",
    "title": "Midterm Exam",
    "description": "Comprehensive examination of course material",
    "type": "exam",
    "max_score": 100,
    "due_date": "2025-04-15T23:59:59Z",
    "created_at": "2025-03-29T14:00:00Z",
    "updated_at": "2025-03-29T14:00:00Z"
  },
  "teacher_name": "John Smith",
  "course_name": "Advanced Computer Science"
}
```

### Get Assessment Submissions

Retrieves all submissions for a specific assessment.

**Endpoint:** `GET /assessments/:id/submissions`

**URL Parameters:**

- `id`: Assessment ID

**Query Parameters:**

- `page` (optional): Page number for pagination (default: 1)
- `limit` (optional): Number of items per page (default: 10)
- `graded` (optional): Filter by graded status (true/false)

**Response:**

Status Code: 200 OK

```json
{
  "submissions": [
    {
      "id": "5e6f7g8h-9i0j-1k2l-3m4n-5o6p7q8r9s0t",
      "assessment_id": "4d5e6f7g-8h9i-0j1k-2l3m-4n5o6p7q8r9s",
      "student_id": "1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p",
      "student_name": "John Smith",
      "content": "Student's exam answers",
      "submitted_at": "2025-04-14T18:30:00Z",
      "is_graded": true,
      "score": 85,
      "feedback": "Good work, but needs improvement in section 3"
    }
  ],
  "pagination": {
    "total": 10,
    "page": 1,
    "limit": 10,
    "pages": 1
  }
}
```

### Get Submission Grade

Retrieves the grade for a specific submission.

**Endpoint:** `GET /submissions/:submissionId/grade`

**URL Parameters:**

- `submissionId`: Submission ID

**Response:**

Status Code: 200 OK

```json
{
  "submission_id": "5e6f7g8h-9i0j-1k2l-3m4n-5o6p7q8r9s0t",
  "assessment_id": "4d5e6f7g-8h9i-0j1k-2l3m-4n5o6p7q8r9s",
  "student_id": "1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p",
  "score": 85,
  "max_score": 100,
  "percentage": 85,
  "feedback": "Good work, but needs improvement in section 3",
  "graded_by": "John Smith",
  "graded_at": "2025-04-15T10:15:00Z"
}
```

## Error Responses

All endpoints may return the following error responses:

Status Code: 401 Unauthorized - Missing or invalid token

```json
{
  "message": "Unauthorized"
}
```

Status Code: 403 Forbidden - Insufficient permissions

```json
{
  "message": "Forbidden: Insufficient permissions"
}
```

Status Code: 404 Not Found - Resource not found

```json
{
  "message": "Resource not found"
}
```

Status Code: 500 Internal Server Error - Server error

```json
{
  "message": "Internal server error"
}
```