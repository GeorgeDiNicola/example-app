package debug

import "net/http"

// intentional error
func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "intentional debug error", http.StatusInternalServerError)
}
