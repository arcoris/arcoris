package objectwriter

import (
	"context"

	"arcoris.dev/apimachinery/api/objectlifecycle"
)

// Applier is the minimal write contract for components that reconcile desired
// object state through the api/objectlifecycle Apply operation.
//
// The interface deliberately exposes objectlifecycle.ApplyRequest and
// objectlifecycle.ApplyResult directly. That keeps field ownership, descriptor
// validation, conflict handling, and error semantics in api/objectlifecycle
// instead of inventing a second runtime write vocabulary.
type Applier interface {
	Apply(context.Context, objectlifecycle.ApplyRequest) (objectlifecycle.ApplyResult, error)
}
