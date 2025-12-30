package requests

import (
	"github.com/IllumiKnowLabs/labstore/client"
)

func HandleList(path string) {
	client.List(path)
}
