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

package objectreconcilewrite

import (
	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/objectlifecycle"
)

// PatchMetadata builds a lifecycle request for generic label and annotation
// patches on the current object.
//
// Expected is Current.Revision. Label and annotation maps are detached,
// including pointed string values, so caller mutation cannot alter the request.
func (c Current) PatchMetadata(
	labels map[string]*string,
	annotations map[string]*string,
	owner fieldownership.Owner,
) (objectlifecycle.PatchMetadataRequest, error) {
	if err := c.validate(); err != nil {
		return objectlifecycle.PatchMetadataRequest{}, err
	}

	return objectlifecycle.PatchMetadataRequest{
		Resource:    c.key.Resource,
		Object:      c.key.Object,
		Labels:      clonePatchMap(labels),
		Annotations: clonePatchMap(annotations),
		Owner:       owner,
		Expected:    c.Revision(),
	}, nil
}

func clonePatchMap(in map[string]*string) map[string]*string {
	if in == nil {
		return nil
	}

	out := make(map[string]*string, len(in))
	for key, value := range in {
		if value == nil {
			out[key] = nil
			continue
		}
		cloned := *value
		out[key] = &cloned
	}

	return out
}
