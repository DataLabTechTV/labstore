package core

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/IllumiKnowLabs/labstore/backend/internal/helper"
	"github.com/google/uuid"
)

func NewRequestID() string {
	return uuid.NewString()
}

func ReadXML(w http.ResponseWriter, r *http.Request, dst any) error {
	decoder := xml.NewDecoder(r.Body)
	defer helper.CloseWithErr(r.Body, nil)

	if err := decoder.Decode(&dst); err != nil {
		return fmt.Errorf("failed to decode XML: %w", err)
	}

	return nil
}
