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

package codec

import (
	"io"

	"arcoris.dev/apimachinery/api/objectownership"
)

var _ ObjectOwnershipStreamDecoder = fakeObjectOwnershipStreamCodec{}
var _ ObjectOwnershipStreamEncoder = fakeObjectOwnershipStreamCodec{}

// ObjectOwnershipStreamCodec is a partial streaming capability for ownership state.
var _ ObjectOwnershipStreamCodec = fakeObjectOwnershipStreamCodec{}

type fakeObjectOwnershipStreamCodec struct {
	fakeBaseCodec
}

func (fakeObjectOwnershipStreamCodec) DecodeObjectOwnershipFrom(io.Reader) (objectownership.State, error) {
	return objectownership.State{}, nil
}

func (fakeObjectOwnershipStreamCodec) EncodeObjectOwnershipTo(io.Writer, objectownership.State) error {
	return nil
}
