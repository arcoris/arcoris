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

package objectenqueue

import (
	"arcoris.dev/apimachinery/api/objectstore"
)

// Mapper maps one committed object change to zero or more reconciliation items.
//
// Mapper owns mapping policy only. It should not consume watch streams, mutate
// caches, run reconciliation, or call queue APIs directly.
type Mapper interface {
	Map(objectstore.Change, EmitFunc) error
}

// MapperFunc adapts a function to Mapper.
type MapperFunc func(objectstore.Change, EmitFunc) error

// Map calls f with change and emit.
//
// A nil MapperFunc is treated as missing wiring and returns ErrNilMapper.
func (f MapperFunc) Map(change objectstore.Change, emit EmitFunc) error {
	if f == nil {
		return ErrNilMapper
	}

	return f(change, emit)
}
