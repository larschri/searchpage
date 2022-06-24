package pestotrap

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) documentLookupHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	index, ok := h.indices[vars["index"]]
	if !ok {
		http.Error(w, "Missing index", http.StatusBadRequest)
		return
	}
	if ixf, ok := index.(interface {
		ServeHTTPx(w http.ResponseWriter, r *http.Request)
	}); ok {
		ixf.ServeHTTPx(w, r)
		return
	}

	rawDocumentHandler(w, r)
}
