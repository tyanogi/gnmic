// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/goyang/pkg/yang"
)

func InitDebugLog() {
	f, err := os.OpenFile("gnmic_prompt_debug.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	f.Close()
}

func DebugLog(format string, a ...interface{}) {
	f, err := os.OpenFile("gnmic_prompt_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, format+"\n", a...)
}

// GetPath returns the absolute gNMI path of the given YANG entry.
// It skips choice and case nodes to produce a path compatible with gNMI.
func GetPath(entry *yang.Entry) string {
	parts := []string{}
	// Stop when Parent is nil (usually the module)
	for e := entry; e != nil && e.Parent != nil; e = e.Parent {
		if e.IsChoice() || e.IsCase() {
			continue
		}
		parts = append(parts, e.Name)
	}
	// Reverse
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return "/" + strings.Join(parts, "/")
}

func (a *App) printCapResponse(printPrefix string, msg *gnmi.CapabilityResponse) {
	sb := strings.Builder{}
	sb.WriteString("gNMI version: ")
	sb.WriteString(msg.GNMIVersion)
	sb.WriteString("\n")
	if a.Config.LocalFlags.CapabilitiesVersion {
		return
	}
	sb.WriteString("supported models:\n")
	for _, sm := range msg.SupportedModels {
		sb.WriteString("  - ")
		sb.WriteString(sm.GetName())
		sb.WriteString(", ")
		sb.WriteString(sm.GetOrganization())
		sb.WriteString(", ")
		sb.WriteString(sm.GetVersion())
		sb.WriteString("\n")
	}
	sb.WriteString("supported encodings:\n")
	for _, se := range msg.SupportedEncodings {
		sb.WriteString("  - ")
		sb.WriteString(se.String())
		sb.WriteString("\n")
	}
	fmt.Fprintf(a.out, "%s\n", indent(printPrefix, sb.String()))
}

func indent(prefix, s string) string {
	if prefix == "" {
		return s
	}
	prefix = "\n" + strings.TrimRight(prefix, "\n")
	lines := strings.Split(s, "\n")
	return strings.TrimLeft(fmt.Sprintf("%s%s", prefix, strings.Join(lines, prefix)), "\n")
}
