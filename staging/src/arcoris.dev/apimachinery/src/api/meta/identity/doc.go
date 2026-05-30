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

// Package identity defines object metadata identity primitives.
//
// The package owns names, namespace markers, UIDs, object names, object
// identities, and object references. It does not define storage keys, route
// keys, cache keys, resource collection keys, or REST paths.
//
// API group, version, kind, and resource identity remains in
// arcoris.dev/apimachinery/api/identity. ObjectReference composes those API
// identity values with metadata identity fields when a metadata reference must
// describe a concrete API object.
//
// Empty Namespace means namespace absence. It does not mean a default namespace,
// and this package never applies defaulting.
package identity
