package request

type Message struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	HttpStatus  int    `json:"httpStatus"`
}
