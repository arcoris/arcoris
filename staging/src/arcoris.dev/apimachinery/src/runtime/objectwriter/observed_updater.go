package objectwriter

import (
	"context"

	"arcoris.dev/apimachinery/api/objectlifecycle"
)

// ObservedUpdater is the minimal write contract for components that replace the
// observed surface of existing live object state.
//
// UpdateObserved is intentionally separate from Applier. Desired-state writes
// and observed-state writes have different ownership in api/objectlifecycle,
// and callers should depend only on the side they need.
type ObservedUpdater interface {
	UpdateObserved(context.Context, objectlifecycle.UpdateObservedRequest) (objectlifecycle.Result, error)
}
