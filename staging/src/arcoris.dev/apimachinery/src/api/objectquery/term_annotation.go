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

// AnnotationExists matches items whose metadata annotations contain key.
func AnnotationExists(key string) (Query, error) {
	return metadataQuery(metadataAnnotations, OperatorExists, key)
}

// AnnotationDoesNotExist matches items whose metadata annotations do not contain key.
func AnnotationDoesNotExist(key string) (Query, error) {
	return metadataQuery(metadataAnnotations, OperatorDoesNotExist, key)
}

// AnnotationEquals matches items whose metadata annotation equals value.
func AnnotationEquals(key string, value string) (Query, error) {
	return metadataQuery(metadataAnnotations, OperatorEquals, key, value)
}

// AnnotationNotEquals matches items whose metadata annotation is absent or differs.
func AnnotationNotEquals(key string, value string) (Query, error) {
	return metadataQuery(metadataAnnotations, OperatorNotEquals, key, value)
}

// AnnotationIn matches items whose metadata annotation value is in values.
func AnnotationIn(key string, values ...string) (Query, error) {
	return metadataQuery(metadataAnnotations, OperatorIn, key, values...)
}

// AnnotationNotIn matches items whose metadata annotation is absent or outside values.
func AnnotationNotIn(key string, values ...string) (Query, error) {
	return metadataQuery(metadataAnnotations, OperatorNotIn, key, values...)
}
