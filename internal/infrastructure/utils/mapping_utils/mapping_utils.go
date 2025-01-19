package mapping_utils

import (
	"encoding/json"
	"errors"
	"log"
	appErrors "messenger/internal/infrastructure/errors"
	"net/http"

	"github.com/go-playground/validator"
)

func ValidateRequestModel(model any) error {
	validate := validator.New()
	if err := validate.Struct(model); err != nil {
		log.Println(err.Error())
		return appErrors.Wrap(
			appErrors.ErrBadRequestWithMessage("validation error"),
			err,
			errors.New("ValidateRequestModel"),
		)
	}
	return nil
}

func FromRequest[T any](request *http.Request) (T, error) {
	model := new(T)
	if err := json.NewDecoder(request.Body).Decode(&model); err != nil {
		return *model, appErrors.Wrap(
			appErrors.ErrBadRequest,
			err,
			errors.New("FromRequest"),
		)
	}
	if err := ValidateRequestModel(*model); err != nil {
		return *model, err
	}
	return *model, nil
}

func FromJsonString[T any](jsonStr string) (T, error) {
	parsed := new(T)
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return *parsed, appErrors.Wrap(
			appErrors.ErrBadRequest,
			err,
			errors.New("FromRequest"),
		)
	}
	return *parsed, nil
}

func ToJsonString(object any) (string, error) {
	jsonStr, err := json.Marshal(object)
	if err != nil {
		return "", appErrors.Wrap(
			appErrors.ErrInternal,
			err,
			errors.New("ToJsonString"),
		)
	}
	return string(jsonStr), nil
}
