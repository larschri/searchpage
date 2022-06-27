package pestotrap

import (
	"net/http"

	"github.com/blevesearch/bleve/v2"
	bhttp "github.com/blevesearch/bleve/v2/http"
	"github.com/gorilla/mux"
)

type Handler struct {
	*mux.Router
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

	h := Handler{
		mux.NewRouter(),
		DefaultConfig,
		m,
		bleve.NewIndexAlias(indices...),
	}

	if cfg != nil {
		h.initFrom(cfg)
	}

	h.Use(h.hxRequestMiddleware)
	h.HandleFunc("/", h.indexHTMLHandler)
	h.HandleFunc("/q", h.searchQueryHandler)
	h.HandleFunc("/d/{index}/{id}", docHandler)
	return &h
}
