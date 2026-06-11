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

package valuevalidation

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// Validator stores reusable validation options.
//
// Validator is immutable by convention. Each Validate or ValidateAt call creates
// a fresh validation run, so per-run mutable state such as reference stacks,
// pattern caches, and collected diagnostics is never shared across calls.
type Validator struct {
	opts Options
}

// New returns a reusable Validator configured with opts.
func New(opts Options) Validator {
	return Validator{opts: opts}
}

// Validate checks val against descriptor starting at the semantic root path.
func (v Validator) Validate(val value.Value, descriptor types.Descriptor) error {
	return v.ValidateAt(fieldpath.Root(), val, descriptor)
}

// ValidateAt checks val against descriptor starting at path.
//
// ValidateAt validates only path's structural well-formedness before starting
// descriptor-aware value validation. Descriptor preparation remains an upstream
// responsibility.
func (v Validator) ValidateAt(path fieldpath.Path, val value.Value, descriptor types.Descriptor) error {
	if err := path.ValidateStructure(); err != nil {
		return ErrorList{wrapAt(
			path,
			ErrInvalidPath,
			ErrorReasonInvalidPath,
			"base field path is invalid",
			err,
		)}
	}

	run := newValidator(v.opts)
	run.validate(path, val, descriptor, 0)

	return run.result()
}
