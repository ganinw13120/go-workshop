package helper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/labstack/echo/v4"
)

func EchoBindErrorTranslator(err error) string {
	fmt.Println(err)
	switch t := err.(*echo.HTTPError).Unwrap().(type) {
	case *json.UnmarshalTypeError:
		switch t.Type.Kind() {
		case reflect.String:
			return fmt.Sprintf("%s must be a string", t.Field)
		case reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
			return fmt.Sprintf("%s must be a number", t.Field)
		default:
			return fmt.Sprintf("%s must be a %s", t.Field, t.Type.Name())
		}
	case *time.ParseError:
		return fmt.Sprintf("%s is not a valid time format", t.Value)
	default:
		return err.Error()
	}
}
