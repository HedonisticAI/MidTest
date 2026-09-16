package http

import (
	"encoding/json"
	"errors"
	"io"
	stdhttp "net/http"
	"strconv"
	"strings"

	"midtest/internal/domain"
	"midtest/internal/usecase"
)

type Handler struct {
	uc usecase.Usecase
}

func NewHandler(uc usecase.Usecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) Register(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	if r.Method != stdhttp.MethodPost {
		h.writeError(w, domain.ErrBadRequestMethod)
		return
	}

	var input usecase.RegisterInput
	if err := decodeRequest(r, &input); err != nil {
		h.writeError(w, domain.ErrBadParameter)
		return
	}

	output, err := h.uc.Register(r.Context(), input)
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, stdhttp.StatusOK, domain.DomainResponse{Response: output})
}

func (h *Handler) Auth(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	if r.Method != stdhttp.MethodPost {
		h.writeError(w, domain.ErrBadRequestMethod)
		return
	}

	var input usecase.LoginInput
	if err := decodeRequest(r, &input); err != nil {
		h.writeError(w, domain.ErrBadParameter)
		return
	}

	output, err := h.uc.LogIn(r.Context(), input)
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, stdhttp.StatusOK, domain.DomainResponse{Response: output})
}

func (h *Handler) EndSession(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	if r.Method != stdhttp.MethodDelete {
		h.writeError(w, domain.ErrBadRequestMethod)
		return
	}

	token := pathToken(r)
	if token == "" {
		token = requestToken(r)
	}
	if token == "" {
		h.writeError(w, domain.ErrBadParameter)
		return
	}

	if err := h.uc.EndSession(token); err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, stdhttp.StatusOK, domain.DomainResponse{
		Response: map[string]bool{token: true},
	})
}

func (h *Handler) CreateDocument(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	if r.Method != stdhttp.MethodPost {
		h.writeError(w, domain.ErrBadRequestMethod)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		h.writeError(w, domain.ErrBadParameter)
		return
	}

	var meta usecase.Meta
	metaRaw := r.FormValue("meta")
	if metaRaw == "" || json.Unmarshal([]byte(metaRaw), &meta) != nil {
		h.writeError(w, domain.ErrBadParameter)
		return
	}

	var document interface{}
	if raw := r.FormValue("json"); raw != "" {
		if err := json.Unmarshal([]byte(raw), &document); err != nil {
			h.writeError(w, domain.ErrBadParameter)
			return
		}
	}

	var body []byte
	file, _, err := r.FormFile("file")
	if err == nil {
		defer file.Close()
		body, err = io.ReadAll(file)
		if err != nil {
			h.writeError(w, domain.ErrInternal)
			return
		}
	} else if !errors.Is(err, stdhttp.ErrMissingFile) {
		h.writeError(w, domain.ErrBadParameter)
		return
	}

	if meta.Token == "" {
		meta.Token = requestToken(r)
	}

	output, err := h.uc.WriteFile(r.Context(), usecase.WriteFileInput{
		Meta: meta,
		Json: document,
		Body: body,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	data := map[string]interface{}{"file": output.Name}
	if len(output.Json) > 0 {
		var value interface{}
		if json.Unmarshal(output.Json, &value) == nil {
			data["json"] = value
		} else {
			data["json"] = json.RawMessage(output.Json)
		}
	}

	h.writeJSON(w, stdhttp.StatusOK, domain.DomainResponse{Data: data})
}

func (h *Handler) ListDocuments(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	if r.Method != stdhttp.MethodGet && r.Method != stdhttp.MethodHead {
		h.writeError(w, domain.ErrBadRequestMethod)
		return
	}

	if requestToken(r) == "" {
		h.writeError(w, domain.ErrBadParameter)
		return
	}

	q := r.URL.Query()
	filters := make(map[string]string)

	if key, value := q.Get("key"), q.Get("value"); key != "" && value != "" {
		filters[key] = value
	}
	if login := q.Get("login"); login != "" {
		filters["login"] = login
	}
	if limit := q.Get("limit"); limit != "" {
		n, err := strconv.Atoi(limit)
		if err != nil || n < 1 {
			h.writeError(w, domain.ErrBadParameter)
			return
		}
		filters["limit"] = strconv.Itoa(n)
	}

	files, err := h.uc.ListFiles(r.Context(), usecase.ListInput{Filters: filters})
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, stdhttp.StatusOK, domain.DomainResponse{
		Data: map[string]interface{}{"docs": files},
	})
}

func (h *Handler) GetDocument(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	if r.Method != stdhttp.MethodGet && r.Method != stdhttp.MethodHead {
		h.writeError(w, domain.ErrBadRequestMethod)
		return
	}

	id := pathID(r)
	token := requestToken(r)
	if id == "" || token == "" {
		h.writeError(w, domain.ErrBadParameter)
		return
	}

	value, err := h.uc.GetFile(r.Context(), id, token)
	if err != nil {
		h.writeError(w, err)
		return
	}

	body, ok := value.([]byte)
	if !ok {
		h.writeJSON(w, stdhttp.StatusOK, domain.DomainResponse{Data: value})
		return
	}

	// The current usecase returns []byte for both files and JSON documents.
	// JSON is wrapped in the common API response; binary content is returned as-is.
	var jsonValue interface{}
	if json.Unmarshal(body, &jsonValue) == nil {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == stdhttp.MethodHead {
			w.WriteHeader(stdhttp.StatusOK)
			return
		}
		h.writeJSON(w, stdhttp.StatusOK, domain.DomainResponse{Data: jsonValue})
		return
	}

	w.Header().Set("Content-Type", stdhttp.DetectContentType(body))
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(stdhttp.StatusOK)
	if r.Method != stdhttp.MethodHead {
		_, _ = w.Write(body)
	}
}

func (h *Handler) DeleteDocument(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	if r.Method != stdhttp.MethodDelete {
		h.writeError(w, domain.ErrBadRequestMethod)
		return
	}

	id := pathID(r)
	token := requestToken(r)
	if id == "" || token == "" {
		h.writeError(w, domain.ErrBadParameter)
		return
	}

	if err := h.uc.DeleteFile(r.Context(), id, token); err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, stdhttp.StatusOK, domain.DomainResponse{
		Response: map[string]bool{id: true},
	})
}

