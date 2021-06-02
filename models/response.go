package models

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	JsonResponse struct {
		Code    int         `json:"code"`
		Content interface{} `json:"content"`
	}
)

func (r *JsonResponse) Send(w http.ResponseWriter) {
	rj, _ := json.Marshal(r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Code)
	fmt.Fprintf(w, "%s", rj)
}
