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

package valueapply

// Apply applies req.Applied to req.Live under descriptor and ownership semantics.
//
// Apply is side-effect free: it does not mutate req.Live, req.Applied, or
// req.Ownership. Successful calls return the merged value and replacement
// ownership state. Conflict errors return pre-force metadata, merge errors
// return pre-merge metadata, and ownership-update errors return the merged
// value without a replacement ownership state.
func Apply(req Request, opts Options) (Result, error) {
	return New(opts).Apply(req)
}
