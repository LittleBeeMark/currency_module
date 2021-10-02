package worker

import (
	"context"
	"time"
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
