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

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/gnmic/pkg/api/target"
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
	target := &target.Target{Client: mockClient}
	engine := &suggestionEngineImpl{
		target: target,
		cache:  mockCache,
	}

	list := ListInfo{
		Path:      "/interfaces/interface",
		Key:       "name",
		StatePath: "/interfaces/interface/state/name",
	}

	// Execute (Red Phase: prefetchList is not implemented)
	ctx := context.Background()
	err := engine.prefetchList(ctx, list)
	if err != nil {
		t.Fatalf("prefetchList failed: %v", err)
	}

	// Verify cache
	got := mockCache.Get("/interfaces/interface")
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
