package request

type Association struct {
	Parent   string `json:"parent"`
	Children string `json:"children"`
}

func (a Association) isReq() {}
