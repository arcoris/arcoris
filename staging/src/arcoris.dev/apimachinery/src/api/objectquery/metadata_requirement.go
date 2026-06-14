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

// metadataRequirement is the canonical shared representation behind label and
// annotation requirements.
//
// The public API intentionally exposes separate LabelRequirement and
// AnnotationRequirement types so callers cannot accidentally mix the two
// metadata domains. Internally, both domains share this representation because
// their operator/value semantics are identical.
type metadataRequirement struct {
	key    string
	op     Operator
	values []string
}

// validatorFunc adapts label and annotation lexical validators.
type validatorFunc func(string) error
