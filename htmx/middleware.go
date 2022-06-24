package htmx

import (
	_ "embed"
	"net/http"
	"strconv"
)

const (
	xNew    = 0
	xActive = 1
	xSkip   = 2
)

var (
	//go:embed head.htmx
	Header []byte

	//go:embed foot.htmx
	Footer []byte
)

type xWrapper struct {
	http.ResponseWriter
	status int
}

func (x *xWrapper) WriteHeader(code int) {
	defer x.ResponseWriter.WriteHeader(code)

	if code != 200 {
		x.status = xSkip
		return
	}

	lns, ok := x.Header()["Content-Length"]
	if !ok {
		return
	}

	ln, err := strconv.Atoi(lns[0])
	if err != nil {
		return
	}

	lns[0] = strconv.Itoa(ln + len(Header) + len(Footer))
}

func (x *xWrapper) Write(b []byte) (int, error) {
	if x.status == xNew {
		if _, err := x.ResponseWriter.Write(Header); err != nil {
			x.status = xSkip
			return 0, err
		}
		x.status = xActive
	}
	return x.ResponseWriter.Write(b)
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.Header["Hx-Request"]) != 0 {
			next.ServeHTTP(w, r)
			return
		}

		x := xWrapper{w, xNew}
		next.ServeHTTP(&x, r)

		if x.status == xActive {
			w.Write(Footer)
		}
	})
}
