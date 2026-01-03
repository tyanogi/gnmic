// © 2026 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"testing"

	"github.com/openconfig/gnmic/pkg/app"
)

func TestCheckValueSuggestions(t *testing.T) {
	// Initialize Mock Store
	gApp = app.New() // Re-initialize gApp for testing
	store := app.NewVariableStore()
	store.Set("interface", []string{"eth0", "eth1"})
	store.Set("netinst", []string{"default", "mgmt"})
	gApp.PromptVariableStore = store

	tests := []struct {
		name     string
		line     string
		want     int // number of suggestions
		category string
	}{
		{
			name: "No match",
			line: "get --path /",
			want: 0,
		},
		{
			name:     "Interface name start",
			line:     "get --path /interfaces/interface[name=",
			want:     2,
			category: "interface",
		},
		{
			name:     "Interface name with quote",
			line:     "get --path /interfaces/interface[name='",
			want:     2,
			category: "interface",
		},
		{
			name:     "Netinst name start",
			line:     "get --path /network-instances/network-instance[name=",
			want:     2,
			category: "netinst",
		},
		{
			name: "Closed bracket",
			line: "get --path /interfaces/interface[name=eth0]",
			want: 0,
		},
		{
			name: "Wrong path",
			line: "get --path /system[name=",
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkValueSuggestions(tt.line)
			if len(got) != tt.want {
				t.Errorf("checkValueSuggestions() got %v suggestions, want %v", len(got), tt.want)
			}
			if tt.want > 0 {
				if got[0].Description != "Discovered "+tt.category {
					t.Errorf("checkValueSuggestions() Description = %v, want Discovered %v", got[0].Description, tt.category)
				}
			}
		})
	}
}

func TestGetStoreSuggestionsNilStore(t *testing.T) {
	// Ensure gApp.PromptVariableStore is nil
	gApp = app.New()
	gApp.PromptVariableStore = nil

	got := checkValueSuggestions("get --path /interfaces/interface[name=")
	if got != nil {
		t.Errorf("expected nil suggestions for nil store, got %v", got)
	}
}

func TestGetStoreSuggestionsEmptyStore(t *testing.T) {
	// Initialize gApp with empty store
	gApp = app.New()
	store := app.NewVariableStore()
	// No values set
	gApp.PromptVariableStore = store

	got := checkValueSuggestions("get --path /interfaces/interface[name=")
	if got != nil {
		t.Errorf("expected nil suggestions for empty store, got %v", got)
	}
}
