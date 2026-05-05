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
	"errors"
	"testing"

	"arcoris.dev/component-base/pkg/health"
	"google.golang.org/grpc"
)

// fakeRegistrar records generated gRPC registration calls.
type fakeRegistrar struct {
	// desc is the registered service descriptor.
	desc *grpc.ServiceDesc

	// impl is the registered service implementation.
	impl any
}

// RegisterService records desc and impl for assertions.
func (r *fakeRegistrar) RegisterService(desc *grpc.ServiceDesc, impl any) {
	r.desc = desc
	r.impl = impl
}

func TestRegisterRejectsNilRegistrar(t *testing.T) {
	t.Parallel()

	server := mustNewServer(t, staticSource{status: health.StatusHealthy})
	err := Register(nil, server)
	if !errors.Is(err, ErrNilRegistrar) {
		t.Fatalf("Register(nil) = %v, want ErrNilRegistrar", err)
	}
}

func TestRegisterRejectsNilServer(t *testing.T) {
	t.Parallel()

	err := Register(&fakeRegistrar{}, nil)
	if !errors.Is(err, ErrNilServer) {
		t.Fatalf("Register(nil server) = %v, want ErrNilServer", err)
	}
}

func TestRegisterInstallsHealthService(t *testing.T) {
	t.Parallel()

	registrar := &fakeRegistrar{}
	server := mustNewServer(t, staticSource{status: health.StatusHealthy})

	if err := Register(registrar, server); err != nil {
		t.Fatalf("Register() = %v, want nil", err)
	}
	if registrar.desc == nil {
		t.Fatal("registered desc is nil")
	}
	if registrar.desc.ServiceName != "grpc.health.v1.Health" {
		t.Fatalf("ServiceName = %q, want grpc.health.v1.Health", registrar.desc.ServiceName)
	}
	if registrar.impl != server {
		t.Fatal("registered implementation is not server")
	}
}
