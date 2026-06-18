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

package objectstorewatch

import "errors"

var (
	// ErrNilBackend reports construction without an objectstore.Store backend.
	ErrNilBackend = errors.New("nil object store backend")
	// ErrInvalidOption reports malformed Store construction options.
	ErrInvalidOption = errors.New("invalid object store watch option")
	// ErrStreamOverflow reports that a stream queue filled before it could
	// consume a matching event.
	ErrStreamOverflow = errors.New("object store watch stream overflow")
)
