-- Drop enrollments table first due to dependencies
DROP TABLE IF EXISTS enrollments;

DROP TABLE IF EXISTS students;
CREATE TABLE students (
    student_id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    gender VARCHAR(10) NOT NULL,
    major VARCHAR(100) NOT NULL
);

DROP TABLE IF EXISTS courses;
CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    max_capacity INT NOT NULL,
    credits INT NOT NULL
);

CREATE TABLE enrollments (
    student_id VARCHAR(50) REFERENCES students(student_id) ON DELETE CASCADE,
    course_id INT REFERENCES courses(id) ON DELETE CASCADE,
    enrollment_date DATE NOT NULL,
    grade FLOAT,
    PRIMARY KEY (student_id, course_id)
);

-- Add indexes for faster lookups
CREATE INDEX idx_student_id ON enrollments(student_id);
CREATE INDEX idx_course_id ON enrollments(course_id);

-- Dummy data for students
INSERT INTO students (student_id, name, email, gender, major) VALUES
('S001', 'John Doe', 'john@example.com', 'Male', 'Computer Science'),
('S002', 'Jane Smith', 'jane@example.com', 'Female', 'Electrical Engineering'),
('S003', 'Michael Johnson', 'michael@example.com', 'Male', 'Physics'),
('S004', 'Emily Brown', 'emily@example.com', 'Female', 'Mathematics'),
('S005', 'David Lee', 'david@example.com', 'Male', 'Chemistry'),
('S006', 'Sarah Taylor', 'sarah@example.com', 'Female', 'Biology'),
('S007', 'Robert Garcia', 'robert@example.com', 'Male', 'History'),
('S008', 'Jennifer Martinez', 'jennifer@example.com', 'Female', 'Sociology'),
('S009', 'Daniel Rodriguez', 'daniel@example.com', 'Male', 'Economics'),
('S010', 'Michelle Wilson', 'michelle@example.com', 'Female', 'Political Science'),
('S011', 'Christopher Anderson', 'christopher@example.com', 'Male', 'English'),
('S012', 'Amanda Hernandez', 'amanda@example.com', 'Female', 'Art'),
('S013', 'Matthew Lopez', 'matthew@example.com', 'Male', 'Music'),
('S014', 'Ashley Gonzalez', 'ashley@example.com', 'Female', 'Psychology'),
('S015', 'Joshua Perez', 'joshua@example.com', 'Male', 'Geography'),
('S016', 'Megan Carter', 'megan@example.com', 'Female', 'Anthropology'),
('S017', 'Kevin Evans', 'kevin@example.com', 'Male', 'Philosophy'),
('S018', 'Rachel Flores', 'rachel@example.com', 'Female', 'Foreign Languages'),
('S019', 'Brandon Torres', 'brandon@example.com', 'Male', 'Engineering'),
('S020', 'Stephanie Murphy', 'stephanie@example.com', 'Female', 'Environmental Science');

-- Dummy data for courses
INSERT INTO courses (name, max_capacity, credits) VALUES
('Introduction to Programming', 30, 3),
('Database Management Systems', 25, 4),
('Linear Algebra', 20, 3),
('Electromagnetism', 15, 4),
('Computer Networks', 10, 5);

-- Dummy data for enrollments
-- Dummy data for enrollments
INSERT INTO enrollments (student_id, course_id, enrollment_date, grade) VALUES
('S001', 5, '2023-09-01', 3.5),
('S002', 5, '2023-09-01', 4.0),
('S003', 5, '2023-09-02', 3.7),
('S004', 5, '2023-09-03', 3.9),
('S005', 5, '2023-09-03', 3.8),
('S006', 5, '2023-09-03', 4.0),
('S007', 5, '2023-09-03', 3.9),
('S008', 5, '2023-09-03', 3.7),
('S009', 5, '2023-09-03', 4.0),
('S010', 5, '2023-09-03', 3.5),
('S011', 4, '2023-09-04', 4.0),
('S012', 2, '2023-09-04', 3.6),
('S013', 1, '2023-09-04', 3.9),
('S014', 3, '2023-09-04', 3.8),
('S015', 4, '2023-09-04', 4.0),
('S016', 2, '2023-09-05', 3.7),
('S017', 3, '2023-09-05', 3.9),
('S018', 4, '2023-09-05', 3.8),
('S019', 1, '2023-09-05', 3.5),
('S020', 2, '2023-09-05', 4.0),
('S001', 3, '2023-09-05', 3.6),
('S002', 4, '2023-09-05', 3.9);