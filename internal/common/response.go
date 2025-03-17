package common

type Message struct {
	Message string `json:"message"`
}

type Metadata struct {
	Service string `json:"service"`
	Version string `json:"version"`
	Staging string `json:"staging"`
	Githash string `json:"githash"`
	Gobuild string `json:"gobuild"`
	Compile string `json:"compile"`
}

type Response struct {
	Content   any    `json:"content"`
	Timestamp string `json:"timestamp"`
}

func NewResponse(content any) *Response {
	return &Response{
		Content:   content,
		Timestamp: TimestampISO3339NS(),
	}
}
