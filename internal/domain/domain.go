package domain

type DomainErrorField struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type AuthCred struct {
	Login    string
	Password string
}

type DomainResponse struct {
	Error    *DomainErrorField `json:"error,omitempty"`
	Response interface{}       `json:"response,omitempty"`
	Data     interface{}       `json:"data,omitempty"`
}

type FileInfo struct {
	Name  string
	Users []string
	ID    string
	Path  string
	File  bool
}
