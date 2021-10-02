package worker

import (
	"fmt"
)

type Worker interface {
	Do(ctx *Context)
}

// BaseWorker doc
type BaseWorker struct {
	Name string
	Age  int
	Sex  string
}

// Do doc
func (bw *BaseWorker) Do(ctx *Context) {
	fmt.Printf("My name is %s i am %d years old \n", bw.Name, bw.Age)
}

// BaseInfo doc
func (bw *BaseWorker) BaseInfo() {
	fmt.Printf("My Base info : %+v", bw)
}
