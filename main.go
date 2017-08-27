package main

import (
	"log"
	"oxylus/handler"
	"oxylus/poller"
	"oxylus/pollerregistry"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/labstack/echo"

	mw "github.com/labstack/echo/middleware"
)

var (
	retries = 5
)

func main() {

	e := echo.New()

	e.Use(mw.Logger())
	e.Use(mw.CORS())
	e.Use(mw.Recover())

	db, err := mgo.Dial("db")
	if err != nil {
		log.Fatal(err)
	}
	// Create indices
	if err = db.Copy().DB("oxylus").C("users").EnsureIndex(mgo.Index{
		Key:    []string{"uuid"},
		Unique: true,
	}); err != nil {
		log.Fatal(err)
	}

	h := handler.Handler{
		PollerRegistry: pollerregistry.New(),
		DB:             db,
	}

	e.GET("/", h.Test)

	e.GET("/users", h.GetUsers)
	e.POST("/users", h.CreateUser)
	e.DELETE("/users/:id", h.DeleteUser)
	// database query
	e.GET("/search/pollers/:id", h.GetAll)

	// events are things that change state
	e.GET("/users/:id/events", h.GetUserEvents)
	e.GET("/users/:id/events/:event", h.GetUserEvent)
	e.POST("/users/:id/events", h.CreateEvent)
	e.DELETE("/users/:id/events/:event", h.DeleteEvent)

	// pollers read sensors and dont change state
	e.GET("/users/:id/pollers", h.GetPollers)
	e.POST("/users/:id/pollers", h.AddPoller)
	e.PUT("/pollers/:id", h.UpdatePoller)
	e.DELETE("/pollers/:id", h.DeletePoller)

	// insert to database
	go func(h *handler.Handler) {
		for {
			select {
			case msg := <-h.PollerRegistry.ToDB:
				db := h.DB.Clone()
				log.Printf("[POLLED] %v\n", msg)
				if err := db.DB("oxylus").C("metrics").Insert(msg); err != nil {
					log.Println(err)
				}
				db.Close()
			}
		}
	}(&h)

	// toggle state off when polling fails
	go func(h *handler.Handler) {
		for {
			select {
			case msg := <-h.PollerRegistry.UpdateStatus:
				db := h.DB.Clone()
				if err := db.DB("oxylus").C("pollers").Update(bson.M{"uuid": msg}, bson.M{"$set": bson.M{"ispolling": false}}); err != nil {
					log.Println(err)
				} else {
					log.Printf("setting isPolling to false for %s\n", msg)
				}
				db.Close()
			}
		}
	}(&h)

	// restart any active pollers
	db = h.DB.Clone()
	var results []*poller.Poller
	if err := db.DB("oxylus").C("pollers").Find(bson.M{"ispolling": true}).All(&results); err != nil {
		log.Println(err)
	}

	for _, val := range results {
		h.PollerRegistry.Add(val.UUID, val)
	}
	db.Close()

	e.Logger.Fatal(e.Start(":1323"))
}
