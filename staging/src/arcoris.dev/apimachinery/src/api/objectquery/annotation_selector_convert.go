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

package objectquery

// annotationRequirementsFromMetadata adapts shared metadata requirements to
// the public annotation requirement wrapper.
func annotationRequirementsFromMetadata(requirements []metadataRequirement) []AnnotationRequirement {
	if len(requirements) == 0 {
		return nil
	}

	out := make([]AnnotationRequirement, len(requirements))
	for i, req := range requirements {
		out[i] = AnnotationRequirement{req: req.clone()}
	}

	return out
}
