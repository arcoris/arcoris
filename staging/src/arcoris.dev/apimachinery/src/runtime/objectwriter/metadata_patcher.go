package objectwriter

import (
	"context"

	"arcoris.dev/apimachinery/api/objectlifecycle"
)

// MetadataPatcher is the minimal write contract for components that patch
// generic labels and annotations on existing live object state.
//
// The interface stays intentionally metadata-specific. Components that only
// patch labels or annotations do not need to depend on the full Writer contract.
type MetadataPatcher interface {
	PatchMetadata(context.Context, objectlifecycle.PatchMetadataRequest) (objectlifecycle.Result, error)
}
