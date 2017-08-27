package particleio

import (
	"testing"
	"time"
)

func TestNewDriver(t *testing.T) {
	p := NewDriver()
	if p.Name() != "ParticleIO" {
		t.Error("failed to instantiate particleio")
	}
	if p.Description() != "ParticleIO gives you everything you need to securely and reliably connect your IoT devices to the web." {
		t.Error("failed to instantiate particleio")
	}
}

func TestRequest(t *testing.T) {
	p := NewDriver()
	p.Authenticate()
	if err := p.Run("fans"); err != nil {
		t.Error("request failed to return 200")
	}
	time.Sleep(3 * time.Second)
	if err := p.Run("fans"); err != nil {
		t.Error("request failed to return 200")
	}
}

func TestReadRequest(t *testing.T) {
	p := NewDriver()
	p.Authenticate()
	var response Response
	if err := p.Get("lightValue", &response); err != nil {
		t.Error("request failed")
	}
	if response.Name != "lightValue" {
		t.Error("request failed")
	}
}
