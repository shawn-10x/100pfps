package handler

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/shawn-10x/100pfps/model"
	"github.com/shawn-10x/100pfps/utils"
	"github.com/shawn-10x/100pfps/validators"
)

func AdminGet(c echo.Context) (err error) {
	return c.Render(http.StatusOK, "admin.html", nil)
}

func AdminSignIn(c echo.Context) (err error) {
	type Form struct {
		User     string `form:"user" validate:"required,max=20"`
		Password string `form:"password" validate:"required,max=255"`
	}

	var form Form

	if err = c.Bind(&form); err != nil {
		return
	}

	showErrors := func(errors utils.Ms) error {
		return c.Render(http.StatusBadRequest, "admin.html", utils.M{
			"form-admin-errors": errors,
			"form-admin-data":   form,
		})
	}

	if err = c.Validate(&form); err != nil {
		errors := validators.ValidationErrors(err.(validator.ValidationErrors))
		return showErrors(errors)
	}

	admin := model.GetAdmin(form.User)
	if admin == nil {
		return showErrors(utils.Ms{
			"kind": "Invalid user",
		})
	}

	if !admin.CheckPassword(form.Password) {
		return showErrors(utils.Ms{
			"kind": "Invalid user or password",
		})
	}

	admin.NewToken()
	admin.Save()

	c.SetCookie(&http.Cookie{
		Name:     "admin_token",
		Value:    admin.Token,
		Expires:  time.Now().AddDate(0, 0, 7),
		HttpOnly: true,
		Path:     "/",
	})

	return c.Redirect(http.StatusSeeOther, "/")
}

func CreateAdmin(c echo.Context) (err error) {
	admin := c.Get("admin").(*model.Admin)
	if admin == nil || admin.Role != model.Owner {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	type Form struct {
		User     string          `form:"user" validate:"required,max=20"`
		Password string          `form:"password" validate:"required,max=255"`
		Role     model.AdminRole `form:"role" validate:"required,lte=2,gte=0"`
	}

	var form Form

	if err = c.Bind(&form); err != nil {
		return
	}

	if err = c.Validate(&form); err != nil {
		return
	}

	new_admin := model.Admin{
		User:     form.User,
		Password: form.Password,
		Role:     form.Role,
	}
	new_admin.Create()

	return c.Redirect(http.StatusSeeOther, "/")
}
