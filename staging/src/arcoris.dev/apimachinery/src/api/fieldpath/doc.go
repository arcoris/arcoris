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

// Package fieldpath defines semantic paths for ARCORIS API payload locations.
//
// A Path identifies a location after concrete payload data has been interpreted
// by descriptor-aware layers. Field elements address fixed object fields, key
// elements address dynamic map entries, index elements address ordered list
// items, and selector elements address associative list entries by stable key
// fields.
//
// For example:
//
//	$.spec.replicas
//	$.metadata.labels["app"]
//	$.containers[0].image
//	$.conditions[{"type":"Ready"}].status
//
// Paths are structured values, not diagnostic strings. This lets validation,
// diff, apply, and future managed-field layers compare, sort, group, and own
// paths without reparsing formatted text.
//
// Package fieldpath does not inspect api/value payloads, read api/types
// descriptors, decode wire formats, validate values, compare payloads, apply
// changes, or manage field ownership. Descriptor-aware callers decide which
// element kind to append while interpreting values through resource schemas.
package fieldpath
