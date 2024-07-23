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

	tags, err := model.GetAvaliableTags()
	if err != nil {
		return err
	}

	if err = c.Validate(&filter); err != nil {
		profiles, err2 := model.GetProfiles(nil)
		if err2 != nil {
			return err2
		}

		return c.Render(http.StatusBadRequest, "board.html", utils.M{
			"profiles": profiles,
			"tags":     tags,
			"tag":      filter.Tag,
		})
	}

	profiles, err := model.GetProfiles(filter.Tag)
	if err != nil {
		return
	}

	tags, err = model.GetAvaliableTags()
	if err != nil {
		return
	}

	return c.Render(http.StatusBadRequest, "board.html", utils.M{
		"profiles": profiles,
		"tags":     tags,
		"tag":      filter.Tag,
	})
}

func PostProfile(c echo.Context) (err error) {
	type Form struct {
		Name        string `form:"name" validate:"required,max=20"`
		Description string `form:"description" validate:"required,max=100"`
		Tags        string `form:"tags" validate:"max=75,tags,tags_max_count=5,tag_length=15"`
	}

	form := Form{}

	if err = c.Bind(&form); err != nil {
		return
	}

	profiles, err := model.GetProfiles(nil)
	if err != nil {
		return
	}

	tags, err := model.GetAvaliableTags()
	if err != nil {
		return
	}

	showErrors := func(errors utils.Ms) error {
		return c.Render(http.StatusBadRequest, "board.html", utils.M{
			"profiles": profiles,
			"tags":     tags,
			"form":     form,
			"errors":   errors,
		})
	}

	if err := c.Validate(&form); err != nil {
		errors := validators.ValidationErrors(err.(validator.ValidationErrors))
		return showErrors(errors)
	}

	image, err := utils.ReadImage(c, "img")
	if err != nil {
		return showErrors(utils.Ms{
			"Image": "Insert an image",
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
			"Image": "Error processing image",
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

	model.InsertProfile(&profile)
	utils.WriteImage(finalImg, model.GetProfileImg(&profile))

	return c.Redirect(http.StatusSeeOther, "/")
}
