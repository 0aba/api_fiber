package validators_response

type EmptySuccessfulResponse struct{}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details"`
}
