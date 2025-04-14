package handlers

import "net/http"

func isFromHtmx(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
