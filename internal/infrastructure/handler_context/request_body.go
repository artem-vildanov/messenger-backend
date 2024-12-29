package handler_context

import (
	"encoding/json"
	"fmt"
	"io"
	"messenger/internal/app/errors"
	"net/http"
)

type requestModel interface {
	FromRequest(body RequestBody) *errors.Error
}

type RequestBody map[string]interface{}

func (r RequestBody) GetString(key string) (string, *errors.Error) {
	value := r[key]
	valueCasted, ok := value.(string)
	if !ok {
		return "", errors.NewValidationError().
			WithFieldName(key).
			WithReason("failed to cast into string").
			BuildError()
	}
	return valueCasted, nil
}

func (r RequestBody) FillObject(key string, object requestModel) *errors.Error {
	value := r[key]
	valueCasted, ok := value.(map[string]interface{})
	if !ok {
		return errors.NewValidationError().
			WithFieldName(key).
			WithReason("failed to cast into object").
			BuildError()
	}

	if err := object.FromRequest(valueCasted); err != nil {
		return err
	}

	return nil
}

func (r *RequestBody) fromRequest(request *http.Request) *errors.Error {
	if request.ContentLength == 0 {
		return nil
	}

	if request.Method == "OPTIONS" {
		return nil
	}

	defer request.Body.Close()
	rawBody, err := io.ReadAll(request.Body)
	if err != nil {
		return errors.InternalError().WithVerbose(fmt.Sprintf("Error reading body: %s", err))
	}

	err = json.Unmarshal(rawBody, r)
	if err != nil {
		return errors.InternalError().WithVerbose(fmt.Sprintf("Error parsing JSON: %s", err))
	}

	return nil
}

func (c *HandlerContext) ErrorResponse(err *errors.Error) {
	c.Response().
		WithCode(err.GetCode()).
		WithContent(err.Error()).
		Json()
}

func (c *HandlerContext) Response() *ResponseBuilder {
	c.responseWriter.Header().Set("Content-Type", "application/json")
	return &ResponseBuilder{
		responseWriter: c.responseWriter,
		code:           http.StatusOK,
		content:        "OK",
	}
}
