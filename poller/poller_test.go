package poller

import (
	"fmt"
	"oxylus/driver/particleio"
	"testing"
	"time"
)

func TestPoller(t *testing.T) {
	p := Poller{}
	p.Action = "test"
	p.Driver = particleio.New()
	p.C = make(chan interface{})
	p.PollInterval = time.Second * 2
	var response particleio.Response
	p.Poll(&response)
	go func() {
		for {
			select {
			case v := <-p.C:
				fmt.Println(v)
			}
		}
	}()
	time.Sleep(time.Second * 10)
}
