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

package admissioncatalog

import "fmt"

// formatPathError builds concise diagnostics while keeping structured details
// available through typed errors.
func formatPathError(kind, path, detail string) string {
	if path != "" && detail != "" {
		return fmt.Sprintf("admissioncatalog: %s at %s: %s", kind, path, detail)
	}
	if path != "" {
		return fmt.Sprintf("admissioncatalog: %s at %s", kind, path)
	}
	if detail != "" {
		return fmt.Sprintf("admissioncatalog: %s: %s", kind, detail)
	}
	return "admissioncatalog: " + kind
}

// descriptorPath returns the stable diagnostic location for a descriptor inside
// a declaration slice.
func descriptorPath(prefix string, index int) string {
	return fmt.Sprintf("%s[%d]", prefix, index)
}
