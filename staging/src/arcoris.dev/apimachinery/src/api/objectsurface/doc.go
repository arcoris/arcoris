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

// Package objectsurface defines the stable object surface taxonomy shared by
// ownership, lifecycle, codecs, and tests.
//
// Desired is declarative user or manager intent. Observed is runtime,
// controller, or agent state. Metadata labels and annotations are key/value
// metadata maps with their own mutation rules.
//
// TypeMeta is not an ownable surface. ObjectMeta identity and system fields,
// including name, namespace, uid, resourceVersion, generation, createdAt, and
// deletion, are not generic ownable surfaces. Finalizers and ownerReferences
// are named here as reserved metadata surfaces for lifecycle/governance layers;
// current ownership state does not model them as generic patch surfaces.
package objectsurface
