package worker

import (
	"fmt"
	"sync"
)

// WorkPool doc
type WorkPool struct {
	WorkCh chan Worker
	Wg     sync.WaitGroup
}

// StartWorkPool doc
func StartWorkPool(ctx *Context, currencyCount int) *WorkPool {
	wp := WorkPool{
		WorkCh: make(chan Worker),
	}

	wp.Wg.Add(currencyCount)
	for i := 0; i < currencyCount; i++ {
		go func(index int) {
			defer wp.Wg.Done()

			for w := range wp.WorkCh {
				w.Do(ctx)
			}
			fmt.Printf("workPool number%d stopped \n", index)
		}(i)
	}

	return &wp
}

// Dispatch 任务派发
func (wp *WorkPool) Dispatch(w Worker) {
	wp.WorkCh <- w
}

// Stop 停止该并发池
func (wp *WorkPool) Stop() {
	close(wp.WorkCh)
	wp.Wg.Wait()
}
