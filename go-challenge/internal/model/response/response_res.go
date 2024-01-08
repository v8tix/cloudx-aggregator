package response

type Response struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Count       int    `json:"count"`
}

func NewResponse(source string, destination string, count int) Response {
	return Response{Source: source, Destination: destination, Count: count}
}
