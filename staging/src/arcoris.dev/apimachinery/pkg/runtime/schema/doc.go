/*
   Copyright 2026 The ARCORIS Authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

// Package schema defines the strict identity primitives used by ARCORIS API
// machinery.
//
// The package models API groups, versions, kinds, resources, subresources, and
// their canonical combinations. It is intentionally small and dependency-free:
// runtime, metadata, discovery, REST mapping, serializers, control-plane APIs,
// and plugin configuration layers must be able to depend on these identities
// without pulling higher-level machinery back into the foundation layer.
//
// Schema distinguishes type identity from resource identity. A
// GroupVersionKind names a concrete object schema. A GroupVersionResource names
// an addressable resource collection. Resource names never include
// subresources; a ResourcePath or GroupVersionResourcePath is used when the
// subresource segment is part of the identity.
//
// The package accepts only canonical ARCORIS forms:
//
//   - GroupVersion: "v1" or "control.arcoris.dev/v1alpha1"
//   - GroupKind: "Pod" or "WorkloadClass.control.arcoris.dev"
//   - GroupResource: "pods" or "workloadclasses.control.arcoris.dev"
//   - GroupVersionKind: "v1, Kind=Pod" or "control.arcoris.dev/v1alpha1, Kind=WorkloadClass"
//   - GroupVersionResource: "v1:pods" or "control.arcoris.dev/v1alpha1:workloadclasses"
//   - ResourcePath: "pods" or "pods/status"
//   - GroupVersionResourcePath: "v1:pods/status" or "control.arcoris.dev/v1alpha1:workloadclasses/status"
//
// This package intentionally does not implement Kubernetes legacy CLI parsing
// behavior. ARCORIS schema accepts only canonical forms.
package schema
