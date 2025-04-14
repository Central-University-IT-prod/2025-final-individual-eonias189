package middlewares

import "net/http"

type Middleware func(next http.Handler) http.Handler

func Apply(handler http.Handler, mwares ...Middleware) http.Handler {
	for i := len(mwares) - 1; i >= 0; i-- {
		handler = mwares[i](handler)
	}
	return handler
}
