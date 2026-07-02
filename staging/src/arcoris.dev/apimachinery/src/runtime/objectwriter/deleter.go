package objectwriter

import (
	"context"

	"arcoris.dev/apimachinery/api/objectlifecycle"
)

// Deleter is the minimal write contract for components that remove committed
// live object state.
//
// DeleteRequest carries the optimistic concurrency information defined by
// api/objectlifecycle. objectwriter does not hide stale revision errors,
// reread state, or repeat a delete attempt.
type Deleter interface {
	Delete(context.Context, objectlifecycle.DeleteRequest) (objectlifecycle.Result, error)
}
