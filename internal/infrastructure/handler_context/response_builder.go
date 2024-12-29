package handler_context

import (
	"encoding/json"
	"fmt"
	"log"
	"messenger/internal/app/errors"
	"net/http"
)

type ResponseBuilder struct {
	responseWriter http.ResponseWriter
	code           int
	content        any
}

func (r *ResponseBuilder) WithCode(code int) *ResponseBuilder {
	r.code = code
	return r
}

func (r *ResponseBuilder) WithContent(content any) *ResponseBuilder {
	r.content = content
	return r
}

func (r *ResponseBuilder) WithHeader(key string, value string) *ResponseBuilder {
	r.responseWriter.Header().Set(key, value)
	return r
}

func (r *ResponseBuilder) WithCookie(sessionId string) *ResponseBuilder {
	r.responseWriter.Header().Set(
		"Set-Cookie",
		fmt.Sprintf("sessionId=%s; HttpOnly; Path=/; SameSite=Strict", sessionId),
	)
	return r
}

func (r *ResponseBuilder) Json() *errors.Error {
	r.responseWriter.Header().Set("Content-Type", "application/json")
	r.content = map[string]interface{}{"message": r.content}
	r.responseWriter.WriteHeader(r.code)
	err := json.NewEncoder(r.responseWriter).Encode(r.content)
	if err != nil {
		log.Printf("Failed to parse into json with error: %s\n", err.Error())
		return errors.InternalError()
	}
	return nil
}

func (r *ResponseBuilder) Empty() *errors.Error {
	r.WithCode(http.StatusNoContent)
	r.responseWriter.Header().Del("Content-Type")
	r.responseWriter.WriteHeader(r.code)
	if _, err := r.responseWriter.Write(nil); err != nil {
		log.Printf("Failed to send empty response: %s\n", err.Error())
		return errors.InternalError()
	}
	return nil
}

func (r *ResponseBuilder) HTML() *errors.Error {
	return r.respondWithText("html")
}

func (r *ResponseBuilder) TextPlain() *errors.Error {
	return r.respondWithText("plain")
}

func (r *ResponseBuilder) respondWithText(textType string) *errors.Error {
	contentType := fmt.Sprintf("text/%s", textType)
	r.responseWriter.Header().Set(
		"Content-Type",
		fmt.Sprintf("%s; charset=utf-8", contentType),
	)

	if r.content == nil {
		r.content = ""
	}

	htmlContent, ok := r.content.(string)
	if !ok {
		return errors.InternalError().
			WithVerbose(fmt.Sprintf("Invalid content type for %s\n", contentType))
	}

	r.responseWriter.WriteHeader(r.code)
	_, err := r.responseWriter.Write([]byte(htmlContent))
	if err != nil {
		return errors.InternalError().
			WithVerbose(fmt.Sprintf("Error writing %s response: %s", contentType, err.Error()))
	}

	return nil
}
