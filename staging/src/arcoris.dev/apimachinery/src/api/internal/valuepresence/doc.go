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

// Package valuepresence distinguishes absent traversal results from present
// concrete API values.
//
// The package is a small internal primitive for descriptor-aware value
// algorithms. It models the result of looking up a field, map key, list index,
// or associative-list selector where "not found" must stay different from a
// present null, scalar zero, or empty composite value.
package valuepresence
