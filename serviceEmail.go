package jdsport_raffle_backend

type ResponseServiceEmail struct {
	Message      string `json:"message"`
	ErrorMessage string `json:"error_message"`
	Data         string `json:"data"`
}
