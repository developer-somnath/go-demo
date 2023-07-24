package middleware

import (
	"github.com/gobuffalo/buffalo"
)

func SessionAuth(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		// Code to be executed before the next handler (if any).
		isAuthenticated := false

		if _, ok := c.Session().Get("UserId").(string); ok {
			isAuthenticated = true
		}
		if isAuthenticated {
			c.Logger().Info("User is authenticated.")
			return next(c)
		}
		return c.Redirect(302, "/login")

	}
}
