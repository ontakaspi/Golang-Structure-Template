package helper

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type SimplifiedError []string

func SimplifyError(errors error) (simplifiedError SimplifiedError) {

	switch errors.(type) {
	case validator.ValidationErrors:
		for _, f := range errors.(validator.ValidationErrors) {

			translatedErr := fmt.Errorf(f.Translate(Translator))

			simplifiedError = append(simplifiedError, "Field "+ToSnakeCase(translatedErr.Error()))
		}

		return simplifiedError

	case *json.UnmarshalTypeError:
		errs := errors.(*json.UnmarshalTypeError)
		err := "Field " + errs.Field + " must be type of " + errs.Type.String()
		simplifiedError = append(simplifiedError, err)
		return simplifiedError

	default:
		err := fmt.Sprintf(errors.Error())
		simplifiedError = append(simplifiedError, err)
		return simplifiedError

	}

}
