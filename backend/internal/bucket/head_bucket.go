package bucket

import (
	"net/http"

	"github.com/DataLabTechTV/labstore/backend/internal/core"
	"github.com/DataLabTechTV/labstore/backend/pkg/iam"
)

func HeadBucket(w http.ResponseWriter, _ *http.Request, accessKey string) {
	// *NOTE: This will likely share code with GET due to using the same headers.
	// TODO: organize shared code somewhere

	if !iam.CheckPolicy(accessKey, "", "ListBuckets") {
		core.WriteS3Error(w, "AccessDenied", "Access Denied", 403)
		return
	}

	w.WriteHeader(http.StatusOK)
}
