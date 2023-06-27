package responses

type ErrorResponse struct {
	Message string `json:"message"`
}

type SendRateResponse struct {
	UnsentEmails []string `json:"unsent"`
}
