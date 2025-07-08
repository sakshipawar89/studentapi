package student

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sakshipawar89/StudentApi/internal/types"
	"github.com/sakshipawar89/StudentApi/internal/utils/response"
)

var validate = validator.New()

func New() http.HandlerFunc {
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

		// Validate fields
		if err := validate.Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			log.Println("[DEBUG] Validation failed")
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		log.Println("[DEBUG] Student parsed:", student)
		response.WriteJSON(w, http.StatusCreated, map[string]string{"success": "Student created successfully"})
	}
}
