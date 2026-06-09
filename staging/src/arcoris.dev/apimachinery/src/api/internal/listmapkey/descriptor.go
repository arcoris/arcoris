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

package listmapkey

import "arcoris.dev/apimachinery/api/types"

// objectDescriptor resolves the list element descriptor to an object view.
func (r selectorRequest) objectDescriptor() (types.ObjectView, error) {
	resolvedDescriptor, err := r.resolver.resolve(r.indexPath, r.element, 0)
	if err != nil {
		return types.ObjectView{}, err
	}

	if resolvedDescriptor.Code() != types.DescriptorObject {
		return types.ObjectView{}, failure(
			r.indexPath,
			FailureInvalidDescriptor,
			"ListMap element descriptor is not an object",
		)
	}

	elementObjectView, ok := resolvedDescriptor.AsObject()
	if !ok {
		return types.ObjectView{}, failure(
			r.indexPath,
			FailureInvalidDescriptor,
			"ListMap element descriptor cannot expose object fields",
		)
	}

	return elementObjectView, nil
}

// fieldDescriptor finds one declared object field descriptor by ListMap key.
func fieldDescriptor(
	objectView types.ObjectView,
	key types.FieldName,
) (types.FieldDescriptor, bool) {
	for _, declaredField := range objectView.Fields() {
		if declaredField.Name() == key {
			return declaredField, true
		}
	}

	return types.FieldDescriptor{}, false
}
