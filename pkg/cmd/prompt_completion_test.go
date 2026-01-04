// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"testing"

	"github.com/openconfig/gnmic/pkg/app"
	"github.com/openconfig/goyang/pkg/yang"
)

func TestFindMatchedXPATH_Subinterface(t *testing.T) {
	// Setup Schema
	// /interfaces/interface[name]/subinterfaces/subinterface[index]

	subintList := &yang.Entry{
		Name: "subinterface",
		Key:  "index",
		Dir:  map[string]*yang.Entry{},
	}

	subintsCont := &yang.Entry{
		Name: "subinterfaces",
		Dir: map[string]*yang.Entry{
			"subinterface": subintList,
		},
	}
	subintList.Parent = subintsCont

	intList := &yang.Entry{
		Name: "interface",
		Key:  "name",
		Dir: map[string]*yang.Entry{
			"subinterfaces": subintsCont,
		},
	}
	subintsCont.Parent = intList

	interfacesCont := &yang.Entry{
		Name: "interfaces",
		Dir: map[string]*yang.Entry{
			"interface": intList,
		},
	}
	intList.Parent = interfacesCont

	root := &yang.Entry{
		Name: "root",
		Dir: map[string]*yang.Entry{
			"interfaces": interfacesCont,
		},
	}
	interfacesCont.Parent = root

	// Mock Cache
	cache := app.NewSuggestionCache()
	// Path for subinterface: /interfaces/interface/subinterfaces/subinterface
	// Key: index
	cache.Add("/interfaces/interface/subinterfaces/subinterface::index", []string{"0", "1"})

	// Setup Global App
	// We need to save the original gApp and restore it after test
	originalGApp := gApp
	defer func() { gApp = originalGApp }()

	gApp = app.New()
	gApp.SuggestionCache = cache
	gApp.Config.LocalFlags.PromptSuggestWithOrigin = false

	// Test Case
	// input: /interfaces/interface[name=eth0]/subinterfaces/subinterface[index=
	input := "/interfaces/interface[name=eth0]/subinterfaces/subinterface[index="

	suggestions := findMatchedXPATH(root, input, false)

	found0 := false
	found1 := false

	for _, s := range suggestions {
		if s.Text == "/interfaces/interface[name=eth0]/subinterfaces/subinterface[index=0]" {
			found0 = true
		}
		if s.Text == "/interfaces/interface[name=eth0]/subinterfaces/subinterface[index=1]" {
			found1 = true
		}
	}

	if !found0 {
		t.Errorf("Expected suggestion ending in index=0, got %v", suggestions)
	}
	if !found1 {
		t.Errorf("Expected suggestion ending in index=1, got %v", suggestions)
	}
}
