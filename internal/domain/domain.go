package domain

import "time"

type DomainErrorField struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type AuthCred struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type DomainResponse struct {
	Error    *DomainErrorField `json:"error,omitempty"`
	Response interface{}       `json:"response,omitempty"`
	Data     interface{}       `json:"data,omitempty"`
}

type FileInfo struct {
	Name       string    `json:"name"`
	Users      []string  `json:"grants"`
	ID         string    `json:"id"`
	Path       string    `json:"path"`
	File       bool      `json:"file"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}
