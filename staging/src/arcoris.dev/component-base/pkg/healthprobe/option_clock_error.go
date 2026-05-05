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

package healthprobe

import (
	"errors"
	"reflect"

	"arcoris.dev/component-base/pkg/clock"
)

var (
	// ErrNilClock identifies a nil clock passed to WithClock or left in Runner
	// configuration.
	//
	// Runner owns a ticker loop and uses clock reads for cache update timestamps
	// and staleness checks. A nil clock would panic during runtime probing.
	ErrNilClock = errors.New("healthprobe: nil clock")
)

// nilClock reports whether clk is nil, including typed nil values stored in the
// clock.Clock interface.
func nilClock(clk clock.Clock) bool {
	if clk == nil {
		return true
	}

	value := reflect.ValueOf(clk)
	switch value.Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Ptr,
		reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
