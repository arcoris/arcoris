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
	"arcoris.dev/apimachinery/api/objectstore"
)

// Extractor extracts zero or more secondary index values from one list item.
//
// Extractors own value normalization. They should emit only stable, comparable
// strings and return errors when an object cannot be indexed safely.
type Extractor interface {
	Extract(objectstore.ListItem, EmitFunc) error
}

// ExtractorFunc adapts a function to Extractor.
type ExtractorFunc func(objectstore.ListItem, EmitFunc) error

// Extract delegates to f.
//
// A nil function returns ErrNilExtractor. Delegated errors are returned
// unchanged.
func (f ExtractorFunc) Extract(item objectstore.ListItem, emit EmitFunc) error {
	if f == nil {
		return ErrNilExtractor
	}

	return f(item, emit)
}

// extractValues runs every registered definition for one object and groups the
// emitted values by index name. Callers use the complete result to preserve
// all-or-nothing Replace and ApplyChange semantics.
func extractValues(names []Name, definitions map[Name]Extractor, item objectstore.ListItem) (map[Name]valueSet, error) {
	valuesByName := make(map[Name]valueSet, len(names))
	for _, name := range names {
		values, err := extractDefinitionValues(definitions[name], item)
		if err != nil {
			return nil, err
		}
		valuesByName[name] = values
	}

	return valuesByName, nil
}

// extractDefinitionValues gives the extractor a cloned list item and records a
// set of emitted values. Duplicate values collapse here before they reach index
// membership state.
func extractDefinitionValues(extractor Extractor, item objectstore.ListItem) (valueSet, error) {
	values := make(valueSet)
	emit := func(value Value) error {
		if value == "" {
			return ErrInvalidIndex
		}
		values[value] = struct{}{}

		return nil
	}
	if err := extractor.Extract(item.Clone(), emit); err != nil {
		return nil, err
	}

	return values, nil
}
