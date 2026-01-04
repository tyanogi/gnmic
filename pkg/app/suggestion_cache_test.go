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
	"sync"
	"testing"
)

func TestSuggestionCache_AddGet(t *testing.T) {
	cache := NewSuggestionCache()

	path := "/interfaces/interface/name"
	candidates := []string{"eth1", "eth2"}

	cache.Add(path, candidates)

	got := cache.Get(path)
	if len(got) != len(candidates) {
		t.Errorf("Expected %d candidates, got %d", len(candidates), len(got))
	}
	for i, v := range candidates {
		if got[i] != v {
			t.Errorf("Expected candidate %s at index %d, got %s", v, i, got[i])
		}
	}

	// Test non-existent path
	if got := cache.Get("/non/existent"); got != nil {
		t.Errorf("Expected nil for non-existent path, got %v", got)
	}
}

func TestSuggestionCache_Clear(t *testing.T) {
	cache := NewSuggestionCache()
	cache.Add("/path1", []string{"a"})
	cache.Add("/path2", []string{"b"})

	cache.Clear()

	if got := cache.Get("/path1"); got != nil {
		t.Errorf("Expected nil after clear, got %v", got)
	}
	if got := cache.Get("/path2"); got != nil {
		t.Errorf("Expected nil after clear, got %v", got)
	}
}

func TestSuggestionCache_Concurrency(t *testing.T) {
	cache := NewSuggestionCache()
	var wg sync.WaitGroup
	numRoutines := 100
	numOperations := 1000

	// Concurrent writes
	wg.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				path := fmt.Sprintf("/path/%d", id)
				cache.Add(path, []string{"val"})
			}
		}(i)
	}

	// Concurrent reads
	wg.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				path := fmt.Sprintf("/path/%d", id)
				cache.Get(path)
			}
		}(i)
	}

	wg.Wait()
}
