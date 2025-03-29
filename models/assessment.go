package models

import (
	"time"
)

// AssessmentType represents the type of assessment
type AssessmentType string

const (
	AssessmentTypeQuiz       AssessmentType = "quiz"
	AssessmentTypeExam       AssessmentType = "exam"
	AssessmentTypeAssignment AssessmentType = "assignment"
	AssessmentTypeProject    AssessmentType = "project"
)

// Assessment represents an assessment that teachers create for courses
type Assessment struct {
	ID          string         `json:"id"`
	CourseID    string         `json:"course_id"`
	TeacherID   string         `json:"teacher_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Type        AssessmentType `json:"type"`
	MaxScore    int            `json:"max_score"`
	DueDate     *time.Time     `json:"due_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// AssessmentSubmission represents a student's submission for an assessment
type AssessmentSubmission struct {
	ID           string    `json:"id"`
	AssessmentID string    `json:"assessment_id"`
	StudentID    string    `json:"student_id"`
	Content      string    `json:"content"`
	SubmittedAt  time.Time `json:"submitted_at"`
}

// Grade represents the grade given to a student's assessment submission
type Grade struct {
	SubmissionID string    `json:"submission_id"`
	Score        float64   `json:"score"`
	Feedback     string    `json:"feedback"`
	GradedBy     string    `json:"graded_by"`
	GradedAt     time.Time `json:"graded_at"`
}

// AssessmentWithSubmissionCount combines an assessment with submission statistics
type AssessmentWithSubmissionCount struct {
	Assessment      *Assessment `json:"assessment"`
	SubmissionCount int         `json:"submission_count"`
	GradedCount     int         `json:"graded_count"`
	UngradedCount   int         `json:"ungraded_count"`
}

// AssessmentWithDetails combines an assessment with teacher and course details
type AssessmentWithDetails struct {
	Assessment  *Assessment `json:"assessment"`
	TeacherName string      `json:"teacher_name"`
	CourseName  string      `json:"course_name"`
}

// StudentAssessmentStatus represents a student's status for an assessment
type StudentAssessmentStatus struct {
	Assessment   *Assessment           `json:"assessment"`
	Submission   *AssessmentSubmission `json:"submission,omitempty"`
	Grade        *Grade                `json:"grade,omitempty"`
	HasSubmitted bool                  `json:"has_submitted"`
	IsGraded     bool                  `json:"is_graded"`
	DaysUntilDue int                   `json:"days_until_due,omitempty"`
	IsOverdue    bool                  `json:"is_overdue"`
}

// CreateAssessmentRequest represents the data needed to create a new assessment
type CreateAssessmentRequest struct {
	CourseID    string         `json:"course_id" validate:"required"`
	Title       string         `json:"title" validate:"required,min=3,max=255"`
	Description string         `json:"description"`
	Type        AssessmentType `json:"type" validate:"required,oneof=quiz exam assignment project"`
	MaxScore    int            `json:"max_score" validate:"required,min=1"`
	DueDate     *time.Time     `json:"due_date"`
}

// UpdateAssessmentRequest represents the data needed to update an assessment
type UpdateAssessmentRequest struct {
	Title       *string         `json:"title" validate:"omitempty,min=3,max=255"`
	Description *string         `json:"description"`
	Type        *AssessmentType `json:"type" validate:"omitempty,oneof=quiz exam assignment project"`
	MaxScore    *int            `json:"max_score" validate:"omitempty,min=1"`
	DueDate     *time.Time      `json:"due_date"`
}

// CreateSubmissionRequest represents the data needed to create a new submission
type CreateSubmissionRequest struct {
	Content string `json:"content" validate:"required"`
}

// GradeSubmissionRequest represents the data needed to grade a submission
type GradeSubmissionRequest struct {
	Score    float64 `json:"score" validate:"required,min=0"`
	Feedback string  `json:"feedback"`
}
