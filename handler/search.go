package handler

import (
	"log"
	"net/http"
	"oxylus/driver/particleio"

	"gopkg.in/mgo.v2/bson"

	"github.com/labstack/echo"
)

// GetAll returns all metric records for a poller id
func (h *Handler) GetAll(c echo.Context) error {
	id := c.Param("id")
	db := h.DB.Clone()
	defer db.Close()

	var results []*particleio.Result
	if err := db.DB("oxylus").C("metrics").Find(bson.M{"uuid": id}).Sort("time").All(&results); err != nil {
		log.Println(err)
	}
	return c.JSON(http.StatusOK, results)
}
