package worker

import (
	"fmt"
	"time"
)

// Worker2 doc
type Worker2 struct {
	BaseWorker
}

// Do doc
func (w2 *Worker2) Do(ctx *Context) {
	err := w2.Core(ctx)
	if err != nil {
		fmt.Println("worker2 has an err: ", err)
		return
	}

	fmt.Println("worker2 has been done ")
}

// Core doc
func (w2 *Worker2) Core(ctx *Context) error {
	time.Sleep(1 * time.Second)
	select {
	case <-ctx.TimeoutCtx.Done():
		return ctx.TimeoutCtx.Err()

	case DealChan <- "worker2":
		return nil
	}
}
