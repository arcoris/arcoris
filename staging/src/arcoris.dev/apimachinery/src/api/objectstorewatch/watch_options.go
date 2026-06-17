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

package objectstorewatch

// WatchOptions contains request-conversion knobs for Boundary.WatchRequest.
//
// It does not configure transports, buffering, retries, filtering, or runtime
// behavior. Those responsibilities belong to future source/reflector layers.
type WatchOptions struct {
	// AllowProgress permits, but does not require, objectwatch.EventProgress
	// markers on the generated watch request.
	AllowProgress bool
}
