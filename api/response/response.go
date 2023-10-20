package response

type ErrorResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
