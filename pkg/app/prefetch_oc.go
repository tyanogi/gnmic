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
	"log"
	"sync"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/gnmic/pkg/api/path"
	"github.com/openconfig/gnmic/pkg/api/target"
)

const (
	ocInterfaceNamePath = "/interfaces/interface/config/name"
	ocNetInstNamePath   = "/network-instances/network-instance/config/name"
)

var prefetchTargets = map[string]string{
	"interface": ocInterfaceNamePath,
	"netinst":   ocNetInstNamePath,
}

// StartOpenConfigPrefetch starts the async prefetch of OpenConfig resources
func StartOpenConfigPrefetch(ctx context.Context, t *target.Target, store *VariableStore, logger *log.Logger) {
	var wg sync.WaitGroup

	for category, pathStr := range prefetchTargets {
		wg.Add(1)
		go func(cat, pStr string) {
			defer wg.Done()
			values, err := fetchValues(ctx, t, pStr)
			if err != nil {
				if logger != nil {
					logger.Printf("failed to prefetch %s: %v", cat, err)
				}
				return
			}
			if len(values) > 0 {
				store.Set(cat, values)
			}
		}(category, pathStr)
	}
	wg.Wait()
}

func fetchValues(ctx context.Context, t *target.Target, pathStr string) ([]string, error) {
	if t.Client == nil {
		return nil, fmt.Errorf("gNMI client not initialized")
	}
	gnmiPath, err := path.ParsePath(pathStr)
	if err != nil {
		return nil, err
	}

	req := &gnmi.GetRequest{
		Path:     []*gnmi.Path{gnmiPath},
		Type:     gnmi.GetRequest_CONFIG,
		Encoding: gnmi.Encoding_JSON,
	}

	resp, err := t.Client.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return extractLeafValues(resp), nil
}

func extractLeafValues(resp *gnmi.GetResponse) []string {
	values := make([]string, 0)
	for _, notif := range resp.Notification {
		for _, update := range notif.Update {
			if update.Val == nil {
				continue
			}
			valStr := getValueString(update.Val)
			if valStr != "" {
				values = append(values, valStr)
			}
		}
	}
	return values
}

func getValueString(tv *gnmi.TypedValue) string {
	if tv == nil {
		return ""
	}
	switch v := tv.Value.(type) {
	case *gnmi.TypedValue_StringVal:
		return v.StringVal
	case *gnmi.TypedValue_IntVal:
		return fmt.Sprintf("%d", v.IntVal)
	case *gnmi.TypedValue_UintVal:
		return fmt.Sprintf("%d", v.UintVal)
	case *gnmi.TypedValue_BoolVal:
		return fmt.Sprintf("%v", v.BoolVal)
	case *gnmi.TypedValue_BytesVal:
		return string(v.BytesVal)
	case *gnmi.TypedValue_FloatVal:
		return fmt.Sprintf("%f", v.FloatVal)
	case *gnmi.TypedValue_DecimalVal:
		return fmt.Sprintf("%d.%d", v.DecimalVal.Digits, v.DecimalVal.Precision)
	case *gnmi.TypedValue_LeaflistVal:
		// Handle leaf-list if necessary, strictly taking first element or joining
		// For simple name extraction, this might not be hit if path points to leaf
		return fmt.Sprintf("%v", v.LeaflistVal)
	case *gnmi.TypedValue_AsciiVal:
		return v.AsciiVal
	case *gnmi.TypedValue_AnyVal:
		return fmt.Sprintf("%v", v.AnyVal)
	case *gnmi.TypedValue_JsonVal:
		return string(v.JsonVal)
	case *gnmi.TypedValue_JsonIetfVal:
		return string(v.JsonIetfVal)
	default:
		return fmt.Sprintf("%v", v)
	}
}
