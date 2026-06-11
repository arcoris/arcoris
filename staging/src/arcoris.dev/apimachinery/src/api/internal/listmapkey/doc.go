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

// Package listmapkey extracts stable semantic keys from concrete API values
// under api/types ListMap descriptor semantics.
//
// The package is intentionally narrow. Its job is deriving fieldpath.Selector
// values for list items whose stable identity comes from one or more declared
// object fields. In api/types this descriptor shape is ListMap.
//
// ListSet scalar uniqueness is intentionally outside this package today:
// api/valuevalidation owns concrete set duplicate checks, while compare, merge,
// and apply currently treat ListSet as a whole-list semantic field.
//
// The package does not validate complete values, extract field sets, compare
// values, apply changes, manage ownership, normalize payloads, or interpret API
// objects/resources. Callers wrap key extraction failures into their own public
// diagnostic models.
package listmapkey
