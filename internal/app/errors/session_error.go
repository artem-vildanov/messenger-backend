package errors

import "fmt"

type SessionError struct {
	*Error
	errorReason string
	sessionId string
	userId int
}

func (e *SessionError) WithId(id string) *SessionError {
	e.sessionId = id
	return e
}

func (e *SessionError) WithReason(reason string) *SessionError {
	e.errorReason = reason
	return e
}

func (e *SessionError) WithUserId(userId int) *SessionError {
	e.userId = userId
	return e
}

func (e *SessionError) BuildError() *Error {
	if e.userId != 0 {
		e.WithVerbose(fmt.Sprintf("with userId: [%d]", e.userId))
	}
	if len(e.sessionId) > 0 {
		e.WithVerbose(fmt.Sprintf("with sessionId: [%s]", e.sessionId))
	}
	if len(e.errorReason) > 0 {
		e.WithVerbose(fmt.Sprintf("for reason: [%s]", e.errorReason))
	}
	return e.Error
}

func FailedToCreateSession() *SessionError {
	return &SessionError{
		Error: InternalError().WithVerbose("failed to create session"),
	}
}

func FailedToFindSession() *SessionError {
	return &SessionError{
		Error: InternalError().WithVerbose("failed to find session"),
	}
}

func FailedToDeleteSession() *SessionError {
	return &SessionError{
		Error: InternalError().WithVerbose("failed to delete session"),
	}
}

func SessionNotFoundError() *SessionError {
	return &SessionError{
		Error: ForbiddenError().WithVerbose("session not found"),
	}
}

func SessionExpiredError() *SessionError {
	return &SessionError{
		Error: ForbiddenError().WithVerbose("session expired"),
	}
}
