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
	// ErrNilMux identifies a nil HTTP mux passed to a health HTTP install
	// helper.
	//
	// Install helpers require a mux because they register handlers immediately.
	// A nil mux cannot store routes and would otherwise panic when Handle is
	// called.
	ErrNilMux = errors.New("healthhttp: nil mux")
)

// Mux is the minimal HTTP route registration capability required by install
// helpers.
//
// The interface intentionally matches the common net/http.ServeMux registration
// shape and many compatible mux implementations:
//
//	mux.Handle(pattern string, handler http.Handler)
//
// healthhttp does not require a concrete *http.ServeMux. Applications may pass
// any router adapter that implements this method.
//
// Mux owns only registration. It does not define route matching semantics,
// duplicate route behavior, method dispatch, middleware, logging, metrics,
// authentication, authorization, or path normalization. Those concerns belong to
// the concrete HTTP server/router used by the application.
//
// Install helpers do not recover panics from Mux.Handle. Duplicate pattern
// handling is router-specific; for example, net/http.ServeMux panics on duplicate
// registrations. Such conflicts should remain visible to server setup code.
type Mux interface {
	Handle(pattern string, handler http.Handler)
}

// nilMux reports whether mux is nil, including typed nil implementations.
//
// A typed nil interface value can occur when a nil pointer to a concrete mux is
// assigned to Mux. Calling Handle on such a value may panic. Install helpers
// reject typed nil muxes at their boundary.
func nilMux(mux Mux) bool {
	if mux == nil {
		return true
	}

	value := reflect.ValueOf(mux)
	switch value.Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Pointer,
		reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
