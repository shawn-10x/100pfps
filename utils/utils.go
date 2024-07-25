package utils

import (
	"image"
	_ "image/jpeg"
	png "image/png"
	"io"
	"os"

	"github.com/h2non/bimg"
	"github.com/labstack/echo/v4"
)

type M map[string]any
type Ms map[string]string

// func BindAndValidate(c echo.Context, obj any) error {
// 	if err := c.Bind(obj); err != nil {
// 		return err
// 	}
// 	if err := c.Validate(obj); err != nil {
// 		switch err.(type) {
// 		case *validator.InvalidValidationError:
// 			return errors.New("invalid form body")
// 		case validator.ValidationErrors:
// 			return ValidationErrors(err.(validator.ValidationErrors))
// 		}
// 		return err
// 	}

// 	return nil
// }

func ReadImage(c echo.Context, name string) (img *bimg.Image, err error) {
	file, err := c.FormFile(name)
	if err != nil {
		return nil, err
	}
	fd, err := file.Open()
	if err != nil {
		return
	}
	defer fd.Close()

	data, err := io.ReadAll(fd)
	if err != nil {
		return
	}
	img = bimg.NewImage(data)

	return img, nil
}

func WriteImage(img image.Image, name string) error {
	fd, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fd.Close()

	return png.Encode(fd, img)
}
