package student

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/as-harudeen/students-api/internal/storage"
	"github.com/as-harudeen/students-api/types"
	"github.com/as-harudeen/students-api/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		}

		if err := validator.New().Struct(student); err != nil {
			validateErr := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErr))
			return
		}

		slog.Info("creating student", student.Name, student.Email)

		id, err := storage.CreateStudent(student.Name, student.Name, int64(student.Age))
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("internal server error")))
		}

		response.WriteJson(w, http.StatusCreated, map[string]int64{"status": id})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id")))
		}

		student, err := storage.GetStudentById(intId)
		if err != nil {
			if err == sql.ErrNoRows {
				response.WriteJson(w, http.StatusNotFound, err)
				return
			}
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]types.Student{"student": student})
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		students, err := storage.GetStudents()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("internal server error")))
			return
		}

		response.WriteJson(w, http.StatusOK, students)
	}
}
