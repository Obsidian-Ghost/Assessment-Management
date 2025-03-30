# Database Schema

This document describes the database schema for the Assessment Management System.

## Overview

The system uses PostgreSQL as its database and follows a relational model. The schema includes tables for organizations, users, courses, assessments, and submissions, with appropriate relationships and constraints.

## Tables

### Organizations

Stores information about organizations in the system.

```sql
CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slogan VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_organizations_name ON organizations(name);
```

### Users

Stores user information including authentication details and role.

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'teacher', 'student')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT users_email_unique UNIQUE (email)
);

CREATE INDEX idx_users_organization_id ON users(organization_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
```

### Courses

Represents courses within an organization.

```sql
CREATE TABLE courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    enrollment_open BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_courses_organization_id ON courses(organization_id);
CREATE INDEX idx_courses_name ON courses(name);
```

### Course_Teachers

Represents the many-to-many relationship between courses and teachers.

```sql
CREATE TABLE course_teachers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT course_teachers_unique UNIQUE (course_id, teacher_id)
);

CREATE INDEX idx_course_teachers_course_id ON course_teachers(course_id);
CREATE INDEX idx_course_teachers_teacher_id ON course_teachers(teacher_id);
```

### Course_Students

Represents the many-to-many relationship between courses and students.

```sql
CREATE TABLE course_students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    enrolled_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT course_students_unique UNIQUE (course_id, student_id)
);

CREATE INDEX idx_course_students_course_id ON course_students(course_id);
CREATE INDEX idx_course_students_student_id ON course_students(student_id);
```

### Assessments

Represents assessments (exams, assignments, projects, etc.) created by teachers.

```sql
CREATE TABLE assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL CHECK (type IN ('quiz', 'exam', 'assignment', 'project')),
    max_score INTEGER NOT NULL CHECK (max_score > 0),
    due_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_assessments_course_id ON assessments(course_id);
CREATE INDEX idx_assessments_teacher_id ON assessments(teacher_id);
CREATE INDEX idx_assessments_due_date ON assessments(due_date);
CREATE INDEX idx_assessments_type ON assessments(type);
```

### Submissions

Represents student submissions for assessments.

```sql
CREATE TABLE submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT,
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_graded BOOLEAN DEFAULT false,
    CONSTRAINT submissions_unique UNIQUE (assessment_id, student_id)
);

CREATE INDEX idx_submissions_assessment_id ON submissions(assessment_id);
CREATE INDEX idx_submissions_student_id ON submissions(student_id);
CREATE INDEX idx_submissions_is_graded ON submissions(is_graded);
```

### Grades

Stores grades for student submissions.

```sql
CREATE TABLE grades (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id UUID NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,
    graded_by UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    score INTEGER NOT NULL CHECK (score >= 0),
    feedback TEXT,
    graded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT grades_submission_unique UNIQUE (submission_id)
);

CREATE INDEX idx_grades_submission_id ON grades(submission_id);
CREATE INDEX idx_grades_graded_by ON grades(graded_by);
```

### Refresh Tokens

Stores refresh tokens for user authentication sessions.

```sql
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked BOOLEAN DEFAULT false,
    CONSTRAINT refresh_tokens_token_unique UNIQUE (token)
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX idx_refresh_tokens_revoked ON refresh_tokens(revoked);
```

## Diagram

Below is a textual representation of the database schema diagram:

```
organizations
    |
    ├── users
    |     |
    |     ├── refresh_tokens
    |     |
    |     ├── course_teachers
    |     |       |
    |     |       └── courses
    |     |             |
    |     |             ├── assessments
    |     |             |       |
    |     |             |       └── submissions
    |     |             |             |
    |     |             |             └── grades
    |     |             |
    |     └── course_students
    |             |
    |             └── courses
    |
    └── courses
