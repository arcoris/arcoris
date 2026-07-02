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
	"context"
	"errors"
	"testing"

	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestHandlerEmitterReturnsFirstQueueError(t *testing.T) {
	wantErr := errors.New("queue failed")
	queue := &recordingQueue{err: wantErr}
	emitter := handlerEmitter{ctx: context.Background(), queue: queue}

	err := emitter.emit(objectworkqueue.Item{Key: testKey(1)})
	if err != wantErr {
		t.Fatalf("first emit error = %v; want %v", err, wantErr)
	}

	err = emitter.emit(objectworkqueue.Item{Key: testKey(2)})
	if err != wantErr {
		t.Fatalf("second emit error = %v; want %v", err, wantErr)
	}
	queue.requireItems(t, objectworkqueue.Item{Key: testKey(1)})
}

func TestHandlerEmitterResultPrefersQueueError(t *testing.T) {
	queueErr := errors.New("queue failed")
	mapErr := errors.New("mapper failed")
	emitter := handlerEmitter{err: queueErr}

	err := emitter.result(mapErr)

	if err != queueErr {
		t.Fatalf("error = %v; want %v", err, queueErr)
	}
}

func TestHandlerEmitterResultReturnsMapperError(t *testing.T) {
	mapErr := errors.New("mapper failed")
	emitter := handlerEmitter{}

	err := emitter.result(mapErr)

	if err != mapErr {
		t.Fatalf("error = %v; want %v", err, mapErr)
	}
}
