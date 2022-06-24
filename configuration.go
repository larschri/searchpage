package pestotrap

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	_ "embed"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search"
)

//go:embed htmx/form.htmx
var formHTML []byte

//go:embed htmx/match.htmx
var matchHTML string

var matchTemplate = template.Must(template.New("t").Parse(matchHTML))

type Config struct {
	FormHTML []byte

	RenderResults func(w io.Writer, matches []*search.DocumentMatch)

	RenderNextPageLink func(w io.Writer, nextPage string)

	PageSize int

	Request func(r *http.Request) *bleve.SearchRequest
}

var DefaultConfig = Config{
	PageSize: 30,

	RenderResults: func(w io.Writer, matches []*search.DocumentMatch) {
		for _, m := range matches {
			renderK8sMatch(w, m)
		}
	},

	RenderNextPageLink: func(w io.Writer, nextPage string) {
		fmt.Fprintf(w, `<div hx-get="%s" hx-trigger="revealed"/>`, nextPage)
	},

	Request: func(r *http.Request) *bleve.SearchRequest {
		srch := ""
		if len(r.Form["search"]) > 0 {
			srch = r.Form["search"][0]
		}

		query := bleve.NewQueryStringQuery(srch)
		request := bleve.NewSearchRequest(query)
		request.Fields = []string{"kind", "metadata.name", "metadata.namespace"}
		return request
	},
}

func init() {
	DefaultConfig.FormHTML = formHTML
}

func renderMatch(w io.Writer, m *search.DocumentMatch) {
	matchTemplate.Execute(w, map[string]string{
		"Name":     m.ID,
		"Type":     "",
		"Taxonomy": m.Index,
		"Url":      "d/" + m.Index + "/" + m.ID,
	})
}

func renderK8sMatch(w io.Writer, m *search.DocumentMatch) {
	matchTemplate.Execute(w, map[string]interface{}{
		"Name":     m.Fields["metadata.name"],
		"Type":     m.Fields["kind"],
		"Taxonomy": m.Index,
		"Url":      "d/" + m.Index + "/" + m.ID,
	})
}
