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

import "arcoris.dev/apimachinery/api/objectstore"

// ListItemMapper maps one live listed object to zero or more reconciliation items.
//
// ListItemMapper is used at replace/relist boundaries where the source provides
// current live state rather than a committed change. It owns mapping policy
// only and should not call queue APIs directly.
type ListItemMapper interface {
	MapListItem(objectstore.ListItem, EmitFunc) error
}

// ListItemMapperFunc adapts a function to ListItemMapper.
type ListItemMapperFunc func(objectstore.ListItem, EmitFunc) error

// MapListItem calls f with item and emit.
//
// A nil ListItemMapperFunc is treated as missing wiring and returns
// ErrNilListItemMapper.
func (f ListItemMapperFunc) MapListItem(item objectstore.ListItem, emit EmitFunc) error {
	if f == nil {
		return ErrNilListItemMapper
	}

	return f(item, emit)
}
