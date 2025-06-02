package web 

type UserResponse struct {
	Email string `json:"email"`
	Name string `json:"name"`
	Picture string `json:"picture"`
}