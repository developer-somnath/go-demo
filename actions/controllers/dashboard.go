package controllers

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// DashboardHandler is a default handler to serve up
// a home page.
func DashboardHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("pages/dashboard.plush.html"))

}
