package validation

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"sync"

	"github.com/arfan21/project-sprint-shopifyx-api/pkg/constant"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/shopspring/decimal"
)

var (
	uni        *ut.UniversalTranslator
	validate   *validator.Validate
	translator ut.Translator

	uniOnce        sync.Once
	validateOnce   sync.Once
	translatorOnce sync.Once
)

func Validate[T any](modelValidate T) error {
	uniOnce.Do(func() {
		en := en.New()
		uni = ut.New(en, en)
	})

	validateOnce.Do(func() {
		validate = validator.New()

		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
			if valuer, ok := field.Interface().(decimal.Decimal); ok {
				return valuer.String()
			}
			return nil
		}, decimal.Decimal{})

		validate.RegisterValidation("stringuint", func(fl validator.FieldLevel) bool {
			if fl.Field().Kind() == reflect.String {
				val, err := decimal.NewFromString(fl.Field().String())
				if err != nil {
					return false
				}

				if val.LessThan(decimal.NewFromInt(0)) {
					return false
				}

				return true
			}
			return true
		})

		validate.RegisterValidation("dgte", func(fl validator.FieldLevel) bool {
			if fl.Field().Kind() == reflect.String {
				val, err := decimal.NewFromString(fl.Field().String())
				if err != nil {
					return false
				}

				param := fl.Param()
				paramVal, err := decimal.NewFromString(param)
				if err != nil {
					return false
				}

				if val.LessThan(paramVal) {
					return false
				}

				return true
			}
			return true
		})

		validate.RegisterValidation("customurl", func(fl validator.FieldLevel) bool {
			if fl.Field().Kind() == reflect.String {
				val := fl.Field().String()
				if val == "" {
					return true
				}

				u, err := url.Parse(val)
				if err != nil {
					fmt.Println("error parse url", err)
					return false
				}

				if !strings.Contains(u.Host, ".") {
					fmt.Println("error has prefix url", err)
					return false
				}

				return true
			}
			return true
		})
	})

	translatorOnce.Do(func() {
		translatorUni, _ := uni.GetTranslator("en")
		translator = translatorUni
		en_translations.RegisterDefaultTranslations(validate, translator)

		addTranslation("stringuint", "{0} must be a positive number")
		addTranslation("dgte", "{0} must be greater than or equal to {1}")
		addTranslation("customurl", "{0} must be a valid url")
	})

	err := validate.Struct(modelValidate)
	if err != nil {
		var messages []map[string]interface{}

		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.Field()

			messages = append(messages, map[string]interface{}{
				"field":   fieldName,
				"message": err.Translate(translator),
			})
		}
		jsonMessage, errJson := json.Marshal(messages)
		if errJson != nil {
			return errJson
		}

		return &constant.ErrValidation{Message: string(jsonMessage)}
	}

	return nil
}

func addTranslation(tag string, errMessage string) {
	registerFn := func(ut ut.Translator) error {
		return ut.Add(tag, errMessage, false)
	}

	transFn := func(ut ut.Translator, fe validator.FieldError) string {
		param := fe.Param()
		tag := fe.Tag()

		t, err := ut.T(tag, fe.Field(), param)
		if err != nil {
			return fe.(error).Error()
		}
		return t
	}

	_ = validate.RegisterTranslation(tag, translator, registerFn, transFn)
}
