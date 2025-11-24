package core

import (
	"bytes"
	"encoding/xml"
	"net/http"
)

func WriteXML(w http.ResponseWriter, status int, v any) {
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)

	if err := enc.Encode(v); err != nil {
		http.Error(w, "Failed to encode XML", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
	w.Write(buf.Bytes())
}
