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
	"errors"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/internal/listmapkey"
	"arcoris.dev/apimachinery/api/internal/valuepresence"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// compareListMap compares items by selector identity, ignoring physical order.
//
// Selector paths are the only successful output form for ListMap items. If any
// item cannot yield a selector, comparison fails instead of falling back to
// unstable physical index comparison.
func (c *comparer) compareListMap(
	path fieldpath.Path,
	oldList value.ListView,
	newList value.ListView,
	element types.Type,
	keys []types.FieldName,
	depth int,
) (Result, error) {
	oldEntries, err := c.listMapEntries(path, oldList, element, keys)
	if err != nil {
		return Result{}, err
	}
	newEntries, err := c.listMapEntries(path, newList, element, keys)
	if err != nil {
		return Result{}, err
	}

	result := EmptyResult()
	for _, key := range unionSortedListMapKeys(oldEntries, newEntries) {
		oldEntry, oldFound := oldEntries[key]
		newEntry, newFound := newEntries[key]
		selector := oldEntry.selector
		if !oldFound {
			selector = newEntry.selector
		}

		child, err := c.compare(
			path.Select(selector),
			listMapOperand(oldEntry, oldFound),
			listMapOperand(newEntry, newFound),
			element,
			depth+1,
		)
		if err != nil {
			return Result{}, err
		}

		result = result.merge(child)
	}

	return result, nil
}

// equalListMap compares ListMap items by selector identity.
//
// Physical order is ignored, matching compareListMap. Iteration uses sorted
// selector strings so traversal and error reporting stay deterministic.
func (c *comparer) equalListMap(
	path fieldpath.Path,
	oldList value.ListView,
	newList value.ListView,
	element types.Type,
	keys []types.FieldName,
	depth int,
) (bool, error) {
	oldEntries, err := c.listMapEntries(path, oldList, element, keys)
	if err != nil {
		return false, err
	}
	newEntries, err := c.listMapEntries(path, newList, element, keys)
	if err != nil {
		return false, err
	}
	if len(oldEntries) != len(newEntries) {
		return false, nil
	}

	for _, key := range sortedMapKeys(oldEntries) {
		oldEntry := oldEntries[key]
		newEntry, found := newEntries[key]
		if !found {
			return false, nil
		}

		equal, err := c.equalValue(path.Select(oldEntry.selector), oldEntry.item, newEntry.item, element, depth+1)
		if err != nil {
			return false, err
		}
		if !equal {
			return false, nil
		}
	}

	return true, nil
}

// listMapEntry stores one concrete ListMap item after selector extraction.
type listMapEntry struct {
	// selector is the stable semantic identity for this item.
	selector fieldpath.Selector
	// indexPath preserves the physical occurrence path for duplicate diagnostics.
	indexPath fieldpath.Path
	// item is the original concrete list item.
	item value.Value
}

// listMapEntries indexes ListMap items by canonical selector string.
//
// The stored key is selector.String() only for deterministic map lookup. The
// original selector is retained for semantic path construction.
func (c *comparer) listMapEntries(
	path fieldpath.Path,
	list value.ListView,
	element types.Type,
	keys []types.FieldName,
) (map[string]listMapEntry, error) {
	if len(keys) == 0 {
		return nil, errorAt(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "list map has no key fields")
	}

	n := list.Len()
	if n == 0 {
		return nil, nil
	}

	entries := make(map[string]listMapEntry, n)
	for i := 0; i < n; i++ {
		entry, err := c.listMapEntryAt(path, list, i, element, keys)
		if err != nil {
			return nil, err
		}

		key := entry.selector.String()
		if previous, exists := entries[key]; exists {
			return nil, duplicateListMapEntryError(path.Select(entry.selector), previous.indexPath, entry.indexPath)
		}

		entries[key] = entry
	}

	return entries, nil
}

