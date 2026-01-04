// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package app

import "sync"

// SuggestionCache defines the interface for a thread-safe cache
// storing suggestions (candidate values) for YANG schema paths.
type SuggestionCache interface {
	// Add stores a list of candidates for a given path.
	Add(path string, candidates []string)
	// Get retrieves the list of candidates for a given path.
	// Returns nil if the path is not found.
	Get(path string) []string
	// Clear removes all entries from the cache.
	Clear()
}

// suggestionCacheImpl is the thread-safe implementation of SuggestionCache
// using sync.Map.
type suggestionCacheImpl struct {
	m sync.Map
}

// NewSuggestionCache creates a new instance of SuggestionCache.
func NewSuggestionCache() SuggestionCache {
	return &suggestionCacheImpl{}
}

// Add stores a list of candidates for a given path.
func (s *suggestionCacheImpl) Add(path string, candidates []string) {
	s.m.Store(path, candidates)
}

// Get retrieves the list of candidates for a given path.
func (s *suggestionCacheImpl) Get(path string) []string {
	if v, ok := s.m.Load(path); ok {
		if candidates, ok := v.([]string); ok {
			return candidates
		}
	}
	return nil
}

// Clear removes all entries from the cache.
func (s *suggestionCacheImpl) Clear() {
	s.m.Range(func(key, value interface{}) bool {
		s.m.Delete(key)
		return true
	})
}
