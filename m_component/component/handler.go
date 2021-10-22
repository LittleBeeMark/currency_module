package subworker

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Context doc
type Context struct {
	TimeoutCtx context.Context
	context.CancelFunc
}

// GetContext doc
func GetContext(d time.Duration) *Context {
	c := &Context{}
	c.TimeoutCtx, c.CancelFunc = context.WithTimeout(context.Background(), d)
	return c
}

// Handler doc
type Handler interface {
	Mount(hs ...Handler)
	Remove(h Handler)
	Do(ctx *Context, currentHandler Handler, wg *sync.WaitGroup)
	BusinessDo(chan interface{}, chan struct{})
	BusinessDealRes(interface{})
	ChildrenDo(ctx *Context)
}

// BaseHandler doc
type BaseHandler struct {
	Wg            *sync.WaitGroup
	ChildHandlers []Handler
	ResChan       chan interface{}
	StopChan      chan struct{}
	Err           error
}

// Mount 挂载任务
func (bh *BaseHandler) Mount(hs ...Handler) {
	bh.ChildHandlers = append(bh.ChildHandlers, hs...)
}

// Remove doc
func (bh *BaseHandler) Remove(h Handler) {
	if len(bh.ChildHandlers) == 0 {
		return
	}

	for i, v := range bh.ChildHandlers {
		if v == h {
			bh.ChildHandlers = append(bh.ChildHandlers[:i], bh.ChildHandlers[i+1:]...)
		}
	}
}

// ChildrenDo doc
func (bh *BaseHandler) ChildrenDo(ctx *Context) {
	if bh.Wg == nil {
		bh.Wg = &sync.WaitGroup{}
	}

	for _, bhc := range bh.ChildHandlers {
		bh.Wg.Add(1)
		go bh.Do(ctx, bhc, bh.Wg)
	}

	bh.Wg.Wait()
}

// Do doc
func (bh *BaseHandler) Do(ctx *Context, current Handler, wg *sync.WaitGroup) {
	defer wg.Done()

	if bh.ResChan == nil {
		bh.ResChan = make(chan interface{})
	}
	if bh.StopChan == nil {
		bh.StopChan = make(chan struct{}, 1)
	}

	go current.BusinessDo(bh.ResChan, bh.StopChan)

	select {
	case res := <-bh.ResChan:
		current.BusinessDealRes(res)
		break

	case <-ctx.TimeoutCtx.Done():
		select {
		case bh.StopChan <- struct{}{}:
		default:
		}

		bh.Err = errors.New("monitor sub work ctx time out ")
		break
	}

}

func (bh *BaseHandler) BusinessDealRes(res interface{}) {
	logrus.Info("this base handler,if you see it means this handler not need BusinessDealRes")
}

func (bh *BaseHandler) BusinessDo(resChan chan interface{}, stopChan chan struct{}) {
	logrus.Info("this base handler,if you see it means this handler not need BusinessDo")
}
