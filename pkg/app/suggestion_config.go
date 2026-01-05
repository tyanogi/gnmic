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
	"gopkg.in/yaml.v2"
)

// SuggestionConfig represents the configuration for the suggestion engine.
type SuggestionConfig struct {
	Suggestions []SuggestionItem `yaml:"suggestions"`
}

// SuggestionItem represents a single suggestion configuration item.
type SuggestionItem struct {
	Path      string `yaml:"path"`
	Key       string `yaml:"key"`
	StatePath string `yaml:"state-path"`
}

// ParseSuggestionConfig parses the suggestion configuration from a byte slice.
func ParseSuggestionConfig(data []byte) (*SuggestionConfig, error) {
	var config SuggestionConfig
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// LoadSuggestionConfigFromFile reads the configuration from the specified file path and parses it.
func LoadSuggestionConfigFromFile(path string) (*SuggestionConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseSuggestionConfig(data)
}