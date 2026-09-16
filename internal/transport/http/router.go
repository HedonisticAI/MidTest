package http

import (
	stdhttp "net/http"

	"midtest/internal/usecase"
)

// NewRouter creates the REST API router.
func NewRouter(uc usecase.Usecase) stdhttp.Handler {
	h := NewHandler(uc)

	mux := stdhttp.NewServeMux()
	mux.HandleFunc("/api/register", h.Register)
	mux.HandleFunc("/api/auth", h.Auth)
	mux.HandleFunc("/api/auth/", h.EndSession)
	mux.HandleFunc("/api/docs", h.CreateOrListDocuments)
	mux.HandleFunc("/api/docs/", h.GetOrDeleteDocument)

	return mux
}
