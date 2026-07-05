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

package objectindex

// Name identifies one secondary index inside an Index.
//
// Names are stable caller-chosen identifiers such as "spec.workerName" or
// "metadata.label.app". The package treats names as opaque strings and only
// rejects the empty name.
type Name string

// Value identifies one extractor-emitted value inside a named secondary index.
//
// Values are opaque normalized strings chosen by the extractor. Empty values are
// rejected because they usually mean "not indexed" and are easy to emit by
// mistake.
type Value string

// EmitFunc emits one secondary index value for the current object.
type EmitFunc func(Value) error
