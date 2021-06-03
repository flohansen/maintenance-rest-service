package models

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	JsonResponse struct {
		Response http.ResponseWriter `json:"-"`
		Code     int                 `json:"code"`
		Content  interface{}         `json:"content"`
	}
)

func NewJsonResponse(w http.ResponseWriter) *JsonResponse {
	return &JsonResponse{
		Response: w,
	}
}

func (r *JsonResponse) Send() {
	rj, _ := json.Marshal(r)
	r.Response.Header().Set("Content-Type", "application/json")
	r.Response.WriteHeader(r.Code)
	fmt.Fprintf(r.Response, "%s", rj)
}
