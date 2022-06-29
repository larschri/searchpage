package pestotrap

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search"
)

//go:embed htmx/*
var htmxDir embed.FS

var DefaultMatch = template.Must(template.ParseFS(htmxDir, "htmx/match.htmx"))

type Config struct {
	FormHTML []byte

	HeaderHTML []byte

	FooterHTML []byte

	RenderMatches func(w io.Writer, matches []*search.DocumentMatch)

	RenderNextPageLink func(w io.Writer, nextPage string)

	PageSize int

	Request func(r *http.Request) *bleve.SearchRequest
}

func (c *Config) initFrom(src *Config) {
	if src.FormHTML != nil {
		c.FormHTML = src.FormHTML
	}

	if src.HeaderHTML != nil {
		c.HeaderHTML = src.HeaderHTML
	}

	if src.FooterHTML != nil {
		c.FooterHTML = src.FooterHTML
	}

	if src.RenderMatches != nil {
		c.RenderMatches = src.RenderMatches
	}

	if src.RenderNextPageLink != nil {
		c.RenderNextPageLink = src.RenderNextPageLink
	}

	if src.PageSize > 0 {
		c.PageSize = src.PageSize
	}

	if src.Request != nil {
		c.Request = src.Request
	}
}

var DefaultConfig = Config{
	PageSize: 30,

	RenderMatches: func(w io.Writer, matches []*search.DocumentMatch) {
		for _, m := range matches {
			DefaultMatch.Execute(w, map[string]interface{}{
				"Name":     m.Fields["metadata.name"],
				"Type":     m.Fields["kind"],
				"Taxonomy": m.Index,
				"Url":      "d/" + m.Index + "/" + m.ID,
			})
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
	DefaultConfig.FormHTML, _ = htmxDir.ReadFile("htmx/form.htmx")
	DefaultConfig.HeaderHTML, _ = htmxDir.ReadFile("htmx/head.htmx")
	DefaultConfig.FooterHTML, _ = htmxDir.ReadFile("htmx/foot.htmx")
}
