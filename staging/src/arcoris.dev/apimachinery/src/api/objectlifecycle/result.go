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
	"arcoris.dev/apimachinery/api/objectapply"
	"arcoris.dev/apimachinery/api/objectstore"
)

// Result is the common successful lifecycle operation result.
//
// State is populated with committed store state for found, created, updated,
// and deleted effects. Delete returns the deleted live state in State and the
// tombstone commit revision in Revision.
type Result struct {
	// Operation is the lifecycle operation that succeeded.
	Operation Operation

	// Effect is the externally visible effect of the successful operation.
	Effect Effect

	// State is the committed or deleted live object state associated with Effect.
	State objectstore.State

	// Revision is the store-local revision associated with this result.
	//
	// For Get, Create, and Apply it is the live State revision. For Delete it
	// is the tombstone commit revision while State.Revision remains the deleted
	// live revision.
	Revision objectstore.Revision
}

// IsValid reports whether r has a known operation and effect.
func (r Result) IsValid() bool {
	return r.Operation.IsValid() && r.Effect.IsValid()
}

// ApplyResult is the successful result of Apply.
//
// Apply is populated only when Apply updated an existing live object through
// objectapply. Missing-object apply uses the create path and leaves Apply zero.
type ApplyResult struct {
	// Result is the common lifecycle result.
	Result

	// Apply is the pure objectapply metadata for existing-object apply.
	Apply objectapply.Result
}
