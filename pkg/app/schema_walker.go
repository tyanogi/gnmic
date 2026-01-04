// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"github.com/openconfig/goyang/pkg/yang"
)

// ListInfo represents the information extracted from a YANG list node
// necessary for dynamic prompt completion.
type ListInfo struct {
	// Path is the absolute path to the list node (e.g., "/interfaces/interface").
	Path string
	// Key is the name of the list key (e.g., "name").
	Key string
	// StatePath is the path used to query the state from the device (e.g., "/interfaces/interface/state/name").
	StatePath string
}

// Walk traverses the YANG schema tree starting from the given entry
// and extracts information about all list nodes.
func Walk(entry *yang.Entry) ([]ListInfo, error) {
	// Implementation will be added in subsequent phases.
	return []ListInfo{}, nil
}
