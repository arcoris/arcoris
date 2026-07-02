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
	"errors"
	"testing"

	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

var _ Enqueuer = (*objectworkqueue.Queue)(nil)

func TestEmitFuncReturnsDelegateError(t *testing.T) {
	wantErr := errors.New("emit failed")
	emit := EmitFunc(func(objectworkqueue.Item) error {
		return wantErr
	})

	err := emit(objectworkqueue.Item{Key: testKey(1)})

	if err != wantErr {
		t.Fatalf("error = %v; want %v", err, wantErr)
	}
}
