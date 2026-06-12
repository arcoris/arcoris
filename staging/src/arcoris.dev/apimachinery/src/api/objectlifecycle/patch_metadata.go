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

package objectlifecycle

import (
	"context"

	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/meta"
	"arcoris.dev/apimachinery/api/meta/annotations"
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/value"
)

// PatchMetadata patches labels and annotations on existing live state.
func (e *Executor) PatchMetadata(ctx context.Context, req PatchMetadataRequest) (Result, error) {
	if err := e.requireExecutor(OperationPatchMetadata); err != nil {
		return Result{}, err
	}
	if err := checkContext(OperationPatchMetadata, ctx); err != nil {
		return Result{}, err
	}
	if err := e.validatePatchMetadataRequest(req); err != nil {
		return Result{}, err
	}

	prepared, err := e.prepareKeyRequest(OperationPatchMetadata, req.Resource, req.Object)
	if err != nil {
		return Result{}, err
	}
	if err := validateExpectedRevision(OperationPatchMetadata, prepared.key, req.Expected); err != nil {
		return Result{}, err
	}

	labelFields, err := labelPatchFields(prepared.key, req.Labels)
	if err != nil {
		return Result{}, err
	}
	annotationFields, err := annotationPatchFields(prepared.key, req.Annotations)
	if err != nil {
		return Result{}, err
	}

	live, ok, err := e.store.Get(ctx, prepared.key)
	if err != nil {
		return Result{}, mapStoreError(OperationPatchMetadata, prepared.key, err)
	}
	if !ok {
		return Result{}, errorFor(OperationPatchMetadata, ErrorReasonNotFound, prepared.key, ErrNotFound, nil)
	}

	ownership, err := stateOwnership(OperationPatchMetadata, prepared.key, live.Ownership)
	if err != nil {
		return Result{}, err
	}
	nextOwnership, err := patchMetadataOwnership(prepared.key, ownership, req.Owner, labelFields, annotationFields)
	if err != nil {
		return Result{}, err
	}

	nextMeta := patchObjectMetadata(live.Object.ObjectMeta, req.Labels, req.Annotations)
	nextObject := object.New[value.Value, value.Value](live.Object.TypeMeta, nextMeta, live.Object.Desired.Clone())
	if live.Object.Observed != nil {
		nextObject = nextObject.WithObserved(live.Object.Observed.Clone())
	}
	next := objectstore.State{
		Object:    nextObject,
		Ownership: nextOwnership,
	}

	committed, err := e.store.Update(ctx, prepared.key, req.Expected, next)
	if err != nil {
		return Result{}, mapStoreError(OperationPatchMetadata, prepared.key, err)
	}

	return Result{
		Operation: OperationPatchMetadata,
		Effect:    EffectUpdated,
		State:     committed,
		Revision:  committed.Revision,
	}, nil
}

func patchObjectMetadata(
	current meta.ObjectMeta,
	labelPatch map[string]*string,
	annotationPatch map[string]*string,
) meta.ObjectMeta {
	out := current.Clone()
	for key, value := range labelPatch {
		labelKey := labels.Key(key)
		if value == nil {
			delete(out.Labels, labelKey)
			continue
		}
		if out.Labels == nil {
			out.Labels = labels.Set{}
		}
		out.Labels[labelKey] = labels.Value(*value)
	}
	for key, value := range annotationPatch {
		annotationKey := annotations.Key(key)
		if value == nil {
			delete(out.Annotations, annotationKey)
			continue
		}
		if out.Annotations == nil {
			out.Annotations = annotations.Set{}
		}
		out.Annotations[annotationKey] = annotations.Value(*value)
	}

	return out
}

func patchMetadataOwnership(
	key objectstore.Key,
	current objectownership.State,
	owner fieldownership.Owner,
	labelFields fieldpath.Set,
	annotationFields fieldpath.Set,
) (objectownership.State, error) {
	metadata := current.Metadata()
	if !labelFields.IsEmpty() {
		labelsState, err := metadata.Labels().AddFields(owner, labelFields)
		if err != nil {
			return objectownership.State{}, errorFor(OperationPatchMetadata, ErrorReasonMetadataOwnershipFailed, key, ErrApplyFailed, err)
		}
		metadata = metadata.WithLabels(labelsState)
	}
	if !annotationFields.IsEmpty() {
		annotationsState, err := metadata.Annotations().AddFields(owner, annotationFields)
		if err != nil {
			return objectownership.State{}, errorFor(OperationPatchMetadata, ErrorReasonMetadataOwnershipFailed, key, ErrApplyFailed, err)
		}
		metadata = metadata.WithAnnotations(annotationsState)
	}

	return current.WithMetadata(metadata), nil
}

func labelPatchFields(key objectstore.Key, patch map[string]*string) (fieldpath.Set, error) {
	paths := make([]fieldpath.Path, 0, len(patch))
	for rawKey, rawValue := range patch {
		if err := labels.Key(rawKey).ValidateLexical(); err != nil {
			return fieldpath.Set{}, errorFor(OperationPatchMetadata, ErrorReasonInvalidMetadataKey, key, ErrInvalidRequest, err)
		}
		if rawValue != nil {
			if err := labels.Value(*rawValue).ValidateLexical(); err != nil {
				return fieldpath.Set{}, errorFor(OperationPatchMetadata, ErrorReasonInvalidMetadataPatch, key, ErrInvalidRequest, err)
			}
		}
		path, err := metadataKeyPath(rawKey)
		if err != nil {
			return fieldpath.Set{}, errorFor(OperationPatchMetadata, ErrorReasonInvalidMetadataKey, key, ErrInvalidRequest, err)
		}
		paths = append(paths, path)
	}

	return fieldpath.NewSet(paths...)
}

func annotationPatchFields(key objectstore.Key, patch map[string]*string) (fieldpath.Set, error) {
	paths := make([]fieldpath.Path, 0, len(patch))
	for rawKey, rawValue := range patch {
		if err := annotations.Key(rawKey).ValidateLexical(); err != nil {
			return fieldpath.Set{}, errorFor(OperationPatchMetadata, ErrorReasonInvalidMetadataKey, key, ErrInvalidRequest, err)
		}
		if rawValue != nil {
			if err := annotations.Value(*rawValue).ValidateLexical(); err != nil {
				return fieldpath.Set{}, errorFor(OperationPatchMetadata, ErrorReasonInvalidMetadataPatch, key, ErrInvalidRequest, err)
			}
		}
		path, err := metadataKeyPath(rawKey)
		if err != nil {
			return fieldpath.Set{}, errorFor(OperationPatchMetadata, ErrorReasonInvalidMetadataKey, key, ErrInvalidRequest, err)
		}
		paths = append(paths, path)
	}

	return fieldpath.NewSet(paths...)
}

func metadataKeyPath(key string) (fieldpath.Path, error) {
	mapKey, err := fieldpath.NewMapKey(key)
	if err != nil {
		return fieldpath.Path{}, err
	}

	return fieldpath.Root().Key(mapKey), nil
}
