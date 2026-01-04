// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) if and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"

	"github.com/openconfig/gnmic/pkg/api/target"
)

// SuggestionEngine defines the interface for prefetching suggestion values
// from targets and managing them in a cache.
type SuggestionEngine interface {
	// Start begins the asynchronous prefetching process for the given lists.
	Start(ctx context.Context, lists []ListInfo)
}

// suggestionEngineImpl is the implementation of SuggestionEngine.
type suggestionEngineImpl struct {
	target *target.Target
	cache  SuggestionCache
}

// NewSuggestionEngine creates a new instance of SuggestionEngine.
func NewSuggestionEngine(t *target.Target, c SuggestionCache) SuggestionEngine {
	return &suggestionEngineImpl{
		target: t,
		cache:  c,
	}
}
