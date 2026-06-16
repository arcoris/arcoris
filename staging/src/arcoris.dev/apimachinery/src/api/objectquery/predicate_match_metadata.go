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

import (
	"slices"

	"arcoris.dev/apimachinery/api/meta/annotations"
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/objectstore"
)

// matchMetadataTerm evaluates label and annotation map terms. Negative
// operators intentionally match absent keys according to objectquery semantics.
func matchMetadataTerm(t term, item objectstore.ListItem) bool {
	actual, present := metadataValue(t, item)
	switch t.operator {
	case OperatorExists:
		return present
	case OperatorDoesNotExist:
		return !present
	case OperatorEquals:
		return present && actual == t.stringValues[0]
	case OperatorNotEquals:
		return !present || actual != t.stringValues[0]
	case OperatorIn:
		return present && metadataValueIn(t.stringValues, actual)
	case OperatorNotIn:
		return !present || !metadataValueIn(t.stringValues, actual)
	default:
		return false
	}
}

// metadataValue returns the stored metadata value and whether the key is
// present. Present empty strings are still present values.
func metadataValue(t term, item objectstore.ListItem) (string, bool) {
	switch t.metadataDomain {
	case metadataLabels:
		value, ok := item.State.Object.ObjectMeta.Labels[labels.Key(t.metadataKey)]
		return value.String(), ok
	case metadataAnnotations:
		value, ok := item.State.Object.ObjectMeta.Annotations[annotations.Key(t.metadataKey)]
		return value.String(), ok
	default:
		return "", false
	}
}

// metadataValueIn searches canonical sorted metadata literals.
func metadataValueIn(values []string, actual string) bool {
	_, found := slices.BinarySearch(values, actual)
	return found
}
