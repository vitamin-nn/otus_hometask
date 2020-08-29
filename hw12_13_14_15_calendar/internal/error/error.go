package error

type OutError string

func (e OutError) Error() string {
	return string(e)
}

var (
	ErrEventNotFound = OutError("event not found")
	ErrDateBusy      = OutError("time is busy")
	ErrUnknownUserID = OutError("unknown user_id")
	ErrEmptyEvent    = OutError("empty event param")
	ErrInternal      = OutError("internal server error")
	ErrInvalidParams = OutError("invalid params")
)
