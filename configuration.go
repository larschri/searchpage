package searchpage

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search"
)

var (
	//DefaultConfig has the settings that will be used unless supplied from
	//the client code.
	DefaultConfig = Config{
		FormHTML:   mustBytes(htmxDir.ReadFile("htmx/form.htmx")),
		HeaderHTML: mustBytes(htmxDir.ReadFile("htmx/head.htmx")),
		FooterHTML: mustBytes(htmxDir.ReadFile("htmx/foot.htmx")),
		PageSize:   30,

		RenderMatches: func(w io.Writer, matches []*search.DocumentMatch) {
			for _, m := range matches {
				DefaultMatch.Execute(w, map[string]interface{}{
					"Name":     m.ID,
					"Type":     "",
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
			return request
		},
	}

	//Defaultmatch is the template used to render a single match in
	//DefaultConfig.RenderMatches. It is public so it can be reused by custom
	//implementations.
	DefaultMatch = template.Must(template.ParseFS(htmxDir, "htmx/match.htmx"))

	//go:embed htmx/*
	htmxDir embed.FS
)

//Config contains customizable fields
type Config struct {
	// HeaderHTML is the html head and the first part of the html body. It
	// is prepended in every response unless the hx-request header is
	// supplied.
	HeaderHTML []byte

	// FooterHTML is the last part of the html body. It is appended in
	// every response where HeaderHTML has been appended.
	FooterHTML []byte

	// FormHTML is the html search form
	FormHTML []byte

	// RenderMatches is the function that render the matches. This part of
	// the page is rendered after the FormHTML.
	RenderMatches func(w io.Writer, matches []*search.DocumentMatch)

	// RenderNextPageLink renders a link to the next page after the
	// matches. It is used for inifinite scrolling.
	RenderNextPageLink func(w io.Writer, nextPage string)

	// PageSize is the number of matches per page.
	PageSize int

	//Request creates a bleve.SearchRequest based on the given
	//http.Request.
	Request func(r *http.Request) *bleve.SearchRequest
}

func mustBytes(b []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return b
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
