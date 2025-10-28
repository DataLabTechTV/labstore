package server

import (
	"encoding/xml"
	"net/http"
)

type S3Error struct {
	XMLName   xml.Name `xml:"Error"`
	Code      string
	Message   string
	RequestId string
	HostId    string
}

func WriteS3Error(w http.ResponseWriter, code, message string, statusCode int) {
	w.WriteHeader(statusCode)
	errResp := S3Error{
		Code:    code,
		Message: message,
	}
	xml.NewEncoder(w).Encode(errResp)
}
