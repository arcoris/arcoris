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

package objectapply

import (
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/value"
)

// buildOutputObject preserves live envelope data and installs merged Desired.
//
// The output intentionally starts from live metadata, not applied metadata.
// Desired and Observed value payloads are cloned because value.Value exposes
// explicit deep-copy semantics for composite payloads.
func buildOutputObject(live ValueObject, desired value.Value) ValueObject {
	out := object.New[value.Value, value.Value](
		live.TypeMeta,
		live.ObjectMeta,
		desired.Clone(),
	)

	if live.Observed != nil {
		// WithObserved stores a fresh pointer; Clone detaches nested value data.
		return out.WithObserved(live.Observed.Clone())
	}

	// Be explicit even though object.New returns no observed payload. This keeps
	// the output policy local to this function.
	return out.WithoutObserved()
}
