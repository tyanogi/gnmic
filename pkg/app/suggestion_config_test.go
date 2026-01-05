// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"os"
	"reflect"
	"testing"
)

func TestParseSuggestionConfig(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected *SuggestionConfig
		wantErr  bool
	}{
		{
			name: "valid config",
			input: []byte(`
suggestions:
  - path: "/interfaces/interface"
    key: "name"
    state-path: "/interfaces/interface[name=*]/state/name"
`),
			expected: &SuggestionConfig{
				Suggestions: []SuggestionItem{
					{
						Path:      "/interfaces/interface",
						Key:       "name",
						StatePath: "/interfaces/interface[name=*]/state/name",
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "empty config",
			input:   []byte(``),
			expected: &SuggestionConfig{},
			wantErr: false,
		},
		{
			name: "invalid yaml",
			input: []byte(`
suggestions:
  - path: "/interfaces/interface"
    key: "name"
     state-path: "/interfaces/interface[name=*]/state/name" # Indentation error
`),
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSuggestionConfig(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSuggestionConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParseSuggestionConfig() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLoadSuggestionConfigFromFile(t *testing.T) {
	// Create a temporary file
	content := []byte(`
suggestions:
  - path: "/interfaces/interface"
    key: "name"
    state-path: "/interfaces/interface[name=*]/state/name"
`)
	tmpfile, err := os.CreateTemp("", "suggestion_config_*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test loading
	config, err := LoadSuggestionConfigFromFile(tmpfile.Name())
	if err != nil {
		t.Errorf("LoadSuggestionConfigFromFile() error = %v", err)
	}
	expected := &SuggestionConfig{
		Suggestions: []SuggestionItem{
			{
				Path:      "/interfaces/interface",
				Key:       "name",
				StatePath: "/interfaces/interface[name=*]/state/name",
			},
		},
	}
	if !reflect.DeepEqual(config, expected) {
		t.Errorf("LoadSuggestionConfigFromFile() = %v, want %v", config, expected)
	}

	// Test non-existent file
	_, err = LoadSuggestionConfigFromFile("non_existent_file.yaml")
	if err == nil {
		t.Error("LoadSuggestionConfigFromFile() expected error for non-existent file, got nil")
	}
}