package pestotrap

import (
	"net/http"
	"strconv"
)

func (h *Handler) searchQueryHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	request := h.Config.Request(r)
	request.Size = h.Config.PageSize

	if len(r.Form["offset"]) > 0 {
		request.From, _ = strconv.Atoi(r.Form["offset"][0])
	}

	result, err := h.alias.Search(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Config.RenderMatches(w, result.Hits)

	nextOffset := uint64(request.From + h.Config.PageSize)
	if nextOffset >= result.Total {
		return
	}

	q := r.URL.Query()
	q.Set("offset", strconv.Itoa(request.From+h.Config.PageSize))
	h.Config.RenderNextPageLink(w, "q?"+q.Encode())
}
