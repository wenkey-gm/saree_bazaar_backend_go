package modal

type ErrorResponse struct {
	ErrorMessage string `json:"message"`
	StatusCode   int    `json:"status"`
}
