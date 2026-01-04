package app

import (
	"testing"

	"github.com/openconfig/goyang/pkg/yang"
)

func TestSchemaWalkerTypes(t *testing.T) {
	// This test is designed to verify the existence of the types/interfaces required by the spec.
	// It will fail compilation if ListInfo or Walk are not defined.

	// Verify ListInfo struct exists and has the required fields
	_ = ListInfo{
		Path:      "/test/path",
		Key:       "name",
		StatePath: "/test/path/state/name",
	}

	// Verify Walk function signature
	// It should take a *yang.Entry and return a slice of ListInfo and an error.
	var _ func(*yang.Entry) ([]ListInfo, error) = Walk
}
