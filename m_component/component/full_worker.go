package subworker

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"

	"github.com/sirupsen/logrus"
)

// FullParam  doc
type FullParam struct {
	Domain string // 域名
	IP     string // ip
	Port   string // 端口
}

// FullWorker doc
type FullWorker struct {
	BaseHandler
	MySSLPartnerID, MySSLSecretKey, UserAgent string
	TaskParam                                 *FullParam
	ResultSet                                 FullRes
}

// 结果集的处理

// FullRes doc
type FullRes struct {
	BaseResults
	FullData []*FullResult
}

// Unmarshal doc
func (fr *FullRes) Unmarshal() {
	for _, d := range fr.Data {
		dp, ok := d.(*FullResult)
		if ok {
			fr.FullData = append(fr.FullData, dp)
		}
	}
}

// 核心逻辑的处理

// Do doc
func (fw *FullWorker) Do(ctx *Context, current Handler, wg *sync.WaitGroup) {
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

// ChildrenDo doc
func (fw *FullWorker) ChildrenDo(ctx *Context) {
	if fw.Wg == nil {
		fw.Wg = &sync.WaitGroup{}
	}

	for _, bhc := range fw.ChildHandlers {
		fw.Wg.Add(1)
		go fw.Do(ctx, bhc, fw.Wg)
	}

	fw.Wg.Wait()
}

// BusinessDo doc
func (fw *FullWorker) BusinessDo(resChan chan interface{}, stopChan chan struct{}) {
	res, err := fw.fullCheckCore()
	if err != nil {
		logrus.Debugf("full check core faild when Check IP (%s) err : %v", fw.TaskParam.IP, err)
		return
	}

	select {
	case resChan <- res:
		return

	case <-stopChan:
		logrus.Debugf("full worker (%s): my work has been stopped ", fw.TaskParam.IP)
		return
	}
}

// BusinessDealRes doc
func (fw *FullWorker) BusinessDealRes(res interface{}) {
	v, ok := res.(*FullResult)
	if ok {
		if len(fw.ResultSet.Data) == 0 {
			fw.ResultSet.RPush(v)
			return
		}

		worst, ok := fw.ResultSet.Data[0].(*FullResult)
		if !ok {
			return
		}

		// 等级好的往后放
		if worst.Status >= v.Status {
			fw.ResultSet.RPush(v)
			return
		}

		fw.ResultSet.LPush(v)
	}
}

// FullResult doc
type FullResult struct {
	Status int
	Result string
}

func (fw *FullWorker) fullCheckCore() (*FullResult, error) {
	logrus.Debugf(" getFullCheckResult domain (%s) IP (%s) ", fw.TaskParam.Domain, fw.TaskParam.IP)
	statusNum := rand.Intn(100)
	return &FullResult{
		Status: statusNum,
		Result: fmt.Sprintf(" domain (%s) IP (%s) result", fw.TaskParam.Domain, fw.TaskParam.IP),
	}, nil
}
