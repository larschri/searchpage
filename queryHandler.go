package pestotrap

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"

	"github.com/blevesearch/bleve/v2/search"
)

func (h *Handler) writeMatch(w http.ResponseWriter, m *search.DocumentMatch) {
	ix := h.indices[m.Index]

	if ixf, ok := ix.(interface {
		RenderFull(m *search.DocumentMatch) template.HTML
	}); ok {
		fmt.Fprint(w, ixf.RenderFull(m))
		return
	}

	renderK8sMatch(w, m)
}

func (h *Handler) searchQueryHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	request := DefaultConfig.Request(r)
	pageSize := DefaultConfig.PageSize
	request.Size = pageSize + 1

	if len(r.Form["offset"]) > 0 {
		request.From, _ = strconv.Atoi(r.Form["offset"][0])
	}

	result, err := h.alias.Search(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var nextPage *url.URL
	if len(result.Hits) == request.Size {
		q := r.URL.Query()
		q.Set("offset", strconv.Itoa(request.From+DefaultConfig.PageSize))
		r.URL.RawQuery = q.Encode()
		nextPage = r.URL
		result.Hits = result.Hits[0:DefaultConfig.PageSize]
	}

	DefaultConfig.RenderPage(w, result.Hits, nextPage)
}
