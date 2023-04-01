DROP TABLE IF EXISTS students;

CREATE TABLE students(
    id VARCHAR(32) primary key ,
    name VARCHAR(255) NOT NULL,
    age INTEGER NOT NULL
);

DROP TABLE IF EXISTS test;

CREATE TABLE test(
  id VARCHAR(32) primary key,
  name VARCHAR(255) NOT NULL
);

DROP TABLE IF EXISTS questions;

CREATE TABLE questions(
    id varchar(32) primary key,
    test_id varchar(32) NOT NULL,
    question varchar(255) NOT NULL,
    answer varchar(255) NOT NULL,
    FOREIGN KEY (test_id) REFERENCES test(id)
);

DROP TABLE IF EXISTS enrollments;

CREATE TABLE enrollments(
  student_id varchar(32) NOT NULL,
  test_id varchar(32) NOT NULL,
  FOREIGN KEY (student_id) REFERENCES students(id),
  FOREIGN KEY (test_id) REFERENCES test(id)
);