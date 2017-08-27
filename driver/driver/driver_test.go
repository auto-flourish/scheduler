package driver

import (
	"errors"
	"net/http"
	"testing"
)

type testDriver struct{}

func (t *testDriver) Run(action string) error {
	switch action {
	case "lightsOn":
		return lightsOn()
	case "lightsOff":
		return nil
	}
	return nil
}

func lightsOn() error {
	res, err := http.Get("http://google.com")
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("http get to google failed")
	}
	return nil
}
func TestDriver(t *testing.T) {
	ta := testDriver{}
	if err := ta.Run("lightsOn"); err != nil {
		t.Error("lights on failed")
	}
}
