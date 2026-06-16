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

// LabelExists matches items whose metadata labels contain key.
func LabelExists(key string) (Query, error) {
	return metadataQuery(metadataLabels, OperatorExists, key)
}

// LabelDoesNotExist matches items whose metadata labels do not contain key.
func LabelDoesNotExist(key string) (Query, error) {
	return metadataQuery(metadataLabels, OperatorDoesNotExist, key)
}

// LabelEquals matches items whose metadata label equals value.
func LabelEquals(key string, value string) (Query, error) {
	return metadataQuery(metadataLabels, OperatorEquals, key, value)
}

// LabelNotEquals matches items whose metadata label is absent or differs.
func LabelNotEquals(key string, value string) (Query, error) {
	return metadataQuery(metadataLabels, OperatorNotEquals, key, value)
}

// LabelIn matches items whose metadata label value is in values.
func LabelIn(key string, values ...string) (Query, error) {
	return metadataQuery(metadataLabels, OperatorIn, key, values...)
}

// LabelNotIn matches items whose metadata label is absent or outside values.
func LabelNotIn(key string, values ...string) (Query, error) {
	return metadataQuery(metadataLabels, OperatorNotIn, key, values...)
}
