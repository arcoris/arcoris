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

package valuecompare

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuefieldset"
)

// addSubtree records every semantic path introduced by a newly present value.
func (c *comparer) addSubtree(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Descriptor,
	result Result,
) (Result, error) {
	set, err := c.extractSubtree(path, val, descriptor)
	if err != nil {
		return Result{}, err
	}

	return result.withAdded(set), nil
}

// removeSubtree records every semantic path lost with a removed value.
func (c *comparer) removeSubtree(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Descriptor,
	result Result,
) (Result, error) {
	set, err := c.extractSubtree(path, val, descriptor)
	if err != nil {
		return Result{}, err
	}

	return result.withRemoved(set), nil
}

// extractSubtree delegates added/removed subtree path discovery to valuefieldset.
//
// valuecompare deliberately does not duplicate field-set traversal. That keeps
// ListMap selectors, empty composites, atomic lists, and unknown-field policy
// aligned across both packages.
func (c *comparer) extractSubtree(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Descriptor,
) (fieldpath.Set, error) {
	set, err := valuefieldset.ExtractAt(
		path,
		val,
		descriptor,
		valuefieldset.Options{Resolver: c.resolver, MaxDepth: c.maxDepth},
	)
	if err != nil {
		return fieldpath.Set{}, compareFieldSetError(path, err)
	}

	return set, nil
}