func decodeRequest(r *stdhttp.Request, dst interface{}) error {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		return json.NewDecoder(r.Body).Decode(dst)
	}

	if err := r.ParseForm(); err != nil {
		return err
	}

	values := make(map[string]string, len(r.Form))
	for key, list := range r.Form {
		if len(list) > 0 {
			values[key] = list[0]
		}
	}

	encoded, err := json.Marshal(values)
	if err != nil {
		return err
	}
	return json.Unmarshal(encoded, dst)
}

func requestToken(r *stdhttp.Request) string {
	if token := r.URL.Query().Get("token"); token != "" {
		return token
	}
	if token := r.FormValue("token"); token != "" {
		return token
	}

	header := r.Header.Get("Authorization")
	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	}
	return ""
}

func pathID(r *stdhttp.Request) string {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) == 3 && parts[0] == "api" && parts[1] == "docs" {
		return parts[2]
	}
	return ""
}

func pathToken(r *stdhttp.Request) string {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) == 3 && parts[0] == "api" && parts[1] == "auth" {
		return parts[2]
	}
	return ""
}

func (h *Handler) writeError(w stdhttp.ResponseWriter, err error) {
	status := stdhttp.StatusInternalServerError
	text := domain.ErrInternal.Error()

	switch {
	case errors.Is(err, domain.ErrBadRequestMethod):
		status, text = stdhttp.StatusMethodNotAllowed, err.Error()
	case errors.Is(err, domain.ErrBadParameter):
		status, text = stdhttp.StatusBadRequest, err.Error()
	case errors.Is(err, domain.ErrUnauthorized):
		status, text = stdhttp.StatusUnauthorized, err.Error()
	case errors.Is(err, domain.ErrActionNotAuthorized):
		status, text = stdhttp.StatusForbidden, err.Error()
	case errors.Is(err, domain.ErrMethodNotReady):
		status, text = stdhttp.StatusNotImplemented, err.Error()
	case errors.Is(err, domain.ErrInternal):
		status, text = stdhttp.StatusInternalServerError, err.Error()
	}

	h.writeJSON(w, status, domain.DomainResponse{
		Error: &domain.DomainErrorField{Code: status, Text: text},
	})
}

func (h *Handler) writeJSON(w stdhttp.ResponseWriter, status int, payload domain.DomainResponse) {
	w.Header().Set("Content-Type", "application/json")
	if status == stdhttp.StatusOK && false {
		// Kept intentionally empty: HEAD is handled before this helper is called
		// for document bodies. net/http also suppresses HEAD response bodies.
	}
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
