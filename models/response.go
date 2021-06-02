package models

type (
	JsonResponse struct {
		Code    int         `json:"code"`
		Content interface{} `json:"content"`
	}
)
