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

package objectstore

// NewCreatedChange validates and returns a detached created transition.
func NewCreatedChange(key Key, after State) (Change, error) {
	change := Change{
		Kind:     ChangeCreated,
		Key:      key,
		Revision: after.Revision,
		After:    after.Clone(),
	}
	if err := ValidateChange(change); err != nil {
		return Change{}, err
	}

	return change, nil
}

// MustCreatedChange validates and returns a created transition or panics.
func MustCreatedChange(key Key, after State) Change {
	change, err := NewCreatedChange(key, after)
	if err != nil {
		panic(err)
	}

	return change
}

// NewUpdatedChange validates and returns a detached updated transition.
func NewUpdatedChange(key Key, before State, after State) (Change, error) {
	change := Change{
		Kind:     ChangeUpdated,
		Key:      key,
		Revision: after.Revision,
		Before:   before.Clone(),
		After:    after.Clone(),
	}
	if err := ValidateChange(change); err != nil {
		return Change{}, err
	}

	return change, nil
}

// MustUpdatedChange validates and returns an updated transition or panics.
func MustUpdatedChange(key Key, before State, after State) Change {
	change, err := NewUpdatedChange(key, before, after)
	if err != nil {
		panic(err)
	}

	return change
}

// NewDeletedChange validates and returns a detached deleted transition.
func NewDeletedChange(key Key, before State, revision Revision) (Change, error) {
	change := Change{
		Kind:     ChangeDeleted,
		Key:      key,
		Revision: revision,
		Before:   before.Clone(),
	}
	if err := ValidateChange(change); err != nil {
		return Change{}, err
	}

	return change, nil
}

// MustDeletedChange validates and returns a deleted transition or panics.
func MustDeletedChange(key Key, before State, revision Revision) Change {
	change, err := NewDeletedChange(key, before, revision)
	if err != nil {
		panic(err)
	}

	return change
}
