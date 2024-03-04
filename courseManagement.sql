DROP TABLE IF EXISTS students;
CREATE TABLE students (
    student_id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100),
    gender VARCHAR(10),
    major VARCHAR(100)
);

-- Insert dummy data into the Student table
INSERT INTO students (student_id, name, email, gender, major) VALUES
('1', 'John Doe', 'john.doe@example.com', 'Male', 'Computer Science'),
('2', 'Jane Smith', 'jane.smith@example.com', 'Female', 'Engineering'),
('3', 'Michael Johnson', 'michael.johnson@example.com', 'Male', 'Biology'),
('4', 'Emily Williams', 'emily.williams@example.com', 'Female', 'Psychology'),
('5', 'David Brown', 'david.brown@example.com', 'Male', 'Business Administration');
