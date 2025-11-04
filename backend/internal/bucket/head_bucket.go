package bucket

import (
	"net/http"
)

func HeadBucketHandler(w http.ResponseWriter, r *http.Request) {
	// *NOTE: This will likely share code with GET due to using the same headers.
	// TODO: organize shared code somewhere

	w.WriteHeader(http.StatusOK)
}
