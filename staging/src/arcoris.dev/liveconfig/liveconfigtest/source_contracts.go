/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package liveconfigtest

import "arcoris.dev/snapshot"

// ControlledSource intentionally implements the same read contracts as
// liveconfig.Holder. Keeping the assertions in their own file makes the package
// boundary explicit: tests can depend on snapshot.Source interfaces instead of a
// concrete helper when they only need read behavior.
var _ snapshot.Source[Config] = (*ControlledSource[Config])(nil)
var _ snapshot.StampedSource[Config] = (*ControlledSource[Config])(nil)
var _ snapshot.RevisionSource = (*ControlledSource[Config])(nil)
