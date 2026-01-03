// © 2026 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"reflect"
	"sync"
	"testing"
)

func TestVariableStore(t *testing.T) {
	store := NewVariableStore()

	// Test Set and Get
	category := "interface"
	values := []string{"eth1", "eth2"}
	store.Set(category, values)

	got := store.Get(category)
	if !reflect.DeepEqual(got, values) {
		t.Errorf("expected %v, got %v", values, got)
	}

	// Test Get non-existent category
	got = store.Get("unknown")
	if got != nil {
		t.Errorf("expected nil for unknown category, got %v", got)
	}
}

func TestVariableStoreConcurrency(t *testing.T) {
	store := NewVariableStore()
	var wg sync.WaitGroup
	numRoutines := 100

	wg.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {
		go func(i int) {
			defer wg.Done()
			store.Set("cat", []string{"val"})
			_ = store.Get("cat")
		}(i)
	}
	wg.Wait()
}
