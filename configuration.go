package pestotrap

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search"
)

type Config struct {
	IndexHTML []byte

	RenderPage func(w io.Writer, matches []*search.DocumentMatch, nextPage *url.URL)

	PageSize int

	Request func(r *http.Request) *bleve.SearchRequest
}

var DefaultConfig = Config{
	PageSize: 30,

	RenderPage: func(w io.Writer, matches []*search.DocumentMatch, nextPage *url.URL) {
		for _, m := range matches {
			renderK8sMatch(w, m)
		}

		if nextPage != nil {
			fmt.Fprintf(w, `<div hx-get="%s" hx-trigger="revealed"/>`, nextPage.String())
		}
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
	DefaultConfig.IndexHTML = indexHtml
}

var matchTemplate = template.Must(template.New("t").Parse(`
<details>
	  <summary>
	    {{ .Name }} <small><i>{{ .Type }}</i></small>
	    <br>
	    <small><i>{{ .Location }}</i></small>
	  </summary>
	  <span hx-get="{{ .Url }}"
		hx-trigger="toggle once from:closest details">
	  </span>
</details>`))

func renderMatch(w io.Writer, m *search.DocumentMatch) {
	matchTemplate.Execute(w, map[string]string{
		"Name":     m.ID,
		"Type":     "",
		"Location": m.Index,
		"Url":      "d/" + m.Index + "/" + m.ID,
	})
}

func renderK8sMatch(w io.Writer, m *search.DocumentMatch) {
	matchTemplate.Execute(w, map[string]interface{}{
		"Name":     m.Fields["metadata.name"],
		"Type":     m.Fields["kind"],
		"Location": m.Index,
		"Url":      "d/" + m.Index + "/" + m.ID,
	})
}
