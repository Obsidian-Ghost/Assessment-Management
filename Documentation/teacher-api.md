# Teacher API

This document describes the Teacher API endpoints for the Assessment Management System.

## Base URL

All teacher endpoints are relative to the base URL `/api/teacher`.

## Authentication

All teacher endpoints require authentication with a valid JWT token and teacher role.

## Course Management

### Get Assigned Courses

Retrieves all courses assigned to the authenticated teacher.

**Endpoint:** `GET /courses`

**Query Parameters:**

- `page` (optional): Page number for pagination (default: 1)
- `limit` (optional): Number of items per page (default: 10)

**Response:**

Status Code: 200 OK

```json
{
  "courses": [
    {
      "course": {
        "id": "3c4d5e6f-7g8h-9i0j-1k2l-3m4n5o6p7q8r",
        "organization_id": "92641a7d-966e-4e29-8d52-1d8ea6cac530",
        "name": "Advanced Computer Science",
        "description": "An in-depth exploration of computer science principles",
        "enrollment_open": true,
        "created_at": "2025-03-29T13:15:45.123456Z",
        "updated_at": "2025-03-29T13:35:45.123456Z"
      },
      "student_count": 15
    }
  ],
  "pagination": {
    "total": 3,
    "page": 1,
    "limit": 10,
    "pages": 1
  }
}
```

### Get Course by ID

Retrieves a specific course by ID (must be assigned to the teacher).

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
    "name": "Advanced Computer Science",
    "description": "An in-depth exploration of computer science principles",
    "enrollment_open": true,
    "created_at": "2025-03-29T13:15:45.123456Z",
    "updated_at": "2025-03-29T13:35:45.123456Z"
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

### Get Course Students

Retrieves all students enrolled in a specific course.

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

### Get Organization Details

Retrieves details about the teacher's organization.

**Endpoint:** `GET /organization`

**Response:**

Status Code: 200 OK

```json
{
  "id": "92641a7d-966e-4e29-8d52-1d8ea6cac530",
  "name": "Example University",
  "slogan": "Excellence in Education",
  "teacher_count": 12,
  "student_count": 30,
  "course_count": 8
}
```

## Assessment Management

### Create Assessment

Creates a new assessment for a course.

**Endpoint:** `POST /assessments`

**Request Body:**

```json
{
  "course_id": "3c4d5e6f-7g8h-9i0j-1k2l-3m4n5o6p7q8r",
  "title": "Final Project",
  "description": "Build a full-stack web application",
  "type": "project",
  "max_score": 100,
  "due_date": "2025-05-15T23:59:59Z"
}
```

**Response:**

Status Code: 201 Created

```json
{
  "id": "6f7g8h9i-0j1k-2l3m-4n5o-6p7q8r9s0t1u",
  "course_id": "3c4d5e6f-7g8h-9i0j-1k2l-3m4n5o6p7q8r",
  "teacher_id": "7f8d4e1c-9b0a-4e2d-8c7f-6b5a3d2e1c0b",
  "title": "Final Project",
  "description": "Build a full-stack web application",
  "type": "project",
  "max_score": 100,
  "due_date": "2025-05-15T23:59:59Z",
  "created_at": "2025-03-29T14:30:00Z",
  "updated_at": "2025-03-29T14:30:00Z"
}
```

### Get Course Assessments

Retrieves all assessments for a specific course.

**Endpoint:** `GET /courses/:courseId/assessments`

**URL Parameters:**

- `courseId`: Course ID

**Query Parameters:**