```

## Relationships

- **Organizations to Users**: One-to-many (an organization has many users, a user belongs to one organization)
- **Organizations to Courses**: One-to-many (an organization has many courses, a course belongs to one organization)
- **Users to Refresh Tokens**: One-to-many (a user can have multiple refresh tokens, a refresh token belongs to one user)
- **Courses to Teachers**: Many-to-many through course_teachers (a course can have multiple teachers, a teacher can teach multiple courses)
- **Courses to Students**: Many-to-many through course_students (a course can have multiple students, a student can enroll in multiple courses)
- **Teachers to Assessments**: One-to-many (a teacher creates many assessments, an assessment is created by one teacher)
- **Courses to Assessments**: One-to-many (a course has many assessments, an assessment belongs to one course)
- **Assessments to Submissions**: One-to-many (an assessment has many submissions, a submission is for one assessment)
- **Students to Submissions**: One-to-many (a student makes many submissions, a submission is by one student)
- **Submissions to Grades**: One-to-one (a submission has one grade, a grade belongs to one submission)
- **Teachers to Grades**: One-to-many (a teacher grades many submissions, a grade is given by one teacher)

## Constraints

1. **Foreign Key Constraints**:
    - User organization_id references organizations(id)
    - Course organization_id references organizations(id)
    - course_teachers course_id references courses(id)
    - course_teachers teacher_id references users(id)
    - course_students course_id references courses(id)
    - course_students student_id references users(id)
    - assessment course_id references courses(id)
    - assessment teacher_id references users(id)
    - submission assessment_id references assessments(id)
    - submission student_id references users(id)
    - grade submission_id references submissions(id)
    - grade graded_by references users(id)
    - refresh_token user_id references users(id)

2. **Unique Constraints**:
    - users(email) - Ensures global email uniqueness across all organizations
    - course_teachers(course_id, teacher_id) - Prevents duplicate teacher assignments
    - course_students(course_id, student_id) - Prevents duplicate student enrollments
    - submissions(assessment_id, student_id) - Ensures one submission per student per assessment
    - grades(submission_id) - Ensures one grade per submission
    - refresh_tokens(token) - Ensures token uniqueness

3. **Check Constraints**:
    - users role must be one of 'admin', 'teacher', 'student'
    - assessments type must be one of 'quiz', 'exam', 'assignment', 'project'
    - assessments max_score must be greater than 0
    - grades score must be greater than or equal to 0

## Cascade Behavior

1. **ON DELETE CASCADE**:
    - When an organization is deleted, all associated users, courses, and related data are deleted
    - When a user is deleted, all associated refresh tokens are deleted
    - When a course is deleted, all associated teachers, students, assessments, and related data are deleted
    - When an assessment is deleted, all associated submissions and grades are deleted
    - When a submission is deleted, its grade is deleted

2. **ON DELETE SET NULL**:
    - When a teacher is deleted, their assessments remain but teacher_id is set to NULL
    - When a grader is deleted, grades remain but graded_by is set to NULL

## Indexes

The schema includes indexes on frequently queried columns to improve performance:

1. **Primary Key Indexes** (automatically created on the id columns)
2. **Foreign Key Indexes** to speed up join operations
3. **Additional Indexes** on columns used in WHERE clauses, ORDER BY clauses, or for lookup operations

## Data Integrity

1. **Referential Integrity** is maintained through foreign key constraints
2. **Domain Integrity** is enforced through check constraints and appropriate data types
3. **Entity Integrity** is preserved by primary key constraints
4. **User-Defined Integrity** is implemented via application-level validation

## Migration Strategy

The database schema is initialized and maintained through SQL migrations in the `migrations` directory:

1. **Initial Migration** creates the base tables and constraints
2. **Subsequent Migrations** implement schema changes while preserving data
3. **Email Uniqueness Migration** ensures global email uniqueness across all organizations
4. **Refresh Token Migration** adds support for refresh tokens and enhanced authentication

## Performance Considerations

1. **Indexes** on frequently queried columns
2. **Appropriate Data Types** chosen for each column
3. **Text** used for potentially large fields (description, content, feedback)
4. **UUID** used for primary keys to allow distributed ID generation

## Security Considerations

1. **No Plain Text Passwords** - Only password hashes are stored
2. **Organization Isolation** - Queries filter by organization_id to ensure data isolation
3. **Role-Based Access** - The role column enables role-based access control
4. **Secure Authentication** - Refresh token system with expiration and revocation capabilities
5. **Token Rotation** - Refresh tokens can be revoked individually or all at once for security