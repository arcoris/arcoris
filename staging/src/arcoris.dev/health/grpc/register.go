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

package healthgrpc

import (
	"reflect"

	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Register installs server as the standard gRPC health service.
//
// The function is a small nil-safe wrapper around
// grpc_health_v1.RegisterHealthServer. It does not start, stop, or otherwise
// own the grpc.Server lifecycle.
func Register(registrar grpc.ServiceRegistrar, server *Server) error {
	if nilRegistrar(registrar) {
		return ErrNilRegistrar
	}
	if server == nil {
		return ErrNilServer
	}

	healthpb.RegisterHealthServer(registrar, server)
	return nil
}

// nilRegistrar recognizes nil grpc.ServiceRegistrar values, including typed nils.
//
// grpc.ServiceRegistrar is commonly implemented by pointer types such as
// *grpc.Server. Detecting typed nils keeps Register's error boundary stable
// instead of panicking inside generated registration code.
func nilRegistrar(registrar grpc.ServiceRegistrar) bool {
	if registrar == nil {
		return true
	}

	value := reflect.ValueOf(registrar)
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
