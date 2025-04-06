package handlers

import (
	"net/http"
	"zapbyte/templates/pages"
)

func TutorialHandler(fromHtmx bool, w http.ResponseWriter, r *http.Request) {
	pages.Tutorial(fromHtmx).Render(r.Context(), w)
}
