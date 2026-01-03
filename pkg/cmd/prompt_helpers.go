// © 2026 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"strings"

	goprompt "github.com/c-bata/go-prompt"
)

func checkValueSuggestions(line string) []goprompt.Suggest {
	// Debug logging
	// fmt.Fprintf(os.Stderr, "[debug] checkValueSuggestions line: '%s'\n", line)

	// split line by space to get the last argument which is likely the path
	args := strings.Fields(line)
	if len(args) == 0 {
		return nil
	}
	lastArg := args[len(args)-1]

	// Heuristic: Check if the last occurrence of "[name=" comes after the last "]"
	idx := strings.LastIndex(lastArg, "[name=")
	if idx == -1 {
		return nil
	}

	// Check if closed
	closingBracket := strings.Index(lastArg[idx:], "]")
	if closingBracket != -1 {
		return nil
	}

	prefix := lastArg[:idx]
	// Remove trailing slash if present for suffix check
	prefix = strings.TrimSuffix(prefix, "/")

	// fmt.Fprintf(os.Stderr, "[debug] checkValueSuggestions prefix: '%s'\n", prefix)

	if strings.HasSuffix(prefix, "interfaces/interface") {
		return getStoreSuggestions("interface")
	}
	if strings.HasSuffix(prefix, "network-instances/network-instance") {
		return getStoreSuggestions("netinst")
	}

	return nil
}

func getStoreSuggestions(category string) []goprompt.Suggest {
	if gApp.PromptVariableStore == nil {
		return nil
	}
	values := gApp.PromptVariableStore.Get(category)
	if len(values) == 0 {
		return nil
	}
	suggs := make([]goprompt.Suggest, len(values))
	for i, v := range values {
		suggs[i] = goprompt.Suggest{
			Text:        v,
			Description: fmt.Sprintf("Discovered %s", category),
		}
	}
	return suggs
}