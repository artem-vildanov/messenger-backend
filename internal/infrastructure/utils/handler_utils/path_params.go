package handler_utils

import (
	"fmt"
	"messenger/internal/infrastructure/errors"
	"strconv"
)

func paramNotProvided(name string) errors.ResponseMessage {
	return errors.ResponseMessage(
		fmt.Sprintf("required param [%s] not provided", name),
	)
}

func failedToCastToInt(value any) errors.ResponseMessage {
	return errors.ResponseMessage(
		fmt.Sprintf("failed to cast to int [%v]", value),
	)
}

type PathParams map[string]string

func (pathParams PathParams) GetString(key string) (string, *errors.Error) {
	value, ok := pathParams[key]
	if !ok {
		return "", errors.NotFoundError().
			WithResponseMessage(paramNotProvided(key))
	}
	return value, nil
}

func (pathParams PathParams) GetInteger(key string) (int, *errors.Error) {
	strValue := pathParams[key]
	if len(strValue) == 0 {
		return 0, errors.NotFoundError().
			WithResponseMessage(paramNotProvided(key))
	}
	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		return 0, errors.InternalError().
			WithLogMessage(string(failedToCastToInt(strValue)))
	}
	return intValue, nil
}
