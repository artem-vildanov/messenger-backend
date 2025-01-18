package mapping_utils

import (
	"encoding/json"
	"fmt"
	"messenger/internal/infrastructure/errors"
	"net/http"
	"reflect"

	"github.com/go-playground/validator"
)

func ValidateRequestModel(model any) *errors.Error {
	validate := validator.New()
	if validationError := validate.Struct(model); validationError != nil {
		resultError := errors.BadRequestError().
			WithResponseMessage("validation error").
			WithLogMessage(validationError.Error(), "got validation error")
		return resultError
	}
	return nil
}

func FromRequest[T any](request *http.Request) (T, *errors.Error) {
	model := new(T)
	if err := json.NewDecoder(request.Body).Decode(&model); err != nil {
		return *model, errors.BadRequestError().
			WithResponseMessage("failed to decode json from request").
			WithLogMessage(err.Error())
	}
	if err := ValidateRequestModel(*model); err != nil {
		return *model, err
	}
	return *model, nil
}

func FromJsonString[T any](jsonStr string) (T, *errors.Error) {
	parsed := new(T)
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return *parsed, errors.BadRequestError().
			WithResponseMessage("failed to decode json from string").
			WithLogMessage(err.Error())
	}
	return *parsed, nil
}

func ToJsonString[T any](object T) (string, *errors.Error) {
	jsonStr, err := json.Marshal(object)
	if err != nil {
		structName := reflect.TypeOf(object).Name()
		return "", errors.InternalError().
			WithLogMessage(
				err.Error(),
				fmt.Sprintf("failed to encode object [%s] into string", structName),
			)
	}
	return string(jsonStr), nil
}
