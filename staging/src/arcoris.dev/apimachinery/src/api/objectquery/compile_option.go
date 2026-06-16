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

// CompileOption configures query compilation.
type CompileOption func(*compileOptions)

// compileOptions is the private resolved option set used throughout the
// compiler. Keeping it private lets new options be added without exposing a
// mutable configuration struct.
type compileOptions struct {
	// fields resolves selectable field definitions for field terms.
	fields SelectableFieldSet
}

// WithSelectableFields supplies the registry used to validate field terms.
func WithSelectableFields(fields SelectableFieldSet) CompileOption {
	return func(opts *compileOptions) {
		opts.fields = fields
	}
}

// newCompileOptions applies options defensively, ignoring nil option
// functions so callers can assemble option slices conditionally.
func newCompileOptions(options []CompileOption) compileOptions {
	var opts compileOptions
	for _, option := range options {
		if option != nil {
			option(&opts)
		}
	}

	return opts
}
