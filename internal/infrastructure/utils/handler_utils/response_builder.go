package handler_utils

import (
	"encoding/json"
	"fmt"
	"log"
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
		fmt.Sprintf("sessionId=%s; HttpOnly; Path=/; SameSite=Lax", sessionId),
	)
	return r
}

func (r *ResponseBuilder) Json() {
	r.responseWriter.Header().Set("Content-Type", "application/json")
	r.content = map[string]interface{}{"message": r.content}
	r.responseWriter.WriteHeader(r.code)
	err := json.NewEncoder(r.responseWriter).Encode(r.content)
	if err != nil {
		log.Printf("Failed to parse into json with error: %s\n", err.Error())
	}
}

func (r *ResponseBuilder) Empty() {
	r.WithCode(http.StatusNoContent)
	r.responseWriter.Header().Del("Content-Type")
	r.responseWriter.WriteHeader(r.code)
	if _, err := r.responseWriter.Write(nil); err != nil {
		log.Printf("Failed to send empty response: %s\n", err.Error())
	}
}

func (r *ResponseBuilder) HTML() {
	r.respondWithText("html")
}

func (r *ResponseBuilder) TextPlain() {
	r.respondWithText("plain")
}

func (r *ResponseBuilder) respondWithText(textType string) {
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
		log.Printf("Invalid content type for %s\n", contentType)
		return
	}

	r.responseWriter.WriteHeader(r.code)
	_, err := r.responseWriter.Write([]byte(htmlContent))
	if err != nil {
		log.Printf("Error writing %s response: %s", contentType, err.Error())
		return
	}
}
