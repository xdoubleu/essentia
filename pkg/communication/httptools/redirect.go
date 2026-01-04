package httptools

import (
	"fmt"
	"net/http"
)

// RedirectWithError redirects to the provided url with an error in the query.
func RedirectWithError(w http.ResponseWriter, r *http.Request, url string, err error) {
	http.Redirect(
		w,
		r,
		fmt.Sprintf("%s?error=%s", url, err.Error()),
		http.StatusSeeOther,
	)
}
