package validators

import (
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidateTags(fl validator.FieldLevel) bool {
	v := fl.Field().String()

	tags := strings.Split(v, " ")

	for _, tag := range tags {
		if !strings.HasPrefix(tag, "#") {
			return false
		}
	}

	return true
}

func TagsMaxCount(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	param_ntags, e := strconv.Atoi(fl.Param())
	if e != nil {
		panic(e)
	}

	tags := strings.Split(v, " ")

	if len(tags) > param_ntags {
		return false
	}

	return true
}

func TagsLength(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	length, e := strconv.Atoi(fl.Param())
	if e != nil {
		panic(e)
	}

	tags := strings.Split(v, " ")
	for _, tag := range tags {
		if len(tag) > length {
			return false
		}
	}

	return true
}
