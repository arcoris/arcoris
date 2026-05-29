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

package identity

// ParseResource parses a canonical API resource collection name.
//
// The parser validates lexical identity only. It does not infer plural forms,
// recognize built-in resources, or split subresource paths.
func ParseResource(value string) (Resource, error) {
	resource := Resource(value)
	if err := resource.Validate(); err != nil {
		return "", err
	}
	return resource, nil
}
