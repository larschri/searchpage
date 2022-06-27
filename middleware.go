package pestotrap

import (
	_ "embed"
	"net/http"
	"strconv"
)

const (
	hxNew    = 0
	hxActive = 1
	hxSkip   = 2
)

type hxWrapper struct {
	http.ResponseWriter
	status int
	cfg    *Config
}

func (x *hxWrapper) WriteHeader(code int) {
	defer x.ResponseWriter.WriteHeader(code)

	if code != http.StatusOK {
		x.status = hxSkip
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

	lns[0] = strconv.Itoa(ln + len(x.cfg.HeaderHTML) + len(x.cfg.FooterHTML))
}

func (x *hxWrapper) Write(b []byte) (int, error) {
	if x.status == hxNew {
		if _, err := x.ResponseWriter.Write(x.cfg.HeaderHTML); err != nil {
			x.status = hxSkip
			return 0, err
		}
		x.status = hxActive
	}
	return x.ResponseWriter.Write(b)
}

func (x *hxWrapper) Close() {
	if x.status == hxActive {
		x.ResponseWriter.Write(x.cfg.FooterHTML)
	}
}

func (h *Handler) hxRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.Header["Hx-Request"]) == 0 {
			x := hxWrapper{w, hxNew, &h.Config}
			defer x.Close()
			w = &x
		}

		next.ServeHTTP(w, r)
	})
}
