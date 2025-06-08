package web

type UserRequest struct {
	Name            string `json:"name" validate:"required,min=1,max=200"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password"`
}
