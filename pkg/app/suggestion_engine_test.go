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
	"testing"
	"time"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/gnmic/pkg/api/target"
	"github.com/openconfig/gnmic/pkg/api/types"
	"google.golang.org/grpc"
)

// mockGNMIClient is a mock implementation of gnmi.GNMIClient
type mockGNMIClient struct {
	gnmi.GNMIClient
	getResponse *gnmi.GetResponse
	err         error
}

func (m *mockGNMIClient) Capabilities(ctx context.Context, in *gnmi.CapabilityRequest, opts ...grpc.CallOption) (*gnmi.CapabilityResponse, error) {
	return nil, nil
}

func (m *mockGNMIClient) Get(ctx context.Context, in *gnmi.GetRequest, opts ...grpc.CallOption) (*gnmi.GetResponse, error) {
	return m.getResponse, m.err
}

func (m *mockGNMIClient) Set(ctx context.Context, in *gnmi.SetRequest, opts ...grpc.CallOption) (*gnmi.SetResponse, error) {
	return nil, nil
}

func (m *mockGNMIClient) Subscribe(ctx context.Context, opts ...grpc.CallOption) (gnmi.GNMI_SubscribeClient, error) {
	return nil, nil
}

// mockSuggestionCache is a mock implementation of SuggestionCache
type mockSuggestionCache struct {
	data map[string][]string
}

func (m *mockSuggestionCache) Add(path string, candidates []string) {
	m.data[path] = candidates
}

func (m *mockSuggestionCache) Get(path string) []string {
	return m.data[path]
}

func (m *mockSuggestionCache) Clear() {
	m.data = make(map[string][]string)
}

func TestSuggestionEngine_PrefetchList(t *testing.T) {
	// Setup
	mockClient := &mockGNMIClient{
		getResponse: &gnmi.GetResponse{
			Notification: []*gnmi.Notification{
				{
					Update: []*gnmi.Update{
						{
							Val: &gnmi.TypedValue{
								Value: &gnmi.TypedValue_StringVal{StringVal: "eth1"},
							},
						},
						{
							Val: &gnmi.TypedValue{
								Value: &gnmi.TypedValue_StringVal{StringVal: "eth2"},
							},
						},
					},
				},
			},
		},
	}
	mockCache := &mockSuggestionCache{data: make(map[string][]string)}

	// Create engine with mock target (we need to inject the mock client)
	target := &target.Target{
		Client: mockClient,
		Config: &types.TargetConfig{},
	}
	engine := &suggestionEngineImpl{
		target: target,
		cache:  mockCache,
	}

	item := SuggestionItem{
		Path:      "/interfaces/interface",
		Key:       "name",
		StatePath: "/interfaces/interface[name=*]/state/name",
	}

	// Execute
	ctx := context.Background()
	err := engine.prefetchList(ctx, item)
	if err != nil {
		t.Fatalf("prefetchList failed: %v", err)
	}

	// Verify cache
	// Cache key should be "Path::Key"
	got := mockCache.Get("/interfaces/interface::name")
	if len(got) != 2 {
		t.Errorf("Expected 2 candidates in cache, got %d", len(got))
	}
	expected := map[string]bool{"eth1": true, "eth2": true}
	for _, v := range got {
		if !expected[v] {
			t.Errorf("Unexpected candidate in cache: %s", v)
		}
	}
}

func TestSuggestionEngine_PrefetchList_UintKey(t *testing.T) {
	// Setup mock with UintVal (typical for subinterface index)
	mockClient := &mockGNMIClient{
		getResponse: &gnmi.GetResponse{
			Notification: []*gnmi.Notification{
				{
					Update: []*gnmi.Update{
						{
							Val: &gnmi.TypedValue{
								Value: &gnmi.TypedValue_UintVal{UintVal: 0},
							},
						},
						{
							Val: &gnmi.TypedValue{
								Value: &gnmi.TypedValue_UintVal{UintVal: 100},
							},
						},
					},
				},
			},
		},
	}
	mockCache := &mockSuggestionCache{data: make(map[string][]string)}

	target := &target.Target{
		Client: mockClient,
		Config: &types.TargetConfig{},
	}
	engine := &suggestionEngineImpl{
		target: target,
		cache:  mockCache,
	}

	item := SuggestionItem{
		Path:      "/interfaces/interface/subinterfaces/subinterface",
		Key:       "index",
		StatePath: "/interfaces/interface/subinterfaces/subinterface[index=*]/state/index",
	}

	// Execute
	ctx := context.Background()
	err := engine.prefetchList(ctx, item)
	if err != nil {
		t.Fatalf("prefetchList failed: %v", err)
	}

	// Verify cache
	got := mockCache.Get("/interfaces/interface/subinterfaces/subinterface::index")
	if len(got) != 2 {
		t.Errorf("Expected 2 candidates in cache, got %d. Cache: %v", len(got), got)
	}
	expected := map[string]bool{"0": true, "100": true}
	for _, v := range got {
		if !expected[v] {
			t.Errorf("Unexpected candidate in cache: %s", v)
		}
	}
}

func TestSuggestionEngine_Start(t *testing.T) {
	// Setup mock cache and client
	mockCache := &mockSuggestionCache{data: make(map[string][]string)}
	mockClient := &mockGNMIClient{
		getResponse: &gnmi.GetResponse{
			Notification: []*gnmi.Notification{
				{
					Update: []*gnmi.Update{
						{Val: &gnmi.TypedValue{Value: &gnmi.TypedValue_StringVal{StringVal: "val"}}},
					},
				},
			},
		},
	}
	target := &target.Target{Client: mockClient, Config: &types.TargetConfig{}}
	engine := NewSuggestionEngine(target, mockCache)

	items := []SuggestionItem{
		{Path: "/p1", Key: "k1", StatePath: "/p1[k1=*]/state/k1"},
		{Path: "/p2", Key: "k2", StatePath: "/p2[k2=*]/state/k2"},
	}

	// Execute
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	engine.Start(ctx, items)

	// Wait for async execution (simple sleep for test)
	time.Sleep(100 * time.Millisecond)

	// Verify
	if len(mockCache.Get("/p1::k1")) == 0 {
		t.Error("Expected suggestions for /p1, got none")
	}
	if len(mockCache.Get("/p2::k2")) == 0 {
		t.Error("Expected suggestions for /p2, got none")
	}
}
