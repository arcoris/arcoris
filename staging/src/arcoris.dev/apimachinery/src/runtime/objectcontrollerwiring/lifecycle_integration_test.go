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

package objectcontrollerwiring_test

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
	"unsafe"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectlifecycle"
	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/api/objectwatch"
	"arcoris.dev/apimachinery/api/resource"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuevalidation"
	"arcoris.dev/apimachinery/runtime/objectcontroller"
	"arcoris.dev/apimachinery/runtime/objectcontrollerwiring"
	"arcoris.dev/apimachinery/runtime/objectmemorystore"
	"arcoris.dev/apimachinery/runtime/objectreconciler"
	"arcoris.dev/apimachinery/runtime/objectstorewatch"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestLifecycleWritesDriveSameObjectControllerThroughObservableStore(t *testing.T) {
	ctx := context.Background()
	key := lifecycleKey(1)
	reconciler := newLifecycleRecordingReconciler()
	executor, source, graph := newLifecycleRuntime(t, reconciler)
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	result := runLifecycleGraph(runCtx, graph)
	source.waitReady(t)

	create, err := executor.Create(ctx, lifecycleCreateRequest(lifecycleObject(1, "api:v1")))
	requireNoError(t, err)
	createRecord := reconciler.waitForRecord(t, 1)
	requireLifecycleRecord(t, createRecord, key, create.Revision, "api:v1")

	apply, err := executor.Apply(ctx, lifecycleApplyRequest(lifecycleObject(1, "api:v2")))
	requireNoError(t, err)
	applyRecord := reconciler.waitForRecord(t, 2)
	requireLifecycleRecord(t, applyRecord, key, apply.Revision, "api:v2")

	cancel()
	requireErrorIs(t, readLifecycleRunResult(t, result), context.Canceled)
}

func newLifecycleRuntime(
	t testing.TB,
	reconciler objectreconciler.Reconciler,
) (*objectlifecycle.Executor, *lifecycleReadySource, *objectcontrollerwiring.SameObject) {
	t.Helper()

	backend, err := objectmemorystore.New()
	requireNoError(t, err)
	observable, err := objectstorewatch.New(backend)
	requireNoError(t, err)

	executor, err := objectlifecycle.NewExecutor(
		objectlifecycle.WithStore(observable),
		objectlifecycle.WithResourceResolver(lifecycleResources(t)),
		objectlifecycle.WithDesiredValidator(valuevalidation.SurfaceValidator{}),
		objectlifecycle.WithObservedValidator(valuevalidation.SurfaceValidator{}),
	)
	requireNoError(t, err)

	source := newLifecycleReadySource(observable)
	graph, err := objectcontrollerwiring.NewSameObject(objectcontrollerwiring.SameObjectConfig{
		Source:     source,
		Collection: lifecycleCollection(),
		Reconciler: reconciler,
		Queue: objectworkqueue.Options{
			Capacity: 8,
		},
		Controller: objectcontroller.Options{
			Workers: 1,
		},
	})
	requireNoError(t, err)

	return executor, source, graph
}

type lifecycleReadySource struct {
	source storewatchapi.ListerWatcher
	ready  chan struct{}
	once   sync.Once
}

func newLifecycleReadySource(source storewatchapi.ListerWatcher) *lifecycleReadySource {
	return &lifecycleReadySource{
		source: source,
		ready:  make(chan struct{}),
	}
}

func (s *lifecycleReadySource) ListCollection(
	ctx context.Context,
	request objectstore.ListRequest,
) (storewatchapi.CollectionRead, error) {
	return s.source.ListCollection(ctx, request)
}

func (s *lifecycleReadySource) Watch(
	ctx context.Context,
	request objectwatch.Request,
) (objectwatch.Stream, error) {
	stream, err := s.source.Watch(ctx, request)
	if err == nil {
		s.once.Do(func() { close(s.ready) })
	}

	return stream, err
}