- `page` (optional): Page number for pagination (default: 1)
- `limit` (optional): Number of items per page (default: 10)
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
    },
    {
      "assessment": {
        "id": "6f7g8h9i-0j1k-2l3m-4n5o-6p7q8r9s0t1u",
        "course_id": "3c4d5e6f-7g8h-9i0j-1k2l-3m4n5o6p7q8r",
        "teacher_id": "7f8d4e1c-9b0a-4e2d-8c7f-6b5a3d2e1c0b",
        "title": "Final Project",
        "description": "Build a full-stack web application",
        "type": "project",
        "max_score": 100,
        "due_date": "2025-05-15T23:59:59Z",
        "created_at": "2025-03-29T14:30:00Z",
        "updated_at": "2025-03-29T14:30:00Z"
      },
      "submission_count": 0,
      "graded_count": 0,
      "ungraded_count": 0
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
  "course_name": "Advanced Computer Science",
  "submission_count": 10,
  "graded_count": 8,
  "ungraded_count": 2
}
```

### Update Assessment

Updates an existing assessment.

**Endpoint:** `PUT /assessments/:id`

**URL Parameters:**

- `id`: Assessment ID

**Request Body:**

```json
{
  "title": "Midterm Examination",
  "description": "Updated comprehensive examination of course material",
  "max_score": 120,
  "due_date": "2025-04-20T23:59:59Z"
}
```

**Response:**

Status Code: 200 OK

```json
{
  "id": "4d5e6f7g-8h9i-0j1k-2l3m-4n5o6p7q8r9s",
  "course_id": "3c4d5e6f-7g8h-9i0j-1k2l-3m4n5o6p7q8r",
  "teacher_id": "7f8d4e1c-9b0a-4e2d-8c7f-6b5a3d2e1c0b",
  "title": "Midterm Examination",
  "description": "Updated comprehensive examination of course material",
  "type": "exam",
  "max_score": 120,
  "due_date": "2025-04-20T23:59:59Z",
  "created_at": "2025-03-29T14:00:00Z",
  "updated_at": "2025-03-29T14:45:00Z"
}
```

### Delete Assessment

Deletes an assessment and all associated submissions and grades.

**Endpoint:** `DELETE /assessments/:id`

**URL Parameters:**

- `id`: Assessment ID

**Response:**

Status Code: 204 No Content

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
      "grade": {
        "score": 85,
        "feedback": "Good work, but needs improvement in section 3",
        "graded_at": "2025-04-15T10:15:00Z"
      }
    },
    {
      "id": "6f7g8h9i-0j1k-2l3m-4n5o-6p7q8r9s0t1u",
      "assessment_id": "4d5e6f7g-8h9i-0j1k-2l3m-4n5o6p7q8r9s",
      "student_id": "2b3c4d5e-6f7g-8h9i-0j1k-2l3m4n5o6p7q",
      "student_name": "Jane Doe",
      "content": "Student's exam answers",
      "submitted_at": "2025-04-14T19:45:00Z",
      "is_graded": false
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

### Grade Submission

Grades a student's submission.

**Endpoint:** `POST /submissions/:submissionId/grade`

**URL Parameters:**

- `submissionId`: Submission ID

**Request Body:**

```json
{
  "score": 90,
  "feedback": "Excellent work! Very thorough and well-reasoned responses."
}
```

**Response:**

Status Code: 200 OK

```json
{
  "submission_id": "6f7g8h9i-0j1k-2l3m-4n5o-6p7q8r9s0t1u",
  "assessment_id": "4d5e6f7g-8h9i-0j1k-2l3m-4n5o6p7q8r9s",
  "student_id": "2b3c4d5e-6f7g-8h9i-0j1k-2l3m4n5o6p7q",
  "score": 90,
  "feedback": "Excellent work! Very thorough and well-reasoned responses.",
  "graded_by": "7f8d4e1c-9b0a-4e2d-8c7f-6b5a3d2e1c0b",
  "graded_at": "2025-03-29T15:00:00Z"
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

Status Code: 403 Forbidden - Insufficient permissions or not assigned to the course

```json
{
  "message": "Forbidden: You are not assigned to this course"
}
```

Status Code: 404 Not Found - Resource not found

```json
{
  "message": "Assessment not found"
}
```

Status Code: 400 Bad Request - Invalid input

```json
{
  "message": "Invalid request: Score must be between 0 and the maximum score"
}
```

Status Code: 500 Internal Server Error - Server error

```json
{
  "message": "Internal server error"
}
```