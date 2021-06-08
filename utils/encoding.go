package utils

import (
	"bytes"
	"encoding/json"
)

func MapToStruct(in, out interface{}) {
  buf := new(bytes.Buffer)
  json.NewEncoder(buf).Encode(in)
  json.NewDecoder(buf).Decode(out)
}
