package routes

import (
	"fmt"
	"net/http"
	"os"
)

func RedirectToWebUiHandler(w http.ResponseWriter, r *http.Request) {

	redirectUrl := fmt.Sprintf("%s?host=%s&port=%s", os.Getenv("LITWEBUI"), r.URL.Query().Get("host"), r.URL.Query().Get("port"))

	http.Redirect(w, r, redirectUrl, 301)
}
