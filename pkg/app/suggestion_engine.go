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
	"strings"

	"github.com/openconfig/gnmi/proto/gnmi"
	gnmicpath "github.com/openconfig/gnmic/pkg/api/path"
	"github.com/openconfig/gnmic/pkg/api/target"
)

// SuggestionEngine defines the interface for prefetching suggestion values
// from targets and managing them in a cache.
// It uses gNMI Get requests to fetch candidate values for YANG list keys.
type SuggestionEngine interface {
	// Start begins the asynchronous prefetching process for the given items.
	// Each item is processed in a separate goroutine.
	Start(ctx context.Context, items []SuggestionItem)
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

// Start begins the asynchronous prefetching process for the given items.
// It iterates over the provided SuggestionItem slice and spawns a goroutine for each item
// to perform the gNMI Get request without blocking the main thread.
func (s *suggestionEngineImpl) Start(ctx context.Context, items []SuggestionItem) {
	DebugLog("[SuggestionEngine] Starting for %d items", len(items))
	for _, item := range items {
		go func(si SuggestionItem) {
			if err := s.prefetchList(ctx, si); err != nil {
				// We log errors but don't stop the engine.
				DebugLog("[SuggestionEngine] Error for path %s: %v", si.Path, err)
			}
		}(item)
	}
}

// prefetchList sends a gNMI Get request for a single item and updates the cache.
// It parses the StatePath, executes the Get request, extracts the values from the response,
// and stores them in the cache using a composite key (SchemaPath::KeyName).
func (s *suggestionEngineImpl) prefetchList(ctx context.Context, item SuggestionItem) error {
	gnmiPath, err := gnmicpath.ParsePath(item.StatePath)
	if err != nil {
		return fmt.Errorf("failed to parse path %s: %v", item.StatePath, err)
	}

	req := &gnmi.GetRequest{
		Path: []*gnmi.Path{gnmiPath},
	}

	DebugLog("[SuggestionEngine] Fetching: %s (StatePath: %s)", item.Path, item.StatePath)
	resp, err := s.target.Get(ctx, req)
	if err != nil {
		return fmt.Errorf("gNMI Get failed for %s: %v", item.StatePath, err)
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
			sVal := typedValueToString(val)

			if sVal != "" && !seen[sVal] {
				candidates = append(candidates, sVal)
				seen[sVal] = true
			}
		}
	}

	if len(candidates) > 0 {
		// Use a composite key to distinguish between multiple keys of the same list
		cacheKey := fmt.Sprintf("%s::%s", item.Path, item.Key)
		DebugLog("[SuggestionEngine] Caching %d values for %s", len(candidates), cacheKey)
		s.cache.Add(cacheKey, candidates)
	}

	return nil
}

func typedValueToString(tv *gnmi.TypedValue) string {
	if tv == nil {
		return ""
	}
	var s string
	switch v := tv.Value.(type) {
	case *gnmi.TypedValue_StringVal:
		s = v.StringVal
	case *gnmi.TypedValue_AsciiVal:
		s = v.AsciiVal
	case *gnmi.TypedValue_IntVal:
		s = fmt.Sprintf("%d", v.IntVal)
	case *gnmi.TypedValue_UintVal:
		s = fmt.Sprintf("%d", v.UintVal)
	case *gnmi.TypedValue_BoolVal:
		s = fmt.Sprintf("%t", v.BoolVal)
	case *gnmi.TypedValue_BytesVal:
		s = string(v.BytesVal)
	case *gnmi.TypedValue_FloatVal:
		s = fmt.Sprintf("%f", v.FloatVal)
	case *gnmi.TypedValue_DecimalVal:
		s = fmt.Sprintf("%d.%d", v.DecimalVal.Digits, v.DecimalVal.Precision)
	case *gnmi.TypedValue_JsonVal:
		s = string(v.JsonVal)
	case *gnmi.TypedValue_JsonIetfVal:
		s = string(v.JsonIetfVal)
	case *gnmi.TypedValue_LeaflistVal:
		// keys shouldn't be leaflists, but for completeness
		vals := make([]string, 0, len(v.LeaflistVal.Element))
		for _, elem := range v.LeaflistVal.Element {
			vals = append(vals, typedValueToString(elem))
		}
		s = fmt.Sprintf("[%s]", strings.Join(vals, ", "))
	default:
		s = fmt.Sprintf("%v", v)
	}
	return strings.Trim(s, "\"")
}