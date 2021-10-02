package worker

import (
	"fmt"
	"sync"
	"time"
)

var (
	DealChan      chan interface{}
	DealResultSet []interface{}
)

func init() {
	DealChan = make(chan interface{})
}

func dealCore(res interface{}) {
	fmt.Println("deal the result : ", res)
	DealResultSet = append(DealResultSet, res)
}

// DealResult doc
func DealResult(ctx *Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case res := <-DealChan:
			dealCore(res)

		case <-ctx.TimeoutCtx.Done():
			fmt.Println("deal result ctx is timeout", ctx.TimeoutCtx.Err())
			return

		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
