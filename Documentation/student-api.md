# Student API

This document describes the Student API endpoints for the Assessment Management System.

## Base URL

All student endpoints are relative to the base URL `/api/student`.

## Authentication

All student endpoints require authentication with a valid JWT token and student role.

## Course Management

### Get Enrolled Courses

Retrieves all courses in which the student is enrolled.

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
      "teachers": [
        {
          "id": "7f8d4e1c-9b0a-4e2d-8c7f-6b5a3d2e1c0b",
          "first_name": "John",
          "last_name": "Smith",
          "email": "john.smith@example.com"
        }
      ],
      "enrolled_at": "2025-03-29T13:35:45.123456Z"
    }
  ],
  "pagination": {
    "total": 4,
    "page": 1,
    "limit": 10,
    "pages": 1
  }
}
```

### Get Course by ID

Retrieves a specific course by ID (must be enrolled).

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
  "teachers": [
    {
      "id": "7f8d4e1c-9b0a-4e2d-8c7f-6b5a3d2e1c0b",
      "first_name": "John",
      "last_name": "Smith",
      "email": "john.smith@example.com"
    }
  ],
  "enrolled_at": "2025-03-29T13:35:45.123456Z",
  "organization_name": "Example University",
  "assessment_count": 2
}
```

### Get Available Courses

Retrieves all courses available for enrollment.

**Endpoint:** `GET /courses/available`

**Query Parameters:**

- `page` (optional): Page number for pagination (default: 1)
- `limit` (optional): Number of items per page (default: 10)
- `search` (optional): Search by course name or description

**Response:**

Status Code: 200 OK

```json
{
  "courses": [
    {
      "course": {
        "id": "4d5e6f7g-8h9i-0j1k-2l3m-4n5o6p7q8r9s",
        "organization_id": "92641a7d-966e-4e29-8d52-1d8ea6cac530",
        "name": "Introduction to Data Science",
        "description": "Learn the fundamentals of data analysis and visualization",
        "enrollment_open": true,
        "created_at": "2025-03-29T13:45:45.123456Z",
        "updated_at": "2025-03-29T13:45:45.123456Z"
      },
      "teachers": [
        {
          "id": "8g9h0i1j-2k3l-4m5n-6o7p-8q9r0s1t2u3v",
          "first_name": "Jane",
          "last_name": "Doe",
          "email": "jane.doe@example.com"
        }
      ],
      "student_count": 8
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

### Enroll in Course

Enrolls the student in a course.

**Endpoint:** `POST /courses/:id/enroll`

**URL Parameters:**

- `id`: Course ID

**Response:**

Status Code: 200 OK

```json
{
  "message": "Successfully enrolled in course",
  "course_id": "4d5e6f7g-8h9i-0j1k-2l3m-4n5o6p7q8r9s",
  "enrolled_at": "2025-03-29T15:30:00Z"
}
```

**Error Responses:**

Status Code: 400 Bad Request - Already enrolled or enrollment closed

```json
{
  "message": "You are already enrolled in this course"
}
```

```json
{
  "message": "Enrollment is currently closed for this course"
}
```

Status Code: 404 Not Found - Course not found or no teachers assigned

```json
{
  "message": "Course not found"
}
```

```json
{
  "message": "No teachers are assigned to this course yet. Enrollment is not available."
}
```

### Get Organization Details

Retrieves details about the student's organization.

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

### Get Course Assessments

Retrieves all assessments for a specific course.

**Endpoint:** `GET /courses/:courseId/assessments`

**URL Parameters:**

- `courseId`: Course ID

**Query Parameters:**

- `page` (optional): Page number for pagination (default: 1)
- `limit` (optional): Number of items per page (default: 10)
- `status` (optional): Filter by status (pending, submitted, graded, overdue)

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
        "title": "Midterm Examination",
        "description": "Updated comprehensive examination of course material",
        "type": "exam",
        "max_score": 120,
        "due_date": "2025-04-20T23:59:59Z",
        "created_at": "2025-03-29T14:00:00Z",
        "updated_at": "2025-03-29T14:45:00Z"
      },
      "has_submitted": true,
      "is_graded": true,
      "days_until_due": 22,
      "is_overdue": false
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
      "has_submitted": false,
      "is_graded": false,
      "days_until_due": 47,
      "is_overdue": false
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
    "title": "Midterm Examination",
    "description": "Updated comprehensive examination of course material",
    "type": "exam",
    "max_score": 120,
    "due_date": "2025-04-20T23:59:59Z",
    "created_at": "2025-03-29T14:00:00Z",
    "updated_at": "2025-03-29T14:45:00Z"
  },
  "course_name": "Advanced Computer Science",
  "teacher_name": "John Smith",
  "has_submitted": true,
  "is_graded": true,
  "days_until_due": 22,
  "is_overdue": false,
  "submission": {
    "id": "5e6f7g8h-9i0j-1k2l-3m4n-5o6p7q8r9s0t",
    "submitted_at": "2025-04-14T18:30:00Z",
    "content": "Student's exam answers"
  },
  "grade": {
    "score": 85,
    "max_score": 120,
    "percentage": 70.83,
    "feedback": "Good work, but needs improvement in section 3",
    "graded_at": "2025-04-15T10:15:00Z"
  }
}
```

