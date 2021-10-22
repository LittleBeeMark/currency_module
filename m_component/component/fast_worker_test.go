package subworker

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestFastWorker_BusinessDo(t *testing.T) {
	params := &FastParam{
		IP:     "106.14.248.134",
		Domain: "marksuper.xyz",
		Port:   "443",
	}
	fw := &FastWorker{
		TaskParam: params,
	}

	if fw.ResChan == nil {
		fw.ResChan = make(chan interface{})
	}
	if fw.StopChan == nil {
		fw.StopChan = make(chan struct{}, 1)
	}

	go fw.BusinessDo(fw.ResChan, fw.StopChan)

	select {
	case res := <-fw.ResChan:
		fmt.Println("res : ", res)
		break

	case <-context.Background().Done():
		select {
		case fw.StopChan <- struct{}{}:
		default:
		}
		fw.Err = errors.New("monitor sub work ctx time out ")
		break
	}
}

func TestFastWorker_Do(t *testing.T) {
	params := &FastParam{
		Domain: "marksuper.xyz",
		Port:   "443",
	}
	ips := []string{
		"106.14.248.134",
	}
	for i := 1; i < 255; i++ {
		for j := 1; j < 255; j++ {
			ip := fmt.Sprintf("220.181.%d.%d", i, j)
			ips = append(ips, ip)
		}
	}

	fastHandler := FastWorker{}
	for _, ip := range ips {
		// 挂载任务
		param := *params
		param.IP = ip
		fastHandler.Mount(&FastWorker{
			TaskParam: &param,
		})
	}

	fastHandler.ChildrenDo(GetContext(180 * time.Second))
	fastHandler.ResultSet.Unmarshal()
	fmt.Println("IPS NUM: ", len(ips))
	fmt.Println("Result NUM : ", len(fastHandler.ResultSet.DomainIPs))
	//fmt.Printf("%+v", fastHandler.ResultSet.DomainIPs)
	for i := 0; i < len(fastHandler.ResultSet.DomainIPs); i++ {
		fmt.Printf("fast result : %+v \n", fastHandler.ResultSet.DomainIPs[i])
	}

	return
}
