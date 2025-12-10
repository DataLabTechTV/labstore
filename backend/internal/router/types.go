package router

import (
	"net/http"
	"sync/atomic"
)

type ServerDescriptor struct {
	Name    string
	Server  *http.Server
	Healthy atomic.Bool
}
