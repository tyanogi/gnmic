// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"fmt"

	"github.com/openconfig/gnmi/proto/gnmi"
	gnmicpath "github.com/openconfig/gnmic/pkg/api/path"
	"github.com/openconfig/gnmic/pkg/api/target"
)

// SuggestionEngine defines the interface for prefetching suggestion values
// from targets and managing them in a cache.
// It uses gNMI Get requests to fetch candidate values for YANG list keys.
type SuggestionEngine interface {
	// Start begins the asynchronous prefetching process for the given lists.
	// Each list is processed in a separate goroutine.
	Start(ctx context.Context, lists []ListInfo)
}

// suggestionEngineImpl is the implementation of SuggestionEngine.
// It holds a reference to a gNMI target and a suggestion cache.
type suggestionEngineImpl struct {
	target *target.Target
	cache  SuggestionCache
}

// NewSuggestionEngine creates a new instance of SuggestionEngine with the provided target and cache.
func NewSuggestionEngine(t *target.Target, c SuggestionCache) SuggestionEngine {
	return &suggestionEngineImpl{
		target: t,
		cache:  c,
	}
}

// Start begins the asynchronous prefetching process for the given lists.
// It iterates over the provided ListInfo slice and spawns a goroutine for each list
// to perform the gNMI Get request without blocking the main thread.
func (s *suggestionEngineImpl) Start(ctx context.Context, lists []ListInfo) {
	for _, list := range lists {
		go func(li ListInfo) {
			err := s.prefetchList(ctx, li)
			if err != nil {
				// We log errors but don't stop the engine as per requirements.
				// In a production environment, this might use a proper logger.
				fmt.Printf("Suggestion engine error for path %s: %v\n", li.Path, err)
			}
		}(list)
	}
}

// prefetchList sends a gNMI Get request for a single list and updates the cache.
// It parses the StatePath, executes the Get request, extracts the values from the response,
// and stores them in the cache using a composite key (SchemaPath::KeyName).
func (s *suggestionEngineImpl) prefetchList(ctx context.Context, list ListInfo) error {
	gnmiPath, err := gnmicpath.ParsePath(list.StatePath)
	if err != nil {
		return fmt.Errorf("failed to parse path %s: %v", list.StatePath, err)
	}

	req := &gnmi.GetRequest{
		Path: []*gnmi.Path{gnmiPath},
	}

	resp, err := s.target.Get(ctx, req)
	if err != nil {
		return fmt.Errorf("gNMI Get failed for %s: %v", list.StatePath, err)
	}

	candidates := []string{}
	seen := make(map[string]bool)

	for _, notif := range resp.GetNotification() {
		for _, update := range notif.GetUpdate() {
			val := update.GetVal()
			if val == nil {
				continue
			}
			// Extract value as string
			var sVal string
			switch v := val.Value.(type) {
			case *gnmi.TypedValue_StringVal:
				sVal = v.StringVal
			case *gnmi.TypedValue_AsciiVal:
				sVal = v.AsciiVal
			default:
				// Fallback to string representation if possible, or skip
				sVal = fmt.Sprintf("%v", val)
			}

			if sVal != "" && !seen[sVal] {
				candidates = append(candidates, sVal)
				seen[sVal] = true
			}
		}
	}

	if len(candidates) > 0 {
		// Use a composite key to distinguish between multiple keys of the same list
		cacheKey := fmt.Sprintf("%s::%s", list.Path, list.Key)
		s.cache.Add(cacheKey, candidates)
	}

	return nil
}
