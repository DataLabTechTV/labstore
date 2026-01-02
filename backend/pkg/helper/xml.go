package helper

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/http"
)

func WriteXMLResponse(w http.ResponseWriter, status int, v any) {
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)

	if err := enc.Encode(v); err != nil {
		http.Error(w, "Failed to encode XML", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")); err != nil {
		slog.Error("write xml", "err", err)
	}

	if _, err := w.Write(buf.Bytes()); err != nil {
		slog.Error("write xml", "err", err)
	}
}

func ReadXMLRequest(r *http.Request, dst any) error {
	decoder := xml.NewDecoder(r.Body)
	defer CloseWithErr(r.Body, nil)

	if err := decoder.Decode(&dst); err != nil {
		return fmt.Errorf("failed to decode XML: %w", err)
	}

	return nil
}
