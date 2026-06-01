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

package diagnostic

// Record stores the common structured diagnostic shape shared by API packages.
//
// R is the package-local reason type. Keeping the reason generic lets each
// public package expose its own typed reason constants while sharing one
// implementation for the repeated path/sentinel/detail/cause fields.
type Record[R ~string] struct {
	// Path identifies the domain-specific location that failed.
	Path string

	// Err is the broad sentinel used with errors.Is.
	Err error

	// Reason gives stable machine-readable detail within Err.
	Reason R

	// Detail gives human-readable diagnostic context.
	Detail string

	// Cause preserves a nested lower-level failure.
	Cause error
}
