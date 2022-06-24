package pestotrap

import (
	"net/http"

	_ "embed"

	"github.com/blevesearch/bleve/v2"
	bhttp "github.com/blevesearch/bleve/v2/http"
	"github.com/gorilla/mux"
)

//go:embed index.html
var indexHtml []byte

type Handler struct {
	indices map[string]bleve.Index
	alias   bleve.IndexAlias
}

func NewHandler(indices ...bleve.Index) *Handler {
	m := make(map[string]bleve.Index)

	for _, ix := range indices {
		m[ix.Name()] = ix
		bhttp.RegisterIndexName(ix.Name(), ix)
	}

	return &Handler{
		m,
		bleve.NewIndexAlias(indices...),
	}
}

func (h *Handler) indexHTMLHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(indexHtml)
}

func New(router *mux.Router, indices ...bleve.Index) *Handler {
	h := NewHandler(indices...)

	router.HandleFunc("/", h.indexHTMLHandler)
	router.HandleFunc("/q", h.searchQueryHandler)
	router.HandleFunc("/d/{index}/{id}", h.documentLookupHandler)
	router.HandleFunc("/raw/{index}/{id}", rawDocumentHandler)
	return h
}
