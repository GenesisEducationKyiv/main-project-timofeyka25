package response

type RateResponse struct {
	Rate float64 `json:"rate"`
}

type SendRateResponse struct {
	UnsentEmails []string `json:"unsent"`
}
