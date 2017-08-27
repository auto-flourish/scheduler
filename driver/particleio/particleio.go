package particleio

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	// ErrActionDoesNotExist means the requested action is not available
	ErrActionDoesNotExist = errors.New("action does not exist")
	// ErrStateChangeRequestFailed means the requested action failed to return 200
	ErrStateChangeRequestFailed = errors.New("particle io service did not return 200")
	// ErrRequestCreationFailed means the request could not be built - probably an error in the request body
	ErrRequestCreationFailed = errors.New("problem creating http request")
	// ErrGetRequestFailed means the request failed
	ErrGetRequestFailed = errors.New("problem sending a get request")
)

// ParticleIO implements driver interface for particleio board
type ParticleIO struct {
	UUID        string `json:"uuid"`
	DeviceID    string `json:"deviceId"`
	AccessToken string `json:"accessToken"`
}

// Authenticate hook performs any auth
// Get a new access token at this step
func (p *ParticleIO) Authenticate() error {
	p.DeviceID = "290018000347353137323334"
	p.AccessToken = "f19c897a5232a616fb611500e8fd62951858566a"
	return nil
}

// Poll returns a response to the client with data
// This is used internally to poll the sensors on a regular interval
// The user pulls out data from the store api
func (p *ParticleIO) Poll(action string) (interface{}, error) {
	u := fmt.Sprintf("https://api.particle.io/v1/devices/%s/%s?access_token=%s", p.DeviceID, action, p.AccessToken)
	resp, err := http.Get(u)
	if err != nil {
		return nil, ErrGetRequestFailed
	}
	// parse the response
	var r Response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Println("failed to decode")
		return nil, err
	}
	// unmarshal the result string into a struct
	var result Result
	if err := json.Unmarshal([]byte(r.Result), &result); err != nil {
		return nil, err
	}
	result.UUID = p.UUID
	result.Time = time.Now()

	defer resp.Body.Close()
	return &result, nil
}

// Action implements the driver interface and is invoked when the user wants to alter state
// Inversely, call Get when you just want to read sensor data
func (p *ParticleIO) Action(action string) error {

	u := fmt.Sprintf("https://api.particle.io/v1/products/4439/devices/%s/control", p.DeviceID)
	m := "POST"
	if err := p.makeRequest(u, m, action); err != nil {
		return err
	}
	return nil
}

// makeRequest is an internal method to change state
func (p *ParticleIO) makeRequest(u, m, action string) error {
	client := &http.Client{}

	data := url.Values{}
	data.Add("args", action)
	data.Add("access_token", p.AccessToken)

	req, err := http.NewRequest(m, u, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return ErrRequestCreationFailed
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(req)
	if err != nil {
		return ErrStateChangeRequestFailed
	}
	if resp.StatusCode != 200 {
		return ErrStateChangeRequestFailed
	}
	return nil
}

// Result represents the data structure inserted into database
type Result struct {
	UUID     string    `json:"uuid"`
	Time     time.Time `json:"time"`
	AirTemp0 float64   `json:"aT0"`
	AirTemp1 float64   `json:"aT1"`
	AirTemp2 float64   `json:"aT2"`
	AirTemp3 float64   `json:"aT3"`
	AirTemp4 float64   `json:"aT4"`
	AirTemp5 float64   `json:"aT5"`
	AirTemp6 float64   `json:"aT6"`
	AirTemp7 float64   `json:"aT7"`
	AirHume0 float64   `json:"aH0"`
	AirHume1 float64   `json:"aH1"`
	AirHume2 float64   `json:"aH2"`
	AirHume3 float64   `json:"aH3"`
	AirHume4 float64   `json:"aH4"`
	AirHume5 float64   `json:"aH5"`
	AirHume6 float64   `json:"aH6"`
	AirHume7 float64   `json:"aH7"`
	SoilMoi0 int       `json:"sM0"`
	SoilMoi1 int       `json:"sM1"`
	SoilMoi2 int       `json:"sM2"`
	SoilMoi3 int       `json:"sM3"`
	SoilMoi4 int       `json:"sM4"`
	SoilMoi5 int       `json:"sM5"`
	SoilMoi6 int       `json:"sM6"`
	SoilMoi7 int       `json:"sM7"`
	Lumens0  int       `json:"lm0"`
	Lumens1  int       `json:"lm1"`
	Lumens2  int       `json:"lm2"`
	Lumens3  int       `json:"lm3"`
	Lumens4  int       `json:"lm4"`
	Lumens5  int       `json:"lm5"`
	Lumens6  int       `json:"lm6"`
	Lumens7  int       `json:"lm7"`
	CO2      float64   `json:"CO2"`
}

// Response data from the particleIO service
type Response struct {
	CMD      string   `json:"cmd"`
	Name     string   `json:"name"`
	Result   string   `json:"result"`
	CoreInfo coreInfo `json:"coreInfo"`
}

type coreInfo struct {
	LastApp         string    `json:"last_app"`
	LastHeard       string    `json:"last_heard"`
	Connected       bool      `json:"connected"`
	LastHandshakeAt time.Time `json:"last_handshake_at"`
	DeviceID        string    `json:"deviceID"`
	ProductID       int       `json:"product_id"`
}

/*
POST response
{
	"id":"290018000347353137323334",
	"name":"test_photon",
	"last_app":"",
	"connected":true,
	"return_value":1
}
*/

/*
GET response
{
  "cmd": "VarReturn",
  "name": "lightValue",
  "result": 17,
  "coreInfo": {
    "last_app": "",
    "last_heard": "2017-06-19T00:23:41.545Z",
    "connected": true,
    "last_handshake_at": "2017-06-18T23:06:14.571Z",
    "deviceID": "290018000347353137323334",
    "product_id": 4439
  }
}
*/
