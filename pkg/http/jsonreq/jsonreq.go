package jsonreq

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Unmarshal(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("read request body: %w", err)
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return fmt.Errorf("unmarshal body: %w", err)
	}

	return nil
}
