package controllers

import (
	"net/http"

	// "demo/helpers"
	"demo/models"

	"errors"

	"github.com/gobuffalo/buffalo"
	// "golang.org/x/crypto/bcrypt"
)

type User struct {
	buffalo.Resource
}

func (e User) MyProfile(c buffalo.Context) error {
	var user models.User
	err := models.DB.Find(&user, c.Session().Get("UserId"))
	if err != nil {
		if errors.Is(err, errRecordNotFound) {
			// User not found
			c.Logger().Info("Record not found.")
		}
		c.Logger().Info("Fetching error.", err.Error())
	}
	c.Set("oldData", user)
	return c.Render(http.StatusOK, r.HTML("pages/my-profile.plush.html"))
}

func (e User) UpdateProfile(c buffalo.Context) error {
	return nil
}