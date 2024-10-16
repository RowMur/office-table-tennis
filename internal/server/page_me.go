package server

import (
	"net/http"

	"github.com/RowMur/office-games/internal/views"
	"github.com/labstack/echo/v4"
)

func mePageHandler(c echo.Context) error {
	user := userFromContext(c)
	if user == nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/sign-in")
	}

	return render(c, http.StatusOK, views.MePage(user, views.UserDetailsFormData{Email: user.Email, Username: user.Username}, views.UserDetailsFormErrors{}))
}

func (s *Server) meUpdateHandler(c echo.Context) error {
	user := userFromContext(c)
	if user == nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/sign-in")
	}

	username := c.FormValue("username")
	email := c.FormValue("email")

	if username == "" || email == "" {
		data := views.UserDetailsFormData{Username: username, Email: email}
		errs := views.UserDetailsFormErrors{}
		if username == "" {
			errs.Username = "Username is required"
		}
		if email == "" {
			errs.Email = "Email is required"
		}
		falseVar := false
		return render(c, http.StatusOK, views.UserDetailsForm(data, errs, &falseVar))
	}

	user, updateErrs, err := s.db.UpdateUser(user.ID, map[string]interface{}{"username": username, "email": email})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if updateErrs != nil {
		formErrs := views.UserDetailsFormErrors{
			Username: updateErrs["username"],
			Email:    updateErrs["email"],
		}
		return render(c, http.StatusOK, views.UserDetailsForm(views.UserDetailsFormData{Username: username, Email: email}, formErrs, nil))
	}

	formData := views.UserDetailsFormData{Email: user.Email, Username: user.Username}
	truePtr := true
	return render(c, http.StatusOK, views.UserDetailsForm(formData, views.UserDetailsFormErrors{}, &truePtr))
}
