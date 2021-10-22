package subworker

import (
	"fmt"
	"testing"
	"time"
)

func TestFullWorker_Do(t *testing.T) {
	params := &FullParam{
		Domain: "marksuper.xyz",
		Port:   "443",
	}

	var ips []string
	for i := 1; i < 255; i++ {
		for j := 1; j < 255; j++ {
			ip := fmt.Sprintf("220.181.%d.%d", i, j)
			ips = append(ips, ip)
		}
	}

	fullHandler := FullWorker{}
	for _, ip := range ips {
		// 挂载任务
		param := *params
		param.IP = ip
		fullHandler.Mount(&FullWorker{
			TaskParam: &param,
		})
	}

	fullHandler.ChildrenDo(GetContext(180 * time.Second))
	fullHandler.ResultSet.Unmarshal()
	fmt.Println("IPS NUM: ", len(ips))
	fmt.Println("Result NUM : ", len(fullHandler.ResultSet.FullData))
	fmt.Printf("%+v", fullHandler.ResultSet.FullData)
	for _, v := range fullHandler.ResultSet.FullData {
		fmt.Printf("fullResult : %+v", v)
	}
	return
}
