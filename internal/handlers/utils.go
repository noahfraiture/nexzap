package handlers

import "net/http"

func fromHtmx(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