### Submit Assessment

Submits an answer for an assessment.

**Endpoint:** `POST /assessments/:id/submit`

**URL Parameters:**

- `id`: Assessment ID

**Request Body:**

```json
{
  "content": "This is my complete exam with all answers..."
}
```

**Response:**

Status Code: 201 Created

```json
{
  "id": "5e6f7g8h-9i0j-1k2l-3m4n-5o6p7q8r9s0t",
  "assessment_id": "4d5e6f7g-8h9i-0j1k-2l3m-4n5o6p7q8r9s",
  "student_id": "1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p",
  "content": "This is my complete exam with all answers...",
  "submitted_at": "2025-04-14T18:30:00Z"
}
```

**Error Responses:**

Status Code: 400 Bad Request - Already submitted or overdue

```json
{
  "message": "You have already submitted this assessment"
}
```

```json
{
  "message": "The due date for this assessment has passed"
}
```

### View Submission

Retrieves the student's submission for a specific assessment.

**Endpoint:** `GET /assessments/:id/submission`

**URL Parameters:**

- `id`: Assessment ID

**Response:**

Status Code: 200 OK

```json
{
  "id": "5e6f7g8h-9i0j-1k2l-3m4n-5o6p7q8r9s0t",
  "assessment_id": "4d5e6f7g-8h9i-0j1k-2l3m-4n5o6p7q8r9s",
  "student_id": "1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p",
  "content": "This is my complete exam with all answers...",
  "submitted_at": "2025-04-14T18:30:00Z"
}
```

Status Code: 404 Not Found - No submission found

```json
{
  "message": "You have not submitted this assessment yet"
}
```

### View Grade

Retrieves the grade for the student's submission.

**Endpoint:** `GET /assessments/:id/grade`

**URL Parameters:**

- `id`: Assessment ID

**Response:**

Status Code: 200 OK

```json
{
  "submission_id": "5e6f7g8h-9i0j-1k2l-3m4n-5o6p7q8r9s0t",
  "score": 85,
  "max_score": 120,
  "percentage": 70.83,
  "feedback": "Good work, but needs improvement in section 3",
  "graded_by": "John Smith",
  "graded_at": "2025-04-15T10:15:00Z"
}
```

Status Code: 404 Not Found - No submission or grade found

```json
{
  "message": "You have not submitted this assessment yet"
}
```

```json
{
  "message": "Your submission has not been graded yet"
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

Status Code: 403 Forbidden - Insufficient permissions or not enrolled in the course

```json
{
  "message": "Forbidden: You are not enrolled in this course"
}
```

Status Code: 404 Not Found - Resource not found

```json
{
  "message": "Course not found"
}
```

Status Code: 500 Internal Server Error - Server error

```json
{
  "message": "Internal server error"
}
```