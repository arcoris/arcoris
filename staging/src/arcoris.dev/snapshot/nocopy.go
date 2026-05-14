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

package snapshot

// noCopy marks holder values as non-copyable for go vet's copylocks checker.
//
// Store contains a sync.RWMutex and Publisher contains atomic state. Copying
// either holder after first use would split synchronization ownership and can
// produce subtle data races or stale publications. noCopy is unexported because
// it is an implementation detail, not a public helper type.
type noCopy struct{}

// Lock is a marker method recognized by go vet's copylocks analyzer.
func (*noCopy) Lock() {}

// Unlock is a marker method recognized by go vet's copylocks analyzer.
func (*noCopy) Unlock() {}
