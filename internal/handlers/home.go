package handlers

import (
	"net/http"
	"zapbyte/templates/pages"
)

func HomeHandler(fromHtmx bool, w http.ResponseWriter, r *http.Request) {
	pages.Home(fromHtmx).Render(r.Context(), w)
}
