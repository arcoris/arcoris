package objectwriter

import "arcoris.dev/apimachinery/api/objectlifecycle"

var _ ObservedUpdater = (*objectlifecycle.Executor)(nil)
