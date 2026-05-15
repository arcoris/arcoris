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

package liveconfig

import "arcoris.dev/snapshot"

// contractValue anchors Holder interface assertions without exposing a public
// test-only type.
//
// These assertions keep the read-side API aligned with package snapshot. Holder
// is intentionally a source, not a publisher; consumers can read snapshots and
// revisions, but they cannot publish around Apply.
type contractValue struct{}

var _ snapshot.Source[contractValue] = (*Holder[contractValue])(nil)
var _ snapshot.StampedSource[contractValue] = (*Holder[contractValue])(nil)
var _ snapshot.RevisionSource = (*Holder[contractValue])(nil)
