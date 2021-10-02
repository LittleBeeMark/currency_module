package main

import (
	"currency_module/m_timeout/worker"
	"fmt"
	"sync"
	"time"
)

func main() {
	//全局控制上下文 及 wait group
	ctx := worker.GetContext(10 * time.Second)
	awg := sync.WaitGroup{}

	// 服务启动
	ServerRun(ctx, &awg)

	awg.Wait()
	fmt.Println("successful")
}
