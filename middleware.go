package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

// Cascade transforms a series of middlewares into a single http.Handler which,
// when invoked, will rescursively call the next in the chain, ending at a maximum
// depth, with finalHandler being called.
func Cascade(finalHandler http.Handler, middlewares ...Middleware) http.Handler {
	var curHandler http.Handler
	for index, _ := range middlewares {
		curIndex := len(middlewares) - index - 1
		if index == 0 {
			// Last middleware in the chain invokes finalHandler
			curHandler = middlewares[curIndex](finalHandler)
		} else {
			curHandler = middlewares[curIndex](curHandler)
		}
	}

	return curHandler
}
