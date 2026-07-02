package objectwriter

import (
	"context"

	"arcoris.dev/apimachinery/api/objectlifecycle"
)

// Creator is the minimal write contract for components that create committed
// live object state.
//
// Create uses the request and result vocabulary owned by api/objectlifecycle.
// objectwriter adds no allocation policy, request cloning, id generation,
// defaulting, or additional validation layer around that operation.
type Creator interface {
	Create(context.Context, objectlifecycle.CreateRequest) (objectlifecycle.Result, error)
}
