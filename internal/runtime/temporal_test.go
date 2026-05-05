package runtime

import (
	"testing"
)

// TestTemporalGridCreation tests basic temporal grid initialization.
func TestTemporalGridCreation(t *testing.T) {
	grid := NewTemporalGrid()

	if grid == nil {
		t.Fatal("NewTemporalGrid() returned nil")
	}

	if grid.CurrentTick != 0 {
		t.Errorf("Expected CurrentTick = 0, got %d", grid.CurrentTick)
	}

	if grid.CurrentBeat != 0 {
		t.Errorf("Expected CurrentBeat = 0, got %d", grid.CurrentBeat)
	}

	if grid.CurrentBar != 0 {
		t.Errorf("Expected CurrentBar = 0, got %d", grid.CurrentBar)
	}

	if grid.TicksPerBeat != 480 {
		t.Errorf("Expected TicksPerBeat = 480, got %d", grid.TicksPerBeat)
	}

	if grid.BeatsPerBar != 4 {
		t.Errorf("Expected BeatsPerBar = 4, got %d", grid.BeatsPerBar)
	}

	if grid.Tempo != 120 {
		t.Errorf("Expected Tempo = 120, got %d", grid.Tempo)
	}
}

// TestTemporalGridTick tests single tick advancement.
func TestTemporalGridTick(t *testing.T) {
	grid := NewTemporalGrid()

	// Tick once
	events := grid.Tick()

	if grid.CurrentTick != 1 {
		t.Errorf("Expected CurrentTick = 1, got %d", grid.CurrentTick)
	}

	if len(events) != 0 {
		t.Errorf("Expected 0 events, got %d", len(events))
	}
}

// TestTemporalGridAdvanceBeat tests beat advancement.
func TestTemporalGridAdvanceBeat(t *testing.T) {
	grid := NewTemporalGrid()
	initialTick := grid.CurrentTick

	events := grid.AdvanceBeat()

	expectedTick := initialTick + grid.TicksPerBeat
	if grid.CurrentTick != expectedTick {
		t.Errorf("Expected CurrentTick = %d, got %d", expectedTick, grid.CurrentTick)
	}

	if len(events) != 0 {
		t.Errorf("Expected 0 events, got %d", len(events))
	}
}

// TestTemporalGridAdvanceBar tests bar advancement.
func TestTemporalGridAdvanceBar(t *testing.T) {
	grid := NewTemporalGrid()
	initialTick := grid.CurrentTick

	events := grid.AdvanceBar()

	expectedTick := initialTick + (grid.TicksPerBeat * grid.BeatsPerBar)
	if grid.CurrentTick != expectedTick {
		t.Errorf("Expected CurrentTick = %d, got %d", expectedTick, grid.CurrentTick)
	}

	if len(events) != 0 {
		t.Errorf("Expected 0 events, got %d", len(events))
	}
}

// TestScheduleEvent tests event scheduling.
func TestScheduleEvent(t *testing.T) {
	grid := NewTemporalGrid()

	// Schedule an event at tick 100
	event := &ScheduledEvent{
		Tick:      100,
		Beat:      0,
		Bar:       0,
		EventType: EventTypeAction,
		Payload:   "test_event",
	}

	grid.ScheduleEvent(event)

	if len(grid.EventQueue) != 1 {
		t.Errorf("Expected 1 event in queue, got %d", len(grid.EventQueue))
	}

	if grid.EventQueue[0].Tick != 100 {
		t.Errorf("Expected event at tick 100, got %d", grid.EventQueue[0].Tick)
	}
}

// TestScheduleAtBeat tests beat-based scheduling.
func TestScheduleAtBeat(t *testing.T) {
	grid := NewTemporalGrid()

	grid.ScheduleAtBeat(2, EventTypeAction, "beat_event")

	if len(grid.EventQueue) != 1 {
		t.Fatalf("Expected 1 event in queue, got %d", len(grid.EventQueue))
	}

	expectedTick := 2 * grid.TicksPerBeat
	if grid.EventQueue[0].Tick != expectedTick {
		t.Errorf("Expected event at tick %d, got %d", expectedTick, grid.EventQueue[0].Tick)
	}
}

// TestScheduleAtBar tests bar-based scheduling.
func TestScheduleAtBar(t *testing.T) {
	grid := NewTemporalGrid()

	grid.ScheduleAtBar(1, EventTypeAction, "bar_event")

	if len(grid.EventQueue) != 1 {
		t.Fatalf("Expected 1 event in queue, got %d", len(grid.EventQueue))
	}

	expectedBeat := 1 * grid.BeatsPerBar
	expectedTick := expectedBeat * grid.TicksPerBeat
	if grid.EventQueue[0].Tick != expectedTick {
		t.Errorf("Expected event at tick %d, got %d", expectedTick, grid.EventQueue[0].Tick)
	}
}

