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

package jsonconfig

// EmptyOwnershipSurfaceMode controls empty desired ownership surface emission.
type EmptyOwnershipSurfaceMode uint8

const (
	// EmptyOwnershipSurfaceDefault defers to the package default during Resolve.
	EmptyOwnershipSurfaceDefault EmptyOwnershipSurfaceMode = iota

	// EmptyOwnershipSurfaceEmit emits desired even when the surface is empty.
	EmptyOwnershipSurfaceEmit

	// EmptyOwnershipSurfaceOmit omits desired when the surface is empty.
	EmptyOwnershipSurfaceOmit
)

// isKnownEmptyOwnershipSurfaceMode reports whether mode is part of the public enum.
func isKnownEmptyOwnershipSurfaceMode(mode EmptyOwnershipSurfaceMode) bool {
	switch mode {
	case EmptyOwnershipSurfaceDefault, EmptyOwnershipSurfaceEmit, EmptyOwnershipSurfaceOmit:
		return true
	default:
		return false
	}
}
