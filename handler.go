package searchpage

import (
	"net/http"

	"github.com/blevesearch/bleve/v2"
	bhttp "github.com/blevesearch/bleve/v2/http"
	"github.com/gorilla/mux"
)

type Handler struct {
	http.Handler
	Config
	indices map[string]bleve.Index
	alias   bleve.IndexAlias
}

func (h *Handler) indexHTMLHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(h.Config.FormHTML)
}

func New(cfg *Config, indices ...bleve.Index) *Handler {
	m := make(map[string]bleve.Index)

	for _, ix := range indices {
		m[ix.Name()] = ix
		bhttp.RegisterIndexName(ix.Name(), ix)
	}

	r := mux.NewRouter()
	h := Handler{
		r,
		DefaultConfig,
		m,
		bleve.NewIndexAlias(indices...),
	}

	if cfg != nil {
		h.initFrom(cfg)
	}

	r.Use(h.hxRequestMiddleware)
	r.HandleFunc("/", h.indexHTMLHandler)
	r.HandleFunc("/q", h.searchQueryHandler)
	r.HandleFunc("/d/{index}/{id}", docHandler)
	return &h
}
