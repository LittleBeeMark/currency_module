package worker

import (
	"fmt"
	"sync"
	"time"
)

var LastHandleChan chan interface{}

func init() {
	LastHandleChan = make(chan interface{})
}

// Receiver doc
type Receiver struct {
	BaseWorker
	Result interface{}
}

// NewReceiver new receiver with a res
func NewReceiver(result interface{}) *Receiver {
	return &Receiver{
		Result: result,
	}
}

// Do doc
func (r *Receiver) Do(ctx *Context) {
	err := r.Core(ctx)
	if err != nil {
		fmt.Println("receiver has an err: ", err)
		return
	}

	fmt.Println("receiver has been done ")
}

// Core doc
func (r *Receiver) Core(ctx *Context) error {
	time.Sleep(1 * time.Second)
	select {
	case <-ctx.TimeoutCtx.Done():
		return ctx.TimeoutCtx.Err()

	case LastHandleChan <- r.Result:
		return nil
	}
}

// LastHandle doc
func LastHandle(ctx *Context, awg *sync.WaitGroup) {
	defer awg.Done()

	for {
		select {
		case <-ctx.TimeoutCtx.Done():
			fmt.Println("last handle ctx timeout ", ctx.TimeoutCtx.Err())
			return

		case res := <-LastHandleChan:
			fmt.Printf("the work : %s  has been done ! \n", res)

		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
