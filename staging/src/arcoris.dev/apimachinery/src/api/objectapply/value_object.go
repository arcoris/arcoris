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
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/value"
)

// ValueObject is the objectapply v1 envelope shape.
//
// Generic typed object apply is intentionally outside this package. Typed apply
// needs conversion adapters between typed payloads and api/value.Value, plus
// mutation and deep-copy contracts for arbitrary user types.
type ValueObject = object.Object[value.Value, value.Value]
