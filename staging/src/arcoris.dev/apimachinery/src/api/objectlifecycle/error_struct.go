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

package objectlifecycle

import (
	"strings"

	"arcoris.dev/apimachinery/api/objectstore"
)

// Error is a structured lifecycle diagnostic.
type Error struct {
	// Operation is the operation that failed, when available.
	Operation Operation

	// Reason is stable machine-readable lifecycle detail.
	Reason ErrorReason

	// Key is the objectstore key involved in the failure, when available.
	Key objectstore.Key

	// Err stores the lifecycle sentinel and nested cause for errors.Is/As.
	Err error
}

// Error returns a stable human-readable lifecycle diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	parts := []string{"objectlifecycle"}
	if e.Operation.IsValid() {
		parts = append(parts, e.Operation.String())
	}
	if e.Reason.IsValid() {
		parts = append(parts, e.Reason.String())
	}
	if !e.Key.Equal(objectstore.Key{}) {
		parts = append(parts, "key="+e.Key.String())
	}
	if e.Err != nil {
		parts = append(parts, e.Err.Error())
	}

	return strings.Join(parts, ": ")
}

// Unwrap exposes lifecycle sentinels and nested lower-layer causes.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}