// TestEventProcessing tests that scheduled events are processed correctly.
func TestEventProcessing(t *testing.T) {
	grid := NewTemporalGrid()

	// Schedule events at different times
	grid.ScheduleAtBeat(1, EventTypeAction, "event1")
	grid.ScheduleAtBeat(2, EventTypeAction, "event2")
	grid.ScheduleAtBeat(3, EventTypeAction, "event3")

	// Advance to beat 2
	events := grid.AdvanceBeat() // Beat 1
	events = append(events, grid.AdvanceBeat()...) // Beat 2

	// Should have processed events for beats 1 and 2
	if len(events) != 2 {
		t.Errorf("Expected 2 events processed, got %d", len(events))
	}

	// One event should remain in queue (beat 3)
	if len(grid.EventQueue) != 1 {
		t.Errorf("Expected 1 event remaining in queue, got %d", len(grid.EventQueue))
	}
}

// TestMotifRegistry tests motif registration and retrieval.
func TestMotifRegistry(t *testing.T) {
	registry := NewMotifRegistry()

	motif := &MotifDefinition{
		Name:       "test_motif",
		Parameters: []string{"param1"},
		Duration:   4,
		Body:       nil,
	}

	// Register motif
	err := registry.RegisterMotif(motif)
	if err != nil {
		t.Fatalf("Failed to register motif: %v", err)
	}

	// Retrieve motif
	retrieved, exists := registry.GetMotif("test_motif")
	if !exists {
		t.Fatal("Motif not found after registration")
	}

	if retrieved.Name != "test_motif" {
		t.Errorf("Expected motif name 'test_motif', got '%s'", retrieved.Name)
	}

	if retrieved.Duration != 4 {
		t.Errorf("Expected duration 4, got %d", retrieved.Duration)
	}
}

// TestMotifRegistryDuplicate tests that duplicate motifs are rejected.
func TestMotifRegistryDuplicate(t *testing.T) {
	registry := NewMotifRegistry()

	motif1 := &MotifDefinition{
		Name:     "duplicate",
		Duration: 4,
	}

	motif2 := &MotifDefinition{
		Name:     "duplicate",
		Duration: 8,
	}

	// Register first motif
	err := registry.RegisterMotif(motif1)
	if err != nil {
		t.Fatalf("Failed to register first motif: %v", err)
	}

	// Try to register duplicate
	err = registry.RegisterMotif(motif2)
	if err == nil {
		t.Error("Expected error when registering duplicate motif, got nil")
	}
}

// TestBarRange tests bar range functionality.
func TestBarRange(t *testing.T) {
	br := NewBarRange(1, 4)

	if br.Start != 1 {
		t.Errorf("Expected start 1, got %d", br.Start)
	}

	if br.End != 4 {
		t.Errorf("Expected end 4, got %d", br.End)
	}

	// Test Contains
	if !br.Contains(2) {
		t.Error("Expected bar 2 to be in range [1,4]")
	}

	if br.Contains(5) {
		t.Error("Expected bar 5 to not be in range [1,4]")
	}

	// Test Duration
	if br.Duration() != 4 {
		t.Errorf("Expected duration 4, got %d", br.Duration())
	}

	// Test ToBeats
	startBeat, endBeat := br.ToBeats(4) // 4 beats per bar
	if startBeat != 4 {
		t.Errorf("Expected start beat 4, got %d", startBeat)
	}
	if endBeat != 19 {
		t.Errorf("Expected end beat 19, got %d", endBeat)
	}
}

// TestSetTempo tests tempo changes.
func TestSetTempo(t *testing.T) {
	grid := NewTemporalGrid()

	grid.SetTempo(140)

	if grid.Tempo != 140 {
		t.Errorf("Expected tempo 140, got %d", grid.Tempo)
	}
}

// TestSetTimeSignature tests time signature changes.
func TestSetTimeSignature(t *testing.T) {
	grid := NewTemporalGrid()

	grid.SetTimeSignature(3, 480) // 3/4 time

	if grid.BeatsPerBar != 3 {
		t.Errorf("Expected 3 beats per bar, got %d", grid.BeatsPerBar)
	}

	if grid.TicksPerBeat != 480 {
		t.Errorf("Expected 480 ticks per beat, got %d", grid.TicksPerBeat)
	}
}

// TestReset tests temporal grid reset.
func TestReset(t *testing.T) {
	grid := NewTemporalGrid()

	// Advance time and schedule events
	grid.AdvanceBar()
	grid.ScheduleAtBeat(10, EventTypeAction, "event")

	// Reset
	grid.Reset()

	if grid.CurrentTick != 0 {
		t.Errorf("Expected CurrentTick = 0 after reset, got %d", grid.CurrentTick)
	}

	if grid.CurrentBeat != 0 {
		t.Errorf("Expected CurrentBeat = 0 after reset, got %d", grid.CurrentBeat)
	}

	if grid.CurrentBar != 0 {
		t.Errorf("Expected CurrentBar = 0 after reset, got %d", grid.CurrentBar)
	}

	if len(grid.EventQueue) != 0 {
		t.Errorf("Expected empty event queue after reset, got %d events", len(grid.EventQueue))
	}
}


