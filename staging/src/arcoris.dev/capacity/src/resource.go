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

package capacity

import (
	"fmt"
	"regexp"
)

// resourcePattern is the closed grammar for stable capacity resource names.
//
// Each segment starts with a lowercase letter and may then contain lowercase
// letters, digits, or single underscores before another alphanumeric byte.
var resourcePattern = regexp.MustCompile(`^[a-z](?:[a-z0-9]|_[a-z0-9])*(?:\.[a-z](?:[a-z0-9]|_[a-z0-9])*)*$`)

// Resource is a stable accounting dimension.
//
// Resource identifies what is counted, not a runtime object instance. Valid
// values are dot-separated lower_snake_case segments such as "worker_slots" or
// "resilience.bulkhead.slots".
type Resource string

// MustResource returns value as a Resource or panics when it is invalid.
func MustResource(value string) Resource {
	resource := Resource(value)
	if !resource.IsValid() {
		panic(errorAt(
			"resource",
			ErrInvalidResource,
			ErrorReasonInvalidResource,
			fmt.Sprintf("resource %q must be dot-separated lower_snake_case", resource),
		))
	}
	return resource
}

// IsValid reports whether r follows the capacity resource grammar.
func (r Resource) IsValid() bool {
	return resourcePattern.MatchString(string(r))
}

// String returns r as its stable identifier string.
func (r Resource) String() string {
	return string(r)
}
