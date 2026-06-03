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

package objectapply

// Apply applies req.Applied.Desired to req.Live.Desired.
//
// Apply is side-effect free: it does not mutate req.Live, req.Applied, or
// req.Ownership. Successful calls preserve live metadata and Observed data,
// return the merged object, and return replacement object ownership state.
//
// The function does not perform admission, authorization, storage access,
// resource lookup, metadata apply, observed apply, defaulting, pruning, or API
// version conversion. Callers must provide a resolved resource.Definition.
func Apply(req Request, opts Options) (Result, error) {
	return newApplier(opts).apply(req)
}
