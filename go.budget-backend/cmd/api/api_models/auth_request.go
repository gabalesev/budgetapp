package api_models

type RegisterUserRequest struct {
	Email     *string `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=6"`
	FirstName string  `json:"first_name" validate:"required"`
	LastName  string  `json:"last_name"`
	Gender    string  `json:"gender"`
}

type LoginRequest struct {
	Email    *string `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6"`
}

type ChangePasswordRequest struct {
	CurrentPassword    string `json:"current_password" validate:"required,min=6"`
	NewPassword        string `json:"new_password" validate:"required,min=6"`
	ConfirmNewPassword string `json:"confirm_new_password" validate:"required,eqfield=NewPassword"`
}
