package validatormsg

import (
	"bytes"
	"strings"

	"github.com/go-playground/validator"
)

func ValidatorMsg(err error) string {
	errorMsg := ""
	for i, err := range err.(validator.ValidationErrors) {
		if i > 0 {
			errorMsg += ", " + SnakeCase(strings.Split(err.StructNamespace(), ".")[1])
		} else {
			errorMsg += ": " + SnakeCase(strings.Split(err.StructNamespace(), ".")[1])
		}
	}
	return errorMsg
}

func SnakeCase(camel string) string {
	var buf bytes.Buffer
	for _, c := range camel {
		if 'A' <= c && c <= 'Z' {
			if buf.Len() > 0 {
				buf.WriteRune('_')
			}
			buf.WriteRune(c - 'A' + 'a')
		} else {
			buf.WriteRune(c)
		}
	}
	return buf.String()
}
