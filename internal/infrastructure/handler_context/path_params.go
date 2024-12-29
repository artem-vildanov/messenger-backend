package handler_context

import (
	"messenger/internal/app/errors"
	"strconv"
)


type PathParams map[string]string

func (pathParams PathParams) GetInteger(key string) (int, *errors.Error) {
	strValue := pathParams[key]
	if len(strValue) == 0 {
		return 0, errors.NotFoundError()
	}
	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		return 0, errors.NotFoundError()
	}
	return intValue, nil
}