// listMapEntryAt extracts selector identity for one physical list item.
//
// indexPath is used only for diagnostics while extracting the selector. Once a
// selector exists, successful comparison paths use path.Select(selector).
func (c *comparer) listMapEntryAt(
	path fieldpath.Path,
	list value.ListView,
	index int,
	element types.Type,
	keys []types.FieldName,
) (listMapEntry, error) {
	item, _ := list.At(index)
	indexPath := path.Index(index)
	selector, err := c.extractListMapSelector(indexPath, item, element, keys)
	if err != nil {
		return listMapEntry{}, err
	}

	return listMapEntry{
		selector:  selector,
		indexPath: indexPath,
		item:      item,
	}, nil
}

// listMapOperand converts a selector lookup into presence-aware compare input.
func listMapOperand(entry listMapEntry, present bool) valuepresence.Operand {
	return valuepresence.From(entry.item, present)
}

// unionSortedListMapKeys returns deterministic traversal order for selector keys.
func unionSortedListMapKeys(left, right map[string]listMapEntry) []string {
	return unionSortedMapKeys(left, right)
}

// extractListMapSelector delegates stable selector construction to listmapkey.
//
// The shared helper owns key-field traversal and TypeRef-aware literal
// extraction. This package only maps helper errors into valuecompare errors.
func (c *comparer) extractListMapSelector(
	indexPath fieldpath.Path,
	item value.Value,
	element types.Type,
	keys []types.FieldName,
) (fieldpath.Selector, error) {
	selector, err := listmapkey.ExtractSelector(
		indexPath,
		item,
		element,
		keys,
		listmapkey.Options{Resolver: c.resolver, MaxDepth: c.maxDepth},
	)
	if err != nil {
		return fieldpath.Selector{}, compareListMapKeyError(err)
	}

	return selector, nil
}

// duplicateListMapEntryError reports repeated selector identity with both indexes.
func duplicateListMapEntryError(selectorPath fieldpath.Path, firstPath fieldpath.Path, duplicatePath fieldpath.Path) error {
	return errorfAt(
		selectorPath,
		ErrDuplicateListKey,
		ErrorReasonDuplicateListKey,
		"duplicate list map key; first occurrence at %s, duplicate at %s",
		firstPath,
		duplicatePath,
	)
}

// compareListMapKeyError maps shared ListMap key failures to compare diagnostics.
//
// valuecompare keeps its own public sentinels while preserving the helper's path,
// detail, and nested cause.
func compareListMapKeyError(err error) error {
	keyError, ok := listMapKeyFailure(err)
	if !ok {
		return err
	}

	sentinel, reason := compareListMapKeyErrorKind(keyError.Kind)
	if keyError.Cause != nil {
		return wrapAt(keyError.Path, sentinel, reason, keyError.Detail, keyError.Cause)
	}

	return errorAt(keyError.Path, sentinel, reason, keyError.Detail)
}

// listMapKeyFailure extracts the shared ListMap key diagnostic when present.
func listMapKeyFailure(err error) (*listmapkey.Error, bool) {
	var keyError *listmapkey.Error
	if errors.As(err, &keyError) {
		return keyError, true
	}

	return nil, false
}

// compareListMapKeyErrorKind maps internal selector failure kinds to compare reasons.
func compareListMapKeyErrorKind(kind listmapkey.FailureKind) (error, ErrorReason) {
	switch kind {
	case listmapkey.FailureInvalidDescriptor:
		return ErrInvalidDescriptor, ErrorReasonInvalidDescriptor
	case listmapkey.FailureUnresolvedRef:
		return ErrUnresolvedRef, ErrorReasonUnresolvedRef
	case listmapkey.FailureReferenceCycle:
		return ErrReferenceCycle, ErrorReasonReferenceCycle
	case listmapkey.FailureMissingKey:
		return ErrInvalidListKey, ErrorReasonMissingListKey
	case listmapkey.FailureNullKey,
		listmapkey.FailureKeyKindMismatch,
		listmapkey.FailureKeyIntegerRange,
		listmapkey.FailureInvalidSelector,
		listmapkey.FailureItemKindMismatch:
		return ErrInvalidListKey, ErrorReasonInvalidListKey
	default:
		return ErrInvalidListKey, ErrorReasonInvalidListKey
	}
}
