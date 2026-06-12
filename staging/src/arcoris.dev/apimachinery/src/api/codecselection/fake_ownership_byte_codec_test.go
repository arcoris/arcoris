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

package codecselection

import (
	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/codecregistry"
	"arcoris.dev/apimachinery/api/objectownership"
)

func testOwnershipByteRegistration(id string, mediaTypes ...codec.MediaType) codecregistry.Registration {
	return testRegistration(id, newOwnershipByteCodec(codec.FormatJSON, mediaTypes...))
}

type fakeOwnershipByteCodec struct {
	fakeBaseCodec
}

func newOwnershipByteCodec(format codec.Format, mediaTypes ...codec.MediaType) *fakeOwnershipByteCodec {
	return &fakeOwnershipByteCodec{fakeBaseCodec: fakeBaseCodec{
		info: testInfo(format, mediaTypes, codec.TargetObjectOwnership),
	}}
}

func (f fakeOwnershipByteCodec) DecodeObjectOwnership(
	[]byte,
) (objectownership.State, error) {
	return objectownership.State{}, nil
}

func (f fakeOwnershipByteCodec) EncodeObjectOwnership(
	objectownership.State,
) ([]byte, error) {
	return nil, nil
}
