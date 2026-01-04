// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"
	"strings"

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
	infos := []ListInfo{}

	var traverse func(e *yang.Entry)
	traverse = func(e *yang.Entry) {
		if e == nil {
			return
		}

		if e.ListAttr != nil {
			p := getPath(e)
			if e.Key != "" {
				keys := strings.Fields(e.Key)
				// Handle multiple keys by creating a ListInfo for each key
				for _, k := range keys {
					statePath := p + "/" + k
					// OpenConfig heuristic: check if there's a 'state' container
					// and if it contains a leaf with the same name as the key.
					if state, ok := e.Dir["state"]; ok && state.Kind == yang.DirectoryEntry {
						if _, ok := state.Dir[k]; ok {
							statePath = p + "/state/" + k
						}
					}

					li := ListInfo{
						Path:      p,
						Key:       k,
						StatePath: statePath,
					}
					infos = append(infos, li)
				}
			} else {
				fmt.Printf("List %s has no key\n", e.Name)
			}
		}

		for _, child := range e.Dir {
			traverse(child)
		}
	}

	traverse(entry)
	return infos, nil
}

func getPath(entry *yang.Entry) string {
	parts := []string{}
	// Stop when Parent is nil (usually the module)
	for e := entry; e != nil && e.Parent != nil; e = e.Parent {
		if e.IsChoice() || e.IsCase() {
			continue
		}
		parts = append(parts, e.Name)
	}
	// Reverse
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return "/" + strings.Join(parts, "/")
}
