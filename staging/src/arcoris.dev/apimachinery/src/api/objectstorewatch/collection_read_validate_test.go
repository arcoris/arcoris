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
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestCollectionReadValidate(t *testing.T) {
	tests := []struct {
		name    string
		read    CollectionRead
		wantErr error
	}{
		{name: "valid", read: testCollectionRead(t)},
		{name: "invalid collection", read: CollectionRead{}, wantErr: ErrInvalidCollectionRead},
		{
			name: "invalid result",
			read: CollectionRead{
				collection: testCollection(),
				result:     testListResult([]objectstore.ListItem{{Key: testKey("system", "main")}}, 1),
			},
			wantErr: ErrInvalidCollectionRead,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.read.Validate()
			if tt.wantErr == nil {
				requireNoError(t, err)
				return
			}
			requireErrorIs(t, err, tt.wantErr)
		})
	}
}
