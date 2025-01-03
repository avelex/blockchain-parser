package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/avelex/blockchain-parser/internal/parser"
)

var ethAddressRegex = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

type Handler struct {
	parser parser.Parser
}

func NewHandler(parser parser.Parser) *Handler {
	return &Handler{
		parser: parser,
	}
}

func (h *Handler) Register(m *http.ServeMux) {
	m.HandleFunc("GET /block", h.showCurrentBlock)
	m.HandleFunc("GET /subscribe", h.subscribeForTransactions)
	m.HandleFunc("GET /transactions", h.showTransactions)
}

func (h *Handler) showCurrentBlock(w http.ResponseWriter, r *http.Request) {
	block := h.parser.GetCurrentBlock()
	renderJSON(w, http.StatusOK, block)
}

func (h *Handler) subscribeForTransactions(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if !ethAddressRegex.MatchString(address) {
		renderJSON(w, http.StatusBadRequest, "invalid address")
		return
	}

	if !h.parser.Subscribe(address) {
		renderJSON(w, http.StatusOK, "already subscribed")
		return
	}

	renderJSON(w, http.StatusOK, "subscribed")
}

func (h *Handler) showTransactions(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if !ethAddressRegex.MatchString(address) {
		renderJSON(w, http.StatusBadRequest, "invalid address")
		return
	}

	transactions := h.parser.GetTransactions(r.Context(), address)

	renderJSON(w, http.StatusOK, transactions)
}

func renderJSON(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