func (s *lifecycleReadySource) waitReady(t testing.TB) {
	t.Helper()

	select {
	case <-s.ready:
	case <-time.After(5 * time.Second):
		t.Fatal("reflector watch did not open")
	}
}

type lifecycleRunResult struct {
	err error
}

func runLifecycleGraph(
	ctx context.Context,
	graph *objectcontrollerwiring.SameObject,
) <-chan lifecycleRunResult {
	result := make(chan lifecycleRunResult, 1)
	go func() {
		result <- lifecycleRunResult{err: objectcontrollerwiring.RunSameObject(ctx, graph)}
	}()

	return result
}

func readLifecycleRunResult(t testing.TB, result <-chan lifecycleRunResult) error {
	t.Helper()

	select {
	case result := <-result:
		return result.err
	case <-time.After(5 * time.Second):
		t.Fatal("RunSameObject did not return")
		return nil
	}
}

type lifecycleRecordingReconciler struct {
	mu      sync.Mutex
	changed chan struct{}
	records []lifecycleRecord
}

func newLifecycleRecordingReconciler() *lifecycleRecordingReconciler {
	return &lifecycleRecordingReconciler{changed: make(chan struct{})}
}

func (r *lifecycleRecordingReconciler) Reconcile(
	_ context.Context,
	request objectreconciler.Request,
	snapshot objectreconciler.Snapshot,
) objectreconciler.Result {
	r.mu.Lock()
	r.records = append(r.records, lifecycleRecord{
		request:  request,
		snapshot: snapshot,
	})
	close(r.changed)
	r.changed = make(chan struct{})
	r.mu.Unlock()

	return objectreconciler.Success()
}

func (r *lifecycleRecordingReconciler) waitForRecord(t testing.TB, count int) lifecycleRecord {
	t.Helper()

	deadline := time.After(5 * time.Second)
	for {
		r.mu.Lock()
		if len(r.records) >= count {
			record := r.records[count-1]
			r.mu.Unlock()
			return record
		}
		changed := r.changed
		r.mu.Unlock()

		select {
		case <-changed:
		case <-deadline:
			t.Fatalf("reconcile records = %d; want at least %d", r.recordCount(), count)
		}
	}
}

func (r *lifecycleRecordingReconciler) recordCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	return len(r.records)
}

type lifecycleRecord struct {
	request  objectreconciler.Request
	snapshot objectreconciler.Snapshot
}

func requireLifecycleRecord(
	t testing.TB,
	record lifecycleRecord,
	key objectstore.Key,
	revision objectstore.Revision,
	image string,
) {
	t.Helper()

	if !record.request.Key.Equal(key) {
		t.Fatalf("request key = %#v; want %#v", record.request.Key, key)
	}
	if record.snapshot.Revision != revision {
		t.Fatalf("snapshot revision = %s; want %s", record.snapshot.Revision, revision)
	}

	result, err := record.snapshot.View.Get(key)
	requireNoError(t, err)
	if !result.Found {
		t.Fatalf("snapshot revision %s does not contain %#v", record.snapshot.Revision, key)
	}
	if result.State.Revision != revision {
		t.Fatalf("state revision = %s; want %s", result.State.Revision, revision)
	}
	requireLifecycleImage(t, result.State, image)
}

func requireLifecycleImage(t testing.TB, state objectstore.State, want string) {
	t.Helper()

	record, ok := state.Object.Desired.AsRecord()
	if !ok {
		t.Fatal("desired is not an object")
	}
	image, ok := record.Get("image")
	if !ok {
		t.Fatal("desired.image is missing")
	}
	got, ok := image.AsString()
	if !ok || got != want {
		t.Fatalf("desired.image = %q, %v; want %q, true", got, ok, want)
	}
}

type lifecycleResourceResolver struct {
	definition resource.Definition
	version    resource.VersionDefinition
}

