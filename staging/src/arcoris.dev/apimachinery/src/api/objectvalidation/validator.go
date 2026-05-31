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

package objectvalidation

import "arcoris.dev/apimachinery/api/types"

// SurfaceValidator bridges typed payload values and structural descriptors.
//
// Implementations validate one concrete desired or observed payload value
// against the descriptor selected from a resource version. They must not mutate
// the value. api/objectvalidation does not assume payloads expose a Validate
// method because descriptor-aware conformance requires the selected types.Type.
// The resolver argument is the exact resolver supplied by Plan and may be nil.
type SurfaceValidator[T any] interface {
	ValidateSurface(value T, descriptor types.Type, resolver types.Resolver) error
}
