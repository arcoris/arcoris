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

import (
	"fmt"
	"strings"
)

// Error is a structured objectstore diagnostic.
type Error struct {
	// Reason gives stable machine-readable detail.
	Reason Reason
	// Key is the object key involved in the error, when available.
	Key Key
	// Expected is the expected revision for optimistic concurrency errors.
	Expected Revision
	// Actual is the current observed revision for optimistic concurrency errors.
	Actual Revision
	// Err is the broad sentinel or nested cause exposed through errors.Is.
	Err error
}

// Error returns a stable human-readable diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	parts := []string{"objectstore"}
	if e.Reason.IsValid() {
		parts = append(parts, e.Reason.String())
	}
	if !e.Key.Equal(Key{}) {
		parts = append(parts, "key="+e.Key.String())
	}
	if e.Expected.IsValid() || e.Actual.IsValid() {
		parts = append(parts, fmt.Sprintf("expected=%s actual=%s", e.Expected, e.Actual))
	}
	if e.Err != nil {
		parts = append(parts, e.Err.Error())
	}

	return strings.Join(parts, ": ")
}

// Unwrap exposes the broad sentinel or nested cause for errors.Is/As.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}
