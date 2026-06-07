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

package health

import (
	"errors"
	"reflect"
)

// ErrNilChecker identifies a nil health checker.
//
// Nil checkers cannot provide stable names or produce health results. This is a
// root checker-contract error, not a registry-specific condition. Registries,
// gates, evaluators, and test helpers may wrap or preserve this sentinel when
// they reject or defensively observe nil checker values.
var ErrNilChecker = errors.New("health: nil checker")

// nilChecker reports whether checker is nil, including typed nil interface
// values.
func nilChecker(checker Checker) bool {
	if checker == nil {
		return true
	}

	val := reflect.ValueOf(checker)
	switch val.Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Pointer,
		reflect.Slice:
		return val.IsNil()
	default:
		return false
	}
}
