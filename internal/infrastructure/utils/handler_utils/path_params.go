package handler_utils

import (
	"errors"
	"fmt"
	appErrors "messenger/internal/infrastructure/errors"
	"strconv"
)

func paramNotProvided(name string) appErrors.ResponseMessage {
	return appErrors.ResponseMessage(
		fmt.Sprintf("required param [%s] not provided", name),
	)
}

func failedToCastToInt(value any) appErrors.ResponseMessage {
	return appErrors.ResponseMessage(
		fmt.Sprintf("failed to cast to int [%v]", value),
	)
}

type PathParams map[string]string

func (pathParams PathParams) GetString(key string) (string, error) {
	value, ok := pathParams[key]
	if !ok {
		return "", appErrors.Wrap(
			requiredParamNotProvided(key),
			errors.New("GetString"),
		)
	}
	return value, nil
}

func (pathParams PathParams) GetInteger(key string) (int, error) {
	strValue := pathParams[key]
	if len(strValue) == 0 {
		return 0, appErrors.Wrap(
			requiredParamNotProvided(key),
			errors.New("GetInteger"),
		)
	}
	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		return 0, appErrors.Wrap(
			appErrors.ErrInternal,
			err,
			errors.New("GetInteger"),
		)
	}
	return intValue, nil
}

func requiredParamNotProvided(param string) error {
	return appErrors.ErrBadRequestWithMessage(
		fmt.Sprintf("required path param not provided: %s", param),
	)
}
