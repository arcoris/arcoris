/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package signals_test

import (
	"context"
	"fmt"
	"os"

	"arcoris.dev/runtime/signals"
)

func ExampleNotifyContext() {
	ctx, stop := signals.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	select {
	case <-ctx.Done():
		_, _ = signals.Cause(ctx)
	default:
		fmt.Println("running")
	}

	// Output:
	// running
}

func ExampleSubscription() {
	sub := signals.Subscribe(os.Interrupt)
	defer sub.Stop()

	select {
	case <-sub.C():
		fmt.Println("signal received")
	default:
		fmt.Println("waiting")
	}

	// Output:
	// waiting
}

func ExampleShutdownController() {
	shutdown := signals.NewShutdownController(context.Background())
	defer shutdown.Stop()

	select {
	case <-shutdown.Done():
		fmt.Println("shutdown")
	default:
		fmt.Println("running")
	}

	// Output:
	// running
}
