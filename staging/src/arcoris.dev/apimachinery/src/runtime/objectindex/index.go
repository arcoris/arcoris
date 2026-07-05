// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package objectindex

import (
	"reflect"
	"sync"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/runtime/objectreflector"
)

var _ objectreflector.Sink = (*Index)(nil)

// Index maintains named secondary indexes for one reflected collection.
//
// Index is safe for concurrent Lookup, Replace, and ApplyChange calls. Lookup
// observes either the state before a successful write operation or the state
// after it; failed write operations leave existing memberships unchanged.
type Index struct {
	mu sync.RWMutex

	names       []Name
	definitions map[Name]Extractor
	state       indexState
}

// New validates definitions and returns an empty Index.
func New(definitions ...Definition) (*Index, error) {
	if len(definitions) == 0 {
		return nil, ErrInvalidDefinition
	}

	names := make([]Name, 0, len(definitions))
	byName := make(map[Name]Extractor, len(definitions))
	for _, definition := range definitions {
		if definition.Name == "" {
			return nil, ErrInvalidDefinition
		}
		if _, exists := byName[definition.Name]; exists {
			return nil, ErrInvalidDefinition
		}
		if isNilExtractor(definition.Extractor) {
			return nil, errorWith(ErrInvalidDefinition, ErrNilExtractor)
		}

		names = append(names, definition.Name)
		byName[definition.Name] = definition.Extractor
	}

	return &Index{
		names:       names,
		definitions: byName,
		state:       newIndexState(names),
	}, nil
}

// validateStatic checks immutable construction-time fields that all operations
// rely on before reading or publishing mutable state.
func (i *Index) validateStatic() error {
	if i == nil || len(i.names) == 0 || len(i.definitions) == 0 {
		return ErrInvalidIndex
	}
	for _, name := range i.names {
		if name == "" || isNilExtractor(i.definitions[name]) {
			return ErrInvalidIndex
		}
	}

	return nil
}

// validateStateLocked checks mutable maps while the caller holds i.mu.
func (i *Index) validateStateLocked() error {
	if i.state.byName == nil || i.state.byObject == nil {
		return ErrInvalidIndex
	}
	for _, name := range i.names {
		if i.state.byName[name] == nil {
			return ErrInvalidIndex
		}
	}

	return nil
}

// isNilExtractor rejects both nil interfaces and typed-nil implementations
// accepted through the Extractor interface.
func isNilExtractor(extractor Extractor) bool {
	if extractor == nil {
		return true
	}

	value := reflect.ValueOf(extractor)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

// indexState stores both reverse lookup state and the per-object memberships
// needed to update or remove one object without scanning every value bucket.
type indexState struct {
	byName   map[Name]map[Value]*keySet
	byObject map[objectstore.Key]map[Name]valueSet
}

// newIndexState allocates empty buckets for every known index name so lookup
// and mutation paths can treat missing names as corruption instead of absence.
func newIndexState(names []Name) indexState {
	byName := make(map[Name]map[Value]*keySet, len(names))
	for _, name := range names {
		byName[name] = make(map[Value]*keySet)
	}

	return indexState{
		byName:   byName,
		byObject: make(map[objectstore.Key]map[Name]valueSet),
	}
}

// addObject installs the memberships extracted for one object during a Replace
// rebuild. Objects that emit no values are intentionally omitted from byObject.
func (s indexState) addObject(key objectstore.Key, valuesByName map[Name]valueSet) {
	hasValues := false
	for name, values := range valuesByName {
		if len(values) == 0 {
			continue
		}
		hasValues = true
		for value := range values {
			s.addMembership(key, name, value)
		}
	}
	if hasValues {
		s.byObject[key] = cloneValuesByName(valuesByName)
	}
}

// removeObject deletes all memberships currently recorded for one object key.
func (s indexState) removeObject(key objectstore.Key) {
	valuesByName := s.byObject[key]
	for name, values := range valuesByName {
		for value := range values {
			s.removeMembership(key, name, value)
		}
	}
	delete(s.byObject, key)
}

// updateObject applies a complete replacement membership set for one object.
// Values that remain unchanged keep their existing lookup order.
func (s indexState) updateObject(key objectstore.Key, next map[Name]valueSet) {
	current := s.byObject[key]
	for name, currentValues := range current {
		nextValues := next[name]
		for value := range currentValues {
			if _, keep := nextValues[value]; !keep {
				s.removeMembership(key, name, value)
			}
		}
	}

	hasValues := false
	for name, nextValues := range next {
		currentValues := current[name]
		for value := range nextValues {
			hasValues = true
			if _, exists := currentValues[value]; !exists {
				s.addMembership(key, name, value)
			}
		}
	}
	if hasValues {
		s.byObject[key] = cloneValuesByName(next)
		return
	}

	delete(s.byObject, key)
}

// addMembership records one name/value -> key relationship.
func (s indexState) addMembership(key objectstore.Key, name Name, value Value) {
	values := s.byName[name]
	keys := values[value]
	if keys == nil {
		keys = newKeySet()
		values[value] = keys
	}
	keys.add(key)
}

// removeMembership removes one name/value -> key relationship and drops the
// value bucket once it becomes empty.
func (s indexState) removeMembership(key objectstore.Key, name Name, value Value) {
	values := s.byName[name]
	keys := values[value]
	if keys == nil {
		return
	}
	keys.remove(key)
	if keys.len() == 0 {
		delete(values, value)
	}
}

// keySet keeps deterministic lookup order while providing duplicate collapse.
type keySet struct {
	order []objectstore.Key
	seen  map[objectstore.Key]struct{}
}

// newKeySet creates an empty ordered set for one index value.
func newKeySet() *keySet {
	return &keySet{seen: make(map[objectstore.Key]struct{})}
}

// add inserts key only once and preserves first-seen order.
func (s *keySet) add(key objectstore.Key) {
	if _, exists := s.seen[key]; exists {
		return
	}
	s.seen[key] = struct{}{}
	s.order = append(s.order, key)
}

// remove deletes key from both membership and ordering state.
func (s *keySet) remove(key objectstore.Key) {
	if _, exists := s.seen[key]; !exists {
		return
	}
	delete(s.seen, key)
	for i, existing := range s.order {
		if existing.Equal(key) {
			copy(s.order[i:], s.order[i+1:])
			s.order = s.order[:len(s.order)-1]
			return
		}
	}
}

// keys returns a detached, ordered snapshot for Lookup callers.
func (s *keySet) keys() []objectstore.Key {
	return append([]objectstore.Key(nil), s.order...)
}

// len reports the number of keys currently in the set.
func (s *keySet) len() int {
	if s == nil {
		return 0
	}

	return len(s.order)
}

// valueSet records the values emitted by one extractor for one object.
type valueSet map[Value]struct{}

// cloneValuesByName detaches per-object membership state before storing it in
// byObject, so later map mutations cannot alter index state.
func cloneValuesByName(valuesByName map[Name]valueSet) map[Name]valueSet {
	clone := make(map[Name]valueSet, len(valuesByName))
	for name, values := range valuesByName {
		if len(values) == 0 {
			continue
		}
		clone[name] = values.clone()
	}

	return clone
}

// clone returns a detached copy of one extractor value set.
func (s valueSet) clone() valueSet {
	clone := make(valueSet, len(s))
	for value := range s {
		clone[value] = struct{}{}
	}

	return clone
}
