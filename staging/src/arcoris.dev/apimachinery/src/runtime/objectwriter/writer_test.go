package objectwriter

import "arcoris.dev/apimachinery/api/objectlifecycle"

var _ Writer = (*objectlifecycle.Executor)(nil)
