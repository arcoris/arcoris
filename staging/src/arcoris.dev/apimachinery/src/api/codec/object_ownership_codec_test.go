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

import "arcoris.dev/apimachinery/api/objectownership"

var _ ObjectOwnershipDecoder = fakeObjectOwnershipCodec{}
var _ ObjectOwnershipEncoder = fakeObjectOwnershipCodec{}

// ObjectOwnershipCodec is a partial byte capability for ownership documents.
var _ ObjectOwnershipCodec = fakeObjectOwnershipCodec{}
var _ ObjectOwnershipCodec = fakeByteCodec{}

type fakeObjectOwnershipCodec struct {
	fakeBaseCodec
}

func (fakeObjectOwnershipCodec) DecodeObjectOwnership([]byte) (objectownership.Document, error) {
	return objectownership.Document{}, nil
}

func (fakeObjectOwnershipCodec) EncodeObjectOwnership(objectownership.Document) ([]byte, error) {
	return nil, nil
}

func (fakeByteCodec) DecodeObjectOwnership(data []byte) (objectownership.Document, error) {
	return fakeObjectOwnershipCodec{}.DecodeObjectOwnership(data)
}

func (fakeByteCodec) EncodeObjectOwnership(doc objectownership.Document) ([]byte, error) {
	return fakeObjectOwnershipCodec{}.EncodeObjectOwnership(doc)
}
