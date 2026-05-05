package runtime

import (
	"testing"
)

// TestEntanglementGroupCreation tests basic entanglement group initialization.
func TestEntanglementGroupCreation(t *testing.T) {
	eg := NewEntanglementGroup("test_group")

	if eg == nil {
		t.Fatal("NewEntanglementGroup() returned nil")
	}

	if eg.ID != "test_group" {
		t.Errorf("Expected ID 'test_group', got '%s'", eg.ID)
	}

	if len(eg.Members) != 0 {
		t.Errorf("Expected 0 members, got %d", len(eg.Members))
	}

	if eg.Value != nil {
		t.Error("Expected nil initial value")
	}
}

// TestAddMember tests adding members to an entanglement group.
func TestAddMember(t *testing.T) {
	eg := NewEntanglementGroup("test_group")

	eg.AddMember("var1")
	eg.AddMember("var2")

	if len(eg.Members) != 2 {
		t.Fatalf("Expected 2 members, got %d", len(eg.Members))
	}

	if eg.Members[0] != "var1" {
		t.Errorf("Expected first member 'var1', got '%s'", eg.Members[0])
	}

	if eg.Members[1] != "var2" {
		t.Errorf("Expected second member 'var2', got '%s'", eg.Members[1])
	}
}

// TestSetGetValue tests setting and getting shared values.
func TestSetGetValue(t *testing.T) {
	eg := NewEntanglementGroup("test_group")

	value := &IntValue{Value: 42}
	eg.SetValue(value)

	retrieved := eg.GetValue()
	if retrieved == nil {
		t.Fatal("GetValue() returned nil")
	}

	intVal, ok := retrieved.(*IntValue)
	if !ok {
		t.Fatal("Retrieved value is not IntValue")
	}

	if intVal.Value != 42 {
		t.Errorf("Expected value 42, got %d", intVal.Value)
	}
}

// TestIsMember tests membership checking.
func TestIsMember(t *testing.T) {
	eg := NewEntanglementGroup("test_group")

	eg.AddMember("varA")
	eg.AddMember("varB")

	if !eg.IsMember("varA") {
		t.Error("Expected varA to be a member")
	}

	if !eg.IsMember("varB") {
		t.Error("Expected varB to be a member")
	}

	if eg.IsMember("varC") {
		t.Error("Expected varC to not be a member")
	}
}

// TestEntanglementManagerCreation tests manager initialization.
func TestEntanglementManagerCreation(t *testing.T) {
	em := NewEntanglementManager()

	if em == nil {
		t.Fatal("NewEntanglementManager() returned nil")
	}

	if em.groups == nil {
		t.Error("groups map is nil")
	}

	if em.varToGroup == nil {
		t.Error("varToGroup map is nil")
	}

	if em.nextGroupID != 0 {
		t.Errorf("Expected nextGroupID 0, got %d", em.nextGroupID)
	}
}

// TestCreateGroup tests group creation via manager.
func TestCreateGroup(t *testing.T) {
	em := NewEntanglementManager()

	group := em.CreateGroup()

	if group == nil {
		t.Fatal("CreateGroup() returned nil")
	}

	if group.ID == "" {
		t.Error("Created group has empty ID")
	}

	// Next group should have incremented ID
	group2 := em.CreateGroup()
	if group2.ID == group.ID {
		t.Error("Second group has same ID as first")
	}
}

// TestEntangleVariables tests variable entanglement.
func TestEntangleVariables(t *testing.T) {
	em := NewEntanglementManager()

	varNames := []string{"x", "y", "z"}
	group, err := em.EntangleVariables(varNames)

	if err != nil {
		t.Fatalf("EntangleVariables() returned error: %v", err)
	}

	if group == nil {
		t.Fatal("EntangleVariables() returned nil group")
	}

	if len(group.Members) != 3 {
		t.Errorf("Expected 3 members, got %d", len(group.Members))
	}

	// Check that all variables are mapped to the group
	for _, varName := range varNames {
		if !em.IsEntangled(varName) {
			t.Errorf("Expected variable '%s' to be entangled", varName)
		}
	}
}

// TestEntangleVariablesAlreadyEntangled tests error when variable is already entangled.
func TestEntangleVariablesAlreadyEntangled(t *testing.T) {
	em := NewEntanglementManager()

	// Entangle first group
	_, err := em.EntangleVariables([]string{"x", "y"})
	if err != nil {
		t.Fatalf("First entanglement failed: %v", err)
	}

	// Try to entangle x again
	_, err = em.EntangleVariables([]string{"x", "z"})
	if err == nil {
		t.Error("Expected error when re-entangling variable, got nil")
	}
}

