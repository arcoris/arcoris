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

package signals

// noCopy marks a struct as unsafe to copy after first use.
//
// It is intentionally small, unexported, and behaviorless. Its only purpose is
// to make accidental value copies visible to static analysis tools such as:
//
//	go vet -copylocks
//
// The Go vet copylocks analyzer treats values with Lock and Unlock methods as
// lock-like values. When Timer contains noCopy, copying a Timer value after
// first use is reported similarly to copying a mutex.
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
