package security

import (
	"testing"
)

func TestTruncateSensitiveHeader(t *testing.T) {
	const original = "Data Tags=dogs;cats;bats, Sensitive=00915d97a7d46e3fbbe81580eb632d69"
	const expected = "Data Tags=dogs;cats;bats, Sensitive=00915d9..."

	truncated := TruncParamHeader(original, "sensitive")

	if truncated != expected {
		t.Error("Truncated value mismatch")
	}
}
