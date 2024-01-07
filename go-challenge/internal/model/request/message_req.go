package request

type Message struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	HTTPStatus  int    `json:"httpStatus"`
}

func (m Message) isReq() {}
