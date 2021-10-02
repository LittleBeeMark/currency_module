package main

import (
	"currency_module/m_timeout/worker"
	"fmt"
	"sync"
)

// 控制超时结束并发模型实现
const maxCurrencyCount = 3

// doc
func dealSendMessage(ress []interface{}) interface{} {
	var s string
	for _, v := range ress {
		s += fmt.Sprintf("last handle results :%v   ", v)
	}
	return s
}

func SenderRun(ctx *worker.Context, rwp *worker.WorkPool,
	awg *sync.WaitGroup, cascade bool) {
	defer awg.Done()

	// 子任务并发
	wp := worker.StartWorkPool(ctx, maxCurrencyCount)
	wp.Dispatch(&worker.Worker1{})
	wp.Dispatch(&worker.Worker2{})
	wp.Dispatch(&worker.Worker3{})
	wp.Stop()

	select {
	case <-ctx.TimeoutCtx.Done():
		fmt.Println("send ctx timeout ", ctx.TimeoutCtx.Err())
		rwp.Stop()
		return

	default:
		if len(worker.DealResultSet) > 0 {
			rwp.Dispatch(worker.NewReceiver(dealSendMessage(worker.DealResultSet)))
		}
		rwp.Stop()
		return
	}
}

// ServerRun doc
func ServerRun(ctx *worker.Context, awg *sync.WaitGroup) {

	awg.Add(3)
	// 启动最后结果处理者
	go worker.LastHandle(ctx, awg)

	// 启动接收者
	receiverWp := worker.StartWorkPool(ctx, maxCurrencyCount)

	// 启动发送数据处理者
	go worker.DealResult(ctx, awg)

	// 启动发送者
	go SenderRun(ctx, receiverWp, awg, false)
}
