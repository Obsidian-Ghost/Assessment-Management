package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"

	"assessment-management-system/db"
	"assessment-management-system/models"
)

// CourseRepository handles database operations for courses
type CourseRepository struct {
	db *db.DB
}

// NewCourseRepository creates a new CourseRepository
func NewCourseRepository(db *db.DB) *CourseRepository {
	return &CourseRepository{
		db: db,
	}
}

// Create creates a new course
func (r *CourseRepository) Create(ctx context.Context, organizationID, name, description string, enrollmentOpen bool) (*models.Course, error) {
	var course models.Course
	err := r.db.Pool.QueryRow(ctx,
		`INSERT INTO courses (organization_id, name, description, enrollment_open) 
                VALUES ($1, $2, $3, $4) 
                RETURNING id, organization_id, name, description, enrollment_open, created_at, updated_at`,
		organizationID, name, description, enrollmentOpen).Scan(
		&course.ID, &course.OrganizationID, &course.Name, &course.Description, &course.EnrollmentOpen, &course.CreatedAt, &course.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &course, nil
}

// FindByID retrieves a course by ID
func (r *CourseRepository) FindByID(ctx context.Context, id string) (*models.Course, error) {
	var course models.Course
	err := r.db.Pool.QueryRow(ctx,
		`SELECT id, organization_id, name, description, enrollment_open, created_at, updated_at 
                FROM courses 
                WHERE id = $1`,
		id).Scan(&course.ID, &course.OrganizationID, &course.Name, &course.Description, &course.EnrollmentOpen, &course.CreatedAt, &course.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &course, nil
}

// FindByNameAndOrganization retrieves a course by name and organization
func (r *CourseRepository) FindByNameAndOrganization(ctx context.Context, name, organizationID string) (*models.Course, error) {
	var course models.Course
	err := r.db.Pool.QueryRow(ctx,
		`SELECT id, organization_id, name, description, enrollment_open, created_at, updated_at 
                FROM courses 
                WHERE name = $1 AND organization_id = $2`,
		name, organizationID).Scan(&course.ID, &course.OrganizationID, &course.Name, &course.Description, &course.EnrollmentOpen, &course.CreatedAt, &course.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &course, nil
}

// FindByOrganization retrieves all courses in an organization
func (r *CourseRepository) FindByOrganization(ctx context.Context, organizationID string) ([]*models.Course, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT id, organization_id, name, description, enrollment_open, created_at, updated_at 
                FROM courses 
                WHERE organization_id = $1
                ORDER BY name`,
		organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*models.Course
	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.ID, &course.OrganizationID, &course.Name, &course.Description, &course.EnrollmentOpen, &course.CreatedAt, &course.UpdatedAt); err != nil {
			return nil, err
		}
		courses = append(courses, &course)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

// Update updates a course
func (r *CourseRepository) Update(ctx context.Context, id string, name *string, description *string, enrollmentOpen *bool) (*models.Course, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Get current course
	var course models.Course
	err = tx.QueryRow(ctx,
		`SELECT id, organization_id, name, description, enrollment_open, created_at, updated_at 
                FROM courses 
                WHERE id = $1`,
		id).Scan(&course.ID, &course.OrganizationID, &course.Name, &course.Description, &course.EnrollmentOpen, &course.CreatedAt, &course.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("course not found")
		}
		return nil, err
	}

	// Update fields that are provided
	if name != nil {
		course.Name = *name
	}
	if description != nil {
		course.Description = *description
	}
	if enrollmentOpen != nil {
		course.EnrollmentOpen = *enrollmentOpen
	}

	// Update in database
	err = tx.QueryRow(ctx,
		`UPDATE courses 
                SET name = $2, description = $3, enrollment_open = $4, updated_at = $5
                WHERE id = $1 
                RETURNING id, organization_id, name, description, enrollment_open, created_at, updated_at`,
		id, course.Name, course.Description, course.EnrollmentOpen, time.Now()).Scan(
		&course.ID, &course.OrganizationID, &course.Name, &course.Description, &course.EnrollmentOpen, &course.CreatedAt, &course.UpdatedAt)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &course, nil
}

// UpdateEnrollmentStatus updates a course's enrollment status
func (r *CourseRepository) UpdateEnrollmentStatus(ctx context.Context, id string, enrollmentOpen bool) error {
	commandTag, err := r.db.Pool.Exec(ctx,
		`UPDATE courses 
                SET enrollment_open = $2, updated_at = $3
                WHERE id = $1`,
		id, enrollmentOpen, time.Now())

	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("course not found")
	}
	return nil
}

// Delete deletes a course
func (r *CourseRepository) Delete(ctx context.Context, id string) error {
	commandTag, err := r.db.Pool.Exec(ctx,
		`DELETE FROM courses WHERE id = $1`,
		id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("course not found")
	}
	return nil
}

// CountByOrganization counts courses in an organization
func (r *CourseRepository) CountByOrganization(ctx context.Context, organizationID string) (int, error) {
	var count int
	err := r.db.Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM courses WHERE organization_id = $1`,
		organizationID).Scan(&count)
	return count, err
}

// AssignTeacher assigns a teacher to a course
func (r *CourseRepository) AssignTeacher(ctx context.Context, courseID, teacherID string) error {
	_, err := r.db.Pool.Exec(ctx,
		`INSERT INTO course_teachers (course_id, teacher_id) 
                VALUES ($1, $2)`,
		courseID, teacherID)
	return err
}

// RemoveTeacher removes a teacher from a course
func (r *CourseRepository) RemoveTeacher(ctx context.Context, courseID, teacherID string) error {
	commandTag, err := r.db.Pool.Exec(ctx,
		`DELETE FROM course_teachers 
                WHERE course_id = $1 AND teacher_id = $2`,
		courseID, teacherID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("teacher not assigned to course")
	}
	return nil
}

// IsTeacherAssigned checks if a teacher is assigned to a course
func (r *CourseRepository) IsTeacherAssigned(ctx context.Context, courseID, teacherID string) (bool, error) {
	var exists bool
	err := r.db.Pool.QueryRow(ctx,
		`SELECT EXISTS(
                        SELECT 1 FROM course_teachers 
                        WHERE course_id = $1 AND teacher_id = $2
                )`,
		courseID, teacherID).Scan(&exists)
	return exists, err
}

// FindTeachersByCourse retrieves all teachers assigned to a course
func (r *CourseRepository) FindTeachersByCourse(ctx context.Context, courseID string) ([]*models.User, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT u.id, u.organization_id, u.email, u.first_name, u.last_name, u.role, u.created_at, u.updated_at
                FROM users u
                JOIN course_teachers ct ON u.id = ct.teacher_id
                WHERE ct.course_id = $1
                ORDER BY u.first_name, u.last_name`,
		courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []*models.User
	for rows.Next() {
		var teacher models.User
		if err := rows.Scan(&teacher.ID, &teacher.OrganizationID, &teacher.Email, &teacher.FirstName, &teacher.LastName, &teacher.Role, &teacher.CreatedAt, &teacher.UpdatedAt); err != nil {
			return nil, err
		}
		teachers = append(teachers, &teacher)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return teachers, nil
}

// CountTeachersByCourse counts teachers assigned to a course
func (r *CourseRepository) CountTeachersByCourse(ctx context.Context, courseID string) (int, error) {
	var count int
	err := r.db.Pool.QueryRow(ctx,
		`SELECT COUNT(*) 
                FROM course_teachers 
                WHERE course_id = $1`,
		courseID).Scan(&count)
	return count, err
}

// EnrollStudent enrolls a student in a course
func (r *CourseRepository) EnrollStudent(ctx context.Context, courseID, studentID string) error {
	_, err := r.db.Pool.Exec(ctx,
		`INSERT INTO course_enrollments (course_id, student_id) 
                VALUES ($1, $2)`,
		courseID, studentID)
	return err
}

// UnenrollStudent removes a student from a course
func (r *CourseRepository) UnenrollStudent(ctx context.Context, courseID, studentID string) error {
	commandTag, err := r.db.Pool.Exec(ctx,
		`DELETE FROM course_enrollments 
                WHERE course_id = $1 AND student_id = $2`,
		courseID, studentID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("student not enrolled in course")
	}
	return nil
}

// IsStudentEnrolled checks if a student is enrolled in a course
func (r *CourseRepository) IsStudentEnrolled(ctx context.Context, courseID, studentID string) (bool, error) {
	var exists bool
	err := r.db.Pool.QueryRow(ctx,
		`SELECT EXISTS(
                        SELECT 1 FROM course_enrollments 
                        WHERE course_id = $1 AND student_id = $2
                )`,
		courseID, studentID).Scan(&exists)
	return exists, err
}

// FindStudentsByCourse retrieves all students enrolled in a course
func (r *CourseRepository) FindStudentsByCourse(ctx context.Context, courseID string) ([]*models.User, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT u.id, u.organization_id, u.email, u.first_name, u.last_name, u.role, u.created_at, u.updated_at
                FROM users u
                JOIN course_enrollments ce ON u.id = ce.student_id
                WHERE ce.course_id = $1
                ORDER BY u.first_name, u.last_name`,
		courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*models.User
	for rows.Next() {
		var student models.User
		if err := rows.Scan(&student.ID, &student.OrganizationID, &student.Email, &student.FirstName, &student.LastName, &student.Role, &student.CreatedAt, &student.UpdatedAt); err != nil {
			return nil, err
		}
		students = append(students, &student)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

// CountStudentsByCourse counts students enrolled in a course
func (r *CourseRepository) CountStudentsByCourse(ctx context.Context, courseID string) (int, error) {
	var count int
	err := r.db.Pool.QueryRow(ctx,
		`SELECT COUNT(*) 
                FROM course_enrollments 
                WHERE course_id = $1`,
		courseID).Scan(&count)
	return count, err
}

// FindByTeacher retrieves all courses assigned to a teacher
func (r *CourseRepository) FindByTeacher(ctx context.Context, teacherID string) ([]*models.Course, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT c.id, c.organization_id, c.name, c.description, c.enrollment_open, c.created_at, c.updated_at
                FROM courses c
                JOIN course_teachers ct ON c.id = ct.course_id
                WHERE ct.teacher_id = $1
                ORDER BY c.name`,
		teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*models.Course
	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.ID, &course.OrganizationID, &course.Name, &course.Description, &course.EnrollmentOpen, &course.CreatedAt, &course.UpdatedAt); err != nil {
			return nil, err
		}
		courses = append(courses, &course)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

// FindByStudent retrieves all courses a student is enrolled in
func (r *CourseRepository) FindByStudent(ctx context.Context, studentID string) ([]*models.Course, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT c.id, c.organization_id, c.name, c.description, c.enrollment_open, c.created_at, c.updated_at
                FROM courses c
                JOIN course_enrollments ce ON c.id = ce.course_id
                WHERE ce.student_id = $1
                ORDER BY c.name`,
		studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*models.Course
	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.ID, &course.OrganizationID, &course.Name, &course.Description, &course.EnrollmentOpen, &course.CreatedAt, &course.UpdatedAt); err != nil {
			return nil, err
		}
		courses = append(courses, &course)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

// FindAvailableCourses retrieves all courses available for enrollment
func (r *CourseRepository) FindAvailableCourses(ctx context.Context, organizationID, studentID string) ([]*models.Course, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT c.id, c.organization_id, c.name, c.description, c.enrollment_open, c.created_at, c.updated_at
                FROM courses c
                WHERE c.organization_id = $1
                AND c.enrollment_open = true
                AND NOT EXISTS (
                        SELECT 1 FROM course_enrollments ce 
                        WHERE ce.course_id = c.id AND ce.student_id = $2
                )
                AND EXISTS (
                        SELECT 1 FROM course_teachers ct
                        WHERE ct.course_id = c.id
                )
                ORDER BY c.name`,
		organizationID, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*models.Course
	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.ID, &course.OrganizationID, &course.Name, &course.Description, &course.EnrollmentOpen, &course.CreatedAt, &course.UpdatedAt); err != nil {
			return nil, err
		}
		courses = append(courses, &course)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

// CountByTeacher counts courses assigned to a teacher
func (r *CourseRepository) CountByTeacher(ctx context.Context, teacherID string) (int, error) {
	var count int
	err := r.db.Pool.QueryRow(ctx,
		`SELECT COUNT(DISTINCT course_id) 
                FROM course_teachers 
                WHERE teacher_id = $1`,
		teacherID).Scan(&count)
	return count, err
}

// CountByStudent counts courses a student is enrolled in
func (r *CourseRepository) CountByStudent(ctx context.Context, studentID string) (int, error) {
	var count int
	err := r.db.Pool.QueryRow(ctx,
		`SELECT COUNT(DISTINCT course_id) 
                FROM course_enrollments 
                WHERE student_id = $1`,
		studentID).Scan(&count)
	return count, err
}

// ExecuteInTransaction executes a function within a transaction
func (r *CourseRepository) ExecuteInTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	return r.db.ExecuteTransaction(ctx, fn)
}
