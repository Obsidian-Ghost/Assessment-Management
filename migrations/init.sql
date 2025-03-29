-- Initial schema for the assessment management system

-- Organizations table
CREATE TABLE IF NOT EXISTS organizations (
                                             id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slogan VARCHAR(1000),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                             );

-- User roles enum
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
CREATE TYPE user_role AS ENUM ('admin', 'teacher', 'student');
END IF;
END $$;

-- Users table
CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role user_role NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                                                   CONSTRAINT unique_email_per_org UNIQUE (email, organization_id)
    );

-- Create index on email for faster lookup
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_organization ON users(organization_id);

-- Courses table
CREATE TABLE IF NOT EXISTS courses (
                                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    enrollment_open BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                                                   CONSTRAINT unique_course_name_per_org UNIQUE (name, organization_id)
    );

CREATE INDEX IF NOT EXISTS idx_courses_organization ON courses(organization_id);

-- Course teachers (many-to-many relationship)
CREATE TABLE IF NOT EXISTS course_teachers (
                                               course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                                                                                  PRIMARY KEY (course_id, teacher_id)
    );

-- Course enrollments (student enrollments)
CREATE TABLE IF NOT EXISTS course_enrollments (
                                                  course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    enrolled_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                                                                                     PRIMARY KEY (course_id, student_id)
    );

-- Assessment types enum
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'assessment_type') THEN
CREATE TYPE assessment_type AS ENUM ('quiz', 'exam', 'assignment', 'project');
END IF;
END $$;

-- Assessments table
CREATE TABLE IF NOT EXISTS assessments (
                                           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    type assessment_type NOT NULL,
    max_score INT NOT NULL,
    due_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                                                       );

CREATE INDEX IF NOT EXISTS idx_assessments_course ON assessments(course_id);
CREATE INDEX IF NOT EXISTS idx_assessments_teacher ON assessments(teacher_id);

-- Assessment submissions
CREATE TABLE IF NOT EXISTS assessment_submissions (
                                                      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT,
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                                               CONSTRAINT unique_student_submission UNIQUE (assessment_id, student_id)
    );

-- Grades
CREATE TABLE IF NOT EXISTS grades (
                                      submission_id UUID PRIMARY KEY REFERENCES assessment_submissions(id) ON DELETE CASCADE,
    score NUMERIC(5, 2) NOT NULL,
    feedback TEXT,
    graded_by UUID NOT NULL REFERENCES users(id),
    graded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                                                                                                               );

-- Triggers for updated_at timestamps
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply triggers to tables
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_organizations_timestamp') THEN
CREATE TRIGGER update_organizations_timestamp
    BEFORE UPDATE ON organizations
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();
END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_users_timestamp') THEN
CREATE TRIGGER update_users_timestamp
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();
END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_courses_timestamp') THEN
CREATE TRIGGER update_courses_timestamp
    BEFORE UPDATE ON courses
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();
END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_assessments_timestamp') THEN
CREATE TRIGGER update_assessments_timestamp
    BEFORE UPDATE ON assessments
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();
END IF;
END $$;
