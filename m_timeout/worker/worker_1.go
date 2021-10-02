package worker

import (
	"fmt"
	"time"
)

// Worker1 doc
type Worker1 struct {
	BaseWorker
}

// Do doc
func (w1 *Worker1) Do(ctx *Context) {
	err := w1.Core(ctx)
	if err != nil {
		fmt.Println("worker1 has an err: ", err)
		return
	}

	fmt.Println("worker1 has been done ")
}

// Core doc
func (w1 *Worker1) Core(ctx *Context) error {
	time.Sleep(1 * time.Second)
	select {
	case <-ctx.TimeoutCtx.Done():
		return ctx.TimeoutCtx.Err()

	case DealChan <- "worker1":
		return nil
	}
}
