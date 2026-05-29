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

package resourcecatalog

import (
	"fmt"

	"arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/resource"
)

// groupResourceOf extracts the primary resource-family key.
func groupResourceOf(def resource.Definition) identity.GroupResource {
	return def.GroupResource()
}

// groupKindOf extracts the version-independent kind key.
func groupKindOf(def resource.Definition) identity.GroupKind {
	return def.GroupKind()
}

// versionResourceKeys extracts concrete version/resource keys in declaration
// order.
func versionResourceKeys(def resource.Definition) []identity.GroupVersionResource {
	versions := def.Versions()
	if len(versions) == 0 {
		return nil
	}

	keys := make([]identity.GroupVersionResource, 0, len(versions))
	for _, version := range versions {
		keys = append(keys, identity.GroupVersionResource{
			Group:    def.Group(),
			Version:  version.Version(),
			Resource: def.Resource(),
		})
	}
	return keys
}

// versionKindKeys extracts concrete version/kind keys in declaration order.
func versionKindKeys(def resource.Definition) []identity.GroupVersionKind {
	versions := def.Versions()
	if len(versions) == 0 {
		return nil
	}

	keys := make([]identity.GroupVersionKind, 0, len(versions))
	for _, version := range versions {
		keys = append(keys, identity.GroupVersionKind{
			Group:   def.Group(),
			Version: version.Version(),
			Kind:    def.Kind(),
		})
	}
	return keys
}

// definitionPath returns a stable path for an incoming definition index.
func definitionPath(index int) string {
	return fmt.Sprintf("definitions[%d]", index)
}

// resourcePath returns a stable path for a GroupResource identity.
func resourcePath(gr identity.GroupResource) string {
	return fmt.Sprintf("definitions[%s]", gr)
}

// kindPath returns a stable path for a GroupKind identity.
func kindPath(gk identity.GroupKind) string {
	return fmt.Sprintf("definitions[%s]", gk)
}

// versionResourcePath returns a stable path for a GroupVersionResource identity.
func versionResourcePath(gvr identity.GroupVersionResource) string {
	return fmt.Sprintf("definitions[%s]", gvr)
}

// versionKindPath returns a stable path for a GroupVersionKind identity.
func versionKindPath(gvk identity.GroupVersionKind) string {
	return fmt.Sprintf("definitions[%s]", gvk)
}
