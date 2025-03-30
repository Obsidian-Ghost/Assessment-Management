# Role-Based Access Control (RBAC)

This document details the Role-Based Access Control system implemented in the Assessment Management System, focusing on permissions, access restrictions, and role hierarchies.

## Role Hierarchy

The system implements three distinct roles with different permission levels:

1. **Admin**: Highest level of access with organization-wide management capabilities
2. **Teacher**: Course-level access with assessment management capabilities
3. **Student**: Limited access focused on course enrollment and assessment submission

## Permission Matrix

Below is a detailed permission matrix for each role:

| Feature/Operation       | Admin | Teacher | Student |
|-------------------------|-------|---------|---------|
| **Organizations**       |       |         |         |
| Create Organization     | ✅     | ❌       | ❌       |
| View All Organizations  | ✅     | ❌       | ❌       |
| View Own Organization   | ✅     | ✅       | ✅       |
| Update Organization     | ✅     | ❌       | ❌       |
| Delete Organization     | ✅     | ❌       | ❌       |
| **Users**               |       |         |         |
| Create User             | ✅     | ❌       | ❌       |
| View All Users          | ✅     | ❌       | ❌       |
| View User Details       | ✅     | ❌       | ❌       |
| Update User             | ✅     | ❌       | ❌       |
| Delete User             | ✅     | ❌       | ❌       |
| Bulk Upload Users       | ✅     | ❌       | ❌       |
| **Courses**             |       |         |         |
| Create Course           | ✅     | ❌       | ❌       |
| View All Courses        | ✅     | ❌       | ❌       |
| View Assigned Courses   | N/A   | ✅       | ❌       |
| View Enrolled Courses   | N/A   | N/A     | ✅       |
| View Available Courses  | N/A   | N/A     | ✅       |
| Update Course           | ✅     | ❌       | ❌       |
| Delete Course           | ✅     | ❌       | ❌       |
| Assign Teachers         | ✅     | ❌       | ❌       |
| Remove Teachers         | ✅     | ❌       | ❌       |
| Toggle Enrollment       | ✅     | ❌       | ❌       |
| Enroll Students         | ✅     | ❌       | ❌       |
| View Course Students    | ✅     | ✅       | ❌       |
| Self-Enroll in Course   | ❌     | ❌       | ✅       |
| **Assessments**         |       |         |         |
| Create Assessment       | ❌     | ✅       | ❌       |
| View All Assessments    | ✅     | ❌       | ❌       |
| View Course Assessments | ✅     | ✅       | ✅       |
| Update Assessment       | ❌     | ✅       | ❌       |
| Delete Assessment       | ❌     | ✅       | ❌       |
| View Submissions        | ✅     | ✅       | ❌       |
| Submit Assessment       | ❌     | ❌       | ✅       |
| View Own Submission     | ❌     | ❌       | ✅       |
| Grade Submission        | ❌     | ✅       | ❌       |
| View Grades (all)       | ✅     | ✅       | ❌       |
| View Own Grades         | ❌     | ❌       | ✅       |

## Permission Implementation

RBAC is implemented through several layers in the application:

### 1. Middleware Layer

The system includes role-specific middleware that checks user roles before allowing access to specific API groups:

- `AdminOnly`: Ensures only users with the admin role can access admin routes
- `TeacherOnly`: Ensures only users with the teacher role can access teacher routes
- `StudentOnly`: Ensures only users with the student role can access student routes

### 2. Service Layer

Beyond middleware checks, the service layer implements additional authorization logic:

#### Admin Service
- Validates organization membership
- Prevents deletion of courses with enrolled students
- Ensures only users within the same organization can be modified

#### Teacher Service
- Validates course assignment before allowing assessment management
- Prevents modification of assessments for unassigned courses
- Ensures teachers can only grade submissions for their assessments

#### Student Service
- Validates course enrollment before allowing assessment submission
- Prevents enrollment in closed courses
- Ensures students can only view their own submissions and grades

### 3. Data Access Layer

The repository layer enforces organization boundaries:

- All queries include organization_id filters to ensure data isolation
- Cross-organization operations are prevented by design
- Specialized queries enforce role-specific access patterns

## Business Rules and Constraints

Beyond basic RBAC, the system enforces additional business rules:

1. **Course Deletion**: Courses can only be deleted if they have no enrolled students
2. **Assessment Management**: Only the teacher who created an assessment can modify it
3. **Enrollment Requirements**: Courses must have at least one teacher assigned for students to enroll
4. **Grade Visibility**: Students can only see their own grades, not those of other students
5. **Organization Boundaries**: Users cannot access data from other organizations

## Implementation Examples

### Middleware Implementation

The middleware checks the user's role in the JWT token:

```go
// AdminOnly middleware ensures that only administrators can access the routes
func AdminOnly() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            user := c.Get("user").(*jwt.Token)
            claims := user.Claims.(jwt.MapClaims)
            role := claims["role"].(string)

            if role != "admin" {
                return echo.ErrForbidden
            }
            return next(c)
        }
    }
}
```

### Service Layer Authorization

The service layer implements authorization checks:

```go
// Example from the AssessmentService
func (s *AssessmentService) UpdateAssessment(teacherID string, assessmentID string, req *models.UpdateAssessmentRequest) (*models.Assessment, error) {
    // Get the assessment
    assessment, err := s.repo.FindByID(assessmentID)
    if err != nil {
        return nil, err
    }

    // Check if the teacher created this assessment
    if assessment.TeacherID != teacherID {
        return nil, errors.New("unauthorized: you can only update your own assessments")
    }

    // Continue with the update...
}
```