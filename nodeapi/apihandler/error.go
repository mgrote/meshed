package apihandler

// JSONError is a json error message wrapper
type JSONError struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}
