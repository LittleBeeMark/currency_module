package subworker

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type FastParam struct {
	IP     string
	Port   string
	Domain string
}

// FastWorker doc
type FastWorker struct {
	BaseHandler
	TaskParam *FastParam
	ResultSet FastRes
}

// 结果集的处理

// FastRes doc
type FastRes struct {
	BaseResults
	DomainIPs []*FastResult
}

// Unmarshal doc
func (fr *FastRes) Unmarshal() {
	for _, d := range fr.Data {
		dp, ok := d.(*FastResult)
		if ok {
			fr.DomainIPs = append(fr.DomainIPs, dp)
		}
	}
}

// 核心逻辑的处理

// ChildrenDo doc
func (fw *FastWorker) ChildrenDo(ctx *Context) {
	if fw.Wg == nil {
		fw.Wg = &sync.WaitGroup{}
	}

	for _, bhc := range fw.ChildHandlers {
		fw.Wg.Add(1)
		go fw.Do(ctx, bhc, fw.Wg)
	}

	fw.Wg.Wait()
}

// Do doc
func (fw *FastWorker) Do(ctx *Context, current Handler, wg *sync.WaitGroup) {
	defer wg.Done()

	if fw.ResChan == nil {
		fw.ResChan = make(chan interface{})
	}
	if fw.StopChan == nil {
		fw.StopChan = make(chan struct{}, 1)
	}

	go current.BusinessDo(fw.ResChan, fw.StopChan)

	select {
	case res := <-fw.ResChan:
		fw.BusinessDealRes(res)
		break

	case <-ctx.TimeoutCtx.Done():
		select {
		case fw.StopChan <- struct{}{}:
		default:
		}
		fw.Err = errors.New("monitor sub work ctx time out ")
		break
	}

}

// BusinessDo doc
func (fw *FastWorker) BusinessDo(resChan chan interface{}, stopChan chan struct{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res := fw.fastCheckCore(ctx)
	select {
	case resChan <- res:
		return

	case <-stopChan:
		logrus.Debugf("my work has stopped")
		return
	}
}

// BusinessDealRes doc
func (fw *FastWorker) BusinessDealRes(res interface{}) {
	v, ok := res.(*FastResult)
	if ok {
		if len(fw.ResultSet.Data) == 0 {
			fw.ResultSet.RPush(v)
			return
		}

		worst, ok := fw.ResultSet.Data[0].(*FastResult)
		if !ok {
			return
		}

		// 状态好的往后放
		if v.Status > worst.Status {
			fw.ResultSet.RPush(v)
			return
		}

		fw.ResultSet.LPush(v)
	}
}

// FastResult doc
type FastResult struct {
	Status int
	Result string
}

// fastCheckCore 快速检测核心
func (fw *FastWorker) fastCheckCore(ctx context.Context) *FastResult {
	var (
		err   error
		infos *FastResult // 临时中间变量
	)

	infos, err = GenerateMultipleCertificates(ctx, fw.TaskParam)
	if err != nil {
		logrus.Debugf("fastCheckCore.GenerateMultipleCertificates failed when check IP(%s) : %v",
			fw.TaskParam.IP, err)
	}

	return infos
}

// GenerateMultipleCertificates 获取证书
func GenerateMultipleCertificates(ctx context.Context, params *FastParam) (*FastResult, error) {
	if params.IP == "" {
		return nil, fmt.Errorf("not found ip: %s", params.IP)
	}
	fastRandNum := rand.Intn(100)

	return &FastResult{
		Status: fastRandNum,
		Result: fmt.Sprintf("Domain : %s,IP : %s ", params.Domain, params.IP),
	}, nil
}
