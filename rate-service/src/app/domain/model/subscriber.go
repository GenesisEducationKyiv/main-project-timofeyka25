package model

type Subscriber struct {
	Email string `json:"email" validate:"required,email"`
}