// TestGetGroup tests group retrieval.
func TestGetGroup(t *testing.T) {
	em := NewEntanglementManager()

	varNames := []string{"a", "b"}
	group, _ := em.EntangleVariables(varNames)

	// Get group by variable name
	retrieved, exists := em.GetGroup("a")
	if !exists {
		t.Fatal("Expected group to exist for variable 'a'")
	}

	if retrieved.ID != group.ID {
		t.Errorf("Expected group ID '%s', got '%s'", group.ID, retrieved.ID)
	}

	// Non-existent variable
	_, exists = em.GetGroup("nonexistent")
	if exists {
		t.Error("Expected group to not exist for nonexistent variable")
	}
}

// TestIsEntangled tests entanglement status checking.
func TestIsEntangled(t *testing.T) {
	em := NewEntanglementManager()

	if em.IsEntangled("x") {
		t.Error("Expected 'x' to not be entangled initially")
	}

	em.EntangleVariables([]string{"x", "y"})

	if !em.IsEntangled("x") {
		t.Error("Expected 'x' to be entangled")
	}

	if !em.IsEntangled("y") {
		t.Error("Expected 'y' to be entangled")
	}

	if em.IsEntangled("z") {
		t.Error("Expected 'z' to not be entangled")
	}
}

// TestDisentangleVariable tests removing a variable from entanglement.
func TestDisentangleVariable(t *testing.T) {
	em := NewEntanglementManager()

	em.EntangleVariables([]string{"x", "y", "z"})

	// Disentangle y
	err := em.DisentangleVariable("y")
	if err != nil {
		t.Fatalf("DisentangleVariable() returned error: %v", err)
	}

	// y should no longer be entangled
	if em.IsEntangled("y") {
		t.Error("Expected 'y' to not be entangled after disentanglement")
	}

	// x and z should still be entangled
	if !em.IsEntangled("x") {
		t.Error("Expected 'x' to still be entangled")
	}

	if !em.IsEntangled("z") {
		t.Error("Expected 'z' to still be entangled")
	}
}

// TestDisentangleLastMember tests that group is removed when last member is disentangled.
func TestDisentangleLastMember(t *testing.T) {
	em := NewEntanglementManager()

	group, _ := em.EntangleVariables([]string{"x", "y"})
	groupID := group.ID

	// Disentangle both members
	em.DisentangleVariable("x")
	em.DisentangleVariable("y")

	// Group should be removed
	_, exists := em.groups[groupID]
	if exists {
		t.Error("Expected group to be removed after all members disentangled")
	}
}

// TestListGroups tests listing all groups.
func TestListGroups(t *testing.T) {
	em := NewEntanglementManager()

	em.EntangleVariables([]string{"x", "y"})
	em.EntangleVariables([]string{"a", "b"})

	groups := em.ListGroups()

	if len(groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(groups))
	}
}

// TestCheckTypeCompatibility tests type checking for entangled variables.
func TestCheckTypeCompatibility(t *testing.T) {
	// All same type - should pass
	varTypes1 := map[string]string{
		"x": "Int",
		"y": "Int",
		"z": "Int",
	}
	err := CheckTypeCompatibility(varTypes1)
	if err != nil {
		t.Errorf("Expected no error for compatible types, got: %v", err)
	}

	// Mixed types - should fail
	varTypes2 := map[string]string{
		"x": "Int",
		"y": "Float",
	}
	err = CheckTypeCompatibility(varTypes2)
	if err == nil {
		t.Error("Expected error for incompatible types, got nil")
	}

	// Empty map - should pass
	varTypes3 := map[string]string{}
	err = CheckTypeCompatibility(varTypes3)
	if err != nil {
		t.Errorf("Expected no error for empty map, got: %v", err)
	}
}

// TestEntangledValueSynchronization tests that all members get the same value.
func TestEntangledValueSynchronization(t *testing.T) {
	em := NewEntanglementManager()

	group, _ := em.EntangleVariables([]string{"x", "y", "z"})

	// Set value
	value := &IntValue{Value: 100}
	group.SetValue(value)

	// All variables should return the same value
	for _, varName := range group.Members {
		val, exists := em.GetEntangledValue(varName)
		if !exists {
			t.Errorf("Expected value for '%s'", varName)
			continue
		}

		intVal, ok := val.(*IntValue)
		if !ok {
			t.Errorf("Expected IntValue for '%s'", varName)
			continue
		}

		if intVal.Value != 100 {
			t.Errorf("Expected value 100 for '%s', got %d", varName, intVal.Value)
		}
	}
}


