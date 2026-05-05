package runtime

import (
	"testing"
)

// TestResonanceGraphCreation tests basic resonance graph initialization.
func TestResonanceGraphCreation(t *testing.T) {
	rg := NewResonanceGraph()

	if rg == nil {
		t.Fatal("NewResonanceGraph() returned nil")
	}

	if rg.links == nil {
		t.Error("Resonance graph links map is nil")
	}

	if len(rg.links) != 0 {
		t.Errorf("Expected 0 links, got %d", len(rg.links))
	}
}

// TestAddLink tests adding resonance links.
func TestAddLink(t *testing.T) {
	rg := NewResonanceGraph()

	link := &ResonanceLink{
		SourceString: "stringA",
		SourceMode:   "modeX",
		TargetString: "stringB",
		TargetMode:   "modeY",
		Strength:     0.8,
		CouplingType: CouplingLinear,
	}

	rg.AddLink(link)

	links := rg.GetLinks("stringA", "modeX")
	if len(links) != 1 {
		t.Fatalf("Expected 1 link, got %d", len(links))
	}

	if links[0].TargetString != "stringB" {
		t.Errorf("Expected target string 'stringB', got '%s'", links[0].TargetString)
	}

	if links[0].Strength != 0.8 {
		t.Errorf("Expected strength 0.8, got %f", links[0].Strength)
	}
}

// TestPropagateChange tests linear coupling propagation.
func TestPropagateChange(t *testing.T) {
	rg := NewResonanceGraph()

	// Add a linear coupling: A.x -> B.y with strength 0.5
	link := &ResonanceLink{
		SourceString: "A",
		SourceMode:   "x",
		TargetString: "B",
		TargetMode:   "y",
		Strength:     0.5,
		CouplingType: CouplingLinear,
	}
	rg.AddLink(link)

	// Propagate a change of delta = 10.0
	affected := rg.PropagateChange("A", "x", 10.0)

	if len(affected) != 1 {
		t.Fatalf("Expected 1 affected field, got %d", len(affected))
	}

	targetKey := makeFieldKey("B", "y")
	delta, ok := affected[targetKey]
	if !ok {
		t.Fatal("Expected 'B.y' to be affected")
	}

	expectedDelta := 10.0 * 0.5
	if delta != expectedDelta {
		t.Errorf("Expected propagated delta %f, got %f", expectedDelta, delta)
	}
}

// TestMultipleCouplings tests multiple outgoing couplings from one field.
func TestMultipleCouplings(t *testing.T) {
	rg := NewResonanceGraph()

	// A.x -> B.y with strength 0.5
	rg.AddLink(&ResonanceLink{
		SourceString: "A",
		SourceMode:   "x",
		TargetString: "B",
		TargetMode:   "y",
		Strength:     0.5,
		CouplingType: CouplingLinear,
	})

	// A.x -> C.z with strength 0.3
	rg.AddLink(&ResonanceLink{
		SourceString: "A",
		SourceMode:   "x",
		TargetString: "C",
		TargetMode:   "z",
		Strength:     0.3,
		CouplingType: CouplingLinear,
	})

	// Propagate change
	affected := rg.PropagateChange("A", "x", 10.0)

	if len(affected) != 2 {
		t.Fatalf("Expected 2 affected fields, got %d", len(affected))
	}

	// Check B.y
	byKey := makeFieldKey("B", "y")
	if delta, ok := affected[byKey]; !ok || delta != 5.0 {
		t.Errorf("Expected B.y delta 5.0, got %f", delta)
	}

	// Check C.z
	czKey := makeFieldKey("C", "z")
	if delta, ok := affected[czKey]; !ok || delta != 3.0 {
		t.Errorf("Expected C.z delta 3.0, got %f", delta)
	}
}

// TestStringStateCreation tests string state initialization.
func TestStringStateCreation(t *testing.T) {
	ss := NewStringState("testString")

	if ss == nil {
		t.Fatal("NewStringState() returned nil")
	}

	if ss.Name != "testString" {
		t.Errorf("Expected name 'testString', got '%s'", ss.Name)
	}

	if ss.Modes == nil {
		t.Error("String state modes map is nil")
	}

	if len(ss.Modes) != 0 {
		t.Errorf("Expected 0 modes, got %d", len(ss.Modes))
	}
}

// TestStringStateSetGetMode tests mode value manipulation.
func TestStringStateSetGetMode(t *testing.T) {
	ss := NewStringState("testString")

	// Set a mode value
	ss.SetMode("frequency", 440.0)

	// Get the mode value
	value, exists := ss.GetMode("frequency")
	if !exists {
		t.Fatal("Expected mode 'frequency' to exist")
	}

	if value != 440.0 {
		t.Errorf("Expected frequency 440.0, got %f", value)
	}
}

