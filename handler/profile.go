package handler

import (
	"net"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/h2non/bimg"
	"github.com/jackc/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/shawn-10x/100pfps/model"
	"github.com/shawn-10x/100pfps/utils"
	"github.com/shawn-10x/100pfps/validators"
)

func GetProfiles(c echo.Context) (err error) {
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
		"admin":    c.Get("admin").(*model.Admin),
	})
}

func CreateProfile(c echo.Context) (err error) {
	type Form struct {
		Name                    string `form:"name" validate:"required,max=20"`
		Description             string `form:"description" validate:"required,max=1000"`
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

	bimage, err := utils.ReadImage(c, "img")
	if err != nil {
		return showErrors(utils.Ms{
			"kind": "Insert an image",
		})
	}

	finalImage, errImage := bimage.Process(bimg.Options{
		Width:   100,
		Height:  100,
		Crop:    true,
		Quality: 95,
		Type:    bimg.WEBP,
	})
	thumbnail, errThumbnail := bimage.Process(bimg.Options{
		Width:   50,
		Height:  50,
		Crop:    true,
		Quality: 100,
		Type:    bimg.WEBP,
	})
	if errImage != nil || errThumbnail != nil {
		return showErrors(utils.Ms{
			"kind": "Error processing image",
		})
	}

	var ip pgtype.Inet
	ip.Set(c.Get("ip").(net.IP))

	profile := model.Profile{
		Name:        form.Name,
		Description: form.Description,
		Tags:        model.StrToTags(strings.ToLower(form.Tags)),
		Ip:          ip,
		Image:       finalImage,
		Thumbnail:   thumbnail,
	}

	profile.Insert()

	return c.Redirect(http.StatusSeeOther, "/")
}

func DeleteProfile(c echo.Context) (err error) {
	admin := c.Get("admin").(*model.Admin)
	if admin == nil {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	type Form struct {
		ID uint `form:"id" validate:"required"`
	}

	form := Form{}

	if err = c.Bind(&form); err != nil {
		return
	}

	if err := c.Validate(&form); err != nil {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	profile := model.GetProfile(form.ID)
	if profile == nil {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	profile.Delete()

	return c.Redirect(http.StatusSeeOther, "/")
}

func BanIP(c echo.Context) (err error) {
	admin := c.Get("admin").(*model.Admin)
	if admin == nil || admin.Role == model.Helper {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	type Form struct {
		ID uint `form:"id" validate:"required"`
	}

	form := Form{}

	if err = c.Bind(&form); err != nil {
		return
	}

	if err := c.Validate(&form); err != nil {
		return c.Redirect(http.StatusSeeOther, "/")
	}
	profile := model.GetProfile(form.ID)
	if profile == nil {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	profile.Delete()
	model.BanIP(profile.Ip.IPNet.IP)

	return c.Redirect(http.StatusSeeOther, "/")
}
