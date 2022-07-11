package searchpage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"net/http"

	bhttp "github.com/blevesearch/bleve/v2/http"
)

// jsonBuffer will buffer if content type is json, so we can apply formatting.
type jsonBuffer struct {
	http.ResponseWriter
	buf bytes.Buffer
}

func (m *jsonBuffer) isJSON() bool {
	return m.Header().Get("Content-Type") == "application/json"
}

func (m *jsonBuffer) Write(b []byte) (int, error) {
	if m.isJSON() {
		return m.buf.Write(b)
	}

	return m.ResponseWriter.Write(b)
}

func (m *jsonBuffer) pretty() string {
	if m.buf.Len() == 0 {
		return ""
	}

	var out bytes.Buffer
	json.Indent(&out, m.buf.Bytes(), "", "  ")
	return out.String()
}

var bleveDocHandler = bhttp.DocGetHandler{
	IndexNameLookup: func(r *http.Request) string {
		return r.URL.Query().Get("index")
	},

	DocIDLookup: func(r *http.Request) string {
		return r.URL.Query().Get("doc")
	},
}

func docHandler(w http.ResponseWriter, r *http.Request) {
	buf := &jsonBuffer{w, bytes.Buffer{}}

	bleveDocHandler.ServeHTTP(buf, r)

	if !buf.isJSON() {
		return
	}

	if r.Header.Get("Hx-Request") == "true" {
		fmt.Fprintf(w, "<pre>%s</pre>", html.EscapeString(buf.pretty()))
		return
	}

	w.Write([]byte(buf.pretty()))
}
