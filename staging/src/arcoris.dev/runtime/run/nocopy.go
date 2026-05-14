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

package run

// noCopy marks a runtime owner as unsafe to copy after first use.
//
// It is intentionally small, unexported, and behaviorless. Its only purpose is
// to make accidental value copies visible to static analysis tools such as:
//
//	go vet -copylocks
//
// The Go vet copylocks analyzer treats values with Lock and Unlock methods as
// lock-like values. Group embeds noCopy because it owns mutable runtime
// orchestration state: the group context, cancellation function, task submission
// state, wait state, task-name registry, task error collection, and WaitGroup
// accounting.
//
// Copying a Group after construction would split one logical task owner into
// independent struct values that still refer to overlapping runtime state. That
// can make cancellation ownership ambiguous, allow one copy to close task
// submission while another still appears open, lose or duplicate task errors, or
// corrupt the relationship between submitted tasks and Wait.
//
// noCopy does not provide runtime protection. It does not lock anything, does
// not allocate, and does not participate in synchronization. It is a
// static-analysis marker only.
type noCopy struct{}

// Lock is a marker method used by go vet's copylocks analyzer.
//
// Lock must never be called by production code. It intentionally has no runtime
// behavior.
func (*noCopy) Lock() {}

// Unlock is a marker method used by go vet's copylocks analyzer.
//
// Unlock must never be called by production code. It intentionally has no
// runtime behavior.
func (*noCopy) Unlock() {}
