package handler

import (
	"net"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
	"github.com/shawn-10x/100pfps/model"
	"github.com/shawn-10x/100pfps/utils"
	"github.com/shawn-10x/100pfps/validators"
)

func GetBoard(c echo.Context) (err error) {
	type Filter struct {
		Tag *string `query:"tag" validate:"omitempty,max=15"`
	}

	filter := Filter{}
	if err = c.Bind(&filter); err != nil {
		return
	}

	if err = c.Validate(&filter); err != nil {
		return c.Render(http.StatusBadRequest, "board.html", utils.M{
			"profiles": model.GetProfiles(nil),
			"tags":     model.GetAvaliableTags(),
			"tag":      filter.Tag,
		})
	}

	return c.Render(http.StatusOK, "board.html", utils.M{
		"profiles": model.GetProfiles(filter.Tag),
		"tags":     model.GetAvaliableTags(),
		"tag":      filter.Tag,
	})
}

func PostProfile(c echo.Context) (err error) {
	type Form struct {
		Name                    string `form:"name" validate:"required,max=20"`
		Description             string `form:"description" validate:"required,max=100"`
		Tags                    string `form:"tags" validate:"required,max=75,tags,tags_max_count=5,tag_length=15"`
		RulesAndPrivacyAccepted bool   `form:"rulesandprivacyaccepted" validate:"required"`
	}

	form := Form{}

	if err = c.Bind(&form); err != nil {
		return
	}

	showErrors := func(errors utils.Ms) error {
		return c.Render(http.StatusBadRequest, "board.html", utils.M{
			"profiles":            model.GetProfiles(nil),
			"tags":                model.GetAvaliableTags(),
			"form-profile-data":   form,
			"form-profile-errors": errors,
		})
	}

	if err := c.Validate(&form); err != nil {
		errors := validators.ValidationErrors(err.(validator.ValidationErrors))
		return showErrors(errors)
	}

	if !form.RulesAndPrivacyAccepted {
		return showErrors(utils.Ms{
			"kind": "You need to accept Rules and Privacy Policy to proceed",
		})
	}

	if model.ExistsProfileWithIP(c.Get("ip").(net.IP)) {
		return showErrors(utils.Ms{
			"kind": "You already posted a profile with this IP",
		})
	}

	image, err := utils.ReadImage(c, "img")
	if err != nil {
		return showErrors(utils.Ms{
			"kind": "Insert an image",
		})
	}

	croppedImg, err := cutter.Crop(image, cutter.Config{
		Width:   1,
		Height:  1,
		Mode:    cutter.Centered,
		Options: cutter.Ratio,
	})

	if err != nil {
		return showErrors(utils.Ms{
			"kind": "Error processing image",
		})
	}

	finalImg := resize.Resize(0, 100, croppedImg, resize.Lanczos3)

	var ip pgtype.Inet
	ip.Set(c.Get("ip").(net.IP))

	profile := model.Profile{
		Name:        form.Name,
		Description: form.Description,
		Tags:        model.StrToTags(form.Tags),
		Ip:          ip,
	}

	profile.Insert()
	utils.WriteImage(finalImg, profile.GetProfileImg())

	return c.Redirect(http.StatusSeeOther, "/")
}
