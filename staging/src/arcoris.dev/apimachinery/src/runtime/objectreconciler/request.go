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

package objectreconciler

import "arcoris.dev/apimachinery/api/objectstore"

// Request identifies the object being reconciled.
//
// Request is the reconciliation boundary form of a queue work item. It carries
// object identity without exposing queue-level state or completion semantics to
// user reconciliation logic.
type Request struct {
	// Key identifies the object being reconciled.
	Key objectstore.Key
}

// Validate checks whether r contains a structurally valid object key.
func (r Request) Validate() error {
	if err := objectstore.ValidateKey(r.Key); err != nil {
		return invalidRequestError(err)
	}

	return nil
}

// IsValid reports whether r passes Validate.
func (r Request) IsValid() bool {
	return r.Validate() == nil
}
