package worker // Worker2 doc
import (
	"fmt"
	"time"
)

// Worker3 doc
type Worker3 struct {
	BaseWorker
}

// Do doc
func (w3 *Worker3) Do(ctx *Context) {
	err := w3.Core(ctx)
	if err != nil {
		fmt.Println("worker3 has an err: ", err)
		return
	}

	fmt.Println("worker3 has been done ")
}

// Core doc
func (w3 *Worker3) Core(ctx *Context) error {
	time.Sleep(1 * time.Second)
	select {
	case <-ctx.TimeoutCtx.Done():
		return ctx.TimeoutCtx.Err()

	case DealChan <- "worker3":
		return nil
	}
}