// TestStringStateUpdateMode tests mode value updates with delta.
func TestStringStateUpdateMode(t *testing.T) {
	ss := NewStringState("testString")

	// Initialize mode
	ss.SetMode("amplitude", 1.0)

	// Update with delta
	delta := ss.UpdateMode("amplitude", 0.5)

	if delta != 0.5 {
		t.Errorf("Expected delta 0.5, got %f", delta)
	}

	// Check new value
	value, _ := ss.GetMode("amplitude")
	if value != 1.5 {
		t.Errorf("Expected amplitude 1.5, got %f", value)
	}
}

// TestBraneContextCreation tests brane context initialization.
func TestBraneContextCreation(t *testing.T) {
	bc := NewBraneContext("testBrane", []string{"dim1", "dim2", "dim3"})

	if bc == nil {
		t.Fatal("NewBraneContext() returned nil")
	}

	if bc.Name != "testBrane" {
		t.Errorf("Expected name 'testBrane', got '%s'", bc.Name)
	}

	if len(bc.Dimensions) != 3 {
		t.Fatalf("Expected 3 dimensions, got %d", len(bc.Dimensions))
	}

	if bc.Dimensions[0] != "dim1" {
		t.Errorf("Expected dimension 'dim1', got '%s'", bc.Dimensions[0])
	}
}

// TestBraneContextAttachString tests string attachment to branes.
func TestBraneContextAttachString(t *testing.T) {
	bc := NewBraneContext("testBrane", []string{"time", "space"})

	bc.AttachString("stringA")
	bc.AttachString("stringB")

	if len(bc.Strings) != 2 {
		t.Fatalf("Expected 2 attached strings, got %d", len(bc.Strings))
	}

	if !bc.IsStringAttached("stringA") {
		t.Error("Expected stringA to be attached")
	}

	if !bc.IsStringAttached("stringB") {
		t.Error("Expected stringB to be attached")
	}

	if bc.IsStringAttached("stringC") {
		t.Error("Expected stringC to not be attached")
	}
}

// TestBraneContextDimensionCount tests dimension counting.
func TestBraneContextDimensionCount(t *testing.T) {
	bc := NewBraneContext("testBrane", []string{"x", "y", "z", "t"})

	count := bc.GetDimensionCount()
	if count != 4 {
		t.Errorf("Expected 4 dimensions, got %d", count)
	}
}

// TestMakeFieldKey tests field key generation.
func TestMakeFieldKey(t *testing.T) {
	key := makeFieldKey("stringA", "modeX")

	expected := "stringA.modeX"
	if key != expected {
		t.Errorf("Expected key '%s', got '%s'", expected, key)
	}
}

// TestResonanceCondition tests resonance condition checking.
func TestResonanceCondition(t *testing.T) {
	ss := NewStringState("testString")
	ss.SetMode("amplitude", 5.0)

	// Create a condition: amplitude > 3.0
	rc := &ResonanceCondition{
		Predicate: func(state *StringState) bool {
			amp, _ := state.GetMode("amplitude")
			return amp > 3.0
		},
		Action: func(state *StringState) {
			// Amplify by 2x
			amp, _ := state.GetMode("amplitude")
			state.SetMode("amplitude", amp*2.0)
		},
	}

	// Check condition
	if !rc.CheckCondition(ss) {
		t.Error("Expected resonance condition to be true")
	}

	// Execute action
	rc.ExecuteAction(ss)

	// Verify action effect
	amp, _ := ss.GetMode("amplitude")
	if amp != 10.0 {
		t.Errorf("Expected amplitude 10.0 after action, got %f", amp)
	}
}

// TestCouplingTypes tests different coupling type behaviors.
func TestCouplingTypes(t *testing.T) {
	tests := []struct {
		name         string
		couplingType CouplingType
		expectedName string
	}{
		{"Linear", CouplingLinear, "linear"},
		{"Nonlinear", CouplingNonlinear, "nonlinear"},
		{"Resonant", CouplingResonant, "resonant"},
	}

	rg := NewResonanceGraph()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link := &ResonanceLink{
				SourceString: "A",
				SourceMode:   "x",
				TargetString: "B",
				TargetMode:   "y",
				Strength:     1.0,
				CouplingType: tt.couplingType,
			}

			rg.AddLink(link)
			links := rg.GetLinks("A", "x")

			found := false
			for _, l := range links {
				if l.CouplingType == tt.couplingType {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("Expected to find coupling type %v", tt.couplingType)
			}
		})
	}
}


