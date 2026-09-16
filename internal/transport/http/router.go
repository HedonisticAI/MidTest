package http

import (
	stdhttp "net/http"

	"midtest/internal/domain"
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

func (h *Handler) CreateOrListDocuments(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	switch r.Method {
	case stdhttp.MethodPost:
		h.CreateDocument(w, r)
	case stdhttp.MethodGet, stdhttp.MethodHead:
		h.ListDocuments(w, r)
	default:
		h.writeError(w, domain.ErrBadRequestMethod)
	}
}

func (h *Handler) GetOrDeleteDocument(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	if pathID(r) == "" {
		h.writeError(w, domain.ErrBadParameter)
		return
	}

	switch r.Method {
	case stdhttp.MethodGet, stdhttp.MethodHead:
		h.GetDocument(w, r)
	case stdhttp.MethodDelete:
		h.DeleteDocument(w, r)
	default:
		h.writeError(w, domain.ErrBadRequestMethod)
	}
}
