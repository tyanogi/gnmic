package app

import (
	"testing"

	"github.com/openconfig/goyang/pkg/yang"
)

func TestSchemaWalkerTypes(t *testing.T) {
	// This test is designed to verify the existence of the types/interfaces required by the spec.
	// It will fail compilation if ListInfo or Walk are not defined.

	// Verify ListInfo struct exists and has the required fields
	_ = ListInfo{
		Path:      "/test/path",
		Key:       "name",
		StatePath: "/test/path/state/name",
	}

	// Verify Walk function signature
	// It should take a *yang.Entry and return a slice of ListInfo and an error.
	var _ func(*yang.Entry) ([]ListInfo, error) = Walk
}

func TestCreateMockEntry(t *testing.T) {
	// Attempt to use the helper function (Red Phase: function does not exist)
	e := createMockEntry("module-name")
	if e.Name != "module-name" {
		t.Errorf("Expected name 'module-name', got %s", e.Name)
	}
}

// createMockEntry creates a *yang.Entry with the given name and children.
func createMockEntry(name string, children ...*yang.Entry) *yang.Entry {
	e := &yang.Entry{
		Name: name,
		Dir:  make(map[string]*yang.Entry),
	}
	for _, c := range children {
		e.Dir[c.Name] = c
		c.Parent = e
	}
	return e
}

func TestWalk_SimpleList(t *testing.T) {
	// Build schema:
	// module: test-module
	//   container: interfaces
	//     list: interface [key=name]
	//       leaf: name

	leafName := &yang.Entry{
		Name: "name",
		Kind: yang.LeafEntry,
		Type: &yang.YangType{Kind: yang.Ystring},
	}
	listInterface := &yang.Entry{
		Name:     "interface",
		Kind:     yang.DirectoryEntry,
		ListAttr: &yang.ListAttr{},
		Dir:      map[string]*yang.Entry{"name": leafName},
		Key:      "name",
	}
	leafName.Parent = listInterface

	containerInterfaces := createMockEntry("interfaces", listInterface)
	// We do not add the module level to path usually in gNMI paths, but goyang might include it if it's not handled.
	// Usually Walk starts at a node. If we start at containerInterfaces, path is /interfaces...
	// If we start at module, path might include module name depending on implementation.
	// But standard gNMI paths don't include module name.
	// Let's assume Walk handles the root properly.
	// For this test, let's treat containerInterfaces as a top-level container in a module.
	module := createMockEntry("test-module", containerInterfaces)
	module.Kind = yang.DirectoryEntry // Modules are directories

	results, err := Walk(module)
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	r := results[0]
	// Path generation depends on logic. Assuming standard string concatenation or entry.Path()
	// yang.Entry.Path() returns /module-name/interfaces/interface if parent is module.
	// Wait, we want gNMI paths: /interfaces/interface
	// So we should probably ignore the module name or handle it.
	// Let's assume for now the expected path excludes module name if it's the root.
	// Actually, if I use goyang's Path(), it includes the module.
	// I'll check my implementation requirement: "リストの絶対パス (例: /interfaces/interface)"
	// This implies gNMI path.

	// For the Red phase, I will assert the gNMI path.
	if r.Path != "/interfaces/interface" {
		t.Errorf("Expected path /interfaces/interface, got %s", r.Path)
	}
	if r.Key != "name" {
		t.Errorf("Expected key name, got %s", r.Key)
	}
	if r.StatePath != "/interfaces/interface/name" {
		t.Errorf("Expected StatePath /interfaces/interface/name, got %s", r.StatePath)
	}
}

func TestWalk_CompositeKey(t *testing.T) {
	// Build schema:
	// list: acl-entry [key="sequence-id name"]

	listAcl := &yang.Entry{
		Name:     "acl-entry",
		Kind:     yang.DirectoryEntry,
		ListAttr: &yang.ListAttr{},
		Key:      "sequence-id name",
	}

	// Wrap in a container to have a valid parent-child structure for getPath
	containerAcl := createMockEntry("acl", listAcl)
	module := createMockEntry("acl-module", containerAcl)
	module.Kind = yang.DirectoryEntry

	results, err := Walk(module)
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	r1 := results[0]
	if r1.Key != "sequence-id" {
		t.Errorf("Expected first key sequence-id, got %s", r1.Key)
	}
	if r1.StatePath != "/acl/acl-entry/sequence-id" {
		t.Errorf("Expected StatePath /acl/acl-entry/sequence-id, got %s", r1.StatePath)
	}

	r2 := results[1]
	if r2.Key != "name" {
		t.Errorf("Expected second key name, got %s", r2.Key)
	}
	if r2.StatePath != "/acl/acl-entry/name" {
		t.Errorf("Expected StatePath /acl/acl-entry/name, got %s", r2.StatePath)
	}
}

