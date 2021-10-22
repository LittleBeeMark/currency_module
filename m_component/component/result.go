package subworker

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type Results interface {
	LPush(interface{})
	RPush(interface{})
	Unmarshal()
}

// BaseResults doc
type BaseResults struct {
	M    sync.Mutex
	Data []interface{}
}

// LPush doc
func (br *BaseResults) LPush(d interface{}) {
	br.M.Lock()
	defer br.M.Unlock()

	if cap(br.Data) == len(br.Data) {
		br.Data = append([]interface{}{d}, br.Data...)
	} else {
		var a interface{}
		br.Data = append(br.Data, a)
		copy(br.Data[1:], br.Data)
		br.Data[0] = d
	}
}

// RPush doc
func (br *BaseResults) RPush(d interface{}) {
	br.M.Lock()
	defer br.M.Unlock()

	if cap(br.Data) == 0 {
		br.Data = make([]interface{}, 0, 5)
	}
	br.Data = append(br.Data, d)
}

// Unmarshal doc
func (br *BaseResults) Unmarshal() {
	logrus.Debugln("this is base result unmarshal")
}
