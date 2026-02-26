package auth

import (
	"regexp"

	"github.com/yancarlodev/workspaces-api/internal/platform/validation"
)

type LoginRequestDTO struct {
	Email    string
	Password string
}

func (d *LoginRequestDTO) Validate() validation.ValidationErrors {
	errors := validation.ValidationErrors{}

	if d.Email == "" {
		errors.Add("email", "required")
	}

	if ok, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, d.Email); !ok {
		errors.Add("email", "invalid")
	}

	if d.Password == "" {
		errors.Add("password", "required")
	}

	return errors
}
