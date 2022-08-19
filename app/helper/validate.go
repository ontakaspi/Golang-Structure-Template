package helper

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en2 "github.com/go-playground/validator/v10/translations/en"
	"github.com/thedevsaddam/govalidator"
	"net/http"
	"net/url"
)

var Translator ut.Translator

func BindAndValidate(c *gin.Context, request interface{}, validate *validator.Validate) (errMessages SimplifiedError) {
	english := en.New()
	uni := ut.New(english, english)
	Translator, _ = uni.GetTranslator("en")
	_ = en2.RegisterDefaultTranslations(validate, Translator)

	errBinding := c.BindJSON(&request)
	if errBinding != nil {
		errMessages = SimplifyError(errBinding)
		return errMessages
	}

	errValidation := validate.Struct(request)
	if errValidation != nil {
		errMessages = SimplifyError(errValidation)
		return errMessages
	}

	return nil

}

type CreateMisconfigFile struct {
	Rules   map[string][]string
	Message map[string][]string
}

func ValidateFormData(c *http.Request, rules map[string][]string, messages map[string][]string) (errMessages url.Values) {
	opts := govalidator.Options{
		Request:         c,
		Rules:           rules,
		Messages:        messages,
		RequiredDefault: true,
	}
	v := govalidator.New(opts)
	e := v.Validate()
	if len(e) > 0 {
		errMessages = e
		return errMessages
	}
	return nil

}
