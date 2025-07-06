package validator

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func MustRegistValidator() {
	var validatorMap = map[string]func(fl validator.FieldLevel) bool{
		"trim_gte": func(fl validator.FieldLevel) bool {
			fieldValue := fl.Field().String()
			paramValue := fl.Param()

			paramInt, err := strconv.Atoi(paramValue)
			if err != nil {
				return false
			}
			return len(strings.TrimSpace(fieldValue)) >= paramInt
		},
		"trim_gt": func(fl validator.FieldLevel) bool {
			fieldValue := fl.Field().String()
			paramValue := fl.Param()

			paramInt, err := strconv.Atoi(paramValue)
			if err != nil {
				return false
			}
			return len(strings.TrimSpace(fieldValue)) > paramInt
		},
		"trim_lte": func(fl validator.FieldLevel) bool {
			fieldValue := fl.Field().String()
			paramValue := fl.Param()

			paramInt, err := strconv.Atoi(paramValue)
			if err != nil {
				return false
			}
			return len(strings.TrimSpace(fieldValue)) <= paramInt
		},
		"trim_lt": func(fl validator.FieldLevel) bool {
			fieldValue := fl.Field().String()
			paramValue := fl.Param()

			paramInt, err := strconv.Atoi(paramValue)
			if err != nil {
				return false
			}
			return len(strings.TrimSpace(fieldValue)) < paramInt
		},
		"trim_no_empty": func(fl validator.FieldLevel) bool {
			fieldValue := fl.Field().String()
			return len(strings.TrimSpace(fieldValue)) > 0
		},
		// todo 支持[]string
		"no_spacing": func(fl validator.FieldLevel) bool {
			fieldValue := fl.Field().String()
			return !strings.Contains(fieldValue, " ")
		},
		"no_lt_spacing": func(fl validator.FieldLevel) bool {
			fieldValue := fl.Field().String()
			fl.Field()
			return !(strings.HasPrefix(fieldValue, " ") || strings.HasSuffix(fieldValue, " "))
		},
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		for name, validatorFunc := range validatorMap {
			v.RegisterValidation(name, validatorFunc)
		}
	}
}
