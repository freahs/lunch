package lunch

import (
	"encoding/json"
	"fmt"
)

// HTTPError is an error with a status
type HTTPError struct {
	status int
	error
}

// MarshalJSON implements the Marshaller interface
func (e HTTPError) MarshalJSON() ([]byte, error) {
	type E struct {
		S int    `json:"status"`
		M string `json:"message"`
	}
	return json.Marshal(E{e.status, e.Error()})
}

// Status returns the status code of the error
func (e HTTPError) Status() int {
	return e.status
}

// NewHTTPError returns a new HTTPError
func NewHTTPError(status int, i interface{}) HTTPError {
	switch t := i.(type) {
	case error:
		return HTTPError{status, t}
	default:
		return HTTPError{status, fmt.Errorf("%v", t)}
	}
}
