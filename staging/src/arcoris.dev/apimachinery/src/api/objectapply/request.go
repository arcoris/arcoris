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
	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/resource"
)

// Request contains all inputs for one pure object-level apply operation.
//
// The request has no context, subject, admission result, storage handle, or
// resource catalog hook. Higher layers own those concerns and pass resolved
// data into objectapply.
type Request struct {
	// Owner is the field owner applying Applied.Desired.
	Owner fieldownership.Owner

	// Live is the current live object. Its TypeMeta, ObjectMeta, and Observed
	// surface are preserved on success.
	Live ValueObject

	// Applied is the requested object carrying Desired intent. Its Observed
	// surface must be absent and its metadata may only identify the live object.
	Applied ValueObject

	// Resource is the resolved resource family definition.
	//
	// It is expected to have been validated at construction, registration, or
	// catalog boundaries.
	Resource resource.Definition

	// Ownership is the current object-level ownership state.
	Ownership objectownership.State
}
