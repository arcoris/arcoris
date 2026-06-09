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

package valueapply

import (
	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// Request describes one pure value-level apply operation.
type Request struct {
	// Path is the semantic base path for this value surface.
	//
	// Use fieldpath.RootPath() for root value apply.
	Path fieldpath.Path

	// Owner is the field owner applying Applied.
	Owner fieldownership.Owner

	// Live is the current value.
	Live value.Value

	// Applied is the new applied value.
	Applied value.Value

	// Descriptor describes both Live and Applied.
	Descriptor types.Descriptor

	// Ownership is the current field ownership state for this value surface.
	Ownership fieldownership.State
}
