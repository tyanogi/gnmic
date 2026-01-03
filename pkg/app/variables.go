// © 2026 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package app

import "sync"

type VariableStore struct {
	m    sync.RWMutex
	vars map[string][]string
}

func NewVariableStore() *VariableStore {
	return &VariableStore{
		vars: make(map[string][]string),
	}
}

func (s *VariableStore) Set(category string, values []string) {
	s.m.Lock()
	defer s.m.Unlock()
	s.vars[category] = values
}

func (s *VariableStore) Get(category string) []string {
	s.m.RLock()
	defer s.m.RUnlock()
	if v, ok := s.vars[category]; ok {
		return v
	}
	return nil
}