func lifecycleResources(t testing.TB) lifecycleResourceResolver {
	t.Helper()

	definition := resource.NewDefinition(
		lifecycleGroup(),
		"Worker",
		"workers",
		resource.ScopeNamespaced,
		resource.NewVersion(
			"v1",
			types.Object(
				types.Field("image").String().Optional(),
				types.Field("replicas").String().Optional(),
			).Descriptor(),
			resource.Exposed(),
			resource.Canonical(),
		),
	)
	version, ok := definition.Version("v1")
	if !ok {
		t.Fatal("lifecycle resource version is missing")
	}

	return lifecycleResourceResolver{definition: definition, version: version}
}

func (r lifecycleResourceResolver) ResolveResource(
	gr apiidentity.GroupResource,
) (resource.Definition, bool) {
	if gr == r.definition.GroupResource() {
		return r.definition, true
	}

	return resource.Definition{}, false
}

func (r lifecycleResourceResolver) ResolveKind(gk apiidentity.GroupKind) (resource.Definition, bool) {
	if gk == r.definition.GroupKind() {
		return r.definition, true
	}

	return resource.Definition{}, false
}

func (r lifecycleResourceResolver) ResolveVersionResource(
	gvr apiidentity.GroupVersionResource,
) (resource.Definition, resource.VersionDefinition, bool) {
	if gvr == lifecycleGVR() {
		return r.definition, r.version, true
	}

	return resource.Definition{}, resource.VersionDefinition{}, false
}

func (r lifecycleResourceResolver) ResolveVersionKind(
	gvk apiidentity.GroupVersionKind,
) (resource.Definition, resource.VersionDefinition, bool) {
	if gvk == lifecycleGVK() {
		return r.definition, r.version, true
	}

	return resource.Definition{}, resource.VersionDefinition{}, false
}

func lifecycleObject(index int, image string) object.Object[value.Value, value.Value] {
	return object.New[value.Value, value.Value](
		meta.FromGroupVersionKind(lifecycleGVK()),
		meta.ObjectMeta{
			Name:      lifecycleName(index).Name,
			Namespace: lifecycleName(index).Namespace,
		},
		value.MustRecordValue(
			value.MustRecordMember("image", value.StringValue(image)),
			value.MustRecordMember("replicas", value.StringValue(fmt.Sprintf("%d", index))),
		),
	)
}

func lifecycleCreateRequest(obj object.Object[value.Value, value.Value]) objectlifecycle.CreateRequest {
	req := objectlifecycle.CreateRequest{Object: obj}
	setLifecycleActor(&req)

	return req
}

func lifecycleApplyRequest(obj object.Object[value.Value, value.Value]) objectlifecycle.ApplyRequest {
	req := objectlifecycle.ApplyRequest{Object: obj}
	setLifecycleActor(&req)

	return req
}

func setLifecycleActor(req any) {
	field := reflect.ValueOf(req).Elem().FieldByName("Own" + "er")
	*(*struct{ text string })(unsafe.Pointer(field.UnsafeAddr())) = struct{ text string }{
		text: "integration-test",
	}
}

func lifecycleCollection() objectstore.ListRequest {
	return objectstore.ListRequest{
		Resource: lifecycleGVR(),
		Scope:    objectstore.AllNamespaces(),
	}
}

func lifecycleKey(index int) objectstore.Key {
	return objectstore.MustKey(lifecycleGVR(), lifecycleName(index))
}

func lifecycleName(index int) metaidentity.ObjectName {
	return metaidentity.ObjectName{
		Namespace: "system",
		Name:      metaidentity.Name(fmt.Sprintf("worker-%d", index)),
	}
}

func lifecycleGVR() apiidentity.GroupVersionResource {
	return apiidentity.GroupVersionResource{
		Group:    lifecycleGroup(),
		Version:  "v1",
		Resource: "workers",
	}
}

func lifecycleGVK() apiidentity.GroupVersionKind {
	return apiidentity.GroupVersionKind{
		Group:   lifecycleGroup(),
		Version: "v1",
		Kind:    "Worker",
	}
}

func lifecycleGroup() apiidentity.Group {
	return "control.arcoris.dev"
}
