# Oxylus Events Backend

This repository represents an api for creating and dispatching `events` on anything that implements the driver interface.

## Installation

[Install Go](https://golang.org/dl/)

Ensure `GOROOT` environment variable is set to go install folder - the installer will usually set this for you.

Ensure `GOPATH` environment variable is set.  Go has this idiom for how to organize your code.  It goes like this.


```
Root (this is what you set the GOPATH variable to)
|_bin (if you create binaries they will go here)
|_pkg (object files - dont worry about it)
|_src (where you put your source)
    |_package1
    |_package2
```

## Install dependencies
`go get github.com/labstack/echo/...`
`go get github.com/satori/go.uuid`

## Usage
Run the server with `go run main.go`.  You can also create a binary with `go build main.go`.

1. Get a user UUID by sending a POST request to http://localhost:1323/user.
2. Take that UUID and create an event by sending the following data structure in a POST request to `http://localhost:1323/user/[the-uuid-from-the-previous-step]/event`

```
{
	"action": "test",
	"repeats": true,
	"driver": "particleio",
	"driverParams": {
		"device_id": "12345",
		"access_token": "12345"
	},
	"date": {
		"year": 2017,
		"month": 6,
		"day": 22,
		"hour": 23,
		"minute": 1,
		"second": 15
	}
}
```
This will return the following response:
```
{
  "event": {
    "uuid": "a934b5bc-2505-4c25-97ae-b55ed4f47923",
    "createdAt": "2017-06-22T23:28:13.052672955-07:00",
    "finishAt": "2017-06-22T23:01:15-07:00",
    "action": "test",
    "driver": {
      "uuid": "fe6a0cf8-6004-498e-8f52-dcbab80e5b03",
      "deviceId": "12345",
      "accessToken": "12345"
    },
    "repeats": true,
    "timeInterval": -1618052676481
  },
  "commands": [
    {
      "cmd": "lights",
      "description": "toggle lights"
    },
    {
      "cmd": "fans",
      "description": "toggle fans"
    },
    {
      "cmd": "fillWaterTank",
      "description": "fills the water tank"
    },
    {
      "cmd": "waterPump",
      "description": "toggle water pump"
    },
    {
      "cmd": "mixValve",
      "description": "toggle mix valve"
    },
    {
      "cmd": "feedValve",
      "description": "toggle feed valve"
    },
    {
      "cmd": "stirNutrients",
      "description": "toggle stir nutrients"
    },
    {
      "cmd": "nutrientOne",
      "description": "toggle nutrient pump 1"
    },
    {
      "cmd": "nutrientTwo",
      "description": "toggle nutrient pump 2"
    },
    {
      "cmd": "nutrientThree",
      "description": "toggle nutrient pump 3"
    }
  ],
  "request": {
    "date": {
      "year": 2017,
      "month": 6,
      "day": 22,
      "hour": 23,
      "minute": 1,
      "second": 15
    },
    "action": "test",
    "repeats": true,
    "driver": "particleio",
    "driverParams": {
      "access_token": "12345",
      "device_id": "12345"
    }
  }
}
```
