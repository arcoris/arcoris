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
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/valueapply"
)

// Result contains the object, ownership state, and Desired apply metadata.
//
// Early validation errors return a zero Result. Desired apply failures return
// the partial valueapply Result when one is available, but no output object or
// post-apply object ownership.
type Result struct {
	// Object is the resulting object. It is populated only after Desired apply
	// succeeds.
	Object ValueObject

	// Ownership is the updated object-level ownership state. It is populated only
	// after Desired apply succeeds.
	Ownership objectownership.State

	// Desired is the value-level apply result for the Desired surface. It may be
	// partially populated when Desired apply returns a conflict or merge error.
	Desired valueapply.Result
}
