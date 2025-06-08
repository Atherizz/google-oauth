package web 

type UserResponse struct {
	Id int `json:"id"`
	GoogleId string `json:"google_id"`
	Email string `json:"email"`
	Name string `json:"name"`
	Picture string `json:"picture"`
	Provider string `json:"provider"`
	Role string `json:"role"`
}