package controllers

import (
	"net/http"

	"demo/helpers"
	"demo/models"

	"errors"

	"github.com/gobuffalo/buffalo"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	buffalo.Resource
}

var BlankArray = []interface{}{}
var errRecordNotFound = errors.New("record not found")

func (e Auth) Index(c buffalo.Context) error {
	if _, ok := c.Session().Get("UserId").(string); ok {
		return c.Redirect(302, "/")
	}
	// Render the "log-in.plush.html" template using c.Render with r.HTML.
	return c.Render(http.StatusOK, r.HTML("pages/log-in.plush.html", "before-login.plush.html"))
}

func (e Auth) UserCheck(c buffalo.Context) error {
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	requestedEmail := c.Params().Get("email")
	requestPassword := c.Params().Get("password")
	var User models.User
	err := models.DB.Where("email = ?", requestedEmail).First(&User)
	// Perform actions based on the data...
	if err != nil {
		if errors.Is(err, errRecordNotFound) {
			// User not found
			return c.Render(http.StatusNotFound, r.JSON(map[string]interface{}{
				"status":   false,
				"message":  "Email does not exists",
				"data":     BlankArray,
				"redirect": "",
				"error":    err,
			}))
		}
		return c.Render(http.StatusInternalServerError, r.JSON(map[string]interface{}{
			"status":  false,
			"message": "Error fetching data",
			"data":    BlankArray,
			"error":   err.Error(), // or custom error message if you don't want to expose the database error details
		}))

	}
	err = bcrypt.CompareHashAndPassword([]byte(*User.Password), []byte(requestPassword))
	if err != nil {
		return c.Render(http.StatusUnauthorized, r.JSON(map[string]interface{}{
			"status":  false,
			"message": "Password does not match",
			"data":    BlankArray,
			"error":   err.Error(), // or custom error message if you don't want to expose the database error details
		}))
	}
	session := c.Session()
	// Set a value in the session
	session.Set("UserId", helpers.IntToString(User.ID))
	session.Set("FirstName", helpers.ValueOrDefault(User.FirstName))
	session.Set("LastName", helpers.ValueOrDefault(User.LastName))
	session.Set("RoleId", helpers.IntToString(User.RoleId))
	session.Set("ProfileImage", helpers.ValueOrDefault(User.ProfileImage))

	uuidString := ""
	if User.UUID != nil {
		uuidString = User.UUID.String()
	}
	session.Set("UUID", uuidString)
	err = session.Save()
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	// Return a JSON response
	return c.Render(http.StatusOK, r.JSON(map[string]interface{}{
		"status":   true,
		"message":  "Login Successfull! Redirecting to Dashboard....",
		"data":     BlankArray,
		"redirect": "/",
		"error":    "",
	}))
}

func (e Auth) LogOut(c buffalo.Context) error {
	// Destroy the session (assuming you're using the default session store)
	c.Session().Clear()
	// Optional: You can also save the session to ensure it gets cleared immediately.
	c.Session().Save()
	// Redirect the user to the home page or login page
	// Replace "/login" with the appropriate URL for your application
	return c.Redirect(302, "/login")
}
