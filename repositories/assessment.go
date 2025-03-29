package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"

	"assessment-management-system/db"
	"assessment-management-system/models"
)

// AssessmentRepository handles database operations for assessments
type AssessmentRepository struct {
	db *db.DB
}

// NewAssessmentRepository creates a new AssessmentRepository
func NewAssessmentRepository(db *db.DB) *AssessmentRepository {
	return &AssessmentRepository{
		db: db,
	}
}

// Create creates a new assessment
func (r *AssessmentRepository) Create(ctx context.Context, courseID, teacherID, title, description, assessmentType string, maxScore int, dueDate *time.Time) (*models.Assessment, error) {
	var assessment models.Assessment
	err := r.db.Pool.QueryRow(ctx,
		`INSERT INTO assessments (course_id, teacher_id, title, description, type, max_score, due_date) 
                VALUES ($1, $2, $3, $4, $5, $6, $7) 
                RETURNING id, course_id, teacher_id, title, description, type, max_score, due_date, created_at, updated_at`,
		courseID, teacherID, title, description, assessmentType, maxScore, dueDate).Scan(
		&assessment.ID, &assessment.CourseID, &assessment.TeacherID, &assessment.Title, &assessment.Description,
		&assessment.Type, &assessment.MaxScore, &assessment.DueDate, &assessment.CreatedAt, &assessment.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &assessment, nil
}

// FindByID retrieves an assessment by ID
func (r *AssessmentRepository) FindByID(ctx context.Context, id string) (*models.Assessment, error) {
	var assessment models.Assessment
	var dueDate *time.Time
	err := r.db.Pool.QueryRow(ctx,
		`SELECT id, course_id, teacher_id, title, description, type, max_score, due_date, created_at, updated_at 
                FROM assessments 
                WHERE id = $1`,
		id).Scan(&assessment.ID, &assessment.CourseID, &assessment.TeacherID, &assessment.Title, &assessment.Description,
		&assessment.Type, &assessment.MaxScore, &dueDate, &assessment.CreatedAt, &assessment.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	assessment.DueDate = dueDate
	return &assessment, nil
}

// FindByCourse retrieves all assessments for a course
func (r *AssessmentRepository) FindByCourse(ctx context.Context, courseID string) ([]*models.Assessment, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT id, course_id, teacher_id, title, description, type, max_score, due_date, created_at, updated_at 
                FROM assessments 
                WHERE course_id = $1
                ORDER BY due_date NULLS LAST, created_at DESC`,
		courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assessments []*models.Assessment
	for rows.Next() {
		var assessment models.Assessment
		var dueDate *time.Time
		if err := rows.Scan(&assessment.ID, &assessment.CourseID, &assessment.TeacherID, &assessment.Title, &assessment.Description,
			&assessment.Type, &assessment.MaxScore, &dueDate, &assessment.CreatedAt, &assessment.UpdatedAt); err != nil {
			return nil, err
		}
		assessment.DueDate = dueDate
		assessments = append(assessments, &assessment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assessments, nil
}

// FindByOrganization retrieves all assessments for an organization
func (r *AssessmentRepository) FindByOrganization(ctx context.Context, organizationID, courseID string) ([]*models.Assessment, error) {
	var query string
	var args []interface{}

	if courseID != "" {
		query = `
                        SELECT a.id, a.course_id, a.teacher_id, a.title, a.description, a.type, a.max_score, a.due_date, a.created_at, a.updated_at 
                        FROM assessments a
                        JOIN courses c ON a.course_id = c.id
                        WHERE c.organization_id = $1 AND a.course_id = $2
                        ORDER BY a.due_date NULLS LAST, a.created_at DESC`
		args = append(args, organizationID, courseID)
	} else {
		query = `
                        SELECT a.id, a.course_id, a.teacher_id, a.title, a.description, a.type, a.max_score, a.due_date, a.created_at, a.updated_at 
                        FROM assessments a
                        JOIN courses c ON a.course_id = c.id
                        WHERE c.organization_id = $1
                        ORDER BY a.due_date NULLS LAST, a.created_at DESC`
		args = append(args, organizationID)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assessments []*models.Assessment
	for rows.Next() {
		var assessment models.Assessment
		var dueDate *time.Time
		if err := rows.Scan(&assessment.ID, &assessment.CourseID, &assessment.TeacherID, &assessment.Title, &assessment.Description,
			&assessment.Type, &assessment.MaxScore, &dueDate, &assessment.CreatedAt, &assessment.UpdatedAt); err != nil {
			return nil, err
		}
		assessment.DueDate = dueDate
		assessments = append(assessments, &assessment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assessments, nil
}

// Update updates an assessment
func (r *AssessmentRepository) Update(ctx context.Context, id string, title *string, description *string, assessmentType *models.AssessmentType, maxScore *int, dueDate *time.Time) (*models.Assessment, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Get current assessment
	var assessment models.Assessment
	var currentDueDate *time.Time
	err = tx.QueryRow(ctx,
		`SELECT id, course_id, teacher_id, title, description, type, max_score, due_date, created_at, updated_at 
                FROM assessments 
                WHERE id = $1`,
		id).Scan(&assessment.ID, &assessment.CourseID, &assessment.TeacherID, &assessment.Title, &assessment.Description,
		&assessment.Type, &assessment.MaxScore, &currentDueDate, &assessment.CreatedAt, &assessment.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("assessment not found")
		}
		return nil, err
	}

	assessment.DueDate = currentDueDate

	// Update fields that are provided
	if title != nil {
		assessment.Title = *title
	}
	if description != nil {
		assessment.Description = *description
	}
	if assessmentType != nil {
		assessment.Type = *assessmentType
	}
	if maxScore != nil {
		assessment.MaxScore = *maxScore
	}
	// Due date handling is special because we need to distinguish between setting to null and not changing
	if dueDate != nil {
		assessment.DueDate = dueDate
	}

	// Update in database
	var updatedDueDate *time.Time
	err = tx.QueryRow(ctx,
		`UPDATE assessments 
                SET title = $2, description = $3, type = $4, max_score = $5, due_date = $6, updated_at = $7
                WHERE id = $1 
                RETURNING id, course_id, teacher_id, title, description, type, max_score, due_date, created_at, updated_at`,
		id, assessment.Title, assessment.Description, assessment.Type, assessment.MaxScore, assessment.DueDate, time.Now()).Scan(
		&assessment.ID, &assessment.CourseID, &assessment.TeacherID, &assessment.Title, &assessment.Description,
		&assessment.Type, &assessment.MaxScore, &updatedDueDate, &assessment.CreatedAt, &assessment.UpdatedAt)

	if err != nil {
		return nil, err
	}

	assessment.DueDate = updatedDueDate

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &assessment, nil
}

// Delete deletes an assessment
func (r *AssessmentRepository) Delete(ctx context.Context, id string) error {
	commandTag, err := r.db.Pool.Exec(ctx,
		`DELETE FROM assessments WHERE id = $1`,
		id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("assessment not found")
	}
	return nil
}

// FindSubmissionsByAssessment retrieves all submissions for an assessment
func (r *AssessmentRepository) FindSubmissionsByAssessment(ctx context.Context, assessmentID string) ([]*models.AssessmentSubmission, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT id, assessment_id, student_id, content, submitted_at 
                FROM assessment_submissions 
                WHERE assessment_id = $1
                ORDER BY submitted_at DESC`,
		assessmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []*models.AssessmentSubmission
	for rows.Next() {
		var submission models.AssessmentSubmission
		if err := rows.Scan(&submission.ID, &submission.AssessmentID, &submission.StudentID, &submission.Content, &submission.SubmittedAt); err != nil {
			return nil, err
		}
		submissions = append(submissions, &submission)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return submissions, nil
}

// FindSubmissionByID retrieves a submission by ID
func (r *AssessmentRepository) FindSubmissionByID(ctx context.Context, id string) (*models.AssessmentSubmission, error) {
	var submission models.AssessmentSubmission
	err := r.db.Pool.QueryRow(ctx,
		`SELECT id, assessment_id, student_id, content, submitted_at 
                FROM assessment_submissions 
                WHERE id = $1`,
		id).Scan(&submission.ID, &submission.AssessmentID, &submission.StudentID, &submission.Content, &submission.SubmittedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &submission, nil
}

// FindSubmissionByStudentAndAssessment retrieves a student's submission for an assessment
func (r *AssessmentRepository) FindSubmissionByStudentAndAssessment(ctx context.Context, assessmentID, studentID string) (*models.AssessmentSubmission, error) {
	var submission models.AssessmentSubmission
	err := r.db.Pool.QueryRow(ctx,
		`SELECT id, assessment_id, student_id, content, submitted_at 
                FROM assessment_submissions 
                WHERE assessment_id = $1 AND student_id = $2`,
		assessmentID, studentID).Scan(&submission.ID, &submission.AssessmentID, &submission.StudentID, &submission.Content, &submission.SubmittedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &submission, nil
}

// CreateSubmission creates a new submission
func (r *AssessmentRepository) CreateSubmission(ctx context.Context, assessmentID, studentID, content string) (*models.AssessmentSubmission, error) {
	var submission models.AssessmentSubmission
	err := r.db.Pool.QueryRow(ctx,
		`INSERT INTO assessment_submissions (assessment_id, student_id, content) 
                VALUES ($1, $2, $3) 
                RETURNING id, assessment_id, student_id, content, submitted_at`,
		assessmentID, studentID, content).Scan(&submission.ID, &submission.AssessmentID, &submission.StudentID, &submission.Content, &submission.SubmittedAt)

	if err != nil {
		return nil, err
	}
	return &submission, nil
}

// FindGradeBySubmission retrieves the grade for a submission
func (r *AssessmentRepository) FindGradeBySubmission(ctx context.Context, submissionID string) (*models.Grade, error) {
	var grade models.Grade
	err := r.db.Pool.QueryRow(ctx,
		`SELECT submission_id, score, feedback, graded_by, graded_at 
                FROM grades 
                WHERE submission_id = $1`,
		submissionID).Scan(&grade.SubmissionID, &grade.Score, &grade.Feedback, &grade.GradedBy, &grade.GradedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &grade, nil
}

// CreateGrade creates a new grade
func (r *AssessmentRepository) CreateGrade(ctx context.Context, submissionID string, score float64, feedback, gradedBy string) (*models.Grade, error) {
	var grade models.Grade
	err := r.db.Pool.QueryRow(ctx,
		`INSERT INTO grades (submission_id, score, feedback, graded_by) 
                VALUES ($1, $2, $3, $4) 
                RETURNING submission_id, score, feedback, graded_by, graded_at`,
		submissionID, score, feedback, gradedBy).Scan(&grade.SubmissionID, &grade.Score, &grade.Feedback, &grade.GradedBy, &grade.GradedAt)

	if err != nil {
		return nil, err
	}
	return &grade, nil
}

// UpdateGrade updates a grade
func (r *AssessmentRepository) UpdateGrade(ctx context.Context, submissionID string, score float64, feedback, gradedBy string) (*models.Grade, error) {
	var grade models.Grade
	err := r.db.Pool.QueryRow(ctx,
		`UPDATE grades 
                SET score = $2, feedback = $3, graded_by = $4, graded_at = $5
                WHERE submission_id = $1 
                RETURNING submission_id, score, feedback, graded_by, graded_at`,
		submissionID, score, feedback, gradedBy, time.Now()).Scan(&grade.SubmissionID, &grade.Score, &grade.Feedback, &grade.GradedBy, &grade.GradedAt)

	if err != nil {
		return nil, err
	}
	return &grade, nil
}

// CountByTeacher counts assessments created by a teacher
func (r *AssessmentRepository) CountByTeacher(ctx context.Context, teacherID string) (int, error) {
	var count int
	err := r.db.Pool.QueryRow(ctx,
		`SELECT COUNT(*) 
                FROM assessments 
                WHERE teacher_id = $1`,
		teacherID).Scan(&count)
	return count, err
}

// CountSubmissionsByStudent counts submissions by a student
func (r *AssessmentRepository) CountSubmissionsByStudent(ctx context.Context, studentID string) (int, error) {
	var count int
	err := r.db.Pool.QueryRow(ctx,
		`SELECT COUNT(*) 
                FROM assessment_submissions 
                WHERE student_id = $1`,
		studentID).Scan(&count)
	return count, err
}

// CountUngradedSubmissionsByTeacher counts ungraded submissions for assessments by a teacher
func (r *AssessmentRepository) CountUngradedSubmissionsByTeacher(ctx context.Context, teacherID string) (int, error) {
	var count int
	err := r.db.Pool.QueryRow(ctx,
		`SELECT COUNT(s.id) 
                FROM assessment_submissions s
                JOIN assessments a ON s.assessment_id = a.id
                LEFT JOIN grades g ON s.id = g.submission_id
                WHERE a.teacher_id = $1 AND g.submission_id IS NULL`,
		teacherID).Scan(&count)
	return count, err
}

// CountPendingAssessmentsByStudent counts pending assessments for a student
func (r *AssessmentRepository) CountPendingAssessmentsByStudent(ctx context.Context, studentID string) (int, error) {
	var count int
	err := r.db.Pool.QueryRow(ctx,
		`SELECT COUNT(a.id)
                FROM assessments a
                JOIN courses c ON a.course_id = c.id
                JOIN course_enrollments ce ON c.id = ce.course_id
                WHERE ce.student_id = $1
                AND NOT EXISTS (
                        SELECT 1 FROM assessment_submissions s
                        WHERE s.assessment_id = a.id AND s.student_id = $1
                )
                AND (a.due_date IS NULL OR a.due_date > CURRENT_TIMESTAMP)`,
		studentID).Scan(&count)
	return count, err
}

// GetAverageGradeForStudent calculates the average grade for a student
func (r *AssessmentRepository) GetAverageGradeForStudent(ctx context.Context, studentID string) (float64, error) {
	var avgGrade float64
	err := r.db.Pool.QueryRow(ctx,
		`SELECT COALESCE(AVG(g.score), 0)
                FROM grades g
                JOIN assessment_submissions s ON g.submission_id = s.id
                WHERE s.student_id = $1`,
		studentID).Scan(&avgGrade)
	return avgGrade, err
}

// ExecuteInTransaction executes a function within a transaction
func (r *AssessmentRepository) ExecuteInTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	return r.db.ExecuteTransaction(ctx, fn)
}
