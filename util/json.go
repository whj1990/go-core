package util

import (
	"bytes"
	"encoding/json"
)

func JsonMarshal(v any) []byte {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(v)
	return bf.Bytes()
}
