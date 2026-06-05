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

package codecregistry

import "fmt"

// registrationPath returns the stable diagnostic path for one registration.
func registrationPath(index int) string {
	return fmt.Sprintf("registrations[%d]", index)
}

// registrationIDPath returns the stable diagnostic path for one entry ID.
func registrationIDPath(index int) string {
	return registrationPath(index) + ".id"
}

// codecPath returns the stable diagnostic path for one registered codec.
func codecPath(index int) string {
	return registrationPath(index) + ".codec"
}

// infoPath returns the stable diagnostic path for codec metadata.
func infoPath(index int) string {
	return registrationPath(index) + ".info"
}

// mediaTypePath returns the stable diagnostic path for one declared media type.
func mediaTypePath(index int, mediaTypeIndex int) string {
	return fmt.Sprintf("%s.mediaTypes[%d]", infoPath(index), mediaTypeIndex)
}

// capabilityPath returns the stable diagnostic path for capability checks.
func capabilityPath(index int) string {
	return registrationPath(index) + ".capabilities"
}
