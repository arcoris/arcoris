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

package lifecycle

import "context"

const (
	errNilContext = "lifecycle: nil context"
	errNilOption  = "lifecycle: nil option"
)

// requireContext enforces the runtime-wide nil-context policy for public APIs.
func requireContext(ctx context.Context) {
	if ctx == nil {
		panic(errNilContext)
	}
}

// requireOption enforces strict option composition for controller construction.
func requireOption(opt Option) {
	if opt == nil {
		panic(errNilOption)
	}
}
