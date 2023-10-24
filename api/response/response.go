package response

type ErrorResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
