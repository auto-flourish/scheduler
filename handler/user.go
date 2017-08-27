package handler

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"oxylus/user"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

var (
	name       = "oxylus"
	usersTable = "users"
)

// CreateUser creates a global uuid that keys all their events
func (h *Handler) CreateUser(c echo.Context) error {
	db := h.DB.Clone()
	defer db.Close()

	u := &user.User{
		UUID:      uuid.NewV4().String(),
		FirstName: "",
		LastName:  "",
	}

	if err := db.DB(name).C(usersTable).Insert(u); err != nil {
		return err
	}

	return c.JSONPretty(http.StatusCreated, u, "  ")
}

// DeleteUser --
func (h *Handler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	db := h.DB.Clone()
	defer db.Close()
	if err := db.DB("oxylus").C("users").Remove(bson.M{"uuid": id}); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

// GetUsers returns all users
func (h *Handler) GetUsers(c echo.Context) error {
	db := h.DB.Clone()
	defer db.Close()

	var results []user.User
	if err := db.DB(name).C(usersTable).Find(bson.M{}).All(&results); err != nil {
		return err
	}
	return c.JSONPretty(http.StatusOK, results, "  ")
}
