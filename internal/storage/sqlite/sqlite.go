package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/as-harudeen/students-api/internal/config"
	"github.com/as-harudeen/students-api/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		slog.Error("failed to connect db")
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        email TEXT,
        age INTEGER
        );
`)

	if err != nil {
		slog.Error("failed to create table")
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int64) (int64, error) {

	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	slog.Info("Student has been save successfully", slog.String("id:", strconv.Itoa(int(id))))

	return id, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, email, name, age FROM students")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Email, &student.Name, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %s", id)
		}
		return types.Student{}, nil
	}

	return student, nil

}

func (s *Sqlite) GetStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, email, name, age FROM students")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()

	if err != nil {
		return nil, err
	}

	var students []types.Student

	for rows.Next() {
		var student types.Student

		err := rows.Scan(&student.Id, &student.Email, &student.Name, &student.Age)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	return students, nil
}
