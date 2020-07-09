package server

type ReplyType string

const (
	ReplyError ReplyType = "Error"
)

type ReplyStructure struct {
	Error *ReplyErrorStructure `json:"error,omitempty"`
	Type  ReplyType            `json:"type,omitempty"`
	Data  interface{}          `json:"data"`
}

type ReplyErrorStructure struct {
	Code    string `json:"code,omitempty"`
	Title   string `json:"title,omitempty"`
	Message string `json:"message,omitempty"`
}
