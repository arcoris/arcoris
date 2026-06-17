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

import "arcoris.dev/apimachinery/api/objectwatch"

// WatchCapabilities returns the wrapper's supported watch contract features.
//
// The wrapper supports historical starts from retained in-memory history and
// StartAtCurrent by observing the current collection under the Store lock.
// Progress events are not emitted by this passive implementation, and terminal
// continuity failures are returned as errors instead of restart-required events.
func (s *Store) WatchCapabilities() objectwatch.Capabilities {
	return objectwatch.Capabilities{
		StartAtCurrent:  true,
		HistoricalStart: true,
		Progress:        false,
		RestartEvents:   false,
	}
}
