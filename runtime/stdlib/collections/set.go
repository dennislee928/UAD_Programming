package collections

import (
	"fmt"

	"github.com/dennislee928/uad-lang/internal/runtime"
)

// SetType represents the Set type
type SetType struct {
	ElementType runtime.Type
}

func (t *SetType) String() string {
	return fmt.Sprintf("Set<%s>", t.ElementType.String())
}

func (t *SetType) Equals(other runtime.Type) bool {
	otherSet, ok := other.(*SetType)
	if !ok {
		return false
	}
	return t.ElementType.Equals(otherSet.ElementType)
}

// SetValue represents a Set value at runtime
type SetValue struct {
	elements map[string]runtime.Value // Key is string representation of value
	Type     *SetType
}

// NewSet creates a new empty set
func NewSet(elementType runtime.Type) *SetValue {
	return &SetValue{
		elements: make(map[string]runtime.Value),
		Type:     &SetType{ElementType: elementType},
	}
}

// String returns string representation
func (s *SetValue) String() string {
	return fmt.Sprintf("Set{%d elements}", len(s.elements))
}

// TypeOf returns the type
func (s *SetValue) TypeOf() runtime.Type {
	return s.Type
}

// Equals checks equality
func (s *SetValue) Equals(other runtime.Value) bool {
	otherSet, ok := other.(*SetValue)
	if !ok {
		return false
	}
	
	if len(s.elements) != len(otherSet.elements) {
		return false
	}
	
	for key := range s.elements {
		if _, exists := otherSet.elements[key]; !exists {
			return false
		}
	}
	
	return true
}

// Add adds an element to the set
func (s *SetValue) Add(value runtime.Value) {
	key := value.String()
	s.elements[key] = value
}

// Remove removes an element from the set
func (s *SetValue) Remove(value runtime.Value) bool {
	key := value.String()
	_, exists := s.elements[key]
	if exists {
		delete(s.elements, key)
	}
	return exists
}

// Contains checks if an element exists in the set
func (s *SetValue) Contains(value runtime.Value) bool {
	key := value.String()
	_, exists := s.elements[key]
	return exists
}

// Size returns the number of elements
func (s *SetValue) Size() int {
	return len(s.elements)
}

// Clear removes all elements
func (s *SetValue) Clear() {
	s.elements = make(map[string]runtime.Value)
}

// ToArray converts the set to an array
func (s *SetValue) ToArray() []runtime.Value {
	result := make([]runtime.Value, 0, len(s.elements))
	for _, v := range s.elements {
		result = append(result, v)
	}
	return result
}

// Union returns the union of two sets
func (s *SetValue) Union(other *SetValue) *SetValue {
	result := NewSet(s.Type.ElementType)
	
	// Add all elements from first set
	for key, value := range s.elements {
		result.elements[key] = value
	}
	
	// Add all elements from second set
	for key, value := range other.elements {
		result.elements[key] = value
	}
	
	return result
}

// Intersection returns the intersection of two sets
func (s *SetValue) Intersection(other *SetValue) *SetValue {
	result := NewSet(s.Type.ElementType)
	
	// Add elements that exist in both sets
	for key, value := range s.elements {
		if _, exists := other.elements[key]; exists {
			result.elements[key] = value
		}
	}
	
	return result
}

// Difference returns the difference of two sets (elements in s but not in other)
func (s *SetValue) Difference(other *SetValue) *SetValue {
	result := NewSet(s.Type.ElementType)
	
	// Add elements that exist in s but not in other
	for key, value := range s.elements {
		if _, exists := other.elements[key]; !exists {
			result.elements[key] = value
		}
	}
	
	return result
}

// IsSubset checks if s is a subset of other
func (s *SetValue) IsSubset(other *SetValue) bool {
	for key := range s.elements {
		if _, exists := other.elements[key]; !exists {
			return false
		}
	}
	return true
}

// IsSuperset checks if s is a superset of other
func (s *SetValue) IsSuperset(other *SetValue) bool {
	return other.IsSubset(s)
}


