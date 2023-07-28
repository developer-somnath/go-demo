package controllers

import (
	"net/http"

	"demo/helpers"
	"demo/models"

	"errors"

	"mime/multipart"

	"github.com/gobuffalo/buffalo"
	"golang.org/x/crypto/bcrypt"
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
	// Get the currently logged-in user
	var user models.User
	err := models.DB.Find(&user, c.Session().Get("UserId"))
	// Retrieve form data
	firstName := c.Request().FormValue("firstName")
	lastName := c.Request().FormValue("lastName")
	email := c.Request().FormValue("email")
	oldPassword := c.Request().FormValue("oldPassword")
	newPassword := c.Request().FormValue("newPassword")
	newPasswordRepeat := c.Request().FormValue("profileNewPasswordRepeat")

	// Validate form data, perform any necessary checks
	// Check if the user uploaded a profile image
	fileHeader, err := c.File("profileImage")
	if err == nil {
		// Create a new *multipart.FileHeader instance from binding.File
		newFileHeader := &multipart.FileHeader{
			Filename: fileHeader.Filename,
			Size:     fileHeader.Size,
		}

		// Pass the newly created *multipart.FileHeader and the file content to the helpers.UploadFileSingle function
		file, err := fileHeader.Open()
		if err != nil {
			// c.Logger().Error("Error opening profile image.", err)
			return c.Render(http.StatusInternalServerError, r.JSON(map[string]interface{}{
				"status":   false,
				"message":  "Error opening profile image",
				"data":     BlankArray,
				"redirect": "",
				"error":    err,
			}))

		} else {
			uploaded := helpers.UploadFileSingle(newFileHeader, file, "admin/uploads/profileImages")
			if uploaded.Err != nil {
				// c.Logger().Error("Error uploading profile image.", uploaded.Err)
				return c.Render(http.StatusInternalServerError, r.JSON(map[string]interface{}{
					"status":   false,
					"message":  "Error uploading profile image",
					"data":     BlankArray,
					"redirect": "",
					"error":    uploaded.Err,
				}))
			} else {
				// Update the user's profile image filename in the database
				user.ProfileImage = &uploaded.FileName
			}
		}
	}
	// Update the user's profile information (including first name, last name, email)
	// and potentially change the password if necessary
	if len(firstName) > 0 {
		user.FirstName = &firstName
	}
	if len(lastName) > 0 {
		user.LastName = &lastName
	}
	if len(email) > 0 {
		user.Email = &email
	}
	if len(newPassword) > 0 && len(oldPassword) > 0 {
		// Check if the old password matches the one in the database
		err := bcrypt.CompareHashAndPassword([]byte(helpers.ValueOrDefault(user.Password)), []byte(oldPassword))
		if err != nil {
			// c.Logger().Error("Old password does not match.", err)
			return c.Render(http.StatusBadRequest, r.JSON(map[string]interface{}{
				"status":   false,
				"message":  "Old password does not match.",
				"data":     BlankArray,
				"redirect": "",
				"error":    err,
			}))
		} else {
			// Hash the new password before storing it in the database
			if newPassword != newPasswordRepeat {
				// c.Logger().Error("New Password & Confirm password does not match.")
				return c.Render(http.StatusInternalServerError, r.JSON(map[string]interface{}{
					"status":   false,
					"message":  "New Password & Confirm password does not match.",
					"data":     BlankArray,
					"redirect": "",
					"error":    "",
				}))
			}
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
			if err != nil {
				// c.Logger().Error("Error hashing the new password.", err)
				return c.Render(http.StatusInternalServerError, r.JSON(map[string]interface{}{
					"status":   false,
					"message":  "Something went wrong.",
					"data":     BlankArray,
					"redirect": "",
					"error":    err,
				}))
			}
			// Update the user's password hash in the database
			hashedPasswordPointer := string(hashedPassword)
			user.Password = &hashedPasswordPointer
		}
	}

	// Save the updated user profile in the database
	err = models.DB.Save(&user)
	if err != nil {
		// c.Logger().Error("Error saving user profile.", err)
		// Handle the error appropriately, e.g., show an error message to the user
		return c.Render(http.StatusInternalServerError, r.JSON(map[string]interface{}{
			"status":   false,
			"message":  "Something went wrong.",
			"data":     BlankArray,
			"redirect": "",
			"error":    err.Error(),
		}))
	}

	// Redirect the user to the profile page or show a success message
	return c.Render(http.StatusOK, r.JSON(map[string]interface{}{
		"status":   true,
		"message":  "Profile updated succefully",
		"data":     BlankArray,
		"redirect": "",
		"error":    "",
	}))
}