func TestWalk_OpenConfigHeuristic(t *testing.T) {
	// Build schema:
	// list interface [key=name]
	//   leaf name
	//   container state
	//     leaf name

	leafName := &yang.Entry{
		Name: "name",
		Kind: yang.LeafEntry,
		Type: &yang.YangType{Kind: yang.Ystring},
	}
	leafStateName := &yang.Entry{
		Name: "name",
		Kind: yang.LeafEntry,
		Type: &yang.YangType{Kind: yang.Ystring},
	}
	containerState := &yang.Entry{
		Name: "state",
		Kind: yang.DirectoryEntry,
		Dir:  map[string]*yang.Entry{"name": leafStateName},
	}
	leafStateName.Parent = containerState

	listInterface := &yang.Entry{
		Name:     "interface",
		Kind:     yang.DirectoryEntry,
		ListAttr: &yang.ListAttr{},
		Dir:      map[string]*yang.Entry{"name": leafName, "state": containerState},
		Key:      "name",
	}
	leafName.Parent = listInterface
	containerState.Parent = listInterface

	module := createMockEntry("oc-interfaces", listInterface)
	module.Kind = yang.DirectoryEntry

	results, err := Walk(module)
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	r := results[0]
	// Heuristic should find /interface/state/name
	if r.StatePath != "/interface/state/name" {
		t.Errorf("Expected StatePath /interface/state/name, got %s", r.StatePath)
	}
}

func TestWalk_ComplexNested(t *testing.T) {
	// Build schema:
	// /network-instances
	//   /network-instance [key=name]
	//     /protocols
	//       /protocol [key="identifier name"]
	//         /bgp
	//           /neighbors
	//             /neighbor [key=neighbor-address]
	//               /state
	//                 /neighbor-address

	neighborAddressLeaf := &yang.Entry{Name: "neighbor-address", Kind: yang.LeafEntry}
	neighborStateAddressLeaf := &yang.Entry{Name: "neighbor-address", Kind: yang.LeafEntry}
	neighborState := &yang.Entry{
		Name: "state",
		Kind: yang.DirectoryEntry,
		Dir:  map[string]*yang.Entry{"neighbor-address": neighborStateAddressLeaf},
	}
	neighborStateAddressLeaf.Parent = neighborState

	listNeighbor := &yang.Entry{
		Name:     "neighbor",
		Kind:     yang.DirectoryEntry,
		ListAttr: &yang.ListAttr{},
		Key:      "neighbor-address",
		Dir:      map[string]*yang.Entry{"neighbor-address": neighborAddressLeaf, "state": neighborState},
	}
	neighborAddressLeaf.Parent = listNeighbor
	neighborState.Parent = listNeighbor

	containerNeighbors := createMockEntry("neighbors", listNeighbor)
	containerBgp := createMockEntry("bgp", containerNeighbors)

	protocolIdLeaf := &yang.Entry{Name: "identifier", Kind: yang.LeafEntry}
	protocolNameLeaf := &yang.Entry{Name: "name", Kind: yang.LeafEntry}
	listProtocol := &yang.Entry{
		Name:     "protocol",
		Kind:     yang.DirectoryEntry,
		ListAttr: &yang.ListAttr{},
		Key:      "identifier name",
		Dir:      map[string]*yang.Entry{"identifier": protocolIdLeaf, "name": protocolNameLeaf, "bgp": containerBgp},
	}
	protocolIdLeaf.Parent = listProtocol
	protocolNameLeaf.Parent = listProtocol
	containerBgp.Parent = listProtocol

	containerProtocols := createMockEntry("protocols", listProtocol)

	niNameLeaf := &yang.Entry{Name: "name", Kind: yang.LeafEntry}
	listNi := &yang.Entry{
		Name:     "network-instance",
		Kind:     yang.DirectoryEntry,
		ListAttr: &yang.ListAttr{},
		Key:      "name",
		Dir:      map[string]*yang.Entry{"name": niNameLeaf, "protocols": containerProtocols},
	}
	niNameLeaf.Parent = listNi
	containerProtocols.Parent = listNi

	containerNi := createMockEntry("network-instances", listNi)
	module := createMockEntry("oc-ni", containerNi)
	module.Kind = yang.DirectoryEntry

	results, err := Walk(module)
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	// Expected results:
	// 1. network-instance (key: name)
	// 2. protocol (key: identifier)
	// 3. protocol (key: name)
	// 4. neighbor (key: neighbor-address)
	if len(results) != 4 {
		t.Fatalf("Expected 4 results, got %d", len(results))
	}

	expected := []struct {
		path      string
		key       string
		statePath string
	}{
		{"/network-instances/network-instance", "name", "/network-instances/network-instance/name"},
		{"/network-instances/network-instance/protocols/protocol", "identifier", "/network-instances/network-instance/protocols/protocol/identifier"},
		{"/network-instances/network-instance/protocols/protocol", "name", "/network-instances/network-instance/protocols/protocol/name"},
		{"/network-instances/network-instance/protocols/protocol/bgp/neighbors/neighbor", "neighbor-address", "/network-instances/network-instance/protocols/protocol/bgp/neighbors/neighbor/state/neighbor-address"},
	}

	for i, exp := range expected {
		if results[i].Path != exp.path {
			t.Errorf("[%d] Expected path %s, got %s", i, exp.path, results[i].Path)
		}
		if results[i].Key != exp.key {
			t.Errorf("[%d] Expected key %s, got %s", i, exp.key, results[i].Key)
		}
		if results[i].StatePath != exp.statePath {
			t.Errorf("[%d] Expected statePath %s, got %s", i, exp.statePath, results[i].StatePath)
		}
	}
}
