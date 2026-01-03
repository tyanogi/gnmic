// © 2026 Nokia.
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
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/gnmic/pkg/api/target"
	"github.com/openconfig/gnmic/pkg/api/types"
	"google.golang.org/grpc"
)

// MockGNMIClient implements gnmi.GNMIClient
type MockGNMIClient struct {
	GetFn func(ctx context.Context, in *gnmi.GetRequest, opts ...grpc.CallOption) (*gnmi.GetResponse, error)
}

func (m *MockGNMIClient) Capabilities(ctx context.Context, in *gnmi.CapabilityRequest, opts ...grpc.CallOption) (*gnmi.CapabilityResponse, error) {
	return nil, nil
}
func (m *MockGNMIClient) Get(ctx context.Context, in *gnmi.GetRequest, opts ...grpc.CallOption) (*gnmi.GetResponse, error) {
	if m.GetFn != nil {
		return m.GetFn(ctx, in, opts...)
	}
	return nil, fmt.Errorf("GetFn not implemented")
}
func (m *MockGNMIClient) Set(ctx context.Context, in *gnmi.SetRequest, opts ...grpc.CallOption) (*gnmi.SetResponse, error) {
	return nil, nil
}
func (m *MockGNMIClient) Subscribe(ctx context.Context, opts ...grpc.CallOption) (gnmi.GNMI_SubscribeClient, error) {
	return nil, nil
}

func TestStartOpenConfigPrefetch(t *testing.T) {
	store := NewVariableStore()
	mockClient := &MockGNMIClient{
		GetFn: func(ctx context.Context, in *gnmi.GetRequest, opts ...grpc.CallOption) (*gnmi.GetResponse, error) {
			path := in.Path[0]
			pathStr := ""
			for _, elem := range path.Elem {
				pathStr += "/" + elem.Name
			}

			// Mock response based on path
			if pathStr == "/interfaces/interface/config/name" {
				return &gnmi.GetResponse{
					Notification: []*gnmi.Notification{
						{
							Update: []*gnmi.Update{
								{
									Val: &gnmi.TypedValue{
										Value: &gnmi.TypedValue_StringVal{StringVal: "Ethernet1/1"},
									},
								},
								{
									Val: &gnmi.TypedValue{
										Value: &gnmi.TypedValue_StringVal{StringVal: "Ethernet1/2"},
									},
								},
							},
						},
					},
				}, nil
			} else if pathStr == "/network-instances/network-instance/config/name" {
				return &gnmi.GetResponse{
					Notification: []*gnmi.Notification{
						{
							Update: []*gnmi.Update{
								{
									Val: &gnmi.TypedValue{
										Value: &gnmi.TypedValue_StringVal{StringVal: "default"},
									},
								},
							},
						},
					},
				}, nil
			}
			return nil, fmt.Errorf("unexpected path: %s", pathStr)
		},
	}

	tg := &target.Target{
		Config: &types.TargetConfig{
			Name:    "test-target",
			Address: "localhost:10161",
		},
		Client: mockClient,
	}

	// Run prefetch with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	StartOpenConfigPrefetch(ctx, tg, store)

	// Verify interfaces
	interfaces := store.Get("interface")
	sort.Strings(interfaces)
	expectedInterfaces := []string{"Ethernet1/1", "Ethernet1/2"}
	if !reflect.DeepEqual(interfaces, expectedInterfaces) {
		t.Errorf("Expected interfaces %v, got %v", expectedInterfaces, interfaces)
	}

	// Test IntVal extraction (simulating a numeric name or ID)
	mockClient.GetFn = func(ctx context.Context, in *gnmi.GetRequest, opts ...grpc.CallOption) (*gnmi.GetResponse, error) {
		path := in.Path[0]
		pathStr := ""
		for _, elem := range path.Elem {
			pathStr += "/" + elem.Name
		}
		if pathStr == "/interfaces/interface/config/name" {
			return &gnmi.GetResponse{
				Notification: []*gnmi.Notification{
					{
						Update: []*gnmi.Update{
							{
								Val: &gnmi.TypedValue{
									Value: &gnmi.TypedValue_IntVal{IntVal: 100},
								},
							},
						},
					},
				},
			}, nil
		}
		return nil, fmt.Errorf("error")
	}

	StartOpenConfigPrefetch(ctx, tg, store)
	interfaces = store.Get("interface")
	// Should be overwritten or merged? Current impl overwrites using Set.
	// Since goroutines run concurrently, order isn't guaranteed between different StartOpenConfigPrefetch calls
	// But here we are calling sequentially.
	expectedInterfaces = []string{"100"}
	if !reflect.DeepEqual(interfaces, expectedInterfaces) {
		t.Errorf("Expected interfaces %v, got %v", expectedInterfaces, interfaces)
	}
}

func TestStartOpenConfigPrefetchError(t *testing.T) {
	store := NewVariableStore()
	mockClient := &MockGNMIClient{
		GetFn: func(ctx context.Context, in *gnmi.GetRequest, opts ...grpc.CallOption) (*gnmi.GetResponse, error) {
			return nil, fmt.Errorf("connection failed")
		},
	}
	tg := &target.Target{
		Config: &types.TargetConfig{Name: "test-target"},
		Client: mockClient,
	}

	// Should not panic and should handle error gracefully
	StartOpenConfigPrefetch(context.Background(), tg, store)

	if store.Get("interface") != nil {
		t.Error("Expected nil store for interface on error")
	}
}
