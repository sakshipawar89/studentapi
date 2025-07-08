package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"github.com/sakshipawar89/StudentApi/internal/storage"
	"github.com/sakshipawar89/StudentApi/internal/types"
	"github.com/sakshipawar89/StudentApi/internal/utils/response"
)

var validate = validator.New()

// POST /api/students
func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("[DEBUG] Student handler hit")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			log.Println("[DEBUG] Empty body received")
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		if err != nil {
			log.Println("[ERROR] Failed to decode body:", err)
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validate.Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			log.Println("[DEBUG] Validation failed")
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		id, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			log.Println("[ERROR] Failed to save student:", err)
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("User created", slog.String("userId", fmt.Sprint(id)))
		response.WriteJSON(w, http.StatusCreated, map[string]int64{"id": id})
	}
}

// GET /api/students/{id}
func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		slog.Info("Fetching student by ID", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id format: %w", err)))
			return
		}

		student, err := storage.GetStudentById(intId)
		if err != nil {
			slog.Error("Error fetching student", slog.String("id", id))
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJSON(w, http.StatusOK, student)
	}
}

// GET /api/students
func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Getting all students")

		students, err := storage.GetStudents()
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJSON(w, http.StatusOK, students)
	}
}
