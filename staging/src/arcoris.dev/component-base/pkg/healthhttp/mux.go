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

package healthhttp

import (
	"errors"
	"net/http"
	"reflect"
)

var (
	// ErrNilMux identifies a nil HTTP mux passed to a health HTTP install helper.
	ErrNilMux = errors.New("healthhttp: nil mux")
)

// Mux is the minimal HTTP route registration capability required by install
// helpers.
type Mux interface {
	Handle(pattern string, handler http.Handler)
}

// nilMux reports whether mux is nil, including typed nil implementations.
func nilMux(mux Mux) bool {
	if mux == nil {
		return true
	}

	value := reflect.ValueOf(mux)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
