package domain

type DomainErrorField struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type DomainResponse struct {
	Error    *DomainErrorField `json:"error,omitempty"`
	Response interface{}       `json:"response,omitempty"`
	Data     interface{}       `json:"data,omitempty"`
}

type FileInfo struct {
	Name  string
	Users []string
	File  bool
}
