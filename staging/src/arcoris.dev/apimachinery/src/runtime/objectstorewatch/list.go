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

import (
	"context"

	"arcoris.dev/apimachinery/api/objectstore"
)

// List delegates raw structural objectstore listing under the observable lock.
//
// List is the plain objectstore.Store operation. It does not create an
// api/objectstorewatch CollectionRead; callers that need a watch boundary
// should use ListCollection instead.
func (s *Store) List(ctx context.Context, request objectstore.ListRequest) (objectstore.ListResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.backend.List(ctx, request)
}
