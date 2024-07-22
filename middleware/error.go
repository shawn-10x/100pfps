package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleError(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			fmt.Println(err.Error())
			return c.Render(http.StatusInternalServerError, "500.html", nil)
		}

		return nil
	}
}

// func HandleError(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		if err := next(c); err != nil {
// 			var status int
// 			r := utils.M{}
// 			switch err.(type) {
// 			case *validator.InvalidValidationError:
// 				status, r["message"] = http.StatusBadRequest, "invalid JSON body"
// 			case validator.ValidationErrors:
// 				status, r["message"] = http.StatusBadRequest, "invalid JSON data"
// 				r["errors"] = validationErrors(err.(validator.ValidationErrors))
// 			case *echo.HTTPError:
// 				err := err.(*echo.HTTPError)
// 				status, r["errors"] = err.Code, err.Message
// 			case *pgconn.ConnectError:
// 				status, r["message"] = http.StatusInternalServerError, "error connecting to the DB"
// 			case *pgconn.PgError:
// 				status, r["message"] = http.StatusInternalServerError, "error executing the query"
// 			default:
// 				status = http.StatusInternalServerError
// 			}

// 			return c.JSON(status, r)
// 		}
// 		return nil
// 	}
// }

// func validationErrors(v_errors validator.ValidationErrors) map[string]string {
// 	errors := map[string]string{}
// 	for _, err := range v_errors {
// 		var msg string
// 		switch err.Tag() {
// 		case "required":
// 			msg = "Missing field"
// 		case "min":
// 			msg = "Too short"
// 		case "max":
// 			msg = "Too long"
// 		default:
// 			msg = "Invalid"
// 			fmt.Println("namespace", err.Namespace())
// 			fmt.Println("field", err.Field())
// 			fmt.Println("struct namespace", err.StructNamespace())
// 			fmt.Println("struct field", err.StructField())
// 			fmt.Println("tag", err.Tag())
// 			fmt.Println("actual tag", err.ActualTag())
// 			fmt.Println("kind", err.Kind())
// 			fmt.Println("type", err.Type())
// 			fmt.Println("value", err.Value())
// 			fmt.Println("param", err.Param())
// 			fmt.Println()
// 		}
// 		errors[err.Field()] = msg
// 	}
// 	return errors
// }
