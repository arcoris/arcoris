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

import (
	"os"
	"reflect"
)

const (
	// errNilSignalSetSignal is the stable diagnostic text used when a public
	// signal set helper receives a nil signal.
	//
	// Signal sets describe process notification configuration. A nil signal cannot
	// be registered or compared meaningfully by this package and is rejected at
	// the value-helper boundary.
	errNilSignalSetSignal = "signals: nil signal in signal set"
)

// Clone returns an independent copy of sigs.
//
// Clone preserves nil input as nil. Non-nil input is copied even when empty.
// Signal values are interface values and are copied by assignment.
func Clone(sigs []os.Signal) []os.Signal {
	if sigs == nil {
		return nil
	}

	return append([]os.Signal(nil), sigs...)
}

// Unique returns sigs with duplicate signals removed.
//
// The first occurrence of each signal is preserved. Unique panics when any
// signal is nil.
func Unique(sigs []os.Signal) []os.Signal {
	if len(sigs) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(sigs))
	unique := make([]os.Signal, 0, len(sigs))
	for _, sig := range sigs {
		requireSignal(sig, errNilSignalSetSignal)

		key := signalKey(sig)
		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}
		unique = append(unique, sig)
	}

	return unique
}

// Merge combines signal sets into one ordered unique set.
//
// Sets are processed from left to right. The first occurrence of each signal is
// preserved. Merge panics when any signal is nil.
func Merge(sets ...[]os.Signal) []os.Signal {
	var merged []os.Signal
	for _, set := range sets {
		merged = append(merged, set...)
	}

	return Unique(merged)
}

// Contains reports whether sigs contains sig.
//
// Contains panics when sig is nil or when sigs contains a nil signal.
func Contains(sigs []os.Signal, sig os.Signal) bool {
	requireSignal(sig, errNilSignalSetSignal)

	key := signalKey(sig)
	for _, candidate := range sigs {
		requireSignal(candidate, errNilSignalSetSignal)
		if signalKey(candidate) == key {
			return true
		}
	}

	return false
}

// signalSetOrShutdown returns a normalized signal set.
//
// An empty set means the package default shutdown signal set. The returned slice
// is unique and independent from caller-owned slices.
func signalSetOrShutdown(sigs []os.Signal) []os.Signal {
	if len(sigs) == 0 {
		return Shutdown()
	}

	return Unique(sigs)
}

// signalKey returns a stable comparison key for sig.
//
// os.Signal does not require the dynamic value to be comparable. The package
// uses a string key instead of interface equality so custom test signals cannot
// panic set helpers by carrying non-comparable fields. Real process signals are
// still distinguished by their concrete type and String value.
func signalKey(sig os.Signal) string {
	return typeName(sig) + ":" + sig.String()
}

// typeName returns a compact package-local type name for signal set keys.
func typeName(value any) string {
	typeOf := reflect.TypeOf(value)
	if typeOf == nil {
		return "<nil>"
	}

	return typeOf.PkgPath() + "." + typeOf.String()
}
